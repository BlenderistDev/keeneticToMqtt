package policylist

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
	"keeneticToMqtt/internal/dto/keeneticdto"
	"keeneticToMqtt/internal/errs"
	mock_policylist "keeneticToMqtt/test/mocks/gomock/clients/keenetic/policylist"
)

func TestPolicyList_GetPolicyList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const (
		host   = "host"
		policy = "policy"
	)

	successRes := map[string]keeneticdto.Policy{
		policy: {},
	}
	someErr := errors.New("some err")

	tests := []struct {
		name             string
		expected         map[string]keeneticdto.Policy
		expectedErr      error
		expectedErrStr   string
		validateRequest  func(req *http.Request)
		getResponse      func() *http.Response
		getResponseError func() error
	}{
		{
			name: "success get policy list",
			validateRequest: func(req *http.Request) {
				assert.Equal(t, host+policyListUrl, req.URL.String())
				assert.Equal(t, "application/json;charset=UTF-8", req.Header.Get("Content-Type"))
				assert.Equal(t, http.MethodGet, req.Method)
			},
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
			expected: successRes,
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
			expectedErrStr: "error in GetPolicyList request, status code: 400",
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
			expectedErrStr: "unmarshal response error in setpolicy request:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := mock_policylist.NewMockclient(ctrl)
			client.EXPECT().Do(gomock.Cond(func(x any) bool {
				req, ok := x.(*http.Request)
				if !ok || req == nil {
					t.Errorf("empty request")
					return false
				}
				tt.validateRequest(req)
				return true
			})).Return(tt.getResponse(), tt.getResponseError())

			policyList := NewPolicyList(host, client)
			res, err := policyList.GetPolicyList()
			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr)
			} else if tt.expectedErrStr != "" {
				assert.Regexp(t, tt.expectedErrStr+".*", err.Error())
			} else {
				assert.Equal(t, tt.expected, res)
				assert.Nil(t, err)
			}
		})
	}
}
