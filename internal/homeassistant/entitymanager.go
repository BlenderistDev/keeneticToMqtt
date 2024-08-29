package homeassistant

import (
	"fmt"
	"sync"
	"time"

	"keeneticToMqtt/internal/dto"
)

//go:generate mockgen -source=entitymanager.go -destination=../../test/mocks/gomock/homeassistant/entitymanager.go

type mqtt interface {
	Subscribe(topic string) chan string
	SendMessage(topic, message string, retained bool)
}

// Entity home assistant entity.
type Entity interface {
	SendDiscoveryMessage(client dto.Client) error
	Consume(client dto.Client, message string) error
	GetCommandTopic(client dto.Client) string
	GetStateTopic(client dto.Client) string
	GetState(client dto.Client) (string, error)
}

type clientList interface {
	GetClientList() ([]dto.Client, error)
}

type logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
}

// EntityManager entity manager for keenetic client entities in home assistant.
type EntityManager struct {
	entities          []Entity
	clientList        clientList
	mqtt              mqtt
	pollingInterval   time.Duration
	logger            logger
	clients           map[string]dto.Client
	entityStates      map[string]map[string]string
	entityStatesMutex sync.RWMutex
}

// NewEntityManager creates new EntityManager.
func NewEntityManager(
	entities []Entity,
	clientList clientList,
	mqtt mqtt,
	pollingInterval time.Duration,
	logger logger,
) *EntityManager {
	return &EntityManager{
		entities:        entities,
		clientList:      clientList,
		mqtt:            mqtt,
		pollingInterval: pollingInterval,
		logger:          logger,
		clients:         map[string]dto.Client{},
		entityStates:    make(map[string]map[string]string),
	}
}

// Run entity updates and command consumer.
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

// updateEntitiesState sends mqtt messages with updates to state topic only if state changes.
func (m *EntityManager) updateEntitiesState(client dto.Client) {
	for _, entity := range m.entities {
		stateTopic := entity.GetStateTopic(client)
		if stateTopic != "" {
			state, err := entity.GetState(client)
			if err != nil {
				m.logger.Error("Entity manager get state error",
					"client", client,
					"entity", entity,
					"error", err,
				)
				return
			}
			m.entityStatesMutex.Lock()
			entityStorage, ok := m.entityStates[stateTopic]
			if ok {
				storageState, ok := entityStorage[client.Mac]
				if ok && storageState == state {
					continue
				}
			}
			if m.entityStates[stateTopic] == nil {
				m.entityStates[stateTopic] = make(map[string]string)
			}
			m.entityStates[stateTopic][client.Mac] = state
			m.entityStatesMutex.Unlock()
			m.mqtt.SendMessage(stateTopic, state, false)
		}
	}
}

func (m *EntityManager) runClient(client dto.Client) {
	for _, entity := range m.entities {
		e := entity
		go m.runClientEntityConsumer(e, client)
		go m.sendDiscovery(client, e)
	}
}

func (m *EntityManager) runClientEntityConsumer(e Entity, client dto.Client) {
	commandTopic := e.GetCommandTopic(client)
	if commandTopic == "" {
		return
	}
	ch := m.mqtt.Subscribe(commandTopic)
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
			"entity", fmt.Sprintf("%v", e),
		)
	}
}
