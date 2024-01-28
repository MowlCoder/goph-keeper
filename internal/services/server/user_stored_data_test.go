package server

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	"github.com/MowlCoder/goph-keeper/internal/domain"
	mock_server "github.com/MowlCoder/goph-keeper/internal/services/server/mocks"
)

type userStoredDataTestSuite struct {
	suite.Suite

	repository *mock_server.MockuserStoredDataRepository
	cryptor    *mock_server.MockcryptorForUserStoredDataService

	service *UserStoredDataService
}

func (suite *userStoredDataTestSuite) SetupSuite() {
}

func (suite *userStoredDataTestSuite) TearDownSuite() {
}

func (suite *userStoredDataTestSuite) SetupTest() {
	ctrl := gomock.NewController(suite.T())

	suite.repository = mock_server.NewMockuserStoredDataRepository(ctrl)
	suite.cryptor = mock_server.NewMockcryptorForUserStoredDataService(ctrl)

	suite.service = NewUserStoredDataService(suite.repository, suite.cryptor)
}

func (suite *userStoredDataTestSuite) TearDownTest() {
}

func TestUserStoredDataSuite(t *testing.T) {
	suite.Run(t, new(userStoredDataTestSuite))
}

func (suite *userStoredDataTestSuite) TestGetAllUserData() {
	testCases := []struct {
		name    string
		err     error
		prepare func() int
	}{
		{
			name: "valid",
			err:  nil,
			prepare: func() int {
				userID := 1
				data := domain.TextData{
					Text: "123",
				}
				b, _ := json.Marshal(data)
				crypted := []byte{1, 2, 3}

				suite.repository.
					EXPECT().
					GetUserAllData(gomock.Any(), userID).
					Return([]domain.UserStoredData{{ID: 1, UserID: userID, CryptedData: crypted, DataType: domain.TextDataType}}, nil)

				suite.cryptor.
					EXPECT().
					DecryptBytes(crypted).
					Return(b, nil)

				return userID
			},
		},
		{
			name: "invalid GetUserAllData",
			err:  domain.ErrInternal,
			prepare: func() int {
				userID := 1

				suite.repository.
					EXPECT().
					GetUserAllData(gomock.Any(), userID).
					Return(nil, domain.ErrInternal)

				return userID
			},
		},
	}

	for _, testCase := range testCases {
		suite.Run(testCase.name, func() {
			userID := testCase.prepare()
			_, err := suite.service.GetAllUserData(context.Background(), userID)
			suite.Equal(testCase.err, err)
		})
	}
}

func (suite *userStoredDataTestSuite) TestGetUserData() {
	testCases := []struct {
		name    string
		err     error
		prepare func() (int, string, *domain.StorageFilters)
	}{
		{
			name: "valid",
			err:  nil,
			prepare: func() (int, string, *domain.StorageFilters) {
				userID := 1
				dataType := domain.TextDataType
				filters := &domain.StorageFilters{
					IsPaginated:    false,
					IsSortedByDate: false,
				}
				data := domain.TextData{
					Text: "123",
				}
				b, _ := json.Marshal(data)
				crypted := []byte{1, 2, 3}

				suite.repository.
					EXPECT().
					GetWithType(gomock.Any(), userID, dataType, filters).
					Return([]domain.UserStoredData{{ID: 1, UserID: userID, CryptedData: crypted, DataType: domain.TextDataType}}, nil)

				suite.repository.
					EXPECT().
					CountUserDataOfType(gomock.Any(), userID, dataType).
					Return(1, nil)

				suite.cryptor.
					EXPECT().
					DecryptBytes(crypted).
					Return(b, nil)

				return userID, dataType, filters
			},
		},
		{
			name: "invalid GetWithType",
			err:  domain.ErrInternal,
			prepare: func() (int, string, *domain.StorageFilters) {
				userID := 1
				dataType := domain.TextDataType
				filters := &domain.StorageFilters{
					IsPaginated:    false,
					IsSortedByDate: false,
				}

				suite.repository.
					EXPECT().
					GetWithType(gomock.Any(), userID, dataType, filters).
					Return(nil, domain.ErrInternal)

				return userID, dataType, filters
			},
		},
		{
			name: "err CountUserDataOfType",
			err:  domain.ErrInternal,
			prepare: func() (int, string, *domain.StorageFilters) {
				userID := 1
				dataType := domain.TextDataType
				filters := &domain.StorageFilters{
					IsPaginated:    false,
					IsSortedByDate: false,
				}
				crypted := []byte{1, 2, 3}

				suite.repository.
					EXPECT().
					GetWithType(gomock.Any(), userID, dataType, filters).
					Return([]domain.UserStoredData{{ID: 1, UserID: userID, CryptedData: crypted}}, nil)

				suite.repository.
					EXPECT().
					CountUserDataOfType(gomock.Any(), userID, dataType).
					Return(0, domain.ErrInternal)

				return userID, dataType, filters
			},
		},
	}

	for _, testCase := range testCases {
		suite.Run(testCase.name, func() {
			userID, dataType, filters := testCase.prepare()
			_, err := suite.service.GetUserData(context.Background(), userID, dataType, filters)
			suite.Equal(testCase.err, err)
		})
	}
}

