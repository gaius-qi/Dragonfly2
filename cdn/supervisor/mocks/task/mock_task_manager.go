// Code generated by MockGen. DO NOT EDIT.
// Source: d7y.io/dragonfly/v2/cdn/supervisor/task (interfaces: Manager)

// Package mock is a generated GoMock package.
package task

import (
	reflect "reflect"

	task "d7y.io/dragonfly/v2/cdn/supervisor/task"
	gomock "github.com/golang/mock/gomock"
)

// MockManager is a mock of Manager interface.
type MockManager struct {
	ctrl     *gomock.Controller
	recorder *MockManagerMockRecorder
}

// MockManagerMockRecorder is the mock recorder for MockManager.
type MockManagerMockRecorder struct {
	mock *MockManager
}

// NewMockManager creates a new mock instance.
func NewMockManager(ctrl *gomock.Controller) *MockManager {
	mock := &MockManager{ctrl: ctrl}
	mock.recorder = &MockManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockManager) EXPECT() *MockManagerMockRecorder {
	return m.recorder
}

// AddOrUpdate mocks base method.
func (m *MockManager) AddOrUpdate(arg0 *task.SeedTask) (*task.SeedTask, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddOrUpdate", arg0)
	ret0, _ := ret[0].(*task.SeedTask)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddOrUpdate indicates an expected call of AddOrUpdate.
func (mr *MockManagerMockRecorder) AddOrUpdate(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddOrUpdate", reflect.TypeOf((*MockManager)(nil).AddOrUpdate), arg0)
}

// Delete mocks base method.
func (m *MockManager) Delete(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Delete", arg0)
}

// Delete indicates an expected call of Delete.
func (mr *MockManagerMockRecorder) Delete(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockManager)(nil).Delete), arg0)
}

// Exist mocks base method.
func (m *MockManager) Exist(arg0 string) (*task.SeedTask, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Exist", arg0)
	ret0, _ := ret[0].(*task.SeedTask)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// Exist indicates an expected call of Exist.
func (mr *MockManagerMockRecorder) Exist(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exist", reflect.TypeOf((*MockManager)(nil).Exist), arg0)
}

// Get mocks base method.
func (m *MockManager) Get(arg0 string) (*task.SeedTask, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0)
	ret0, _ := ret[0].(*task.SeedTask)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockManagerMockRecorder) Get(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockManager)(nil).Get), arg0)
}

// Update mocks base method.
func (m *MockManager) Update(arg0 string, arg1 *task.SeedTask) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockManagerMockRecorder) Update(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockManager)(nil).Update), arg0, arg1)
}
