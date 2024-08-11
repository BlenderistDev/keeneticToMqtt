package discovery

import (
	"encoding/json"
	"fmt"

	"keeneticToMqtt/internal/clients/mqtt"
)

const (
	defaultDiscoveryPrefix = "homeassistant"
)

type Device struct {
	Manufacturer string `json:"manufacturer"`
	Name         string `json:"name"`
}

type Discovery struct {
	discoveryPrefix string
	mqtt            *mqtt.Client
}

func NewDiscovery(
	discoveryPrefix string,
	mqtt *mqtt.Client,
) *Discovery {
	if discoveryPrefix == "" {
		discoveryPrefix = defaultDiscoveryPrefix
	}

	return &Discovery{
		discoveryPrefix: discoveryPrefix,
		mqtt:            mqtt,
	}
}

func (d *Discovery) SendDiscoverySelect(commandTopic, stateTopic, mac, name string, options []string) error {
	config := struct {
		CommandTopic string   `json:"command_topic"`
		StateTopic   string   `json:"state_topic"`
		Name         string   `json:"name"`
		Options      []string `json:"options"`
		Device       Device
	}{
		CommandTopic: commandTopic,
		StateTopic:   stateTopic,
		Name:         name,
		Options:      options,
		Device: Device{
			Manufacturer: "BlenderistDev keeneticToMqtt",
			Name:         mac,
		},
	}

	configStr, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("error while marshal select discovery config: %w", err)
	}
	d.SendDiscovery("select", mac, string(configStr))

	return nil
}

func (d *Discovery) SendDiscovery(component, deviceID, config string) {
	d.mqtt.SendMessage(
		d.buildDiscoveryTopic(component, deviceID),
		config,
	)
}

func (d Discovery) buildDiscoveryTopic(component, deviceID string) string {
	return d.discoveryPrefix + "/" + component + "/" + deviceID + "/" + "config"
}
