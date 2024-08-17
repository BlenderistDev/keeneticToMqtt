package accessupdate

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"keeneticToMqtt/internal/dto/homeassistantdto"
	"keeneticToMqtt/internal/errs"
	mock_accessupdate "keeneticToMqtt/test/mocks/gomock/clients/keenetic/accessupdate"
)

func TestAccessUpdate_SetPermit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const (
		host = "host"
		mac  = "mac"
	)

	successRes := response{
		"key": {},
	}
	someErr := errors.New("some err")

	tests := []struct {
		name             string
		expectedErr      error
		expectedErrStr   string
		validateRequest  func(req *http.Request)
		getResponse      func() *http.Response
		getResponseError func() error
		permit           bool
	}{
		{
			name: "success set permit to true",
			validateRequest: func(req *http.Request) {
				assert.Equal(t, host+ipHotspotHostURL, req.URL.String())
				assert.Equal(t, "application/json;charset=UTF-8", req.Header.Get("Content-Type"))
				assert.Equal(t, http.MethodPost, req.Method)

				b, err := io.ReadAll(req.Body)
				assert.Nil(t, err)
				var body permitTrueReq
				err = json.Unmarshal(b, &body)
				assert.Nil(t, err)
				assert.Equal(t, body.Mac, mac)
				assert.Equal(t, body.Permit, true)
			},
			permit: true,
			getResponse: func() *http.Response {
				body := successRes
				bodyStr, err := json.Marshal(body)
				assert.Nil(t, err)

				bytesReader := bytes.NewReader(bodyStr)
				resp := http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytesReader),
				}
				return &resp
			},
			getResponseError: func() error {
				return nil
			},
		},
		{
			name: "success set permit to false",
			validateRequest: func(req *http.Request) {
				assert.Equal(t, host+ipHotspotHostURL, req.URL.String())
				assert.Equal(t, "application/json;charset=UTF-8", req.Header.Get("Content-Type"))
				assert.Equal(t, http.MethodPost, req.Method)

				b, err := io.ReadAll(req.Body)
				assert.Nil(t, err)
				var body permitFalseReq
				err = json.Unmarshal(b, &body)
				assert.Nil(t, err)
				assert.Equal(t, body.Mac, mac)
				assert.Equal(t, body.Deny, true)
			},
			permit: false,
			getResponse: func() *http.Response {
				body := successRes
				bodyStr, err := json.Marshal(body)
				assert.Nil(t, err)

				bytesReader := bytes.NewReader(bodyStr)
				resp := http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytesReader),
				}
				return &resp
			},
			getResponseError: func() error {
				return nil
			},
		},
		{
			name: "empty resp",
			validateRequest: func(req *http.Request) {
				assert.Equal(t, host+ipHotspotHostURL, req.URL.String())
				assert.Equal(t, "application/json;charset=UTF-8", req.Header.Get("Content-Type"))
				assert.Equal(t, http.MethodPost, req.Method)

				b, err := io.ReadAll(req.Body)
				assert.Nil(t, err)
				var body permitTrueReq
				err = json.Unmarshal(b, &body)
				assert.Nil(t, err)
				assert.Equal(t, body.Mac, mac)
				assert.Equal(t, body.Permit, true)
			},
			permit: true,
			getResponse: func() *http.Response {
				bodyStr, err := json.Marshal(nil)
				assert.Nil(t, err)

				bytesReader := bytes.NewReader(bodyStr)
				resp := http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytesReader),
				}
				return &resp
			},
			getResponseError: func() error {
				return nil
			},
			expectedErrStr: "no status in setaccess response",
		},
		{
			name:            "error from client",
			validateRequest: func(req *http.Request) {},
			getResponse: func() *http.Response {
				return nil
			},
			getResponseError: func() error {
				return someErr
			},
			expectedErr: someErr,
		},
		{
			name:            "http.StatusUnauthorized status code",
			validateRequest: func(req *http.Request) {},
			getResponse: func() *http.Response {
				bytesReader := strings.NewReader("")
				resp := http.Response{
					StatusCode: http.StatusUnauthorized,
					Body:       io.NopCloser(bytesReader),
				}
				return &resp
			},
			getResponseError: func() error {
				return nil
			},
			expectedErr: errs.ErrUnauthorized,
		},
		{
			name:            "status code not 200",
			validateRequest: func(req *http.Request) {},
			getResponse: func() *http.Response {
				bytesReader := strings.NewReader("")
				resp := http.Response{
					StatusCode: http.StatusBadRequest,
					Body:       io.NopCloser(bytesReader),
				}
				return &resp
			},
			getResponseError: func() error {
				return nil
			},
			expectedErrStr: "error in setaccess request, status code: 400",
		},
		{
			name:            "error while unmarshal body",
			validateRequest: func(req *http.Request) {},
			getResponse: func() *http.Response {
				stringReader := strings.NewReader("")
				resp := http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(stringReader),
				}
				return &resp
			},
			getResponseError: func() error {
				return nil
			},
			expectedErrStr: "unmarshal response error in setaccess request:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := mock_accessupdate.NewMockclient(ctrl)
			client.EXPECT().Do(gomock.Cond(func(x any) bool {
				req, ok := x.(*http.Request)
				if !ok || req == nil {
					t.Errorf("empty request")
					return false
				}
				tt.validateRequest(req)
				return true
			})).Return(tt.getResponse(), tt.getResponseError())

			accessUpdate := NewAccessUpdate(host, client)
			err := accessUpdate.SetPermit(mac, tt.permit)
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

