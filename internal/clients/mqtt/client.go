package mqtt

import (
	"fmt"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type logger interface {
	Error(msg string, args ...any)
	Info(msg string, args ...any)
	Debug(msg string, args ...any)
}

type Client struct {
	topicPrefix string
	client      mqtt.Client
	logger      logger
}

func NewClient(broker, clientID, username, password string, log logger) *Client {
	var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("TOPIC: %s\n", msg.Topic())
		fmt.Printf("MSG: %s\n", msg.Payload())
	}

	opts := mqtt.NewClientOptions().AddBroker(broker).SetClientID(clientID)
	opts.SetKeepAlive(2 * time.Second)
	opts.SetDefaultPublishHandler(f)
	opts.SetPingTimeout(1 * time.Second)
	opts.SetUsername(username)
	opts.SetPassword(password)

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	if token := c.Subscribe("go-mqtt/sample", 0, func(client mqtt.Client, message mqtt.Message) {
		fmt.Printf(string(message.Payload()))
	}); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	return &Client{
		client: c,
		logger: log,
	}
}

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

func (c *Client) Subscribe(topic string) chan string {
	ch := make(chan string)
	c.client.Subscribe(topic, 0, func(client mqtt.Client, message mqtt.Message) {
		fmt.Println(topic, message)
		ch <- string(message.Payload())
	})

	return ch
}
