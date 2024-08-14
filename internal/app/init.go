package app

import (
	"log/slog"
	"net/http/cookiejar"
	"time"

	"keeneticToMqtt/internal/clients/keenetic"
	"keeneticToMqtt/internal/clients/keenetic/auth"
	"keeneticToMqtt/internal/clients/keenetic/list"
	"keeneticToMqtt/internal/clients/keenetic/policy"
	"keeneticToMqtt/internal/clients/mqtt"
	configs "keeneticToMqtt/internal/config"
	"keeneticToMqtt/internal/homeassistant"
	"keeneticToMqtt/internal/homeassistant/clientpolicy"
	"keeneticToMqtt/internal/services/clientlist"
	"keeneticToMqtt/internal/services/discovery"
	policy2 "keeneticToMqtt/internal/storages/policy"
)

type Container struct {
	Logger            *slog.Logger
	Config            *configs.Config
	ClientListService *clientlist.ClientList
	DiscoveryService  *discovery.Discovery
	EntityManager     *homeassistant.EntityManager
}

func NewContainer() (*Container, error) {
	cont := Container{}
	cont.Logger = slog.Default()

	conf, err := configs.NewDefaultConfig()
	if err != nil {
		return nil, err
	}

	cont.Config = conf

	cookie, _ := cookiejar.New(&cookiejar.Options{})

	authClient := auth.NewAuth(cont.Config.Keenetic.Host, cont.Config.Keenetic.Login, cont.Config.Keenetic.Password, cookie)
	keeneticClient := keenetic.NewKeenetic(authClient, cookie, cont.Config.Keenetic.Host, cont.Config.Keenetic.Login, cont.Config.Keenetic.Password)
	policyClient := policy.NewPolicy(cont.Config.Keenetic.Host, keeneticClient)
	listClient := list.NewList(cont.Config.Keenetic.Host, keeneticClient)
	mqttClient := mqtt.NewClient(cont.Config.Mqtt.Host, cont.Config.Mqtt.ClientID, cont.Config.Mqtt.Login, cont.Config.Mqtt.Password)

	policyStorage := policy2.NewStorage(policyClient, time.Second*10, cont.Logger)

	cont.ClientListService = clientlist.NewClientList(listClient, cont.Config.Homeassistant.WhiteList)
	cont.DiscoveryService = discovery.NewDiscovery("", cont.Config.Homeassistant.DeviceID, mqttClient)

	clientPolicy := clientpolicy.NewClientPolicy(cont.Config.Mqtt.BaseTopic, cont.DiscoveryService, mqttClient, policyClient, policyStorage, cont.Logger)

	cont.EntityManager = homeassistant.NewEntityManager(
		[]homeassistant.Entity{
			clientPolicy,
		},
		cont.ClientListService,
		cont.Config.Homeassistant.UpdateInterval,
		cont.Logger,
	)

	return &cont, nil
}
