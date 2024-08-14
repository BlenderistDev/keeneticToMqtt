package keenetic

import (
	"log/slog"
	"net/http"
	"net/http/cookiejar"

	"keeneticToMqtt/internal/clients/keenetic/auth"
)

type Keenetic struct {
	host, login, password string
	client                *http.Client
}

func NewKeenetic(auth *auth.Auth, cookiejar *cookiejar.Jar, host, login, password string, log *slog.Logger) *Keenetic {
	keenetic := &Keenetic{
		host:     host,
		login:    login,
		password: password,
	}

	t := &http.Transport{}

	var rt http.RoundTripper
	rt = &AuthRoundTripper{
		proxied: t,
		auth:    auth,
	}
	//rt = &logger.RoundTripper{
	//	Proxied:    rt,
	//	Log:        log,
	//	ClientName: "keenetic",
	//}

	client := &http.Client{
		Transport: rt,
		Jar:       cookiejar,
	}

	keenetic.client = client
	return keenetic
}

func (k *Keenetic) Do(req *http.Request) (*http.Response, error) {
	return k.client.Do(req)
}
