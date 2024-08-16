package policy

import (
	"time"

	"keeneticToMqtt/internal/dto/homeassistantdto"
	"keeneticToMqtt/internal/dto/keeneticdto"
)

//go:generate mockgen -source=policy.go -destination=../../../test/mocks/gomock/storages/policy/policy.go

type (
	policyClient interface {
		GetPolicyList() (map[string]keeneticdto.Policy, error)
	}
	logger interface {
		Error(msg string, args ...any)
		Info(msg string, args ...any)
	}

	// Storage store policies in-memory.
	Storage struct {
		policyClient    policyClient
		refreshInterval time.Duration
		policies        []string
		logger          logger
	}
)

// NewStorage creates Storage
func NewStorage(policyClient policyClient, refreshInterval time.Duration, logger logger) *Storage {
	s := &Storage{
		policyClient:    policyClient,
		refreshInterval: refreshInterval,
		logger:          logger,
	}

	return s
}

// Run start storage updates.
func (s *Storage) Run() chan struct{} {
	done := make(chan struct{})
	ticker := time.NewTicker(s.refreshInterval)

	go func() {
		for {
			select {
			case <-done:
				s.logger.Info("shutdown policy storage")
				return
			case _ = <-ticker.C:
				s.refresh()
			}
		}
	}()

	return done
}

// GetPolicyList returns policy list.
func (s *Storage) GetPolicyList() []string {
	if len(s.policies) == 0 {
		s.refresh()
	}
	return s.policies
}

func (s *Storage) refresh() {
	resp, err := s.policyClient.GetPolicyList()
	if err != nil {
		s.logger.Error("error while refresh policies storage", "error", err)
		return
	}
	policies := make([]string, 0, len(resp)+1)

	policies = append(policies, homeassistantdto.NonePolicy)
	for k := range resp {
		policies = append(policies, k)
	}

	s.policies = policies
	s.logger.Info("update policies", "policies", policies)
}
