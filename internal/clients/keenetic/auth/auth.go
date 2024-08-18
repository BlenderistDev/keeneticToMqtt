package auth

import (
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"keeneticToMqtt/internal/errs"
)

//go:generate mockgen -source=auth.go -destination=../../../../test/mocks/gomock/clients/keenetic/auth/auth.go

const (
	authUrl            = "/auth"
	ndmChallengeHeader = "X-NDM-Challenge"
	ndmRealmHeader     = "X-NDM-Realm"
)

type (
	client interface {
		Do(req *http.Request) (*http.Response, error)
	}
)

// Auth struct for keenetic auth.
type Auth struct {
	login, password, host string
	client                client
}

// NewAuth creates new Auth.
func NewAuth(host, login, password string, cookiejar http.CookieJar) *Auth {
	return &Auth{
		login:    login,
		password: password,
		host:     host,
		client:   &http.Client{Jar: cookiejar},
	}
}

// RefreshAuth refresh keenetic auth.
// try GET /auth. If 200 - OK.
// IF 401 - try POST /auth with headers as password salt.
func (a *Auth) RefreshAuth() error {
	realm, challenge, err := a.checkAuth()
	switch {
	case errors.Is(err, errs.ErrUnauthorized):
		if err := a.auth(realm, challenge); err != nil {
			return fmt.Errorf("error while keenetic auth: %w", err)
		}
	case err != nil:
		return fmt.Errorf("error while keenetic check auth: %w", err)
	}

	return nil
}

func (a *Auth) checkAuth() (realm, challenge string, err error) {
	req, err := http.NewRequest(http.MethodGet, a.host+authUrl, nil)
	if err != nil {
		err = fmt.Errorf("build request error in checkauth request: %w", err)
		return
	}

	resp, err := a.client.Do(req)
	if err != nil {
		err = fmt.Errorf("send error in checkauth request: %w", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		challenge = resp.Header.Get(ndmChallengeHeader)
		realm = resp.Header.Get(ndmRealmHeader)
		err = errs.ErrUnauthorized
		return
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("error in checkauth request, status code: %d", resp.StatusCode)
		return
	}

	return
}

type authRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (a *Auth) auth(realm, challenge string) error {
	hashMd5 := md5.Sum([]byte(a.login + ":" + realm + ":" + a.password))
	hashMd5Str := hex.EncodeToString(hashMd5[:])

	hashSha256 := sha256.Sum256([]byte(challenge + hashMd5Str))
	pass := hex.EncodeToString(hashSha256[:])

	body := authRequest{
		Login:    a.login,
		Password: pass,
	}

	b, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("auth body marshal error: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, a.host+"/auth", bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("build request error in checkauth request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json;charset=UTF-8")

	resp, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("send error in checkauth request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return errs.ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error in checkauth request, status code: %d", resp.StatusCode)
	}

	return nil
}
