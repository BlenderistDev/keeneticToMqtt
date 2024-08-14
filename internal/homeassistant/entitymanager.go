package homeassistant

import (
	"fmt"
	"log/slog"
	"time"

	"keeneticToMqtt/internal/dto"
	"keeneticToMqtt/internal/services/clientlist"
)

type Entity interface {
	RunConsumer(mac string)
	SendDiscoveryMessage(mac, name string) error
	SendState(client dto.Client)
}

type EntityManager struct {
	entities        []Entity
	clientList      *clientlist.ClientList
	pollingInterval time.Duration
	logger          *slog.Logger
	clients         map[string]dto.Client
}

func NewEntityManager(
	entities []Entity,
	clientList *clientlist.ClientList,
	pollingInterval time.Duration,
	logger *slog.Logger,
) *EntityManager {
	return &EntityManager{
		entities:        entities,
		clientList:      clientList,
		pollingInterval: pollingInterval,
		logger:          logger,
		clients:         map[string]dto.Client{},
	}
}

func (m *EntityManager) update() error {
	clients, err := m.clientList.GetClientList()

	m.logger.Info("Entity manager update", "clients", clients)

	if err != nil {
		return fmt.Errorf("error while get clieent list update: %s", err.Error())
	}

	for _, client := range clients {

		_, ok := m.clients[client.Mac]
		if !ok {
			for _, entity := range m.entities {
				e := entity
				go e.RunConsumer(client.Mac)
				go func() {
					if err := e.SendDiscoveryMessage(client.Mac, client.Name); err != nil {
						m.logger.Error("Entity manager update error while sending discovery message",
							"error", err,
							"client", client,
							"Entity", fmt.Sprintf("%v", e),
						)
					}
				}()
			}
		}
		m.clients[client.Mac] = client
		for _, entity := range m.entities {

			entity.SendState(client)
		}

	}

	return nil
}

func (m *EntityManager) Run() chan bool {
	ticker := time.NewTicker(m.pollingInterval)

	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case _ = <-ticker.C:
				if err := m.update(); err != nil {
					m.logger.Error("error while update entitymanager", "error", err.Error())
				}
			}
		}
	}()

	return done
}
