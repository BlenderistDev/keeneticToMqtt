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
	policyStorage interface {
		GetPolicyList() []string
	}
)

// ClientPolicy struct for handle home assistant client policy entities.
type ClientPolicy struct {
	basetopic       string
	discoveryClient discovery
	accessUpdate    accessUpdate
	policyStorage   policyStorage
}

// NewClientPolicy creates new ClientPolicy.
func NewClientPolicy(
	basetopic string,
	discoveryClient discovery,
	accessUpdate accessUpdate,
	policyStorage policyStorage,
) *ClientPolicy {
	return &ClientPolicy{
		basetopic:       basetopic,
		discoveryClient: discoveryClient,
		accessUpdate:    accessUpdate,
		policyStorage:   policyStorage,
	}
}

// SendDiscoveryMessage sends homeassistant discovery message.
func (p *ClientPolicy) SendDiscoveryMessage(client dto.Client) error {
	commandTopic := p.GetCommandTopic(client)
	stateTopic := p.GetStateTopic(client)
	policies := p.policyStorage.GetPolicyList()

	if err := p.discoveryClient.SendDiscoverySelect(commandTopic, stateTopic, client.Name, client.Name+"_"+entityTypeName, policies); err != nil {
		return fmt.Errorf("ClientPolicy SendDiscoveryMessage error: %w", err)
	}

	return nil
}

// GetState returns entity state.
func (p *ClientPolicy) GetState(client dto.Client) (string, error) {
	return client.Policy, nil
}

// Consume consumes message.
func (p *ClientPolicy) Consume(client dto.Client, message string) error {
	if err := p.accessUpdate.SetPolicy(client.Mac, message); err != nil {
		return fmt.Errorf("client error while setting policy: %w", err)
	}
	return nil
}

// GetStateTopic returns state topic.
func (p *ClientPolicy) GetStateTopic(client dto.Client) string {
	mac := strings.Replace(client.Mac, ":", "_", -1)
	return fmt.Sprintf("%s/%s_%s/state", p.basetopic, mac, entityTypeName)
}

// GetCommandTopic returns command topic.
func (p *ClientPolicy) GetCommandTopic(client dto.Client) string {
	mac := strings.Replace(client.Mac, ":", "_", -1)
	return fmt.Sprintf("%s/%s_%s/command", p.basetopic, mac, entityTypeName)
}
