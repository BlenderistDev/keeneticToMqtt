package accessupdate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"keeneticToMqtt/internal/dto/homeassistantdto"
	"keeneticToMqtt/internal/errs"
)

type client interface {
	Do(req *http.Request) (*http.Response, error)
}

type (
	AccessUpdate struct {
		host   string
		client client
	}

	responseStatus struct {
		Status  string `json:"status"`
		Code    string `json:"code"`
		Ident   string `json:"ident"`
		Message string `json:"message"`
	}

	responseClient struct {
		Status []responseStatus `json:"status"`
	}

	response map[string]responseClient
)

func NewAccessUpdate(host string, client client) *AccessUpdate {
	return &AccessUpdate{
		host:   host,
		client: client,
	}
}

func (p *AccessUpdate) SetPolicy(mac, policy string) error {
	var body interface{}
	if policy == homeassistantdto.NonePolicy {
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

	return p.ipHotspotHostRequest(body)
}

func (p *AccessUpdate) SetPermit(mac string, permit bool) error {
	var body interface{}
	if permit {
		body = struct {
			Mac    string `json:"mac"`
			Permit bool   `json:"permit"`
		}{
			Mac:    mac,
			Permit: permit,
		}
	} else {
		body = struct {
			Mac  string `json:"mac"`
			Deny bool   `json:"deny"`
		}{
			Mac:  mac,
			Deny: !permit,
		}
	}

	return p.ipHotspotHostRequest(body)
}

func (p *AccessUpdate) ipHotspotHostRequest(body interface{}) error {
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

	var res response

	if err := json.Unmarshal(resBytes, &res); err != nil {
		return fmt.Errorf("unmarshal response error in setpolicy request: %w", err)
	}

	if len(res) == 0 {
		return fmt.Errorf("no status in setpolicy response: %w", err)
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return errs.ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error in setpolicy request, status code: %d", resp.StatusCode)
	}

	return nil
}
