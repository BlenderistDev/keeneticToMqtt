package mqtt

import (
	"fmt"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Client struct {
	topicPrefix string
	client      mqtt.Client
}

func NewClient(broker, clientID, username, password string) *Client {
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
	}
}

func (c *Client) SendMessage(topic, message string) {
	token := c.client.Publish(topic, 0, false, message)
	token.Wait()
}

func (c *Client) Subscribe(topic string) chan string {
	ch := make(chan string)
	c.client.Subscribe(topic, 0, func(client mqtt.Client, message mqtt.Message) {
		fmt.Println(topic, message)
		ch <- string(message.Payload())
	})

	return ch
}
