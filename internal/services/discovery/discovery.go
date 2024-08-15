package discovery

import (
	"encoding/json"
	"fmt"
)

//go:generate mockgen -source=discovery.go -destination=../../../test/mocks/gomock/services/discovery/discovery.go

const (
	defaultDiscoveryPrefix = "homeassistant"
)

type (
	mqttClient interface {
		SendMessage(topic, message string)
	}
	device struct {
		Manufacturer string `json:"manufacturer"`
		Name         string `json:"name"`
	}
	// Discovery struct to send home assistant discovery messages.
	Discovery struct {
		discoveryPrefix, deviceID string
		mqtt                      mqttClient
	}
)

// NewDiscovery creates new Discovery struct.
func NewDiscovery(
	discoveryPrefix, deviceID string,
	mqtt mqttClient,
) *Discovery {
	if discoveryPrefix == "" {
		discoveryPrefix = defaultDiscoveryPrefix
	}

	return &Discovery{
		discoveryPrefix: discoveryPrefix,
		deviceID:        deviceID,
		mqtt:            mqtt,
	}
}

// SendDiscoverySelect sends home assistant discovery message for switch.
func (d *Discovery) SendDiscoverySelect(commandTopic, stateTopic, deviceName, name string, options []string) error {
	config := struct {
		CommandTopic string   `json:"command_topic"`
		StateTopic   string   `json:"state_topic"`
		Name         string   `json:"name"`
		Options      []string `json:"options"`
		Device       device
	}{
		CommandTopic: commandTopic,
		StateTopic:   stateTopic,
		Name:         name,
		Options:      options,
		Device: device{
			Manufacturer: "BlenderistDev keeneticToMqtt",
			Name:         deviceName,
		},
	}

	configStr, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("error while marshal select discovery config: %w", err)
	}
	d.sendDiscovery("select", d.deviceID+name, string(configStr))

	return nil
}

// SendDiscoverySwitch sends home assistant discovery message for switch.
func (d *Discovery) SendDiscoverySwitch(commandTopic, stateTopic, deviceName, name string) error {
	config := struct {
		CommandTopic string `json:"command_topic"`
		StateTopic   string `json:"state_topic"`
		Name         string `json:"name"`
		Device       device
	}{
		CommandTopic: commandTopic,
		StateTopic:   stateTopic,
		Name:         name,
		Device: device{
			Manufacturer: "BlenderistDev keeneticToMqtt",
			Name:         deviceName,
		},
	}

	configStr, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("error while marshal select discovery config: %w", err)
	}
	d.sendDiscovery("switch", d.deviceID+name, string(configStr))

	return nil
}

func (d *Discovery) sendDiscovery(component, deviceID, config string) {
	d.mqtt.SendMessage(
		d.buildDiscoveryTopic(component, deviceID),
		config,
	)
}

func (d Discovery) buildDiscoveryTopic(component, deviceID string) string {
	return d.discoveryPrefix + "/" + component + "/" + deviceID + "/" + "config"
}
