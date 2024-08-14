// Code generated by MockGen. DO NOT EDIT.
// Source: delivery_end.go

// Package mock_usecase is a generated GoMock package.
package mock_usecase

import (
	context "context"
	reflect "reflect"
	time "time"
	models "touchgift-job-manager/domain/models"
	repository "touchgift-job-manager/domain/repository"

	gomock "github.com/golang/mock/gomock"
)

// MockDeliveryEnd is a mock of DeliveryEnd interface.
type MockDeliveryEnd struct {
	ctrl     *gomock.Controller
	recorder *MockDeliveryEndMockRecorder
}

// MockDeliveryEndMockRecorder is the mock recorder for MockDeliveryEnd.
type MockDeliveryEndMockRecorder struct {
	mock *MockDeliveryEnd
}

// NewMockDeliveryEnd creates a new mock instance.
func NewMockDeliveryEnd(ctrl *gomock.Controller) *MockDeliveryEnd {
	mock := &MockDeliveryEnd{ctrl: ctrl}
	mock.recorder = &MockDeliveryEndMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDeliveryEnd) EXPECT() *MockDeliveryEndMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockDeliveryEnd) Close() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close.
func (mr *MockDeliveryEndMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockDeliveryEnd)(nil).Close))
}

// CreateWorker mocks base method.
func (m *MockDeliveryEnd) CreateWorker(ctx context.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "CreateWorker", ctx)
}

// CreateWorker indicates an expected call of CreateWorker.
func (mr *MockDeliveryEndMockRecorder) CreateWorker(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateWorker", reflect.TypeOf((*MockDeliveryEnd)(nil).CreateWorker), ctx)
}

// Delete mocks base method.
func (m *MockDeliveryEnd) Delete(ctx context.Context, campaignID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, campaignID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockDeliveryEndMockRecorder) Delete(ctx, campaignID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockDeliveryEnd)(nil).Delete), ctx, campaignID)
}

// ExecuteNow mocks base method.
func (m *MockDeliveryEnd) ExecuteNow(campaign *models.Campaign) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ExecuteNow", campaign)
}

// ExecuteNow indicates an expected call of ExecuteNow.
func (mr *MockDeliveryEndMockRecorder) ExecuteNow(campaign interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExecuteNow", reflect.TypeOf((*MockDeliveryEnd)(nil).ExecuteNow), campaign)
}

// GetDeliveryDataCampaigns mocks base method.
func (m *MockDeliveryEnd) GetDeliveryDataCampaigns(ctx context.Context, to time.Time, status []string, limit int) ([]*models.Campaign, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeliveryDataCampaigns", ctx, to, status, limit)
	ret0, _ := ret[0].([]*models.Campaign)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDeliveryDataCampaigns indicates an expected call of GetDeliveryDataCampaigns.
func (mr *MockDeliveryEndMockRecorder) GetDeliveryDataCampaigns(ctx, to, status, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeliveryDataCampaigns", reflect.TypeOf((*MockDeliveryEnd)(nil).GetDeliveryDataCampaigns), ctx, to, status, limit)
}

// Reserve mocks base method.
func (m *MockDeliveryEnd) Reserve(ctx context.Context, endAt time.Time, campaign *models.Campaign) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Reserve", ctx, endAt, campaign)
}

// Reserve indicates an expected call of Reserve.
func (mr *MockDeliveryEndMockRecorder) Reserve(ctx, endAt, campaign interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Reserve", reflect.TypeOf((*MockDeliveryEnd)(nil).Reserve), ctx, endAt, campaign)
}

// Stop mocks base method.
func (m *MockDeliveryEnd) Stop(ctx context.Context, tx repository.Transaction, campaign *models.Campaign, status string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stop", ctx, tx, campaign, status)
	ret0, _ := ret[0].(error)
	return ret0
}

// Stop indicates an expected call of Stop.
func (mr *MockDeliveryEndMockRecorder) Stop(ctx, tx, campaign, status interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockDeliveryEnd)(nil).Stop), ctx, tx, campaign, status)
}

// Terminate mocks base method.
func (m *MockDeliveryEnd) Terminate(ctx context.Context, tx repository.Transaction, campaignID int, updatedAt time.Time) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Terminate", ctx, tx, campaignID, updatedAt)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Terminate indicates an expected call of Terminate.
func (mr *MockDeliveryEndMockRecorder) Terminate(ctx, tx, campaignID, updatedAt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Terminate", reflect.TypeOf((*MockDeliveryEnd)(nil).Terminate), ctx, tx, campaignID, updatedAt)
}
