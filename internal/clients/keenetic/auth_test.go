package keenetic

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	mock_keenetic "keeneticToMqtt/test/mocks/gomock/clients/keenetic"
)

func TestAuthRoundTripper_RoundTrip(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := &http.Request{}
	resp := &http.Response{}
	someErr := errors.New("some err")

	tests := []struct {
		name    string
		auth    func() authClient
		proxied func() roundTripper
		resp    *http.Response
		err     error
	}{
		{
			name: "success",
			auth: func() authClient {
				auth := mock_keenetic.NewMockauthClient(ctrl)
				auth.EXPECT().RefreshAuth().Return(nil)
				return auth
			},
			proxied: func() roundTripper {
				proxied := mock_keenetic.NewMockroundTripper(ctrl)
				proxied.EXPECT().RoundTrip(req).Return(resp, nil)
				return proxied
			},
			resp: resp,
		},
		{
			name: "auth error",
			auth: func() authClient {
				auth := mock_keenetic.NewMockauthClient(ctrl)
				auth.EXPECT().RefreshAuth().Return(someErr)
				return auth
			},
			proxied: func() roundTripper {
				proxied := mock_keenetic.NewMockroundTripper(ctrl)
				return proxied
			},
			err: someErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt := authRoundTripper{
				proxied: tt.proxied(),
				auth:    tt.auth(),
			}

			result, err := rt.RoundTrip(req)
			if tt.err != nil {
				assert.Nil(t, result)
				assert.ErrorIs(t, err, tt.err)
			} else {
				assert.Nil(t, err)
				assert.Same(t, tt.resp, result)
			}

		})
	}
}
