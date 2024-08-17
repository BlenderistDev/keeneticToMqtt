package clientlist

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"keeneticToMqtt/internal/dto"
	"keeneticToMqtt/internal/dto/homeassistantdto"
	"keeneticToMqtt/internal/dto/keeneticdto"
	mock_clientlist "keeneticToMqtt/test/mocks/gomock/services/clientlist"

	"go.uber.org/mock/gomock"
)

func TestClientList_GetClientList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const (
		mac1  = "mac1"
		name1 = "name1"
	)

	var (
		policy      = "policy"
		emptyPolicy = ""
		someErr     = errors.New("some err")
	)

	tests := []struct {
		name        string
		listClient  func() listClient
		whitelist   []string
		expected    []dto.Client
		expectedErr error
	}{
		{
			name: "success list building",
			listClient: func() listClient {
				listClient := mock_clientlist.NewMocklistClient(ctrl)
				listClient.EXPECT().GetDeviceList().Return([]keeneticdto.DeviceInfoResponse{
					{
						Mac:  mac1,
						Name: name1,
					},
				}, nil)

				listClient.EXPECT().GetClientPolicyList().Return([]keeneticdto.DevicePolicy{
					{
						Mac:    mac1,
						Policy: &policy,
						Permit: true,
					},
				}, nil)

				return listClient
			},
			whitelist: []string{mac1},
			expected: []dto.Client{
				{
					Mac:    mac1,
					Policy: policy,
					Name:   name1,
					Permit: true,
				},
			},
		},
		{
			name: "mac is not in white list",
			listClient: func() listClient {
				listClient := mock_clientlist.NewMocklistClient(ctrl)
				listClient.EXPECT().GetDeviceList().Return([]keeneticdto.DeviceInfoResponse{
					{
						Mac:  mac1,
						Name: name1,
					},
				}, nil)

				listClient.EXPECT().GetClientPolicyList().Return([]keeneticdto.DevicePolicy{
					{
						Mac:    mac1,
						Policy: &policy,
						Permit: true,
					},
				}, nil)

				return listClient
			},
			whitelist: []string{},
			expected:  []dto.Client{},
		},
		{
			name: "GetDeviceList error",
			listClient: func() listClient {
				listClient := mock_clientlist.NewMocklistClient(ctrl)
				listClient.EXPECT().GetDeviceList().Return(nil, someErr)

				return listClient
			},
			expectedErr: someErr,
		},
		{
			name: "GetClientPolicyList error",
			listClient: func() listClient {
				listClient := mock_clientlist.NewMocklistClient(ctrl)
				listClient.EXPECT().GetDeviceList().Return([]keeneticdto.DeviceInfoResponse{
					{
						Mac:  mac1,
						Name: name1,
					},
				}, nil)

				listClient.EXPECT().GetClientPolicyList().Return(nil, someErr)

				return listClient
			},
			expectedErr: someErr,
		},
		{
			name: "empty policy",
			listClient: func() listClient {
				listClient := mock_clientlist.NewMocklistClient(ctrl)
				listClient.EXPECT().GetDeviceList().Return([]keeneticdto.DeviceInfoResponse{
					{
						Mac:  mac1,
						Name: name1,
					},
				}, nil)

				listClient.EXPECT().GetClientPolicyList().Return([]keeneticdto.DevicePolicy{
					{
						Mac:    mac1,
						Policy: &emptyPolicy,
						Permit: true,
					},
				}, nil)

				return listClient
			},
			whitelist: []string{mac1},
			expected: []dto.Client{
				{
					Mac:    mac1,
					Policy: homeassistantdto.NonePolicy,
					Name:   name1,
					Permit: true,
				},
			},
		},
		{
			name: "nil policy",
			listClient: func() listClient {
				listClient := mock_clientlist.NewMocklistClient(ctrl)
				listClient.EXPECT().GetDeviceList().Return([]keeneticdto.DeviceInfoResponse{
					{
						Mac:  mac1,
						Name: name1,
					},
				}, nil)

				listClient.EXPECT().GetClientPolicyList().Return([]keeneticdto.DevicePolicy{
					{
						Mac:    mac1,
						Policy: nil,
						Permit: true,
					},
				}, nil)

				return listClient
			},
			whitelist: []string{mac1},
			expected: []dto.Client{
				{
					Mac:    mac1,
					Policy: homeassistantdto.NonePolicy,
					Name:   name1,
					Permit: true,
				},
			},
		},
		{
			name: "no policy in map",
			listClient: func() listClient {
				listClient := mock_clientlist.NewMocklistClient(ctrl)
				listClient.EXPECT().GetDeviceList().Return([]keeneticdto.DeviceInfoResponse{
					{
						Mac:  mac1,
						Name: name1,
					},
				}, nil)

				listClient.EXPECT().GetClientPolicyList().Return([]keeneticdto.DevicePolicy{}, nil)

				return listClient
			},
			whitelist: []string{mac1},
			expected: []dto.Client{
				{
					Mac:    mac1,
					Policy: homeassistantdto.NonePolicy,
					Name:   name1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientList := NewClientList(
				tt.listClient(),
				tt.whitelist,
			)
			res, err := clientList.GetClientList()
			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr)
				assert.True(t, len(res) == 0)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.expected, res)
			}
		})
	}
}
