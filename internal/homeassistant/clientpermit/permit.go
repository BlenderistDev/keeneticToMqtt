package clientpermit

import (
	"fmt"
	"strings"

	"keeneticToMqtt/internal/dto"
)

//go:generate mockgen -source=permit.go -destination=../../../test/mocks/gomock/homeassistant/clientpermit/permit.go

const (
	entityTypeName = "permit"
	offPayload     = "OFF"
	onPayload      = "ON"
)

type (
	discovery interface {
		SendDiscoverySwitch(commandTopic, stateTopic, deviceName, name string) error
	}
	accessUpdate interface {
		SetPermit(mac string, permit bool) error
	}
)

// ClientPermit struct for handle home assistant client permit entities.
type ClientPermit struct {
	basetopic       string
	discoveryClient discovery
	accessUpdate    accessUpdate
}

// NewClientPermit creates new ClientPermit.
func NewClientPermit(
	basetopic string,
	discoveryClient discovery,
	accessUpdate accessUpdate,
) *ClientPermit {
	return &ClientPermit{
		basetopic:       basetopic,
		discoveryClient: discoveryClient,
		accessUpdate:    accessUpdate,
	}
}

// SendDiscoveryMessage sending home assistant discovery message.
func (p *ClientPermit) SendDiscoveryMessage(client dto.Client) error {
	commandTopic := p.GetCommandTopic(client)
	stateTopic := p.GetStateTopic(client)

	if err := p.discoveryClient.SendDiscoverySwitch(commandTopic, stateTopic, client.Name, client.Name+"_"+entityTypeName); err != nil {
		return fmt.Errorf("ClientPermit SendDiscoveryMessage error: %w", err)
	}

	return nil
}

// GetStateTopic returns state topic.
func (p *ClientPermit) GetStateTopic(client dto.Client) string {
	mac := strings.Replace(client.Mac, ":", "_", -1)
	return fmt.Sprintf("%s/%s_%s/state", p.basetopic, mac, entityTypeName)
}

// GetCommandTopic returns command topic.
func (p *ClientPermit) GetCommandTopic(client dto.Client) string {
	mac := strings.Replace(client.Mac, ":", "_", -1)
	return fmt.Sprintf("%s/%s_%s/command", p.basetopic, mac, entityTypeName)
}

// GetState returns client permit state.
func (p *ClientPermit) GetState(client dto.Client) (string, error) {
	msg := offPayload
	if client.Permit {
		msg = onPayload
	}
	return msg, nil
}

// Consume consumes message.
func (p *ClientPermit) Consume(client dto.Client, message string) error {
	permit := true
	if message == offPayload {
		permit = false
	}
	if err := p.accessUpdate.SetPermit(client.Mac, permit); err != nil {
		return fmt.Errorf("client error while setting permit: %w", err)
	}

	return nil
}
