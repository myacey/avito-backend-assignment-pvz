// Code generated by MockGen. DO NOT EDIT.
// Source: ./reception_handler.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
	request "github.com/myacey/avito-backend-assignment-pvz/internal/models/dto/request"
	entity "github.com/myacey/avito-backend-assignment-pvz/internal/models/entity"
)

// MockReceptionService is a mock of ReceptionService interface.
type MockReceptionService struct {
	ctrl     *gomock.Controller
	recorder *MockReceptionServiceMockRecorder
}

// MockReceptionServiceMockRecorder is the mock recorder for MockReceptionService.
type MockReceptionServiceMockRecorder struct {
	mock *MockReceptionService
}

// NewMockReceptionService creates a new mock instance.
func NewMockReceptionService(ctrl *gomock.Controller) *MockReceptionService {
	mock := &MockReceptionService{ctrl: ctrl}
	mock.recorder = &MockReceptionServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockReceptionService) EXPECT() *MockReceptionServiceMockRecorder {
	return m.recorder
}

// AddProductToReception mocks base method.
func (m *MockReceptionService) AddProductToReception(arg0 context.Context, arg1 *request.AddProduct) (*entity.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddProductToReception", arg0, arg1)
	ret0, _ := ret[0].(*entity.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddProductToReception indicates an expected call of AddProductToReception.
func (mr *MockReceptionServiceMockRecorder) AddProductToReception(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddProductToReception", reflect.TypeOf((*MockReceptionService)(nil).AddProductToReception), arg0, arg1)
}

// CreateReception mocks base method.
func (m *MockReceptionService) CreateReception(arg0 context.Context, arg1 *request.CreateReception) (*entity.Reception, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateReception", arg0, arg1)
	ret0, _ := ret[0].(*entity.Reception)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateReception indicates an expected call of CreateReception.
func (mr *MockReceptionServiceMockRecorder) CreateReception(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateReception", reflect.TypeOf((*MockReceptionService)(nil).CreateReception), arg0, arg1)
}

// DeleteLastProduct mocks base method.
func (m *MockReceptionService) DeleteLastProduct(arg0 context.Context, arg1 uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteLastProduct", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteLastProduct indicates an expected call of DeleteLastProduct.
func (mr *MockReceptionServiceMockRecorder) DeleteLastProduct(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteLastProduct", reflect.TypeOf((*MockReceptionService)(nil).DeleteLastProduct), arg0, arg1)
}

// FinishReception mocks base method.
func (m *MockReceptionService) FinishReception(arg0 context.Context, arg1 uuid.UUID) (*entity.Reception, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FinishReception", arg0, arg1)
	ret0, _ := ret[0].(*entity.Reception)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FinishReception indicates an expected call of FinishReception.
func (mr *MockReceptionServiceMockRecorder) FinishReception(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FinishReception", reflect.TypeOf((*MockReceptionService)(nil).FinishReception), arg0, arg1)
}

// SearchReceptions mocks base method.
func (m *MockReceptionService) SearchReceptions(arg0 context.Context, arg1 *request.SearchPvz) ([]*entity.PvzWithReception, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchReceptions", arg0, arg1)
	ret0, _ := ret[0].([]*entity.PvzWithReception)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchReceptions indicates an expected call of SearchReceptions.
func (mr *MockReceptionServiceMockRecorder) SearchReceptions(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchReceptions", reflect.TypeOf((*MockReceptionService)(nil).SearchReceptions), arg0, arg1)
}
