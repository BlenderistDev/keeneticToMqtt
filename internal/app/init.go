package app

import (
	"log/slog"
	"net/http/cookiejar"
	"time"

	"keeneticToMqtt/internal/clients/keenetic"
	"keeneticToMqtt/internal/clients/keenetic/accessupdate"
	"keeneticToMqtt/internal/clients/keenetic/auth"
	"keeneticToMqtt/internal/clients/keenetic/list"
	"keeneticToMqtt/internal/clients/keenetic/policylist"
	"keeneticToMqtt/internal/clients/mqtt"
	"keeneticToMqtt/internal/config"
	"keeneticToMqtt/internal/homeassistant"
	"keeneticToMqtt/internal/homeassistant/clientpermit"
	"keeneticToMqtt/internal/homeassistant/clientpolicy"
	"keeneticToMqtt/internal/homeassistant/rxbytes"
	"keeneticToMqtt/internal/homeassistant/txbytes"
	"keeneticToMqtt/internal/logger"
	"keeneticToMqtt/internal/services/clientlist"
	"keeneticToMqtt/internal/services/discovery"
	"keeneticToMqtt/internal/storages/policy"
)

// Container with dependencies.
type Container struct {
	Logger            *slog.Logger
	Config            *config.Config
	ClientListService *clientlist.ClientList
	DiscoveryService  *discovery.Discovery
	EntityManager     *homeassistant.EntityManager
	PolicyStorage     *policy.Storage
	Mqtt              *mqtt.Client
}

// NewContainer creates new Container.
func NewContainer() (*Container, error) {
	cont := Container{}

	conf, err := config.NewDefaultConfig()
	if err != nil {
		return nil, err
	}
	cont.Config = conf

	cont.Logger = logger.NewLogger(cont.Config.LogLevel)

	cookie, _ := cookiejar.New(&cookiejar.Options{})

	cont.Mqtt = mqtt.NewClient(cont.Config.Mqtt.Host, cont.Config.Mqtt.ClientID, cont.Config.Mqtt.Login, cont.Config.Mqtt.Password, cont.Logger)

	authClient := auth.NewAuth(cont.Config.Keenetic.Host, cont.Config.Keenetic.Login, cont.Config.Keenetic.Password, cookie)
	keeneticClient := keenetic.NewKeenetic(authClient, cookie, cont.Config.Keenetic.Host, cont.Config.Keenetic.Login, cont.Config.Keenetic.Password, cont.Logger)
	policyClient := accessupdate.NewAccessUpdate(cont.Config.Keenetic.Host, keeneticClient)
	policyList := policylist.NewPolicyList(cont.Config.Keenetic.Host, keeneticClient)
	listClient := list.NewList(cont.Config.Keenetic.Host, keeneticClient)

	cont.PolicyStorage = policy.NewStorage(policyList, time.Second*10, cont.Logger)

	cont.ClientListService = clientlist.NewClientList(listClient, cont.Config.Homeassistant.WhiteList)
	cont.DiscoveryService = discovery.NewDiscovery("", cont.Config.Homeassistant.DeviceID, cont.Mqtt)

	clientPolicy := clientpolicy.NewClientPolicy(cont.Config.Mqtt.BaseTopic, cont.DiscoveryService, policyClient, cont.PolicyStorage)
	clientPermit := clientpermit.NewClientPermit(cont.Config.Mqtt.BaseTopic, cont.DiscoveryService, policyClient)
	txBytes := txbytes.NewTxBytes(cont.Config.Mqtt.BaseTopic, cont.DiscoveryService)
	rxBytes := rxbytes.NewRxBytes(cont.Config.Mqtt.BaseTopic, cont.DiscoveryService)

	cont.EntityManager = homeassistant.NewEntityManager(
		[]homeassistant.Entity{
			clientPolicy,
			clientPermit,
			txBytes,
			rxBytes,
		},
		cont.ClientListService,
		cont.Mqtt,
		cont.Config.Homeassistant.UpdateInterval,
		cont.Logger,
	)

	return &cont, nil
}
