package discovery

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	mock_discovery "keeneticToMqtt/test/mocks/gomock/services/discovery"
)

func TestDiscovery_SendDiscoverySelect(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const (
		commandTopic    = "commandTopic"
		stateTopic      = "stateTopic"
		deviceName      = "deviceName"
		entityName      = "entityName"
		discoveryPrefix = "discoveryPrefix"
		deviceID        = "deviceID"
	)
	options := []string{"option1", "option2"}

	tests := []struct {
		name                                             string
		commandTopic, stateTopic, deviceName, entityName string
		discoveryPrefix, deviceID                        string
		options                                          []string
		mqttClient                                       func() mqttClient
		expectedErr                                      error
	}{
		{
			name: "success sending select discovery message",
			mqttClient: func() mqttClient {
				client := mock_discovery.NewMockmqttClient(ctrl)
				client.EXPECT().SendMessage(
					gomock.Eq("discoveryPrefix/select/deviceIDentityName/config"),
					gomock.Eq("{\"command_topic\":\"commandTopic\",\"state_topic\":\"stateTopic\",\"name\":\"entityName\",\"options\":[\"option1\",\"option2\"],\"Device\":{\"manufacturer\":\"BlenderistDev keeneticToMqtt\",\"name\":\"deviceName\"}}"),
					gomock.Eq(true),
				)

				return client
			},
			commandTopic:    commandTopic,
			stateTopic:      stateTopic,
			deviceName:      deviceName,
			entityName:      entityName,
			options:         options,
			deviceID:        deviceID,
			discoveryPrefix: discoveryPrefix,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			discovery := NewDiscovery(tt.discoveryPrefix, tt.deviceID, tt.mqttClient())
			err := discovery.SendDiscoverySelect(tt.commandTopic, tt.stateTopic, tt.deviceName, tt.entityName, tt.options)
			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestDiscovery_SendDiscoverySwitch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const (
		commandTopic    = "commandTopic"
		stateTopic      = "stateTopic"
		deviceName      = "deviceName"
		entityName      = "entityName"
		discoveryPrefix = "discoveryPrefix"
		deviceID        = "deviceID"
	)
	options := []string{"option1", "option2"}

	tests := []struct {
		name                                             string
		commandTopic, stateTopic, deviceName, entityName string
		discoveryPrefix, deviceID                        string
		options                                          []string
		mqttClient                                       func() mqttClient
		expectedErr                                      error
	}{
		{
			name: "success sending switch discovery message",
			mqttClient: func() mqttClient {
				client := mock_discovery.NewMockmqttClient(ctrl)
				client.EXPECT().SendMessage(
					gomock.Eq("discoveryPrefix/switch/deviceIDentityName/config"),
					gomock.Eq("{\"command_topic\":\"commandTopic\",\"state_topic\":\"stateTopic\",\"name\":\"entityName\",\"Device\":{\"manufacturer\":\"BlenderistDev keeneticToMqtt\",\"name\":\"deviceName\"}}"),
					gomock.Eq(true),
				)

				return client
			},
			commandTopic:    commandTopic,
			stateTopic:      stateTopic,
			deviceName:      deviceName,
			entityName:      entityName,
			options:         options,
			deviceID:        deviceID,
			discoveryPrefix: discoveryPrefix,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			discovery := NewDiscovery(tt.discoveryPrefix, tt.deviceID, tt.mqttClient())
			err := discovery.SendDiscoverySwitch(tt.commandTopic, tt.stateTopic, tt.deviceName, tt.entityName)
			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestDiscovery_SendDiscoverySensor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const (
		commandTopic    = "commandTopic"
		stateTopic      = "stateTopic"
		deviceName      = "deviceName"
		entityName      = "entityName"
		discoveryPrefix = "discoveryPrefix"
		deviceID        = "deviceID"
		unit            = "unit"
	)
	options := []string{"option1", "option2"}

	tests := []struct {
		name                                             string
		commandTopic, stateTopic, deviceName, entityName string
		discoveryPrefix, deviceID                        string
		options                                          []string
		mqttClient                                       func() mqttClient
		expectedErr                                      error
	}{
		{
			name: "success sending sensor discovery message",
			mqttClient: func() mqttClient {
				client := mock_discovery.NewMockmqttClient(ctrl)
				client.EXPECT().SendMessage(
					gomock.Eq("discoveryPrefix/sensor/deviceIDentityName/config"),
					gomock.Eq("{\"state_topic\":\"stateTopic\",\"name\":\"entityName\",\"Device\":{\"manufacturer\":\"BlenderistDev keeneticToMqtt\",\"name\":\"deviceName\"},\"unit_of_measurement\":\"unit\"}"),
					gomock.Eq(true),
				)

				return client
			},
			commandTopic:    commandTopic,
			stateTopic:      stateTopic,
			deviceName:      deviceName,
			entityName:      entityName,
			options:         options,
			deviceID:        deviceID,
			discoveryPrefix: discoveryPrefix,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			discovery := NewDiscovery(tt.discoveryPrefix, tt.deviceID, tt.mqttClient())
			err := discovery.SendDiscoverySensor(tt.stateTopic, tt.deviceName, tt.entityName, unit)
			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestNewDiscovery_emptyDiscoveryPrefix(t *testing.T) {
	discovery := NewDiscovery("", "", nil)
	assert.Equal(t, defaultDiscoveryPrefix, discovery.discoveryPrefix)
}
