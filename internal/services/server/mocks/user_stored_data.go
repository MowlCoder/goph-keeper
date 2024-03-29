// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/services/server/user_stored_data.go
//
// Generated by this command:
//
//	mockgen -source=./internal/services/server/user_stored_data.go -destination=./internal/services/server/mocks/user_stored_data.go
//
// Package mock_server is a generated GoMock package.
package mock_server

import (
	context "context"
	reflect "reflect"

	domain "github.com/MowlCoder/goph-keeper/internal/domain"
	gomock "go.uber.org/mock/gomock"
)

// MockcryptorForUserStoredDataService is a mock of cryptorForUserStoredDataService interface.
type MockcryptorForUserStoredDataService struct {
	ctrl     *gomock.Controller
	recorder *MockcryptorForUserStoredDataServiceMockRecorder
}

// MockcryptorForUserStoredDataServiceMockRecorder is the mock recorder for MockcryptorForUserStoredDataService.
type MockcryptorForUserStoredDataServiceMockRecorder struct {
	mock *MockcryptorForUserStoredDataService
}

// NewMockcryptorForUserStoredDataService creates a new mock instance.
func NewMockcryptorForUserStoredDataService(ctrl *gomock.Controller) *MockcryptorForUserStoredDataService {
	mock := &MockcryptorForUserStoredDataService{ctrl: ctrl}
	mock.recorder = &MockcryptorForUserStoredDataServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockcryptorForUserStoredDataService) EXPECT() *MockcryptorForUserStoredDataServiceMockRecorder {
	return m.recorder
}

// DecryptBytes mocks base method.
func (m *MockcryptorForUserStoredDataService) DecryptBytes(crypted []byte) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DecryptBytes", crypted)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DecryptBytes indicates an expected call of DecryptBytes.
func (mr *MockcryptorForUserStoredDataServiceMockRecorder) DecryptBytes(crypted any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DecryptBytes", reflect.TypeOf((*MockcryptorForUserStoredDataService)(nil).DecryptBytes), crypted)
}

// EncryptBytes mocks base method.
func (m *MockcryptorForUserStoredDataService) EncryptBytes(raw []byte) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EncryptBytes", raw)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EncryptBytes indicates an expected call of EncryptBytes.
func (mr *MockcryptorForUserStoredDataServiceMockRecorder) EncryptBytes(raw any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EncryptBytes", reflect.TypeOf((*MockcryptorForUserStoredDataService)(nil).EncryptBytes), raw)
}

// MockuserStoredDataRepository is a mock of userStoredDataRepository interface.
type MockuserStoredDataRepository struct {
	ctrl     *gomock.Controller
	recorder *MockuserStoredDataRepositoryMockRecorder
}

// MockuserStoredDataRepositoryMockRecorder is the mock recorder for MockuserStoredDataRepository.
type MockuserStoredDataRepositoryMockRecorder struct {
	mock *MockuserStoredDataRepository
}

// NewMockuserStoredDataRepository creates a new mock instance.
func NewMockuserStoredDataRepository(ctrl *gomock.Controller) *MockuserStoredDataRepository {
	mock := &MockuserStoredDataRepository{ctrl: ctrl}
	mock.recorder = &MockuserStoredDataRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockuserStoredDataRepository) EXPECT() *MockuserStoredDataRepositoryMockRecorder {
	return m.recorder
}

// AddData mocks base method.
func (m *MockuserStoredDataRepository) AddData(ctx context.Context, userID int, dataType string, data []byte, meta string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddData", ctx, userID, dataType, data, meta)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddData indicates an expected call of AddData.
func (mr *MockuserStoredDataRepositoryMockRecorder) AddData(ctx, userID, dataType, data, meta any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddData", reflect.TypeOf((*MockuserStoredDataRepository)(nil).AddData), ctx, userID, dataType, data, meta)
}

// CountUserDataOfType mocks base method.
func (m *MockuserStoredDataRepository) CountUserDataOfType(ctx context.Context, userID int, dataType string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CountUserDataOfType", ctx, userID, dataType)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CountUserDataOfType indicates an expected call of CountUserDataOfType.
func (mr *MockuserStoredDataRepositoryMockRecorder) CountUserDataOfType(ctx, userID, dataType any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CountUserDataOfType", reflect.TypeOf((*MockuserStoredDataRepository)(nil).CountUserDataOfType), ctx, userID, dataType)
}

// DeleteBatch mocks base method.
func (m *MockuserStoredDataRepository) DeleteBatch(ctx context.Context, userID int, id []int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteBatch", ctx, userID, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteBatch indicates an expected call of DeleteBatch.
func (mr *MockuserStoredDataRepositoryMockRecorder) DeleteBatch(ctx, userID, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteBatch", reflect.TypeOf((*MockuserStoredDataRepository)(nil).DeleteBatch), ctx, userID, id)
}

// GetByID mocks base method.
func (m *MockuserStoredDataRepository) GetByID(ctx context.Context, id int) (*domain.UserStoredData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(*domain.UserStoredData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockuserStoredDataRepositoryMockRecorder) GetByID(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockuserStoredDataRepository)(nil).GetByID), ctx, id)
}

// GetUserAllData mocks base method.
func (m *MockuserStoredDataRepository) GetUserAllData(ctx context.Context, userID int) ([]domain.UserStoredData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserAllData", ctx, userID)
	ret0, _ := ret[0].([]domain.UserStoredData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserAllData indicates an expected call of GetUserAllData.
func (mr *MockuserStoredDataRepositoryMockRecorder) GetUserAllData(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserAllData", reflect.TypeOf((*MockuserStoredDataRepository)(nil).GetUserAllData), ctx, userID)
}

// GetWithType mocks base method.
func (m *MockuserStoredDataRepository) GetWithType(ctx context.Context, userID int, dataType string, filters *domain.StorageFilters) ([]domain.UserStoredData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWithType", ctx, userID, dataType, filters)
	ret0, _ := ret[0].([]domain.UserStoredData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWithType indicates an expected call of GetWithType.
func (mr *MockuserStoredDataRepositoryMockRecorder) GetWithType(ctx, userID, dataType, filters any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWithType", reflect.TypeOf((*MockuserStoredDataRepository)(nil).GetWithType), ctx, userID, dataType, filters)
}

// UpdateUserData mocks base method.
func (m *MockuserStoredDataRepository) UpdateUserData(ctx context.Context, userID, dataID int, data any, meta string) (*domain.UserStoredData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserData", ctx, userID, dataID, data, meta)
	ret0, _ := ret[0].(*domain.UserStoredData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateUserData indicates an expected call of UpdateUserData.
func (mr *MockuserStoredDataRepositoryMockRecorder) UpdateUserData(ctx, userID, dataID, data, meta any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserData", reflect.TypeOf((*MockuserStoredDataRepository)(nil).UpdateUserData), ctx, userID, dataID, data, meta)
}
