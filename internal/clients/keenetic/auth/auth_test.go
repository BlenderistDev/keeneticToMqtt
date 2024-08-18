package auth

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"keeneticToMqtt/internal/errs"
	mock_auth "keeneticToMqtt/test/mocks/gomock/clients/keenetic/auth"
)

func TestAuth_RefreshAuth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const (
		host      = "host"
		challenge = "challenge"
		ndmRealm  = "realm"
		login     = "login"
		password  = "password"
	)

	someErr := errors.New("some error")

	tests := []struct {
		name                      string
		validateCheckAuthRequest  func(req *http.Request)
		validateAuthRequest       func(req *http.Request)
		getCheckAuthResponse      func() *http.Response
		getAuthResponse           func() *http.Response
		getCheckAuthResponseError func() error
		getAuthResponseError      func() error
		expectedErr               error
		expectedErrStr            string
	}{
		{
			name: "401 check auth, then success auth",
			validateCheckAuthRequest: func(req *http.Request) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, host+authUrl, req.URL.String())
			},
			getCheckAuthResponse: func() *http.Response {
				resp := http.Response{
					StatusCode: http.StatusUnauthorized,
					Body:       io.NopCloser(strings.NewReader("")),
				}
				resp.Header = make(http.Header)
				resp.Header.Add(ndmChallengeHeader, challenge)
				resp.Header.Add(ndmRealmHeader, ndmRealm)
				return &resp
			},
			getCheckAuthResponseError: func() error {
				return nil
			},
			validateAuthRequest: func(req *http.Request) {
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, host+authUrl, req.URL.String())

				b, err := io.ReadAll(req.Body)
				assert.Nil(t, err)
				var body authRequest
				err = json.Unmarshal(b, &body)
				assert.Nil(t, err)
				assert.Equal(t, body.Login, login)
				assert.Equal(t, body.Password, "c4889195dd1015de53272fdca72714374b1d7f19b74698857bcecac36c433c70")
			},
			getAuthResponse: func() *http.Response {
				resp := http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader("")),
				}
				return &resp
			},
			getAuthResponseError: func() error {
				return nil
			},
		},
		{
			name: "200 check auth",
			validateCheckAuthRequest: func(req *http.Request) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, host+authUrl, req.URL.String())
			},
			getCheckAuthResponse: func() *http.Response {
				resp := http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader("")),
				}
				return &resp
			},
			getCheckAuthResponseError: func() error {
				return nil
			},
			validateAuthRequest: func(req *http.Request) {},
			getAuthResponse: func() *http.Response {
				return nil
			},
			getAuthResponseError: func() error {
				return nil
			},
		},
		{
			name:                     "check auth error",
			validateCheckAuthRequest: func(req *http.Request) {},
			getCheckAuthResponse: func() *http.Response {
				return nil
			},
			getCheckAuthResponseError: func() error {
				return someErr
			},
			validateAuthRequest: func(req *http.Request) {},
			getAuthResponse: func() *http.Response {
				return nil
			},
			getAuthResponseError: func() error {
				return nil
			},
			expectedErr: someErr,
		},
		{
			name:                     "401 check auth, then auth error",
			validateCheckAuthRequest: func(req *http.Request) {},
			getCheckAuthResponse: func() *http.Response {
				resp := http.Response{
					StatusCode: http.StatusUnauthorized,
					Body:       io.NopCloser(strings.NewReader("")),
				}
				resp.Header = make(http.Header)
				resp.Header.Add(ndmChallengeHeader, challenge)
				resp.Header.Add(ndmRealmHeader, ndmRealm)
				return &resp
			},
			getCheckAuthResponseError: func() error {
				return nil
			},
			validateAuthRequest: func(req *http.Request) {},
			getAuthResponse: func() *http.Response {
				return nil
			},
			getAuthResponseError: func() error {
				return someErr
			},
			expectedErr: someErr,
		},
		{
			name: "check auth wrong status code",
			validateCheckAuthRequest: func(req *http.Request) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, host+authUrl, req.URL.String())
			},
			getCheckAuthResponse: func() *http.Response {
				resp := http.Response{
					StatusCode: http.StatusBadRequest,
					Body:       io.NopCloser(strings.NewReader("")),
				}
				resp.Header = make(http.Header)
				return &resp
			},
			getCheckAuthResponseError: func() error {
				return nil
			},
			validateAuthRequest: func(req *http.Request) {},
			getAuthResponse: func() *http.Response {
				return nil
			},
			getAuthResponseError: func() error {
				return nil
			},
			expectedErrStr: "error while keenetic check auth: error in checkauth request, status code: 400",
		},
		{
			name: "401 check auth, then 401 auth",
			validateCheckAuthRequest: func(req *http.Request) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, host+authUrl, req.URL.String())
			},
			getCheckAuthResponse: func() *http.Response {
				resp := http.Response{
					StatusCode: http.StatusUnauthorized,
					Body:       io.NopCloser(strings.NewReader("")),
				}
				resp.Header = make(http.Header)
				resp.Header.Add(ndmChallengeHeader, challenge)
				resp.Header.Add(ndmRealmHeader, ndmRealm)
				return &resp
			},
			getCheckAuthResponseError: func() error {
				return nil
			},
			validateAuthRequest: func(req *http.Request) {},
			getAuthResponse: func() *http.Response {
				resp := http.Response{
					StatusCode: http.StatusUnauthorized,
					Body:       io.NopCloser(strings.NewReader("")),
				}
				return &resp
			},
			getAuthResponseError: func() error {
				return nil
			},
			expectedErr: errs.ErrUnauthorized,
		},
		{
			name: "401 check auth, then wrong status auth",
			validateCheckAuthRequest: func(req *http.Request) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, host+authUrl, req.URL.String())
			},
			getCheckAuthResponse: func() *http.Response {
				resp := http.Response{
					StatusCode: http.StatusUnauthorized,
					Body:       io.NopCloser(strings.NewReader("")),
				}
				resp.Header = make(http.Header)
				resp.Header.Add(ndmChallengeHeader, challenge)
				resp.Header.Add(ndmRealmHeader, ndmRealm)
				return &resp
			},
			getCheckAuthResponseError: func() error {
				return nil
			},
			validateAuthRequest: func(req *http.Request) {},
			getAuthResponse: func() *http.Response {
				resp := http.Response{
					StatusCode: http.StatusBadRequest,
					Body:       io.NopCloser(strings.NewReader("")),
				}
				return &resp
			},
			getAuthResponseError: func() error {
				return nil
			},
			expectedErrStr: "error while keenetic auth: error in checkauth request, status code: 400",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := mock_auth.NewMockclient(ctrl)
			client.EXPECT().Do(gomock.Cond(func(x any) bool {
				req, ok := x.(*http.Request)
				if !ok || req == nil {
					t.Errorf("empty request")
					return false
				}
				tt.validateCheckAuthRequest(req)
				return true
			})).Return(tt.getCheckAuthResponse(), tt.getCheckAuthResponseError())

			if tt.getAuthResponse() != nil || tt.getAuthResponseError() != nil {
				client.EXPECT().Do(gomock.Cond(func(x any) bool {
					req, ok := x.(*http.Request)
					if !ok || req == nil {
						t.Errorf("empty request")
						return false
					}
					tt.validateAuthRequest(req)
					return true
				})).Return(tt.getAuthResponse(), tt.getAuthResponseError())
			}

			cookie, _ := cookiejar.New(&cookiejar.Options{})
			auth := NewAuth(host, login, password, cookie)
			auth.client = client
			err := auth.RefreshAuth()
			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr)
			} else if tt.expectedErrStr != "" {
				assert.Regexp(t, tt.expectedErrStr+".*", err.Error())
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
