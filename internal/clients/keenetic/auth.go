package keenetic

import (
	"net/http"

	"keeneticToMqtt/internal/clients/keenetic/auth"
)

type roundTripper interface {
	RoundTrip(*http.Request) (*http.Response, error)
}

type authRoundTripper struct {
	proxied roundTripper
	auth    *auth.Auth
}

func (rt *authRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if err := rt.auth.RefreshAuth(); err != nil {
		return nil, err
	}
	res, err := rt.proxied.RoundTrip(req)

	return res, err
}
