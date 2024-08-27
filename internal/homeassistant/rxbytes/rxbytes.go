package rxbytes

import (
	"fmt"
	"strconv"
	"strings"

	"keeneticToMqtt/internal/dto"
)

//go:generate mockgen -source=rxbytes.go -destination=../../../test/mocks/gomock/homeassistant/rxbytes/rxbytes.go

const (
	entityTypeName = "rxbytes"
	unit           = "bytes"
)

type (
	discovery interface {
		SendDiscoverySensor(stateTopic, deviceName, name, unit string) error
	}
)

// RxBytes struct for handle home assistant client rxbytes entities.
type RxBytes struct {
	basetopic       string
	discoveryClient discovery
}

// NewRxBytes creates new RxBytes.
func NewRxBytes(
	basetopic string,
	discoveryClient discovery,
) *RxBytes {
	return &RxBytes{
		basetopic:       basetopic,
		discoveryClient: discoveryClient,
	}
}

// SendDiscoveryMessage sends homeassistant discovery message.
func (b *RxBytes) SendDiscoveryMessage(client dto.Client) error {
	stateTopic := b.GetStateTopic(client)
	if err := b.discoveryClient.SendDiscoverySensor(stateTopic, client.Name, client.Name+"_"+entityTypeName, unit); err != nil {
		return fmt.Errorf("RxBytes SendDiscoveryMessage error: %w", err)
	}

	return nil
}

// GetState returns entity state.
func (b *RxBytes) GetState(client dto.Client) (string, error) {
	return strconv.Itoa(int(client.RxBytes)), nil
}

// Consume consumes message.
func (b *RxBytes) Consume(_ dto.Client, _ string) error {
	return nil
}

// GetStateTopic returns state topic.
func (b *RxBytes) GetStateTopic(client dto.Client) string {
	mac := strings.Replace(client.Mac, ":", "_", -1)
	return fmt.Sprintf("%s/%s_%s/state", b.basetopic, mac, entityTypeName)
}

// GetCommandTopic returns command topic.
func (b *RxBytes) GetCommandTopic(_ dto.Client) string {
	return ""
}
