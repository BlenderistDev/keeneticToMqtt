package txbytes

import (
	"errors"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	mock_txbytes "keeneticToMqtt/test/mocks/gomock/homeassistant/txbytes"

	"keeneticToMqtt/internal/dto"
)

func TestTxBytes_SendDiscoveryMessage(t *testing.T) {
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
				discovery := mock_txbytes.NewMockdiscovery(ctrl)
				discovery.EXPECT().
					SendDiscoverySensor(
						gomock.Eq("basetopic/mac_txbytes/state"),
						gomock.Eq(name),
						gomock.Eq("name_txbytes"),
						gomock.Eq(unit),
					).
					Return(nil)

				return discovery
			},
		},
		{
			name: "error while send discovery message",
			discovery: func() discovery {
				discovery := mock_txbytes.NewMockdiscovery(ctrl)
				discovery.EXPECT().
					SendDiscoverySensor(
						gomock.Eq("basetopic/mac_txbytes/state"),
						gomock.Eq(name),
						gomock.Eq("name_txbytes"),
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
			txBytes := NewTxBytes(basetopic, tt.discovery())
			err := txBytes.SendDiscoveryMessage(client)
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
		mac     = "mac"
		txBytes = 123
		name    = "name"
	)

	tests := []struct {
		name     string
		expected string
		client   dto.Client
	}{
		{
			name: "success txbytes get",
			client: dto.Client{
				Mac:     mac,
				TxBytes: txBytes,
				Name:    name,
			},
			expected: strconv.Itoa(txBytes),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			txBytes := TxBytes{}
			res, err := txBytes.GetState(tt.client)
			assert.Nil(t, err)
			assert.Equal(t, tt.expected, res)
		})
	}
}

func TestTxBytes_Consume(t *testing.T) {
	txBytes := TxBytes{}

	err := txBytes.Consume(dto.Client{}, "")
	assert.Nil(t, err)
}

func TestTxBytes_GetCommandTopic(t *testing.T) {
	txBytes := TxBytes{}
	assert.Empty(t, txBytes.GetCommandTopic(dto.Client{}))
}
