package homeassistant

import (
	"fmt"
	"log/slog"
	"time"

	"keeneticToMqtt/internal/dto"
	"keeneticToMqtt/internal/homeassistant/clientpolicy"
	"keeneticToMqtt/internal/services/clientlist"
)

type EntityManager struct {
	clientPolicy    *clientpolicy.ClientPolicy
	clientList      *clientlist.ClientList
	pollingInterval time.Duration
	logger          *slog.Logger
	clients         map[string]dto.Client
}

func NewEntityManager(
	policy *clientpolicy.ClientPolicy,
	clientList *clientlist.ClientList,
	pollingInterval time.Duration,
	logger *slog.Logger,
) *EntityManager {
	return &EntityManager{
		clientPolicy:    policy,
		clientList:      clientList,
		pollingInterval: pollingInterval,
		logger:          logger,
		clients:         map[string]dto.Client{},
	}
}

func (m *EntityManager) update() error {
	clients, err := m.clientList.GetClientList()

	m.logger.Info("entity manager update", "clients", clients)

	if err != nil {
		return fmt.Errorf("error while get clieent list update: %s", err.Error())
	}

	for _, client := range clients {
		_, ok := m.clients[client.Mac]
		if !ok {
			go m.clientPolicy.RunConsumer(client.Mac)
			go func() {
				if err := m.clientPolicy.SendDiscoveryMessage(client.Mac, client.Name); err != nil {
					m.logger.Error("entity manager update error while sending discovery message", "error", err, "client", client)
				}
			}()
		}
		m.clients[client.Mac] = client
		m.clientPolicy.SendState(client)
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
