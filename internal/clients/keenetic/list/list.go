package list

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"keeneticToMqtt/internal/dto/keeneticdto"
	"keeneticToMqtt/internal/errs"
)

const (
	clientPolicyListUrl = "/rci/show/rc/ip/hotspot/host"
	deviceListUrl       = "/rci/show/ip/hotspot/host"
)

type (
	client interface {
		Do(req *http.Request) (*http.Response, error)
	}
)

// List struct for get client lists from keenetic.
type List struct {
	host   string
	client client
}

// NewList creates new List.
func NewList(host string, client client) *List {
	return &List{
		host:   host,
		client: client,
	}
}

// GetDeviceList returns keenetic device list.
func (l *List) GetDeviceList() ([]keeneticdto.DeviceInfoResponse, error) {
	req, err := http.NewRequest(http.MethodGet, l.host+deviceListUrl, nil)
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
		return nil, errs.ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error in GetDeviceList request, status code: %d", resp.StatusCode)
	}

	resBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body error in GetDeviceList request: %w", err)
	}
	var res []keeneticdto.DeviceInfoResponse

	if err := json.Unmarshal(resBytes, &res); err != nil {
		return nil, fmt.Errorf("unmarshal response error in GetDeviceList request: %w", err)
	}

	return res, nil
}

// GetClientPolicyList returns keenetic policy list.
func (l *List) GetClientPolicyList() ([]keeneticdto.DevicePolicy, error) {
	req, err := http.NewRequest(http.MethodGet, l.host+clientPolicyListUrl, nil)
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
		return nil, errs.ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error in GetClientPolicyList request, status code: %d", resp.StatusCode)
	}

	resBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body error in GetClientPolicyList request: %w", err)
	}
	var res []keeneticdto.DevicePolicy

	if err := json.Unmarshal(resBytes, &res); err != nil {
		return nil, fmt.Errorf("unmarshal response error in GetClientPolicyList request: %w", err)
	}

	return res, nil
}
