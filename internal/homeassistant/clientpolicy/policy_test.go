package clientpolicy

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"keeneticToMqtt/internal/dto"
	mock_clientpolicy "keeneticToMqtt/test/mocks/gomock/homeassistant/clientpolicy"
)

func TestClientPolicy_SendDiscoveryMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const (
		mac       = "mac"
		name      = "name"
		basetopic = "basetopic"
	)
	someErr := errors.New("some error")
	policies := []string{name}

	tests := []struct {
		name          string
		expectedErr   error
		discovery     func() discovery
		policyStorage func() policyStorage
	}{
		{
			name: "success send discovery message",
			discovery: func() discovery {
				discovery := mock_clientpolicy.NewMockdiscovery(ctrl)
				discovery.EXPECT().
					SendDiscoverySelect(
						gomock.Eq("basetopic/mac_policy/command"),
						gomock.Eq("basetopic/mac_policy/state"),
						gomock.Eq(name),
						gomock.Eq("name_policy"),
						gomock.Eq(policies),
					).
					Return(nil)

				return discovery
			},
			policyStorage: func() policyStorage {
				policyStorage := mock_clientpolicy.NewMockpolicyStorage(ctrl)
				policyStorage.EXPECT().GetPolicyList().Return(policies)

				return policyStorage
			},
		},
		{
			name: "error while send discovery message",
			discovery: func() discovery {
				discovery := mock_clientpolicy.NewMockdiscovery(ctrl)
				discovery.EXPECT().
					SendDiscoverySelect(
						gomock.Eq("basetopic/mac_policy/command"),
						gomock.Eq("basetopic/mac_policy/state"),
						gomock.Eq(name),
						gomock.Eq("name_policy"),
						gomock.Eq(policies),
					).
					Return(someErr)

				return discovery
			},
			policyStorage: func() policyStorage {
				policyStorage := mock_clientpolicy.NewMockpolicyStorage(ctrl)
				policyStorage.EXPECT().GetPolicyList().Return(policies)

				return policyStorage
			},
			expectedErr: someErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			policy := ClientPolicy{
				basetopic:       basetopic,
				discoveryClient: tt.discovery(),
				policyStorage:   tt.policyStorage(),
			}

			err := policy.SendDiscoveryMessage(mac, name)
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
			name: "success policy send",
			mqtt: func() mqtt {
				mqtt := mock_clientpolicy.NewMockmqtt(ctrl)
				mqtt.EXPECT().
					SendMessage(gomock.Eq("basetopic/mac_policy/state"), gomock.Eq(policy))

				return mqtt
			},
			client: dto.Client{
				Mac:    mac,
				Policy: policy,
				Name:   name,
				Permit: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientPolicy := ClientPolicy{
				basetopic: basetopic,
				mqtt:      tt.mqtt(),
			}

			clientPolicy.SendState(tt.client)
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
			name: "run policy consumer",
			mqtt: func(ch chan string) mqtt {
				mqtt := mock_clientpolicy.NewMockmqtt(ctrl)
				mqtt.EXPECT().
					Subscribe(gomock.Eq("basetopic/mac_policy/command")).
					Return(ch)
				return mqtt
			},
			client: dto.Client{
				Mac:    mac,
				Policy: policy,
				Name:   name,
			},
			accessUpdate: func() accessUpdate {
				accessUpdate := mock_clientpolicy.NewMockaccessUpdate(ctrl)
				accessUpdate.EXPECT().
					SetPolicy(gomock.Eq(mac), gomock.Eq(name)).
					Return(nil)
				return accessUpdate
			},
			logger: func() logger {
				logger := mock_clientpolicy.NewMocklogger(ctrl)
				return logger
			},
			payload: name,
		},
		{
			name: "error while setting permit",
			mqtt: func(ch chan string) mqtt {
				mqtt := mock_clientpolicy.NewMockmqtt(ctrl)
				mqtt.EXPECT().
					Subscribe(gomock.Eq("basetopic/mac_policy/command")).
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
				accessUpdate := mock_clientpolicy.NewMockaccessUpdate(ctrl)
				accessUpdate.EXPECT().
					SetPolicy(gomock.Eq(mac), gomock.Eq(name)).
					Return(someErr)
				return accessUpdate
			},
			logger: func() logger {
				logger := mock_clientpolicy.NewMocklogger(ctrl)
				logger.EXPECT().Error(gomock.Eq("client error while setting policy"), gomock.Eq("error"), gomock.Eq(someErr))
				return logger
			},
			payload:     name,
			expectedErr: someErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			discovery := mock_clientpolicy.NewMockdiscovery(ctrl)
			policyStorage := mock_clientpolicy.NewMockpolicyStorage(ctrl)
			ch := make(chan string)
			permit := NewClientPolicy(
				basetopic,
				discovery,
				tt.mqtt(ch),
				tt.accessUpdate(),
				policyStorage,
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
