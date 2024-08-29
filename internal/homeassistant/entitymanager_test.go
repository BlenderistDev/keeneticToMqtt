package homeassistant

import (
	"errors"
	"testing"
	"time"

	"go.uber.org/mock/gomock"
	"keeneticToMqtt/internal/dto"
	mock_homeassistant "keeneticToMqtt/test/mocks/gomock/homeassistant"
)

func TestEntityManager_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const (
		mac             = "mac"
		macNew          = "macNew"
		stateTopic      = "stateTopic"
		stateTopicNew   = "stateTopicNew"
		commandTopicNew = "commandTopicNew"
		storageState    = "storageState"
		state           = "state"
		stateNew        = "stateNew"
		command         = "command"
	)

	clientDto := dto.Client{Mac: mac}
	clientDtoNew := dto.Client{Mac: macNew}

	clients := []dto.Client{
		clientDto,
	}
	clientsWithNew := []dto.Client{
		clientDto, clientDtoNew,
	}
	clientsOnlyNew := []dto.Client{
		clientDtoNew,
	}

	someErr := errors.New("some error")
	tests := []struct {
		name         string
		entities     func() []Entity
		clientList   func() clientList
		mqtt         func() mqtt
		logger       func() logger
		clients      map[string]dto.Client
		entityStates map[string]map[string]string
	}{
		{
			name: "success update",
			entities: func() []Entity {
				entity := mock_homeassistant.NewMockEntity(ctrl)
				entity.EXPECT().GetStateTopic(clientDto).Return(stateTopic)
				entity.EXPECT().GetState(clientDto).Return(state, nil)
				return []Entity{entity}
			},
			clientList: func() clientList {
				clientList := mock_homeassistant.NewMockclientList(ctrl)
				clientList.EXPECT().GetClientList().Return(clients, nil)
				return clientList
			},
			mqtt: func() mqtt {
				mqtt := mock_homeassistant.NewMockmqtt(ctrl)
				mqtt.EXPECT().SendMessage(stateTopic, state, false)
				return mqtt
			},
			logger: func() logger {
				logger := mock_homeassistant.NewMocklogger(ctrl)
				logger.EXPECT().Info("Entity manager update", "clients", clients)
				logger.EXPECT().Info("shutdown entitymanager")
				return logger
			},
			clients: map[string]dto.Client{
				mac: clientDto,
			},
			entityStates: map[string]map[string]string{},
		},
		{
			name: "state is same as in storage",
			entities: func() []Entity {
				entity := mock_homeassistant.NewMockEntity(ctrl)
				entity.EXPECT().GetStateTopic(clientDto).Return(stateTopic)
				entity.EXPECT().GetState(clientDto).Return(storageState, nil)
				return []Entity{entity}
			},
			clientList: func() clientList {
				clientList := mock_homeassistant.NewMockclientList(ctrl)
				clientList.EXPECT().GetClientList().Return(clients, nil)
				return clientList
			},
			mqtt: func() mqtt {
				mqtt := mock_homeassistant.NewMockmqtt(ctrl)
				return mqtt
			},
			logger: func() logger {
				logger := mock_homeassistant.NewMocklogger(ctrl)
				logger.EXPECT().Info("Entity manager update", "clients", clients)
				logger.EXPECT().Info("shutdown entitymanager")
				return logger
			},
			clients: map[string]dto.Client{
				mac: clientDto,
			},
			entityStates: map[string]map[string]string{
				stateTopic: {mac: storageState},
			},
		},
		{
			name: "error while getting state from entity",
			entities: func() []Entity {
				entity := mock_homeassistant.NewMockEntity(ctrl)
				entity.EXPECT().GetStateTopic(clientDto).Return(stateTopic)
				entity.EXPECT().GetState(clientDto).Return("", someErr)
				return []Entity{entity}
			},
			clientList: func() clientList {
				clientList := mock_homeassistant.NewMockclientList(ctrl)
				clientList.EXPECT().GetClientList().Return(clients, nil)
				return clientList
			},
			mqtt: func() mqtt {
				mqtt := mock_homeassistant.NewMockmqtt(ctrl)
				return mqtt
			},
			logger: func() logger {
				logger := mock_homeassistant.NewMocklogger(ctrl)
				logger.EXPECT().Info("Entity manager update", "clients", clients)
				logger.EXPECT().Error("Entity manager get state error", "client", clientDto, "entity", gomock.Any(), "error", someErr)
				logger.EXPECT().Info("shutdown entitymanager")
				return logger
			},
			clients: map[string]dto.Client{
				mac: clientDto,
			},
			entityStates: map[string]map[string]string{
				stateTopic: {mac: storageState},
			},
		},
		{
			name: "error while getting client list",
			entities: func() []Entity {
				entity := mock_homeassistant.NewMockEntity(ctrl)
				return []Entity{entity}
			},
			clientList: func() clientList {
				clientList := mock_homeassistant.NewMockclientList(ctrl)
				clientList.EXPECT().GetClientList().Return(nil, someErr)
				return clientList
			},
			mqtt: func() mqtt {
				mqtt := mock_homeassistant.NewMockmqtt(ctrl)
				return mqtt
			},
			logger: func() logger {
				logger := mock_homeassistant.NewMocklogger(ctrl)
				logger.EXPECT().Error("Entity manager get state error", "error", someErr)
				logger.EXPECT().Info("shutdown entitymanager")
				return logger
			},
			clients: map[string]dto.Client{
				mac: clientDto,
			},
			entityStates: map[string]map[string]string{
				stateTopic: {mac: storageState},
			},
		},
		{
			name: "add second client and success update",
			entities: func() []Entity {
				entity := mock_homeassistant.NewMockEntity(ctrl)
				entity.EXPECT().GetStateTopic(clientDto).Return(stateTopic)
				entity.EXPECT().GetStateTopic(clientDtoNew).Return(stateTopicNew)
				entity.EXPECT().GetCommandTopic(clientDtoNew).Return(commandTopicNew)
				entity.EXPECT().GetState(clientDto).Return(state, nil)
				entity.EXPECT().GetState(clientDtoNew).Return(stateNew, nil)
				entity.EXPECT().SendDiscoveryMessage(clientDtoNew).Return(nil)
				return []Entity{entity}
			},
			clientList: func() clientList {
				clientList := mock_homeassistant.NewMockclientList(ctrl)
				clientList.EXPECT().GetClientList().Return(clientsWithNew, nil)
				return clientList
			},
			mqtt: func() mqtt {
				ch := make(chan string)
				mqtt := mock_homeassistant.NewMockmqtt(ctrl)
				mqtt.EXPECT().SendMessage(stateTopic, state, false)
				mqtt.EXPECT().SendMessage(stateTopicNew, stateNew, false)
				mqtt.EXPECT().Subscribe(commandTopicNew).Return(ch)
				return mqtt
			},
			logger: func() logger {
				logger := mock_homeassistant.NewMocklogger(ctrl)
				logger.EXPECT().Info("Entity manager update", "clients", clientsWithNew)
				logger.EXPECT().Info("shutdown entitymanager")
				return logger
			},
			clients: map[string]dto.Client{
				mac: clientDto,
			},
			entityStates: map[string]map[string]string{},
		},
		{
			name: "add new and consume",
			entities: func() []Entity {
				entity := mock_homeassistant.NewMockEntity(ctrl)
				entity.EXPECT().GetStateTopic(clientDtoNew).Return(stateTopicNew)
				entity.EXPECT().GetStateTopic(clientDtoNew).Return(stateTopicNew)
				entity.EXPECT().GetCommandTopic(clientDtoNew).Return(commandTopicNew)
				entity.EXPECT().GetState(clientDtoNew).Return(stateNew, nil)
				entity.EXPECT().GetState(clientDtoNew).Return(stateNew, nil)
				entity.EXPECT().SendDiscoveryMessage(clientDtoNew).Return(nil)
				entity.EXPECT().Consume(clientDtoNew, command).Return(nil)
				return []Entity{entity}
			},
			clientList: func() clientList {
				clientList := mock_homeassistant.NewMockclientList(ctrl)
				clientList.EXPECT().GetClientList().Return(clientsOnlyNew, nil)
				clientList.EXPECT().GetClientList().Return(clientsOnlyNew, nil)
				return clientList
			},
			mqtt: func() mqtt {
				ch := make(chan string)
				go func() {
					ch <- command
				}()
				mqtt := mock_homeassistant.NewMockmqtt(ctrl)
				mqtt.EXPECT().SendMessage(stateTopicNew, stateNew, false)
				mqtt.EXPECT().Subscribe(commandTopicNew).Return(ch)
				return mqtt
			},
			logger: func() logger {
				logger := mock_homeassistant.NewMocklogger(ctrl)
				logger.EXPECT().Info("Entity manager update", "clients", clientsOnlyNew)
				logger.EXPECT().Info("Entity manager update", "clients", clientsOnlyNew)
				logger.EXPECT().Info("shutdown entitymanager")
				return logger
			},
			clients:      map[string]dto.Client{},
			entityStates: map[string]map[string]string{},
		},
		{
			name: "add new and error in consume",
			entities: func() []Entity {
				entity := mock_homeassistant.NewMockEntity(ctrl)
				entity.EXPECT().GetStateTopic(clientDtoNew).Return(stateTopicNew)
				entity.EXPECT().GetStateTopic(clientDtoNew).Return(stateTopicNew)
				entity.EXPECT().GetCommandTopic(clientDtoNew).Return(commandTopicNew)
				entity.EXPECT().GetState(clientDtoNew).Return(stateNew, nil)
				entity.EXPECT().GetState(clientDtoNew).Return(stateNew, nil)
				entity.EXPECT().SendDiscoveryMessage(clientDtoNew).Return(nil)
				entity.EXPECT().Consume(clientDtoNew, command).Return(someErr)
				return []Entity{entity}
			},
			clientList: func() clientList {
				clientList := mock_homeassistant.NewMockclientList(ctrl)
				clientList.EXPECT().GetClientList().Return(clientsOnlyNew, nil)
				clientList.EXPECT().GetClientList().Return(clientsOnlyNew, nil)
				return clientList
			},
			mqtt: func() mqtt {
				ch := make(chan string)
				go func() {
					ch <- command
				}()
				mqtt := mock_homeassistant.NewMockmqtt(ctrl)
				mqtt.EXPECT().SendMessage(stateTopicNew, stateNew, false)
				mqtt.EXPECT().Subscribe(commandTopicNew).Return(ch)
				return mqtt
			},
			logger: func() logger {
				logger := mock_homeassistant.NewMocklogger(ctrl)
				logger.EXPECT().Info("Entity manager update", "clients", clientsOnlyNew)
				logger.EXPECT().Info("Entity manager update", "clients", clientsOnlyNew)
				logger.EXPECT().Error("error while entity consume", "client", clientDtoNew, "message", command, "entity", gomock.Any(), "error", someErr)
				logger.EXPECT().Info("shutdown entitymanager")
				return logger
			},
			clients:      map[string]dto.Client{},
			entityStates: map[string]map[string]string{},
		},
		{
			name: "add second client without command topic and success update",
			entities: func() []Entity {
				entity := mock_homeassistant.NewMockEntity(ctrl)
				entity.EXPECT().GetStateTopic(clientDto).Return(stateTopic)
				entity.EXPECT().GetStateTopic(clientDtoNew).Return(stateTopicNew)
				entity.EXPECT().GetCommandTopic(clientDtoNew).Return("")
				entity.EXPECT().GetState(clientDto).Return(state, nil)
				entity.EXPECT().GetState(clientDtoNew).Return(stateNew, nil)
				entity.EXPECT().SendDiscoveryMessage(clientDtoNew).Return(nil)
				return []Entity{entity}
			},
			clientList: func() clientList {
				clientList := mock_homeassistant.NewMockclientList(ctrl)
				clientList.EXPECT().GetClientList().Return(clientsWithNew, nil)
				return clientList
			},
			mqtt: func() mqtt {
				mqtt := mock_homeassistant.NewMockmqtt(ctrl)
				mqtt.EXPECT().SendMessage(stateTopic, state, false)
				mqtt.EXPECT().SendMessage(stateTopicNew, stateNew, false)
				return mqtt
			},
			logger: func() logger {
				logger := mock_homeassistant.NewMocklogger(ctrl)
				logger.EXPECT().Info("Entity manager update", "clients", clientsWithNew)
				logger.EXPECT().Info("shutdown entitymanager")
				return logger
			},
			clients: map[string]dto.Client{
				mac: clientDto,
			},
			entityStates: map[string]map[string]string{},
		},
		{
			name: "add second client, error in discovery message and success update",
			entities: func() []Entity {
				entity := mock_homeassistant.NewMockEntity(ctrl)
				entity.EXPECT().GetStateTopic(clientDto).Return(stateTopic)
				entity.EXPECT().GetStateTopic(clientDtoNew).Return(stateTopicNew)
				entity.EXPECT().GetCommandTopic(clientDtoNew).Return(commandTopicNew)
				entity.EXPECT().GetState(clientDto).Return(state, nil)
				entity.EXPECT().GetState(clientDtoNew).Return(stateNew, nil)
				entity.EXPECT().SendDiscoveryMessage(clientDtoNew).Return(someErr)
				return []Entity{entity}
			},
			clientList: func() clientList {
				clientList := mock_homeassistant.NewMockclientList(ctrl)
				clientList.EXPECT().GetClientList().Return(clientsWithNew, nil)
				return clientList
			},
			mqtt: func() mqtt {
				ch := make(chan string)
				mqtt := mock_homeassistant.NewMockmqtt(ctrl)
				mqtt.EXPECT().SendMessage(stateTopic, state, false)
				mqtt.EXPECT().SendMessage(stateTopicNew, stateNew, false)
				mqtt.EXPECT().Subscribe(commandTopicNew).Return(ch)
				return mqtt
			},
			logger: func() logger {
				logger := mock_homeassistant.NewMocklogger(ctrl)
				logger.EXPECT().Info("Entity manager update", "clients", clientsWithNew)
				logger.EXPECT().Info("shutdown entitymanager")
				logger.EXPECT().Error("Entity manager update error while sending discovery message", "error", someErr, "client", clientDtoNew, "entity", gomock.Any())
				return logger
			},
			clients: map[string]dto.Client{
				mac: clientDto,
			},
			entityStates: map[string]map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			manager := NewEntityManager(
				tt.entities(),
				tt.clientList(),
				tt.mqtt(),
				100*time.Millisecond,
				tt.logger(),
			)

			manager.clients = tt.clients
			manager.entityStates = tt.entityStates

			var stopChan chan struct{}
			go func() {
				stopChan = manager.Run()
			}()
			time.Sleep(time.Millisecond * 120)

			go func() {
				stopChan <- struct{}{}
			}()
			time.Sleep(time.Millisecond * 100)
			return
		})
	}
}
