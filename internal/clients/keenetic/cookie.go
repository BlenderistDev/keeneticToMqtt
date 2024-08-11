package keenetic

import (
	"net/http"
)

type cookieStorage struct {
	cookie []*http.Cookie
}

func (s *cookieStorage) Get() []*http.Cookie {
	return s.cookie
}

func (s *cookieStorage) Set(cookie []*http.Cookie) {
	s.cookie = cookie
}

type cookieRoundTripper struct {
	proxied       http.RoundTripper
	cookieStorage *cookieStorage
}

func (rt cookieRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	for _, cookie := range rt.cookieStorage.Get() {
		req.AddCookie(cookie)
	}
	res, err := rt.proxied.RoundTrip(req)

	if res.Header.Get("Set-Cookie") != "" {
		cookies := make([]*http.Cookie, len(res.Cookies()))
		for i, cookie := range res.Cookies() {
			cookies[i] = cookie
		}
		rt.cookieStorage.Set(cookies)
	}

	return res, err
}
