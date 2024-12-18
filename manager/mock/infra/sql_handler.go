// Code generated by MockGen. DO NOT EDIT.
// Source: sql_handler.go

// Package mock_infra is a generated GoMock package.
package mock_infra

import (
	context "context"
	reflect "reflect"
	repository "touchgift-job-manager/domain/repository"

	gomock "github.com/golang/mock/gomock"
	sqlx "github.com/jmoiron/sqlx"
)

// MockSQLHandler is a mock of SQLHandler interface.
type MockSQLHandler struct {
	ctrl     *gomock.Controller
	recorder *MockSQLHandlerMockRecorder
}

// MockSQLHandlerMockRecorder is the mock recorder for MockSQLHandler.
type MockSQLHandlerMockRecorder struct {
	mock *MockSQLHandler
}

// NewMockSQLHandler creates a new mock instance.
func NewMockSQLHandler(ctrl *gomock.Controller) *MockSQLHandler {
	mock := &MockSQLHandler{ctrl: ctrl}
	mock.recorder = &MockSQLHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSQLHandler) EXPECT() *MockSQLHandlerMockRecorder {
	return m.recorder
}

// Begin mocks base method.
func (m *MockSQLHandler) Begin(ctx context.Context) (repository.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Begin", ctx)
	ret0, _ := ret[0].(repository.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Begin indicates an expected call of Begin.
func (mr *MockSQLHandlerMockRecorder) Begin(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Begin", reflect.TypeOf((*MockSQLHandler)(nil).Begin), ctx)
}

// Close mocks base method.
func (m *MockSQLHandler) Close() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close.
func (mr *MockSQLHandlerMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockSQLHandler)(nil).Close))
}

// In mocks base method.
func (m *MockSQLHandler) In(query string, arg interface{}) (*string, []interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "In", query, arg)
	ret0, _ := ret[0].(*string)
	ret1, _ := ret[1].([]interface{})
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// In indicates an expected call of In.
func (mr *MockSQLHandlerMockRecorder) In(query, arg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "In", reflect.TypeOf((*MockSQLHandler)(nil).In), query, arg)
}

// PrepareContext mocks base method.
func (m *MockSQLHandler) PrepareContext(ctx context.Context, query string) (*sqlx.Stmt, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PrepareContext", ctx, query)
	ret0, _ := ret[0].(*sqlx.Stmt)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PrepareContext indicates an expected call of PrepareContext.
func (mr *MockSQLHandlerMockRecorder) PrepareContext(ctx, query interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PrepareContext", reflect.TypeOf((*MockSQLHandler)(nil).PrepareContext), ctx, query)
}

// PrepareNamedContext mocks base method.
func (m *MockSQLHandler) PrepareNamedContext(ctx context.Context, query string) (*sqlx.NamedStmt, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PrepareNamedContext", ctx, query)
	ret0, _ := ret[0].(*sqlx.NamedStmt)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PrepareNamedContext indicates an expected call of PrepareNamedContext.
func (mr *MockSQLHandlerMockRecorder) PrepareNamedContext(ctx, query interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PrepareNamedContext", reflect.TypeOf((*MockSQLHandler)(nil).PrepareNamedContext), ctx, query)
}

// Select mocks base method.
func (m *MockSQLHandler) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, dest, query}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Select", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Select indicates an expected call of Select.
func (mr *MockSQLHandlerMockRecorder) Select(ctx, dest, query interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, dest, query}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Select", reflect.TypeOf((*MockSQLHandler)(nil).Select), varargs...)
}
