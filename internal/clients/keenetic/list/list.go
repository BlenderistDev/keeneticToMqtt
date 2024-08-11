package list

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"keeneticToMqtt/internal/clients/keenetic/auth"
	"keeneticToMqtt/internal/dto/keeneticdto"
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

func (l *List) GetDeviceList() ([]*keeneticdto.DeviceInfoResponse, error) {
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
		return nil, fmt.Errorf("error in GetDeviceList request, status code: %d", resp.StatusCode)
	}

	resBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body error in GetDeviceList request: %w", err)
	}
	var res []*keeneticdto.DeviceInfoResponse

	if err := json.Unmarshal(resBytes, &res); err != nil {
		return nil, fmt.Errorf("unmarshal response error in setpolicy request: %w", err)
	}

	return res, nil
}

func (l *List) GetClientPolicyList() ([]*keeneticdto.DevicePolicy, error) {
	req, err := http.NewRequest(http.MethodGet, l.host+"/rci/show/rc/ip/hotspot/host", nil)
	if err != nil {
		return nil, fmt.Errorf("build request error in GetClientPolicyList request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json;charset=UTF-8")

	resp, err := l.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send error in GetClientPolicyList request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, auth.ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error in GetClientPolicyList request, status code: %d", resp.StatusCode)
	}

	resBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body error in GetClientPolicyList request: %w", err)
	}
	var res []*keeneticdto.DevicePolicy

	if err := json.Unmarshal(resBytes, &res); err != nil {
		return nil, fmt.Errorf("unmarshal response error in setpolicy request: %w", err)
	}

	return res, nil
}

func (l *List) GetPolicyList() (keeneticdto.PolicyResponse, error) {
	req, err := http.NewRequest(http.MethodGet, l.host+"/rci/show/rc/ip/policy", nil)
	if err != nil {
		return nil, fmt.Errorf("build request error in GetPolicyList request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json;charset=UTF-8")

	resp, err := l.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send error in GetPolicyList request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, auth.ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error in GetPolicyList request, status code: %d", resp.StatusCode)
	}

	resBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body error in GetPolicyList request: %w", err)
	}
	var res keeneticdto.PolicyResponse

	if err := json.Unmarshal(resBytes, &res); err != nil {
		return nil, fmt.Errorf("unmarshal response error in setpolicy request: %w", err)
	}

	return res, nil
}
