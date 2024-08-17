package policylist

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"keeneticToMqtt/internal/dto/keeneticdto"
	"keeneticToMqtt/internal/errs"
)

//go:generate mockgen -source=policy.go -destination=../../../../test/mocks/gomock/clients/keenetic/policylist/policy.go

const policyListUrl = "/rci/show/rc/ip/policy"

type (
	client interface {
		Do(req *http.Request) (*http.Response, error)
	}
)

// PolicyList struct to get keenetic policy list.
type PolicyList struct {
	host   string
	client client
}

// NewPolicyList creates new PolicyList.
func NewPolicyList(host string, client client) *PolicyList {
	return &PolicyList{
		host:   host,
		client: client,
	}
}

// GetPolicyList return map of policies. Key of map is name of policy.
func (l *PolicyList) GetPolicyList() (map[string]keeneticdto.Policy, error) {
	req, err := http.NewRequest(http.MethodGet, l.host+policyListUrl, nil)
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
		return nil, errs.ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error in GetPolicyList request, status code: %d", resp.StatusCode)
	}

	resBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body error in GetPolicyList request: %w", err)
	}
	var res map[string]keeneticdto.Policy

	if err := json.Unmarshal(resBytes, &res); err != nil {
		return nil, fmt.Errorf("unmarshal response error in setpolicy request: %w", err)
	}

	return res, nil
}