func (suite *userStoredDataTestSuite) TestGetUserDataByID() {
	testCases := []struct {
		name    string
		err     error
		prepare func() (int, int)
	}{
		{
			name: "valid",
			err:  nil,
			prepare: func() (int, int) {
				userID := 1
				id := 1
				data := domain.TextData{
					Text: "text",
				}
				b, _ := json.Marshal(data)
				crypted := []byte("123")

				suite.repository.
					EXPECT().
					GetByID(gomock.Any(), id).
					Return(&domain.UserStoredData{ID: id, UserID: userID, CryptedData: crypted, DataType: domain.TextDataType}, nil)

				suite.cryptor.
					EXPECT().
					DecryptBytes(crypted).
					Return(b, nil)

				return userID, id
			},
		},
		{
			name: "not found",
			err:  domain.ErrUserStoredDataNotFound,
			prepare: func() (int, int) {
				userID := 1
				id := 1

				suite.repository.
					EXPECT().
					GetByID(gomock.Any(), id).
					Return(nil, domain.ErrUserStoredDataNotFound)

				return userID, id
			},
		},
		{
			name: "not found (user id not equal)",
			err:  domain.ErrUserStoredDataNotFound,
			prepare: func() (int, int) {
				userID := 1
				id := 1
				crypted := []byte("123")

				suite.repository.
					EXPECT().
					GetByID(gomock.Any(), id).
					Return(&domain.UserStoredData{ID: id, UserID: 2, CryptedData: crypted}, nil)

				return userID, id
			},
		},
		{
			name: "error when decrypting",
			err:  domain.ErrInternal,
			prepare: func() (int, int) {
				userID := 1
				id := 1
				crypted := []byte("123")

				suite.repository.
					EXPECT().
					GetByID(gomock.Any(), id).
					Return(&domain.UserStoredData{ID: id, UserID: userID, CryptedData: crypted}, nil)

				suite.cryptor.
					EXPECT().
					DecryptBytes(crypted).
					Return(nil, domain.ErrInternal)

				return userID, id
			},
		},
	}

	for _, testCase := range testCases {
		suite.Run(testCase.name, func() {
			userID, id := testCase.prepare()
			_, err := suite.service.GetUserDataByID(context.Background(), userID, id)
			suite.Equal(testCase.err, err)
		})
	}
}

func (suite *userStoredDataTestSuite) UpdateUserData() {
	testCases := []struct {
		name    string
		err     error
		prepare func() (int, int, interface{}, string)
	}{
		{
			name: "valid",
			err:  nil,
			prepare: func() (int, int, interface{}, string) {
				userID := 1
				id := 1
				data := domain.TextData{
					Text: "text",
				}
				b, _ := json.Marshal(data)
				meta := "meta"
				encrypted := []uint8{1, 2, 3}

				suite.cryptor.
					EXPECT().
					EncryptBytes(b).
					Return(encrypted, nil)

				suite.repository.
					EXPECT().
					UpdateUserData(gomock.Any(), userID, id, encrypted, meta).
					Return(&domain.UserStoredData{}, nil)

				return userID, id, data, meta
			},
		},
		{
			name: "error when encrypted",
			err:  domain.ErrInternal,
			prepare: func() (int, int, interface{}, string) {
				userID := 1
				id := 1
				data := domain.TextData{
					Text: "text",
				}
				b, _ := json.Marshal(data)
				meta := "meta"

				suite.cryptor.
					EXPECT().
					EncryptBytes(b).
					Return(nil, domain.ErrInternal)

				return userID, id, data, meta
			},
		},
		{
			name: "invalid update",
			err:  domain.ErrInternal,
			prepare: func() (int, int, interface{}, string) {
				userID := 1
				id := 1
				data := domain.TextData{
					Text: "text",
				}
				b, _ := json.Marshal(data)
				meta := "meta"
				encrypted := []uint8{1, 2, 3}

				suite.cryptor.
					EXPECT().
					EncryptBytes(b).
					Return(encrypted, nil)

				suite.repository.
					EXPECT().
					UpdateUserData(gomock.Any(), userID, id, encrypted, meta).
					Return(nil, domain.ErrInternal)

				return userID, id, data, meta
			},
		},
	}

	for _, testCase := range testCases {
		suite.Run(testCase.name, func() {
			userID, id, data, meta := testCase.prepare()
			_, err := suite.service.UpdateUserData(context.Background(), userID, id, data, meta)
			suite.Equal(testCase.err, err)
		})
	}
}

