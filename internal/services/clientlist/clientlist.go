package clientlist

import (
	"fmt"

	"keeneticToMqtt/internal/dto"
	"keeneticToMqtt/internal/dto/homeassistantdto"
	"keeneticToMqtt/internal/dto/keeneticdto"
)

//go:generate mockgen -source=clientlist.go -destination=../../../test/mocks/gomock/services/clientlist/clientlist.go

type listClient interface {
	GetDeviceList() ([]keeneticdto.DeviceInfoResponse, error)
	GetClientPolicyList() ([]keeneticdto.DevicePolicy, error)
}

// ClientList struct for building keenetic client list.
type ClientList struct {
	listClient   listClient
	macWhiteList map[string]bool
}

// NewClientList creates new ClientList.
func NewClientList(listClient listClient, macWhiteList []string) *ClientList {
	macMap := make(map[string]bool, len(macWhiteList))
	for _, mac := range macWhiteList {
		macMap[mac] = true
	}
	return &ClientList{
		listClient:   listClient,
		macWhiteList: macMap,
	}
}

// GetClientList returns list of dto.Client.
func (l *ClientList) GetClientList() ([]dto.Client, error) {
	deviceList, err := l.listClient.GetDeviceList()
	if err != nil {
		return nil, fmt.Errorf("ClientList client error while getting device list: %w", err)
	}
	policyList, err := l.listClient.GetClientPolicyList()
	if err != nil {
		return nil, fmt.Errorf("ClientList client error while getting policy list: %w", err)
	}

	policyMap := make(map[string]keeneticdto.DevicePolicy, len(policyList))
	for _, policy := range policyList {
		policyMap[policy.Mac] = policy
	}

	clientList := make([]dto.Client, 0)
	for _, device := range deviceList {
		if !l.macWhiteList[device.Mac] {
			continue
		}
		client := dto.Client{
			Mac:  device.Mac,
			Name: device.Name,
		}

		policy := policyMap[device.Mac]
		if policy.Policy != nil && *policy.Policy != "" {
			client.Policy = *policy.Policy
		} else {
			client.Policy = homeassistantdto.NonePolicy
		}

		client.Permit = policy.Permit

		clientList = append(clientList, client)
	}

	return clientList, nil
}
