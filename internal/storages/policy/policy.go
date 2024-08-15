package policy

import (
	"log/slog"
	"time"

	"keeneticToMqtt/internal/clients/keenetic/policylist"
	"keeneticToMqtt/internal/dto/homeassistantdto"
)

type Storage struct {
	policyClient    *policylist.PolicyList
	refreshInterval time.Duration
	policies        []string
	logger          *slog.Logger
}

func NewStorage(policyClient *policylist.PolicyList, refreshInterval time.Duration, logger *slog.Logger) *Storage {
	s := &Storage{
		policyClient:    policyClient,
		refreshInterval: refreshInterval,
		logger:          logger,
	}

	return s
}

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
}
