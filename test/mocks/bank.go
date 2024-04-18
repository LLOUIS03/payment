// Code generated by MockGen. DO NOT EDIT.
// Source: domain/clients/bank.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

// MockBank is a mock of Bank interface.
type MockBank struct {
	ctrl     *gomock.Controller
	recorder *MockBankMockRecorder
}

// MockBankMockRecorder is the mock recorder for MockBank.
type MockBankMockRecorder struct {
	mock *MockBank
}

// NewMockBank creates a new mock instance.
func NewMockBank(ctrl *gomock.Controller) *MockBank {
	mock := &MockBank{ctrl: ctrl}
	mock.recorder = &MockBankMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBank) EXPECT() *MockBankMockRecorder {
	return m.recorder
}

// Place mocks base method.
func (m *MockBank) Place(arg0 context.Context, arg1 uuid.UUID, arg2 float64, arg3 uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Place", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// Place indicates an expected call of Place.
func (mr *MockBankMockRecorder) Place(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Place", reflect.TypeOf((*MockBank)(nil).Place), arg0, arg1, arg2, arg3)
}

// Refund mocks base method.
func (m *MockBank) Refund(arg0 context.Context, arg1 uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Refund", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Refund indicates an expected call of Refund.
func (mr *MockBankMockRecorder) Refund(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Refund", reflect.TypeOf((*MockBank)(nil).Refund), arg0, arg1)
}