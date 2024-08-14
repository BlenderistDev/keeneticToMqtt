package clientpolicy

import (
	"fmt"
	"log/slog"
	"strings"

	"keeneticToMqtt/internal/clients/keenetic/policy"
	"keeneticToMqtt/internal/clients/mqtt"
	"keeneticToMqtt/internal/dto"
	discoverycilient "keeneticToMqtt/internal/services/discovery"
	policy2 "keeneticToMqtt/internal/storages/policy"
)

const (
	entityTypeName = "policy"
)

type ClientPolicy struct {
	basetopic       string
	discoveryClient *discoverycilient.Discovery
	mqtt            *mqtt.Client
	policyClient    *policy.Policy
	policyStorage   *policy2.Storage
	logger          *slog.Logger
}

func NewClientPolicy(
	basetopic string,
	discoveryClient *discoverycilient.Discovery,
	mqtt *mqtt.Client,
	policyClient *policy.Policy,
	policyStorage *policy2.Storage,
	logger *slog.Logger,
) *ClientPolicy {

	return &ClientPolicy{
		basetopic:       basetopic,
		discoveryClient: discoveryClient,
		mqtt:            mqtt,
		policyClient:    policyClient,
		policyStorage:   policyStorage,
		logger:          logger,
	}
}

func (p *ClientPolicy) SendDiscoveryMessage(mac, name string) error {
	commandTopic := p.getCommandTopic(mac)
	stateTopic := p.getStateTopic(mac)
	policies := p.policyStorage.GetPolicyList()

	if err := p.discoveryClient.SendDiscoverySelect(commandTopic, stateTopic, name, name+"_"+entityTypeName, policies); err != nil {
		return fmt.Errorf("ClientPolicy SendDiscoveryMessage error: %w", err)
	}

	return nil
}

func (p *ClientPolicy) getStateTopic(mac string) string {
	mac = strings.Replace(mac, ":", "_", -1)
	return fmt.Sprintf("%s/%s_%s/state", p.basetopic, mac, entityTypeName)
}

func (p *ClientPolicy) getCommandTopic(mac string) string {
	mac = strings.Replace(mac, ":", "_", -1)
	return fmt.Sprintf("%s/%s_%s/command", p.basetopic, mac, entityTypeName)
}

func (p *ClientPolicy) SendState(client dto.Client) {
	p.mqtt.SendMessage(p.getStateTopic(client.Mac), client.Policy)
}

func (p *ClientPolicy) RunConsumer(mac string) {
	ch := p.mqtt.Subscribe(p.getCommandTopic(mac))

	for {
		message := <-ch
		if err := p.policyClient.SetPolicy(mac, message); err != nil {
			p.logger.Error("client error while setting policy", "error", err.Error())
		}
	}
}
