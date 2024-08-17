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

const (
	ipHotspotHostURL = "/rci/ip/hotspot/host"
)

type (
	client interface {
		Do(req *http.Request) (*http.Response, error)
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

	permitTrueReq struct {
		Mac    string `json:"mac"`
		Permit bool   `json:"permit"`
	}
	permitFalseReq struct {
		Mac  string `json:"mac"`
		Deny bool   `json:"deny"`
	}
	setEmptyPolicyReq struct {
		Mac    string `json:"mac"`
		Policy bool   `json:"policy"`
	}
	setPolicyReq struct {
		Mac    string `json:"mac"`
		Policy string `json:"policy"`
	}
)

// NewAccessUpdate creates new AccessUpdate.
func NewAccessUpdate(host string, client client) *AccessUpdate {
	return &AccessUpdate{
		host:   host,
		client: client,
	}
}

// AccessUpdate struct for controlling keenetic client access.
type AccessUpdate struct {
	host   string
	client client
}

// SetPolicy set keenetic client policy.
func (p *AccessUpdate) SetPolicy(mac, policy string) error {
	var body interface{}
	if policy == homeassistantdto.NonePolicy {
		body = setEmptyPolicyReq{
			Mac:    mac,
			Policy: false,
		}
	} else {
		body = setPolicyReq{
			Mac:    mac,
			Policy: policy,
		}
	}

	return p.ipHotspotHostRequest(body)
}

// SetPermit set keenetic client internet permit.
func (p *AccessUpdate) SetPermit(mac string, permit bool) error {
	var body interface{}
	if permit {
		body = permitTrueReq{
			Mac:    mac,
			Permit: permit,
		}
	} else {
		body = permitFalseReq{
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

	req, err := http.NewRequest(http.MethodPost, p.host+ipHotspotHostURL, bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("build request error in setaccess request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json;charset=UTF-8")

	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("send error in setaccess request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return errs.ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error in setaccess request, status code: %d", resp.StatusCode)
	}

	resBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body error in access request: %w", err)
	}

	var res response
	if err := json.Unmarshal(resBytes, &res); err != nil {
		return fmt.Errorf("unmarshal response error in setaccess request: %w", err)
	}

	if len(res) == 0 {
		return fmt.Errorf("no status in setaccess response: %w", err)
	}

	return nil
}
