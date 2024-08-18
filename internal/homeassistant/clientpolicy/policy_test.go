package clientpolicy

import (
	"errors"
	"testing"

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

	client := dto.Client{Mac: mac, Name: name}

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

			err := policy.SendDiscoveryMessage(client)
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

	const (
		mac    = "mac"
		policy = "policy"
		name   = "name"
	)

	tests := []struct {
		name     string
		expected string
		client   dto.Client
	}{
		{
			name: "success policy get",
			client: dto.Client{
				Mac:    mac,
				Policy: policy,
				Name:   name,
			},
			expected: policy,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientPolicy := ClientPolicy{}

			res, err := clientPolicy.GetState(tt.client)
			assert.Nil(t, err)
			assert.Equal(t, tt.expected, res)
		})
	}
}

func TestClientPermit_Consume(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const (
		mac       = "mac"
		name      = "name"
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
			name:   "run policy consumer",
			client: dto.Client{Mac: mac},
			accessUpdate: func() accessUpdate {
				accessUpdate := mock_clientpolicy.NewMockaccessUpdate(ctrl)
				accessUpdate.EXPECT().
					SetPolicy(gomock.Eq(mac), gomock.Eq(name)).
					Return(nil)
				return accessUpdate
			},
			payload: name,
		},
		{
			name:   "error while setting permit",
			client: dto.Client{Mac: mac},
			accessUpdate: func() accessUpdate {
				accessUpdate := mock_clientpolicy.NewMockaccessUpdate(ctrl)
				accessUpdate.EXPECT().
					SetPolicy(gomock.Eq(mac), gomock.Eq(name)).
					Return(someErr)
				return accessUpdate
			},
			payload:     name,
			expectedErr: someErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			discovery := mock_clientpolicy.NewMockdiscovery(ctrl)
			policyStorage := mock_clientpolicy.NewMockpolicyStorage(ctrl)
			policy := NewClientPolicy(
				basetopic,
				discovery,
				tt.accessUpdate(),
				policyStorage,
			)

			err := policy.Consume(tt.client, tt.payload)
			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
