package mqtt

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	mock_mqtt "keeneticToMqtt/test/mocks/gomock/clients/mqtt"
	mock_mqtttoken "keeneticToMqtt/test/mocks/gomock/clients/mqtt/token"
)

func TestClient_SendMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const (
		prefix   = "prefix"
		topic    = "topic"
		message  = "message"
		retained = true
	)

	someErr := errors.New("some err")

	tests := []struct {
		name       string
		mqttClient func() mqttClient
		logger     func() logger
	}{
		{
			name: "success publish",
			logger: func() logger {
				log := mock_mqtt.NewMocklogger(ctrl)
				log.EXPECT().Debug("start sending mqtt message",
					"topic", topic,
					"message", message,
					"retained", retained,
				)
				log.EXPECT().Info("sending mqtt message",
					"topic", topic,
					"message", message,
					"retained", retained,
				)

				return log
			},
			mqttClient: func() mqttClient {
				ch := make(chan struct{})
				go func() {
					ch <- struct{}{}
				}()

				token := mock_mqtttoken.NewMockToken(ctrl)
				token.EXPECT().Done().Return(ch)
				token.EXPECT().Error().Return(nil)

				client := mock_mqtt.NewMockmqttClient(ctrl)
				client.EXPECT().Publish(topic, byte(0), retained, message).Return(token)

				return client
			},
		},
		{
			name: "publish error",
			logger: func() logger {
				log := mock_mqtt.NewMocklogger(ctrl)
				log.EXPECT().Debug("start sending mqtt message",
					"topic", topic,
					"message", message,
					"retained", retained,
				)
				log.EXPECT().Error("error sending mqtt message",
					"error", someErr,
					"topic", topic,
					"message", message,
					"retained", retained,
				)

				return log
			},
			mqttClient: func() mqttClient {
				ch := make(chan struct{})
				go func() {
					ch <- struct{}{}
				}()

				token := mock_mqtttoken.NewMockToken(ctrl)
				token.EXPECT().Done().Return(ch)
				token.EXPECT().Error().Return(someErr)

				client := mock_mqtt.NewMockmqttClient(ctrl)
				client.EXPECT().Publish(topic, byte(0), retained, message).Return(token)

				return client
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mqtt := Client{
				topicPrefix: prefix,
				client:      tt.mqttClient(),
				logger:      tt.logger(),
				broker:      "",
			}

			mqtt.SendMessage(topic, message, retained)
		})
	}
}

func TestClient_Connect(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const (
		broker = "broker"
	)

	someErr := errors.New("some err")

	tests := []struct {
		name        string
		mqttClient  func() mqttClient
		logger      func() logger
		expectedErr error
	}{
		{
			name: "success connect",
			logger: func() logger {
				log := mock_mqtt.NewMocklogger(ctrl)
				log.EXPECT().Info("connected to mqtt", "broker", broker)

				return log
			},
			mqttClient: func() mqttClient {
				ch := make(chan struct{})
				go func() {
					ch <- struct{}{}
				}()

				token := mock_mqtttoken.NewMockToken(ctrl)
				token.EXPECT().Wait().Return(true)
				token.EXPECT().Error().Return(nil)

				client := mock_mqtt.NewMockmqttClient(ctrl)
				client.EXPECT().Connect().Return(token)

				return client
			},
		},
		{
			name: "error while connect",
			logger: func() logger {
				log := mock_mqtt.NewMocklogger(ctrl)
				return log
			},
			mqttClient: func() mqttClient {
				ch := make(chan struct{})
				go func() {
					ch <- struct{}{}
				}()

				token := mock_mqtttoken.NewMockToken(ctrl)
				token.EXPECT().Wait().Return(true)
				token.EXPECT().Error().Return(someErr)
				token.EXPECT().Error().Return(someErr)

				client := mock_mqtt.NewMockmqttClient(ctrl)
				client.EXPECT().Connect().Return(token)

				return client
			},
			expectedErr: someErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mqtt := Client{
				client: tt.mqttClient(),
				logger: tt.logger(),
				broker: broker,
			}

			err := mqtt.Connect()
			if tt.expectedErr == nil {
				assert.Nil(t, err)
			} else {
				assert.ErrorIs(t, err, tt.expectedErr)
			}
		})
	}
}

func TestClient_Subscribe(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const (
		topic = "topic"
	)

	client := mock_mqtt.NewMockmqttClient(ctrl)
	client.EXPECT().Subscribe(topic, byte(0), gomock.Any())

	mqtt := Client{
		client: client,
	}

	mqtt.Subscribe(topic)
}
