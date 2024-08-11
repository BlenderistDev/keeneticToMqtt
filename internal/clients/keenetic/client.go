package keenetic

import (
	"net/http"
)

type Keenetic struct {
	host, login, password string
	client                *http.Client
}

func NewKeenetic(host, login, password string) *Keenetic {
	keenetic := &Keenetic{
		host:     host,
		login:    login,
		password: password,
	}

	cookieStorage := &cookieStorage{}
	t := &http.Transport{}

	var rt http.RoundTripper
	rt = cookieRoundTripper{
		proxied:       t,
		cookieStorage: cookieStorage,
	}

	client := &http.Client{
		Transport: rt,
	}

	keenetic.client = client
	return keenetic
}

func (k *Keenetic) Do(req *http.Request) (*http.Response, error) {
	return k.client.Do(req)
}

func (k *Keenetic) GetHost() string {
	return k.host
}
