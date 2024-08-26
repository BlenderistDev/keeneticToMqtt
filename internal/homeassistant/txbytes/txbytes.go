package txbytes

import (
	"fmt"
	"strconv"
	"strings"

	"keeneticToMqtt/internal/dto"
)

//go:generate mockgen -source=txbytes.go -destination=../../../test/mocks/gomock/homeassistant/txbytes/txbytes.go

const (
	entityTypeName = "txbytes"
	unit           = "bytes"
)

type (
	discovery interface {
		SendDiscoverySensor(stateTopic, deviceName, name, unit string) error
	}
)

// TxBytes struct for handle home assistant client txbytes entities.
type TxBytes struct {
	basetopic       string
	discoveryClient discovery
}

// NewTxBytes creates new TxBytes.
func NewTxBytes(
	basetopic string,
	discoveryClient discovery,
) *TxBytes {
	return &TxBytes{
		basetopic:       basetopic,
		discoveryClient: discoveryClient,
	}
}

// SendDiscoveryMessage sends homeassistant discovery message.
func (p *TxBytes) SendDiscoveryMessage(client dto.Client) error {
	stateTopic := p.GetStateTopic(client)
	if err := p.discoveryClient.SendDiscoverySensor(stateTopic, client.Name, client.Name+"_"+entityTypeName, unit); err != nil {
		return fmt.Errorf("TxBytes SendDiscoveryMessage error: %w", err)
	}

	return nil
}

// GetState returns entity state.
func (p *TxBytes) GetState(client dto.Client) (string, error) {
	return strconv.Itoa(int(client.TxBytes)), nil
}

// Consume consumes message.
func (p *TxBytes) Consume(_ dto.Client, _ string) error {
	return nil
}

// GetStateTopic returns state topic.
func (p *TxBytes) GetStateTopic(client dto.Client) string {
	mac := strings.Replace(client.Mac, ":", "_", -1)
	return fmt.Sprintf("%s/%s_%s/state", p.basetopic, mac, entityTypeName)
}

// GetCommandTopic returns command topic.
func (p *TxBytes) GetCommandTopic(_ dto.Client) string {
	return ""
}
