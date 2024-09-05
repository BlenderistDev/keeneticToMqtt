package mqtt

import (
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

//go:generate mockgen -source=client.go -destination=../../../test/mocks/gomock/clients/mqtt/client.go

type logger interface {
	Error(msg string, args ...any)
	Info(msg string, args ...any)
	Debug(msg string, args ...any)
}

type mqttClient interface {
	Connect() mqtt.Token
	Publish(topic string, qos byte, retained bool, payload interface{}) mqtt.Token
	Subscribe(topic string, qos byte, callback mqtt.MessageHandler) mqtt.Token
}

// Client mqtt client.
type Client struct {
	topicPrefix string
	client      mqttClient
	logger      logger
	broker      string
}

// NewClient creates new Client.
func NewClient(broker, clientID, username, password string, log logger) *Client {
	opts := mqtt.
		NewClientOptions().
		AddBroker(broker).
		SetClientID(clientID).
		SetKeepAlive(2 * time.Second).
		SetPingTimeout(1 * time.Second).
		SetUsername(username).
		SetPassword(password)

	c := mqtt.NewClient(opts)

	return &Client{
		client: c,
		logger: log,
		broker: broker,
	}
}

// Connect connection to mqtt broker.
func (c *Client) Connect() error {
	if token := c.client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	c.logger.Info("connected to mqtt", "broker", c.broker)

	return nil
}

// SendMessage sends mqtt message.
func (c *Client) SendMessage(topic, message string, retained bool) {
	c.logger.Debug("start sending mqtt message",
		"topic", topic,
		"message", message,
		"retained", retained,
	)

	token := c.client.Publish(topic, 0, retained, message)
	<-token.Done()
	if err := token.Error(); err != nil {
		c.logger.Error("error sending mqtt message",
			"error", err,
			"topic", topic,
			"message", message,
			"retained", retained,
		)
		return
	}

	c.logger.Info("sending mqtt message",
		"topic", topic,
		"message", message,
		"retained", retained,
	)
}

// Subscribe subscribes to topic.
func (c *Client) Subscribe(topic string) chan string {
	ch := make(chan string)
	c.client.Subscribe(topic, 0, func(client mqtt.Client, message mqtt.Message) {
		ch <- string(message.Payload())
	})

	return ch
}
