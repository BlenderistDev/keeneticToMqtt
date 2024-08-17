package clientpolicy

import (
	"fmt"
	"strings"

	"keeneticToMqtt/internal/dto"
)

//go:generate mockgen -source=policy.go -destination=../../../test/mocks/gomock/homeassistant/clientpolicy/policy.go

const (
	entityTypeName = "policy"
)

type (
	discovery interface {
		SendDiscoverySelect(commandTopic, stateTopic, deviceName, name string, options []string) error
	}
	accessUpdate interface {
		SetPolicy(mac, policy string) error
	}
	mqtt interface {
		Subscribe(topic string) chan string
		SendMessage(topic, message string)
	}
	policyStorage interface {
		GetPolicyList() []string
	}

	logger interface {
		Error(msg string, args ...any)
	}
)

// ClientPolicy struct for handle home assistant client policy entities.
type ClientPolicy struct {
	basetopic       string
	discoveryClient discovery
	mqtt            mqtt
	accessUpdate    accessUpdate
	policyStorage   policyStorage
	logger          logger
}

// NewClientPolicy creates new ClientPolicy.
func NewClientPolicy(
	basetopic string,
	discoveryClient discovery,
	mqtt mqtt,
	accessUpdate accessUpdate,
	policyStorage policyStorage,
	logger logger,
) *ClientPolicy {
	return &ClientPolicy{
		basetopic:       basetopic,
		discoveryClient: discoveryClient,
		mqtt:            mqtt,
		accessUpdate:    accessUpdate,
		policyStorage:   policyStorage,
		logger:          logger,
	}
}

// SendDiscoveryMessage sends homeassistant discovery message.
func (p *ClientPolicy) SendDiscoveryMessage(mac, name string) error {
	commandTopic := p.getCommandTopic(mac)
	stateTopic := p.getStateTopic(mac)
	policies := p.policyStorage.GetPolicyList()

	if err := p.discoveryClient.SendDiscoverySelect(commandTopic, stateTopic, name, name+"_"+entityTypeName, policies); err != nil {
		return fmt.Errorf("ClientPolicy SendDiscoveryMessage error: %w", err)
	}

	return nil
}

// SendState sends state to mqtt.
func (p *ClientPolicy) SendState(client dto.Client) {
	p.mqtt.SendMessage(p.getStateTopic(client.Mac), client.Policy)
}

// RunConsumer runs mqtt consumer.
func (p *ClientPolicy) RunConsumer(mac string) {
	ch := p.mqtt.Subscribe(p.getCommandTopic(mac))

	for {
		message := <-ch
		if err := p.accessUpdate.SetPolicy(mac, message); err != nil {
			p.logger.Error("client error while setting policy", "error", err)
		}
	}
}

func (p *ClientPolicy) getStateTopic(mac string) string {
	mac = strings.Replace(mac, ":", "_", -1)
	return fmt.Sprintf("%s/%s_%s/state", p.basetopic, mac, entityTypeName)
}

func (p *ClientPolicy) getCommandTopic(mac string) string {
	mac = strings.Replace(mac, ":", "_", -1)
	return fmt.Sprintf("%s/%s_%s/command", p.basetopic, mac, entityTypeName)
}
