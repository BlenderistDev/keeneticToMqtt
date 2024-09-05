package keenetic

import (
	"log/slog"
	"net/http"
	"net/http/cookiejar"

	"keeneticToMqtt/internal/logger"
)

const clientName = "keenetic"

// Keenetic client for keenetic.
type Keenetic struct {
	host, login, password string
	client                *http.Client
}

// NewKeenetic creates new Keenetic.
func NewKeenetic(auth authClient, cookiejar *cookiejar.Jar, host, login, password string, log *slog.Logger) *Keenetic {
	keenetic := &Keenetic{
		host:     host,
		login:    login,
		password: password,
	}

	t := &http.Transport{}

	var rt http.RoundTripper
	rt = &authRoundTripper{
		proxied: t,
		auth:    auth,
	}

	rt = &logger.RoundTripper{
		Proxied:    rt,
		Log:        log,
		ClientName: clientName,
	}

	client := &http.Client{
		Transport: rt,
		Jar:       cookiejar,
	}

	keenetic.client = client
	return keenetic
}

// Do request.
func (k *Keenetic) Do(req *http.Request) (*http.Response, error) {
	return k.client.Do(req)
}
