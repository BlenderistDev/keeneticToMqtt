package keenetic

import (
	"net/http"

	"keeneticToMqtt/internal/clients/keenetic/auth"
)

type AuthRoundTripper struct {
	proxied http.RoundTripper
	auth    *auth.Auth
}

func (k *Keenetic) NewAuthRoundTripper(proxied http.RoundTripper, auth *auth.Auth) http.RoundTripper {
	return &AuthRoundTripper{
		proxied: proxied,
		auth:    auth,
	}
}

func (rt *AuthRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if err := rt.auth.RefreshAuth(); err != nil {
		return nil, err
	}
	res, err := rt.proxied.RoundTrip(req)

	return res, err
}
