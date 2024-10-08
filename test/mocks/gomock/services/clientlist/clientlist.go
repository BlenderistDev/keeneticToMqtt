// Code generated by MockGen. DO NOT EDIT.
// Source: clientlist.go
//
// Generated by this command:
//
//	mockgen -source=clientlist.go -destination=../../../test/mocks/gomock/services/clientlist/clientlist.go
//
// Package mock_clientlist is a generated GoMock package.
package mock_clientlist

import (
	keeneticdto "keeneticToMqtt/internal/dto/keeneticdto"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MocklistClient is a mock of listClient interface.
type MocklistClient struct {
	ctrl     *gomock.Controller
	recorder *MocklistClientMockRecorder
}

// MocklistClientMockRecorder is the mock recorder for MocklistClient.
type MocklistClientMockRecorder struct {
	mock *MocklistClient
}

// NewMocklistClient creates a new mock instance.
func NewMocklistClient(ctrl *gomock.Controller) *MocklistClient {
	mock := &MocklistClient{ctrl: ctrl}
	mock.recorder = &MocklistClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MocklistClient) EXPECT() *MocklistClientMockRecorder {
	return m.recorder
}

// GetClientPolicyList mocks base method.
func (m *MocklistClient) GetClientPolicyList() ([]keeneticdto.DevicePolicy, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetClientPolicyList")
	ret0, _ := ret[0].([]keeneticdto.DevicePolicy)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetClientPolicyList indicates an expected call of GetClientPolicyList.
func (mr *MocklistClientMockRecorder) GetClientPolicyList() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetClientPolicyList", reflect.TypeOf((*MocklistClient)(nil).GetClientPolicyList))
}

// GetDeviceList mocks base method.
func (m *MocklistClient) GetDeviceList() ([]keeneticdto.DeviceInfoResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeviceList")
	ret0, _ := ret[0].([]keeneticdto.DeviceInfoResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDeviceList indicates an expected call of GetDeviceList.
func (mr *MocklistClientMockRecorder) GetDeviceList() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeviceList", reflect.TypeOf((*MocklistClient)(nil).GetDeviceList))
}
