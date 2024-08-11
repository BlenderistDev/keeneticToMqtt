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
)

type client interface {
	Do(req *http.Request) (*http.Response, error)
}

var (
	ErrUnauthorized = errors.New("unauthorized")
)

type Auth struct {
	login, password, host string
	client                client
}

func NewAuth(host, login, password string, client client) *Auth {
	return &Auth{
		login:    login,
		password: password,
		host:     host,
		client:   client,
	}
}

func (a *Auth) CheckAuth() (realm, challenge string, err error) {
	req, err := http.NewRequest(http.MethodGet, a.host+"/auth", nil)
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
		challenge = resp.Header.Get("X-NDM-Challenge")
		realm = resp.Header.Get("X-NDM-Realm")
		err = ErrUnauthorized
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

func (a *Auth) Auth(realm, challenge string) error {
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
		return ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error in checkauth request, status code: %d", resp.StatusCode)
	}

	return nil
}
