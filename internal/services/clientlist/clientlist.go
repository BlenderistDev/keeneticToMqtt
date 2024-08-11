package clientlist

import (
	"fmt"

	"keeneticToMqtt/internal/dto"
	"keeneticToMqtt/internal/dto/keeneticdto"
)

type listClient interface {
	GetDeviceList() ([]*keeneticdto.DeviceInfoResponse, error)
	GetClientPolicyList() ([]*keeneticdto.DevicePolicy, error)
}

type ClientList struct {
	listClient listClient
}

func NewClientList(listClient listClient) *ClientList {
	return &ClientList{
		listClient: listClient,
	}
}

func (l *ClientList) GetClientList() ([]dto.Client, error) {
	deviceList, err := l.listClient.GetDeviceList()
	if err != nil {
		return nil, fmt.Errorf("ClientList client error while getting device list: %w", err)
	}
	policyList, err := l.listClient.GetClientPolicyList()
	if err != nil {
		return nil, fmt.Errorf("ClientList client error while getting policy list: %w", err)
	}

	policyMap := make(map[string]*keeneticdto.DevicePolicy, len(policyList))
	for _, policy := range policyList {
		policyMap[policy.Mac] = policy
	}

	clientList := make([]dto.Client, len(deviceList))
	for i, device := range deviceList {
		clientList[i] = dto.Client{
			Mac:  device.Mac,
			Name: device.Name,
		}
		policy := policyMap[device.Mac]
		if policy != nil && policy.Policy != nil {
			clientList[i].Policy = *policy.Policy
		}
	}

	return clientList, nil
}
