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
	mqtt interface {
		Subscribe(topic string) chan string
		SendMessage(topic, message string)
	}
	logger interface {
		Error(msg string, args ...any)
	}
)

// ClientPermit struct for handle home assistant client permit entities.
type ClientPermit struct {
	basetopic       string
	discoveryClient discovery
	accessUpdate    accessUpdate
	mqtt            mqtt
	logger          logger
}

// NewClientPermit creates new ClientPermit.
func NewClientPermit(
	basetopic string,
	discoveryClient discovery,
	mqtt mqtt,
	accessUpdate accessUpdate,
	logger logger,
) *ClientPermit {
	return &ClientPermit{
		basetopic:       basetopic,
		discoveryClient: discoveryClient,
		accessUpdate:    accessUpdate,
		mqtt:            mqtt,
		logger:          logger,
	}
}

// SendDiscoveryMessage sending home assistant discovery message.
func (p *ClientPermit) SendDiscoveryMessage(mac, name string) error {
	commandTopic := p.getCommandTopic(mac)
	stateTopic := p.getStateTopic(mac)

	if err := p.discoveryClient.SendDiscoverySwitch(commandTopic, stateTopic, name, name+"_"+entityTypeName); err != nil {
		return fmt.Errorf("ClientPermit SendDiscoveryMessage error: %w", err)
	}

	return nil
}

func (p *ClientPermit) getStateTopic(mac string) string {
	mac = strings.Replace(mac, ":", "_", -1)
	return fmt.Sprintf("%s/%s_%s/state", p.basetopic, mac, entityTypeName)
}

func (p *ClientPermit) getCommandTopic(mac string) string {
	mac = strings.Replace(mac, ":", "_", -1)
	return fmt.Sprintf("%s/%s_%s/command", p.basetopic, mac, entityTypeName)
}

// SendState sends client permit state.
func (p *ClientPermit) SendState(client dto.Client) {
	msg := offPayload
	if client.Permit {
		msg = onPayload
	}
	p.mqtt.SendMessage(p.getStateTopic(client.Mac), msg)
}

// RunConsumer runs client permit consumer.
func (p *ClientPermit) RunConsumer(mac string) {
	ch := p.mqtt.Subscribe(p.getCommandTopic(mac))

	for {
		message := <-ch
		permit := true
		if message == offPayload {
			permit = false
		}
		if err := p.accessUpdate.SetPermit(mac, permit); err != nil {
			p.logger.Error("client error while setting permit", "error", err)
		}
	}
}
