// Code generated by MockGen. DO NOT EDIT.
// Source: timer.go

// Package mock_usecase is a generated GoMock package.
package mock_usecase

import (
	context "context"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
)

// MockTimer is a mock of Timer interface.
type MockTimer struct {
	ctrl     *gomock.Controller
	recorder *MockTimerMockRecorder
}

// MockTimerMockRecorder is the mock recorder for MockTimer.
type MockTimerMockRecorder struct {
	mock *MockTimer
}

// NewMockTimer creates a new mock instance.
func NewMockTimer(ctrl *gomock.Controller) *MockTimer {
	mock := &MockTimer{ctrl: ctrl}
	mock.recorder = &MockTimerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTimer) EXPECT() *MockTimerMockRecorder {
	return m.recorder
}

// ExecuteAtTime mocks base method.
func (m *MockTimer) ExecuteAtTime(ctx context.Context, specifiedTime time.Time, process func()) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ExecuteAtTime", ctx, specifiedTime, process)
}

// ExecuteAtTime indicates an expected call of ExecuteAtTime.
func (mr *MockTimerMockRecorder) ExecuteAtTime(ctx, specifiedTime, process interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExecuteAtTime", reflect.TypeOf((*MockTimer)(nil).ExecuteAtTime), ctx, specifiedTime, process)
}