func TestAccessUpdate_SetPolicy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const (
		host   = "host"
		mac    = "mac"
		policy = "policy"
	)

	successRes := response{
		"key": {},
	}
	someErr := errors.New("some err")

	tests := []struct {
		name             string
		expectedErr      error
		expectedErrStr   string
		validateRequest  func(req *http.Request)
		getResponse      func() *http.Response
		getResponseError func() error
		policy           string
	}{
		{
			name: "success set policy to some policy",
			validateRequest: func(req *http.Request) {
				assert.Equal(t, host+ipHotspotHostURL, req.URL.String())
				assert.Equal(t, "application/json;charset=UTF-8", req.Header.Get("Content-Type"))
				assert.Equal(t, http.MethodPost, req.Method)

				b, err := io.ReadAll(req.Body)
				assert.Nil(t, err)
				var body setPolicyReq
				err = json.Unmarshal(b, &body)
				assert.Nil(t, err)
				assert.Equal(t, body.Mac, mac)
				assert.Equal(t, body.Policy, policy)
			},
			policy: policy,
			getResponse: func() *http.Response {
				body := successRes
				bodyStr, err := json.Marshal(body)
				assert.Nil(t, err)

				bytesReader := bytes.NewReader(bodyStr)
				resp := http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytesReader),
				}
				return &resp
			},
			getResponseError: func() error {
				return nil
			},
		},
		{
			name: "success set policy to none",
			validateRequest: func(req *http.Request) {
				assert.Equal(t, host+ipHotspotHostURL, req.URL.String())
				assert.Equal(t, "application/json;charset=UTF-8", req.Header.Get("Content-Type"))
				assert.Equal(t, http.MethodPost, req.Method)

				b, err := io.ReadAll(req.Body)
				assert.Nil(t, err)
				var body setEmptyPolicyReq
				err = json.Unmarshal(b, &body)
				assert.Nil(t, err)
				assert.Equal(t, body.Mac, mac)
				assert.Equal(t, body.Policy, false)
			},
			policy: homeassistantdto.NonePolicy,
			getResponse: func() *http.Response {
				body := successRes
				bodyStr, err := json.Marshal(body)
				assert.Nil(t, err)

				bytesReader := bytes.NewReader(bodyStr)
				resp := http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytesReader),
				}
				return &resp
			},
			getResponseError: func() error {
				return nil
			},
		},

		{
			name:            "error from client",
			validateRequest: func(req *http.Request) {},
			getResponse: func() *http.Response {
				return nil
			},
			getResponseError: func() error {
				return someErr
			},
			expectedErr: someErr,
		},
		{
			name:            "http.StatusUnauthorized status code",
			validateRequest: func(req *http.Request) {},
			getResponse: func() *http.Response {
				bytesReader := strings.NewReader("")
				resp := http.Response{
					StatusCode: http.StatusUnauthorized,
					Body:       io.NopCloser(bytesReader),
				}
				return &resp
			},
			getResponseError: func() error {
				return nil
			},
			expectedErr: errs.ErrUnauthorized,
		},
		{
			name:            "status code not 200",
			validateRequest: func(req *http.Request) {},
			getResponse: func() *http.Response {
				bytesReader := strings.NewReader("")
				resp := http.Response{
					StatusCode: http.StatusBadRequest,
					Body:       io.NopCloser(bytesReader),
				}
				return &resp
			},
			getResponseError: func() error {
				return nil
			},
			expectedErrStr: "error in setaccess request, status code: 400",
		},
		{
			name:            "error while unmarshal body",
			validateRequest: func(req *http.Request) {},
			getResponse: func() *http.Response {
				stringReader := strings.NewReader("")
				resp := http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(stringReader),
				}
				return &resp
			},
			getResponseError: func() error {
				return nil
			},
			expectedErrStr: "unmarshal response error in setaccess request:",
		},
		{
			name:            "empty response",
			validateRequest: func(req *http.Request) {},
			policy:          policy,
			getResponse: func() *http.Response {
				bodyStr, err := json.Marshal(nil)
				assert.Nil(t, err)

				bytesReader := bytes.NewReader(bodyStr)
				resp := http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytesReader),
				}
				return &resp
			},
			getResponseError: func() error {
				return nil
			},
			expectedErrStr: "no status in setaccess response",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := mock_accessupdate.NewMockclient(ctrl)
			client.EXPECT().Do(gomock.Cond(func(x any) bool {
				req, ok := x.(*http.Request)
				if !ok || req == nil {
					t.Errorf("empty request")
					return false
				}
				tt.validateRequest(req)
				return true
			})).Return(tt.getResponse(), tt.getResponseError())

			accessUpdate := NewAccessUpdate(host, client)
			err := accessUpdate.SetPolicy(mac, tt.policy)
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
