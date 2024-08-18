package clientpermit

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"keeneticToMqtt/internal/dto"
	mock_clientpermit "keeneticToMqtt/test/mocks/gomock/homeassistant/clientpermit"
)

func TestClientPermit_SendDiscoveryMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const (
		mac       = "mac"
		name      = "name"
		basetopic = "basetopic"
	)
	client := dto.Client{Mac: mac, Name: name}
	someErr := errors.New("some error")

	tests := []struct {
		name        string
		expectedErr error
		discovery   func() discovery
	}{
		{
			name: "success send discovery message",
			discovery: func() discovery {
				discovery := mock_clientpermit.NewMockdiscovery(ctrl)
				discovery.EXPECT().
					SendDiscoverySwitch(gomock.Eq("basetopic/mac_permit/command"), gomock.Eq("basetopic/mac_permit/state"), gomock.Eq("name"), gomock.Eq("name_permit")).
					Return(nil)

				return discovery
			},
		},
		{
			name: "error while send discovery message",
			discovery: func() discovery {
				discovery := mock_clientpermit.NewMockdiscovery(ctrl)
				discovery.EXPECT().
					SendDiscoverySwitch(gomock.Eq("basetopic/mac_permit/command"), gomock.Eq("basetopic/mac_permit/state"), gomock.Eq("name"), gomock.Eq("name_permit")).
					Return(someErr)

				return discovery
			},
			expectedErr: someErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			perimt := ClientPermit{
				basetopic:       basetopic,
				discoveryClient: tt.discovery(),
			}

			err := perimt.SendDiscoveryMessage(client)
			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestClientPermit_GetState(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name        string
		expected    string
		expectedErr error
		client      dto.Client
	}{
		{
			name: "success get permit false",
			client: dto.Client{
				Permit: false,
			},
			expected: "OFF",
		},
		{
			name: "success get permit true",
			client: dto.Client{
				Permit: true,
			},
			expected: "ON",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			permit := ClientPermit{}
			state, err := permit.GetState(tt.client)
			assert.Nil(t, err)
			assert.Equal(t, tt.expected, state)
		})
	}
}

func TestClientPermit_Consume(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const (
		mac       = "mac"
		basetopic = "basetopic"
	)

	someErr := errors.New("some error")

	tests := []struct {
		name         string
		expectedErr  error
		client       dto.Client
		accessUpdate func() accessUpdate
		payload      string
	}{
		{
			name:   "set permit to false",
			client: dto.Client{Mac: mac},
			accessUpdate: func() accessUpdate {
				accessUpdate := mock_clientpermit.NewMockaccessUpdate(ctrl)
				accessUpdate.EXPECT().
					SetPermit(gomock.Eq(mac), gomock.Eq(false)).
					Return(nil)
				return accessUpdate
			},
			payload: offPayload,
		},
		{
			name:   "set permit to true",
			client: dto.Client{Mac: mac},
			accessUpdate: func() accessUpdate {
				accessUpdate := mock_clientpermit.NewMockaccessUpdate(ctrl)
				accessUpdate.EXPECT().
					SetPermit(gomock.Eq(mac), gomock.Eq(true)).
					Return(nil)
				return accessUpdate
			},
			payload: onPayload,
		},
		{
			name:   "error while setting permit",
			client: dto.Client{Mac: mac},
			accessUpdate: func() accessUpdate {
				accessUpdate := mock_clientpermit.NewMockaccessUpdate(ctrl)
				accessUpdate.EXPECT().
					SetPermit(gomock.Eq(mac), gomock.Eq(false)).
					Return(someErr)
				return accessUpdate
			},
			payload:     offPayload,
			expectedErr: someErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			discovery := mock_clientpermit.NewMockdiscovery(ctrl)
			permit := NewClientPermit(
				basetopic,
				discovery,
				tt.accessUpdate(),
			)

			err := permit.Consume(tt.client, tt.payload)
			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
