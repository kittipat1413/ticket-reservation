// Code generated by MockGen. DO NOT EDIT.
// Source: ./health_check_repository.go

// Package repository_mocks is a generated GoMock package.
package repository_mocks

import (
	context "context"
	reflect "reflect"
	repository "ticket-reservation/internal/domain/repository"
	db "ticket-reservation/internal/infra/db"

	gomock "github.com/golang/mock/gomock"
)

// MockHealthCheckRepository is a mock of HealthCheckRepository interface.
type MockHealthCheckRepository struct {
	ctrl     *gomock.Controller
	recorder *MockHealthCheckRepositoryMockRecorder
}

// MockHealthCheckRepositoryMockRecorder is the mock recorder for MockHealthCheckRepository.
type MockHealthCheckRepositoryMockRecorder struct {
	mock *MockHealthCheckRepository
}

// NewMockHealthCheckRepository creates a new mock instance.
func NewMockHealthCheckRepository(ctrl *gomock.Controller) *MockHealthCheckRepository {
	mock := &MockHealthCheckRepository{ctrl: ctrl}
	mock.recorder = &MockHealthCheckRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHealthCheckRepository) EXPECT() *MockHealthCheckRepositoryMockRecorder {
	return m.recorder
}

// CheckDatabaseReadiness mocks base method.
func (m *MockHealthCheckRepository) CheckDatabaseReadiness(ctx context.Context) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckDatabaseReadiness", ctx)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckDatabaseReadiness indicates an expected call of CheckDatabaseReadiness.
func (mr *MockHealthCheckRepositoryMockRecorder) CheckDatabaseReadiness(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckDatabaseReadiness", reflect.TypeOf((*MockHealthCheckRepository)(nil).CheckDatabaseReadiness), ctx)
}

// WithTx mocks base method.
func (m *MockHealthCheckRepository) WithTx(tx db.SqlExecer) repository.HealthCheckRepository {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithTx", tx)
	ret0, _ := ret[0].(repository.HealthCheckRepository)
	return ret0
}

// WithTx indicates an expected call of WithTx.
func (mr *MockHealthCheckRepositoryMockRecorder) WithTx(tx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithTx", reflect.TypeOf((*MockHealthCheckRepository)(nil).WithTx), tx)
}