func (suite *userStoredDataTestSuite) TestAdd() {
	testCases := []struct {
		name    string
		err     error
		prepare func() (int, string, interface{}, string)
	}{
		{
			name: "valid",
			err:  nil,
			prepare: func() (int, string, interface{}, string) {
				userID := 1
				dataType := domain.TextDataType
				data := domain.TextData{
					Text: "text",
				}
				b, _ := json.Marshal(data)
				meta := "meta"
				encrypted := []uint8{1, 2, 3}

				suite.cryptor.
					EXPECT().
					EncryptBytes(b).
					Return(encrypted, nil)

				suite.repository.
					EXPECT().
					AddData(gomock.Any(), userID, dataType, encrypted, meta).
					Return(int64(1), nil)

				return userID, dataType, data, meta
			},
		},
		{
			name: "error when encrypt",
			err:  domain.ErrInternal,
			prepare: func() (int, string, interface{}, string) {
				userID := 1
				dataType := domain.TextDataType
				data := domain.TextData{
					Text: "text",
				}
				b, _ := json.Marshal(data)
				meta := "meta"

				suite.cryptor.
					EXPECT().
					EncryptBytes(b).
					Return(nil, domain.ErrInternal)

				return userID, dataType, data, meta
			},
		},
		{
			name: "error when adding",
			err:  domain.ErrInternal,
			prepare: func() (int, string, interface{}, string) {
				userID := 1
				dataType := domain.TextDataType
				data := domain.TextData{
					Text: "text",
				}
				b, _ := json.Marshal(data)
				meta := "meta"
				encrypted := []uint8{1, 2, 3}

				suite.cryptor.
					EXPECT().
					EncryptBytes(b).
					Return(encrypted, nil)

				suite.repository.
					EXPECT().
					AddData(gomock.Any(), userID, dataType, encrypted, meta).
					Return(int64(0), domain.ErrInternal)

				return userID, dataType, data, meta
			},
		},
	}

	for _, testCase := range testCases {
		suite.Run(testCase.name, func() {
			userID, dataType, data, meta := testCase.prepare()
			_, err := suite.service.Add(context.Background(), userID, dataType, data, meta)
			suite.Equal(testCase.err, err)
		})
	}
}

func (suite *userStoredDataTestSuite) TestDeleteBatch() {
	testCases := []struct {
		name    string
		err     error
		prepare func() (int, []int)
	}{
		{
			name: "valid",
			err:  nil,
			prepare: func() (int, []int) {
				userID := 1
				ids := []int{1}

				suite.repository.
					EXPECT().
					DeleteBatch(gomock.Any(), userID, ids).
					Return(nil)

				return userID, ids
			},
		},
		{
			name: "error",
			err:  domain.ErrInternal,
			prepare: func() (int, []int) {
				userID := 1
				ids := []int{1}

				suite.repository.
					EXPECT().
					DeleteBatch(gomock.Any(), userID, ids).
					Return(domain.ErrInternal)

				return userID, ids
			},
		},
	}

	for _, testCase := range testCases {
		suite.Run(testCase.name, func() {
			userID, ids := testCase.prepare()
			err := suite.service.DeleteBatch(context.Background(), userID, ids)
			suite.Equal(testCase.err, err)
		})
	}
}
