package policy

import (
	"bytes"
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
	Policy struct {
		host   string
		client client
	}

	ResponseStatus struct {
		Status  string `json:"status"`
		Code    string `json:"code"`
		Ident   string `json:"ident"`
		Message string `json:"message"`
	}

	ResponsePolicy struct {
		Status []ResponseStatus `json:"status"`
	}

	SetPolicyResponse struct {
		Policy ResponsePolicy `json:"policy"`
	}
)

func NewPolicy(host string, client client) *Policy {
	return &Policy{
		host:   host,
		client: client,
	}
}

func (p *Policy) SetPolicy(mac, policy string) error {

	var body interface{}
	if policy == "" {
		body = struct {
			Mac    string `json:"mac"`
			Policy bool   `json:"policy"`
		}{
			Mac:    mac,
			Policy: false,
		}
	} else {
		body = struct {
			Mac    string `json:"mac"`
			Policy string `json:"policy"`
		}{
			Mac:    mac,
			Policy: policy,
		}
	}
	b, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("RciIpHotspotHost error: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, p.host+"/rci/ip/hotspot/host", bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("build request error in setpolicy request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json;charset=UTF-8")

	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("send error in setpolicy request: %w", err)
	}
	defer resp.Body.Close()

	resBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body error in setpolicy request: %w", err)
	}

	var res SetPolicyResponse

	if err := json.Unmarshal(resBytes, &res); err != nil {
		return fmt.Errorf("unmarshal response error in setpolicy request: %w", err)
	}

	if len(res.Policy.Status) == 0 {
		return fmt.Errorf("no status in setpolicy response: %w", err)
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return auth.ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error in setpolicy request, status code: %d", resp.StatusCode)
	}

	return nil
}
