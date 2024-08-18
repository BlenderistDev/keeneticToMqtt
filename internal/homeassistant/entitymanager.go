package homeassistant

import (
	"fmt"
	"log/slog"
	"time"

	"keeneticToMqtt/internal/dto"
	"keeneticToMqtt/internal/services/clientlist"
)

type mqtt interface {
	Subscribe(topic string) chan string
	SendMessage(topic, message string)
}

type Entity interface {
	SendDiscoveryMessage(client dto.Client) error
	Consume(client dto.Client, message string) error
	GetCommandTopic(client dto.Client) string
	GetStateTopic(client dto.Client) string
	GetState(client dto.Client) (string, error)
}

type EntityManager struct {
	entities        []Entity
	clientList      *clientlist.ClientList
	mqtt            mqtt
	pollingInterval time.Duration
	logger          *slog.Logger
	clients         map[string]dto.Client
}

func NewEntityManager(
	entities []Entity,
	clientList *clientlist.ClientList,
	mqtt mqtt,
	pollingInterval time.Duration,
	logger *slog.Logger,
) *EntityManager {
	return &EntityManager{
		entities:        entities,
		clientList:      clientList,
		mqtt:            mqtt,
		pollingInterval: pollingInterval,
		logger:          logger,
		clients:         map[string]dto.Client{},
	}
}

func (m *EntityManager) Run() chan struct{} {
	done := make(chan struct{})
	ticker := time.NewTicker(m.pollingInterval)

	go func() {
		for {
			select {
			case <-done:
				m.logger.Info("shutdown entitymanager")
				return
			case _ = <-ticker.C:
				m.update()
			}
		}
	}()

	return done
}

func (m *EntityManager) update() {
	clients, err := m.clientList.GetClientList()
	if err != nil {
		m.logger.Error("Entity manager get state error", "error", err)
		return
	}
	m.logger.Info("Entity manager update", "clients", clients)

	for _, client := range clients {
		_, ok := m.clients[client.Mac]
		if !ok {
			m.runClient(client)
		}
		// update client because it can change
		m.clients[client.Mac] = client

		m.updateEntitiesState(client)
	}
	return
}

func (m *EntityManager) updateEntitiesState(client dto.Client) {
	for _, entity := range m.entities {
		if entity.GetStateTopic(client) != "" {
			state, err := entity.GetState(client)
			if err != nil {
				m.logger.Error("Entity manager get state error",
					"client", client,
					"entity", entity,
					"error", err,
				)
			}
			m.mqtt.SendMessage(entity.GetStateTopic(client), state)
		}
	}
}

func (m *EntityManager) runClient(client dto.Client) {
	for _, entity := range m.entities {
		e := entity
		if e.GetCommandTopic(client) != "" {
			go m.runClientEntityConsumer(e, client)
		}
		go m.sendDiscovery(client, e)
	}
}

func (m *EntityManager) runClientEntityConsumer(e Entity, client dto.Client) {
	ch := m.mqtt.Subscribe(e.GetCommandTopic(client))
	for {
		message := <-ch
		err := e.Consume(client, message)
		if err != nil {
			m.logger.Error("error while entity consume",
				"client", client,
				"message", message,
				"entity", e,
				"error", err,
			)
		}
		m.update()
	}
}

func (m *EntityManager) sendDiscovery(client dto.Client, e Entity) {
	if err := e.SendDiscoveryMessage(client); err != nil {
		m.logger.Error("Entity manager update error while sending discovery message",
			"error", err,
			"client", client,
			"Entity", fmt.Sprintf("%v", e),
		)
	}
}
