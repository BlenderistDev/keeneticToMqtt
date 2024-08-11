package policy

import (
	"log/slog"
	"time"

	"keeneticToMqtt/internal/clients/keenetic/policy"
)

type Storage struct {
	policyClient    *policy.Policy
	refreshInterval time.Duration
	policies        []string
	logger          *slog.Logger
}

func NewStorage(policyClient *policy.Policy, refreshInterval time.Duration) *Storage {
	s := &Storage{
		policyClient:    policyClient,
		refreshInterval: refreshInterval,
	}

	s.refresh()
	ticker := time.NewTicker(refreshInterval)

	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case _ = <-ticker.C:
				s.refresh()
			}
		}
	}()

	return s
}

func (s *Storage) GetPolicyList() []string {
	return s.policies
}

func (s *Storage) refresh() {
	resp, err := s.policyClient.GetPolicyList()
	if err != nil {
		s.logger.Error("error while refresh policies storage", "error", err)
	}
	policies := make([]string, 0, len(resp))
	for k := range resp {
		policies = append(policies, k)
	}

	s.policies = policies
}
