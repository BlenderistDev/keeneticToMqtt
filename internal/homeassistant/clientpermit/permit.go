package clientpolicy

import (
	"fmt"
	"log/slog"
	"strings"

	"keeneticToMqtt/internal/clients/keenetic/accessupdate"
	"keeneticToMqtt/internal/clients/mqtt"
	"keeneticToMqtt/internal/dto"
	discoverycilient "keeneticToMqtt/internal/services/discovery"
)

const (
	entityTypeName = "permit"
	offPayload     = "OFF"
	onPayload      = "ON"
)

type ClientPermit struct {
	basetopic       string
	discoveryClient *discoverycilient.Discovery
	accessUpdate    *accessupdate.AccessUpdate
	mqtt            *mqtt.Client
	logger          *slog.Logger
}

func NewClientPermit(
	basetopic string,
	discoveryClient *discoverycilient.Discovery,
	mqtt *mqtt.Client,
	accessUpdate *accessupdate.AccessUpdate,
	logger *slog.Logger,
) *ClientPermit {
	return &ClientPermit{
		basetopic:       basetopic,
		discoveryClient: discoveryClient,
		accessUpdate:    accessUpdate,
		mqtt:            mqtt,
		logger:          logger,
	}
}

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

func (p *ClientPermit) SendState(client dto.Client) {
	msg := offPayload
	if client.Permit {
		msg = onPayload
	}
	p.mqtt.SendMessage(p.getStateTopic(client.Mac), msg)
}

func (p *ClientPermit) RunConsumer(mac string) {
	ch := p.mqtt.Subscribe(p.getCommandTopic(mac))

	for {
		message := <-ch
		permit := true
		if message == offPayload {
			permit = false
		}
		if err := p.accessUpdate.SetPermit(mac, permit); err != nil {
			p.logger.Error("client error while setting permit", "error", err.Error())
		}
	}
}
