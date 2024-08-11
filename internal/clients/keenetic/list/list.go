package list

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"keeneticToMqtt/internal/clients/keenetic/auth"
)

type client interface {
	Do(req *http.Request) (*http.Response, error)
}

type (
	List struct {
		host   string
		client client
	}
)

func NewList(host string, client client) *List {
	return &List{
		host:   host,
		client: client,
	}
}

func (l *List) GetDeviceList() ([]*DeviceInfo, error) {
	req, err := http.NewRequest(http.MethodGet, l.host+"/rci/show/ip/hotspot/host", nil)
	if err != nil {
		return nil, fmt.Errorf("build request error in GetDeviceList request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json;charset=UTF-8")

	resp, err := l.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send error in GetDeviceList request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, auth.ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error in setpolicy request, status code: %d", resp.StatusCode)
	}

	resBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body error in GetDeviceList request: %w", err)
	}
	var res []*DeviceInfo

	if err := json.Unmarshal(resBytes, &res); err != nil {
		return nil, fmt.Errorf("unmarshal response error in setpolicy request: %w", err)
	}

	return res, nil
}

func (l *List) GetPolicyList() ([]*DevicePolicy, error) {
	req, err := http.NewRequest(http.MethodGet, l.host+"/rci/show/rc/ip/hotspot/host", nil)
	if err != nil {
		return nil, fmt.Errorf("build request error in GetDeviceList request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json;charset=UTF-8")

	resp, err := l.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send error in GetDeviceList request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, auth.ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error in setpolicy request, status code: %d", resp.StatusCode)
	}

	resBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body error in GetDeviceList request: %w", err)
	}
	var res []*DevicePolicy

	if err := json.Unmarshal(resBytes, &res); err != nil {
		return nil, fmt.Errorf("unmarshal response error in setpolicy request: %w", err)
	}

	return res, nil
}
