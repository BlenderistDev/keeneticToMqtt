package policy

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"keeneticToMqtt/internal/dto/keeneticdto"
	mock_policy "keeneticToMqtt/test/mocks/gomock/storages/policy"
)

func TestStorage_GetPolicyList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const (
		name = "name"
	)
	someErr := errors.New("some err")
	policies := []string{"none", name}

	tests := []struct {
		name            string
		policyClient    func() policyClient
		refreshInterval time.Duration
		policies        []string
		logger          func() logger
		expectedRes     []string
	}{
		{
			name: "success get policy list with refresh",
			policyClient: func() policyClient {
				policyClient := mock_policy.NewMockpolicyClient(ctrl)
				policyClient.EXPECT().GetPolicyList().Return(map[string]keeneticdto.Policy{
					name: {},
				}, nil)

				return policyClient
			},
			logger: func() logger {
				logger := mock_policy.NewMocklogger(ctrl)
				logger.EXPECT().Info(gomock.Eq("update policies"), gomock.Eq("policies"), gomock.Eq(policies))
				return logger
			},
			expectedRes: policies,
		},
		{
			name:     "success get policy list without refresh",
			policies: policies,
			policyClient: func() policyClient {
				policyClient := mock_policy.NewMockpolicyClient(ctrl)
				return policyClient
			},
			logger: func() logger {
				logger := mock_policy.NewMocklogger(ctrl)
				return logger
			},
			expectedRes: policies,
		},
		{
			name: "get policy list with error while refresh",
			policyClient: func() policyClient {
				policyClient := mock_policy.NewMockpolicyClient(ctrl)
				policyClient.EXPECT().GetPolicyList().Return(nil, someErr)
				return policyClient
			},
			logger: func() logger {
				logger := mock_policy.NewMocklogger(ctrl)
				logger.EXPECT().Error(
					gomock.Eq("error while refresh policies storage"),
					gomock.Eq("error"),
					gomock.Eq(someErr),
				)
				return logger
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			policy := NewStorage(tt.policyClient(), tt.refreshInterval, tt.logger())
			if len(tt.policies) > 0 {
				policy.policies = tt.policies
			}
			res := policy.GetPolicyList()
			assert.Equal(t, tt.expectedRes, res)
		})
	}
}

func TestStorage_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const (
		name = "name"
	)
	policies := []string{"none", name}

	policyClient := mock_policy.NewMockpolicyClient(ctrl)
	policyClient.EXPECT().GetPolicyList().Return(map[string]keeneticdto.Policy{
		name: {},
	}, nil)

	logger := mock_policy.NewMocklogger(ctrl)
	logger.EXPECT().Info(gomock.Eq("update policies"), gomock.Eq("policies"), gomock.Eq(policies))
	logger.EXPECT().Info("shutdown policy storage")
	storage := NewStorage(policyClient, 10*time.Millisecond, logger)
	done := storage.Run()

	ticker := time.NewTicker(15 * time.Millisecond)
	<-ticker.C
	done <- struct{}{}
	<-ticker.C
	return
}
