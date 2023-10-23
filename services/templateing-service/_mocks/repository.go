// Code generated by MockGen. DO NOT EDIT.
// Source: templates/repository/repository_interface.go
//
// Generated by this command:
//
//	mockgen -package=mocks -destination=_mocks/repository.go -source=templates/repository/repository_interface.go
//
// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	model "github.com/hingew/hsfl-master-ai-cloud-engineering/lib/model"
	gomock "go.uber.org/mock/gomock"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// CreateTemplate mocks base method.
func (m *MockRepository) CreateTemplate(data model.PdfTemplate) (*uint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTemplate", data)
	ret0, _ := ret[0].(*uint)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTemplate indicates an expected call of CreateTemplate.
func (mr *MockRepositoryMockRecorder) CreateTemplate(data any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTemplate", reflect.TypeOf((*MockRepository)(nil).CreateTemplate), data)
}

// DeleteTemplate mocks base method.
func (m *MockRepository) DeleteTemplate(id uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteTemplate", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteTemplate indicates an expected call of DeleteTemplate.
func (mr *MockRepositoryMockRecorder) DeleteTemplate(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTemplate", reflect.TypeOf((*MockRepository)(nil).DeleteTemplate), id)
}

// GetAllTemplates mocks base method.
func (m *MockRepository) GetAllTemplates() (*[]model.PdfTemplate, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllTemplates")
	ret0, _ := ret[0].(*[]model.PdfTemplate)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllTemplates indicates an expected call of GetAllTemplates.
func (mr *MockRepositoryMockRecorder) GetAllTemplates() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllTemplates", reflect.TypeOf((*MockRepository)(nil).GetAllTemplates))
}

// GetTemplateById mocks base method.
func (m *MockRepository) GetTemplateById(id uint) (*model.PdfTemplate, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTemplateById", id)
	ret0, _ := ret[0].(*model.PdfTemplate)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTemplateById indicates an expected call of GetTemplateById.
func (mr *MockRepositoryMockRecorder) GetTemplateById(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTemplateById", reflect.TypeOf((*MockRepository)(nil).GetTemplateById), id)
}

// UpdateTemplate mocks base method.
func (m *MockRepository) UpdateTemplate(id uint, data model.PdfTemplate) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTemplate", id, data)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateTemplate indicates an expected call of UpdateTemplate.
func (mr *MockRepositoryMockRecorder) UpdateTemplate(id, data any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTemplate", reflect.TypeOf((*MockRepository)(nil).UpdateTemplate), id, data)
}
