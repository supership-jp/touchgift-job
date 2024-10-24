// Code generated by MockGen. DO NOT EDIT.
// Source: delivery_operation.go

// Package mock_usecase is a generated GoMock package.
package mock_usecase

import (
	context "context"
	reflect "reflect"
	time "time"
	models "touchgift-job-manager/domain/models"

	gomock "github.com/golang/mock/gomock"
)

// MockDeliveryOperation is a mock of DeliveryOperation interface.
type MockDeliveryOperation struct {
	ctrl     *gomock.Controller
	recorder *MockDeliveryOperationMockRecorder
}

// MockDeliveryOperationMockRecorder is the mock recorder for MockDeliveryOperation.
type MockDeliveryOperationMockRecorder struct {
	mock *MockDeliveryOperation
}

// NewMockDeliveryOperation creates a new mock instance.
func NewMockDeliveryOperation(ctrl *gomock.Controller) *MockDeliveryOperation {
	mock := &MockDeliveryOperation{ctrl: ctrl}
	mock.recorder = &MockDeliveryOperationMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDeliveryOperation) EXPECT() *MockDeliveryOperationMockRecorder {
	return m.recorder
}

// Process mocks base method.
func (m *MockDeliveryOperation) Process(ctx context.Context, current time.Time, campaignLog *models.CampaignLog) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Process", ctx, current, campaignLog)
	ret0, _ := ret[0].(error)
	return ret0
}

// Process indicates an expected call of Process.
func (mr *MockDeliveryOperationMockRecorder) Process(ctx, current, campaignLog interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Process", reflect.TypeOf((*MockDeliveryOperation)(nil).Process), ctx, current, campaignLog)
}
