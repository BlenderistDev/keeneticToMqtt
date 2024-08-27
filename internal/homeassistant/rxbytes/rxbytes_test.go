package rxbytes

import (
	"errors"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"keeneticToMqtt/internal/dto"
	mock_rxbytes "keeneticToMqtt/test/mocks/gomock/homeassistant/rxbytes"
)

func TestRxBytes_SendDiscoveryMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const (
		mac       = "mac"
		name      = "name"
		basetopic = "basetopic"
	)
	someErr := errors.New("some error")

	client := dto.Client{Mac: mac, Name: name}

	tests := []struct {
		name        string
		expectedErr error
		discovery   func() discovery
	}{
		{
			name: "success send discovery message",
			discovery: func() discovery {
				discovery := mock_rxbytes.NewMockdiscovery(ctrl)
				discovery.EXPECT().
					SendDiscoverySensor(
						gomock.Eq("basetopic/mac_rxbytes/state"),
						gomock.Eq(name),
						gomock.Eq("name_rxbytes"),
						gomock.Eq(unit),
					).
					Return(nil)

				return discovery
			},
		},
		{
			name: "error while send discovery message",
			discovery: func() discovery {
				discovery := mock_rxbytes.NewMockdiscovery(ctrl)
				discovery.EXPECT().
					SendDiscoverySensor(
						gomock.Eq("basetopic/mac_rxbytes/state"),
						gomock.Eq(name),
						gomock.Eq("name_rxbytes"),
						gomock.Eq(unit),
					).
					Return(someErr)

				return discovery
			},
			expectedErr: someErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rxBytes := NewRxBytes(basetopic, tt.discovery())
			err := rxBytes.SendDiscoveryMessage(client)
			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestRxBytes_GetState(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const (
		mac     = "mac"
		rxBytes = 123
		name    = "name"
	)

	tests := []struct {
		name     string
		expected string
		client   dto.Client
	}{
		{
			name: "success rxbytes get",
			client: dto.Client{
				Mac:     mac,
				RxBytes: rxBytes,
				Name:    name,
			},
			expected: strconv.Itoa(rxBytes),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rxBytes := RxBytes{}
			res, err := rxBytes.GetState(tt.client)
			assert.Nil(t, err)
			assert.Equal(t, tt.expected, res)
		})
	}
}

func TestRxBytes_Consume(t *testing.T) {
	rxBytes := RxBytes{}

	err := rxBytes.Consume(dto.Client{}, "")
	assert.Nil(t, err)
}

func TestRxBytes_GetCommandTopic(t *testing.T) {
	rxBytes := RxBytes{}
	assert.Empty(t, rxBytes.GetCommandTopic(dto.Client{}))
}
