package clientpermit

import (
	"errors"
	"testing"
	"time"

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

			err := perimt.SendDiscoveryMessage(mac, name)
			if tt.expectedErr != nil {
				assert.True(t, errors.Is(err, tt.expectedErr))
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestClientPermit_SendState(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const (
		mac       = "mac"
		policy    = "policy"
		name      = "name"
		basetopic = "basetopic"
	)

	tests := []struct {
		name        string
		expectedErr error
		mqtt        func() mqtt
		client      dto.Client
	}{
		{
			name: "success send permit false",
			mqtt: func() mqtt {
				mqtt := mock_clientpermit.NewMockmqtt(ctrl)
				mqtt.EXPECT().
					SendMessage(gomock.Eq("basetopic/mac_permit/state"), gomock.Eq(offPayload))

				return mqtt
			},
			client: dto.Client{
				Mac:    mac,
				Policy: policy,
				Name:   name,
				Permit: false,
			},
		},
		{
			name: "success send permit true",
			mqtt: func() mqtt {
				mqtt := mock_clientpermit.NewMockmqtt(ctrl)
				mqtt.EXPECT().
					SendMessage(gomock.Eq("basetopic/mac_permit/state"), gomock.Eq(onPayload))

				return mqtt
			},
			client: dto.Client{
				Mac:    mac,
				Policy: policy,
				Name:   name,
				Permit: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			permit := ClientPermit{
				basetopic: basetopic,
				mqtt:      tt.mqtt(),
			}

			permit.SendState(tt.client)
		})
	}
}

func TestClientPermit_RunConsumer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const (
		mac       = "mac"
		policy    = "policy"
		name      = "name"
		basetopic = "basetopic"
	)

	someErr := errors.New("some error")

	tests := []struct {
		name         string
		expectedErr  error
		mqtt         func(chan string) mqtt
		logger       func() logger
		client       dto.Client
		accessUpdate func() accessUpdate
		payload      string
	}{
		{
			name: "set permit to false",
			mqtt: func(ch chan string) mqtt {
				mqtt := mock_clientpermit.NewMockmqtt(ctrl)
				mqtt.EXPECT().
					Subscribe(gomock.Eq("basetopic/mac_permit/command")).
					Return(ch)
				return mqtt
			},
			client: dto.Client{
				Mac:    mac,
				Policy: policy,
				Name:   name,
				Permit: false,
			},
			accessUpdate: func() accessUpdate {
				accessUpdate := mock_clientpermit.NewMockaccessUpdate(ctrl)
				accessUpdate.EXPECT().
					SetPermit(gomock.Eq(mac), gomock.Eq(false)).
					Return(nil)
				return accessUpdate
			},
			logger: func() logger {
				logger := mock_clientpermit.NewMocklogger(ctrl)
				return logger
			},
			payload: offPayload,
		},
		{
			name: "set permit to true",
			mqtt: func(ch chan string) mqtt {
				mqtt := mock_clientpermit.NewMockmqtt(ctrl)
				mqtt.EXPECT().
					Subscribe(gomock.Eq("basetopic/mac_permit/command")).
					Return(ch)
				return mqtt
			},
			client: dto.Client{
				Mac:    mac,
				Policy: policy,
				Name:   name,
				Permit: false,
			},
			accessUpdate: func() accessUpdate {
				accessUpdate := mock_clientpermit.NewMockaccessUpdate(ctrl)
				accessUpdate.EXPECT().
					SetPermit(gomock.Eq(mac), gomock.Eq(true)).
					Return(nil)
				return accessUpdate
			},
			logger: func() logger {
				logger := mock_clientpermit.NewMocklogger(ctrl)
				return logger
			},
			payload: onPayload,
		},
		{
			name: "error while setting permit",
			mqtt: func(ch chan string) mqtt {
				mqtt := mock_clientpermit.NewMockmqtt(ctrl)
				mqtt.EXPECT().
					Subscribe(gomock.Eq("basetopic/mac_permit/command")).
					Return(ch)
				return mqtt
			},
			client: dto.Client{
				Mac:    mac,
				Policy: policy,
				Name:   name,
				Permit: false,
			},
			accessUpdate: func() accessUpdate {
				accessUpdate := mock_clientpermit.NewMockaccessUpdate(ctrl)
				accessUpdate.EXPECT().
					SetPermit(gomock.Eq(mac), gomock.Eq(false)).
					Return(someErr)
				return accessUpdate
			},
			logger: func() logger {
				logger := mock_clientpermit.NewMocklogger(ctrl)
				logger.EXPECT().Error(gomock.Eq("client error while setting permit"), gomock.Eq("error"), someErr)
				return logger
			},
			payload:     offPayload,
			expectedErr: someErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			discovery := mock_clientpermit.NewMockdiscovery(ctrl)
			ch := make(chan string)
			permit := NewClientPermit(
				basetopic,
				discovery,
				tt.mqtt(ch),
				tt.accessUpdate(),
				tt.logger(),
			)

			go permit.RunConsumer(mac)

			ch <- tt.payload

			ticker := time.NewTicker(time.Millisecond)
			<-ticker.C
			return
		})
	}
}
