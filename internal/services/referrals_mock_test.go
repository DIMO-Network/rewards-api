// Code generated by MockGen. DO NOT EDIT.
// Source: ./referrals.go
//
// Generated by this command:
//
//	mockgen -source=./referrals.go -destination=referrals_mock_test.go -package=services
//

// Package services is a generated GoMock package.
package services

import (
	context "context"
	reflect "reflect"

	common "github.com/ethereum/go-ethereum/common"
	gomock "go.uber.org/mock/gomock"
)

// MockMobileAPIClient is a mock of MobileAPIClient interface.
type MockMobileAPIClient struct {
	ctrl     *gomock.Controller
	recorder *MockMobileAPIClientMockRecorder
	isgomock struct{}
}

// MockMobileAPIClientMockRecorder is the mock recorder for MockMobileAPIClient.
type MockMobileAPIClientMockRecorder struct {
	mock *MockMobileAPIClient
}

// NewMockMobileAPIClient creates a new mock instance.
func NewMockMobileAPIClient(ctrl *gomock.Controller) *MockMobileAPIClient {
	mock := &MockMobileAPIClient{ctrl: ctrl}
	mock.recorder = &MockMobileAPIClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMobileAPIClient) EXPECT() *MockMobileAPIClientMockRecorder {
	return m.recorder
}

// GetReferrer mocks base method.
func (m *MockMobileAPIClient) GetReferrer(ctx context.Context, addr common.Address) (common.Address, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetReferrer", ctx, addr)
	ret0, _ := ret[0].(common.Address)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetReferrer indicates an expected call of GetReferrer.
func (mr *MockMobileAPIClientMockRecorder) GetReferrer(ctx, addr any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetReferrer", reflect.TypeOf((*MockMobileAPIClient)(nil).GetReferrer), ctx, addr)
}
