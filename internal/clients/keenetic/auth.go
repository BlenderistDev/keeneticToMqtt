package keenetic

import (
	"net/http"
)

//go:generate mockgen -source=auth.go -destination=../../../test/mocks/gomock/clients/keenetic/auth.go

type roundTripper interface {
	RoundTrip(*http.Request) (*http.Response, error)
}

type authClient interface {
	RefreshAuth() error
}

type authRoundTripper struct {
	proxied roundTripper
	auth    authClient
}

// RoundTrip checks and refresh auth before requests.
func (rt *authRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if err := rt.auth.RefreshAuth(); err != nil {
		return nil, err
	}
	res, err := rt.proxied.RoundTrip(req)

	return res, err
}
