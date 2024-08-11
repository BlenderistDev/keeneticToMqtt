package clientlist

import (
	"fmt"

	"keeneticToMqtt/internal/dto"
)

type listClient interface {
	GetDeviceList() ([]*dto.KeeneticDeviceInfo, error)
	GetPolicyList() ([]*dto.KeeneticDevicePolicy, error)
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
	policyList, err := l.listClient.GetPolicyList()
	if err != nil {
		return nil, fmt.Errorf("ClientList client error while getting policy list: %w", err)
	}

	policyMap := make(map[string]*dto.KeeneticDevicePolicy, len(policyList))
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
