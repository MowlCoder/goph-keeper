package client

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	"github.com/MowlCoder/goph-keeper/internal/domain"
	mock_client "github.com/MowlCoder/goph-keeper/internal/services/client/mocks"
)

type userStoredDataTestSuite struct {
	suite.Suite

	repository *mock_client.MockuserStoredDataRepository
	cryptor    *mock_client.MockcryptorForUserStoredDataService

	service *UserStoredDataService
}

func (suite *userStoredDataTestSuite) SetupSuite() {
}

func (suite *userStoredDataTestSuite) TearDownSuite() {
}

func (suite *userStoredDataTestSuite) SetupTest() {
	ctrl := gomock.NewController(suite.T())

	suite.repository = mock_client.NewMockuserStoredDataRepository(ctrl)
	suite.cryptor = mock_client.NewMockcryptorForUserStoredDataService(ctrl)

	suite.service = NewUserStoredDataService(suite.repository, suite.cryptor)
}

func (suite *userStoredDataTestSuite) TearDownTest() {
}

func TestUserStoredDataSuite(t *testing.T) {
	suite.Run(t, new(userStoredDataTestSuite))
}

func (suite *userStoredDataTestSuite) TestGetByID() {
	testCases := []struct {
		name    string
		err     error
		prepare func() int
	}{
		{
			name: "valid",
			err:  nil,
			prepare: func() int {
				id := 1
				data := domain.LogPassData{
					Login:    "Test",
					Password: "test",
				}
				b, _ := json.Marshal(data)

				suite.repository.
					EXPECT().
					GetByID(gomock.Any(), id).
					Return(&domain.UserStoredData{ID: id, DataType: domain.LogPassDataType}, nil)

				suite.cryptor.
					EXPECT().
					DecryptBytes(gomock.Any()).
					Return(b, nil)

				return id
			},
		},
		{
			name: "invalid data type",
			err:  domain.ErrInvalidDataType,
			prepare: func() int {
				id := 1
				data := domain.LogPassData{
					Login:    "Test",
					Password: "test",
				}
				b, _ := json.Marshal(data)

				suite.repository.
					EXPECT().
					GetByID(gomock.Any(), id).
					Return(&domain.UserStoredData{ID: id, DataType: "test"}, nil)

				suite.cryptor.
					EXPECT().
					DecryptBytes(gomock.Any()).
					Return(b, nil)

				return id
			},
		},
		{
			name: "user data not found",
			err:  domain.ErrUserStoredDataNotFound,
			prepare: func() int {
				id := 1
				suite.repository.
					EXPECT().
					GetByID(gomock.Any(), id).
					Return(nil, domain.ErrUserStoredDataNotFound)

				return id
			},
		},
		{
			name: "decryption err",
			err:  domain.ErrInternal,
			prepare: func() int {
				id := 1
				suite.repository.
					EXPECT().
					GetByID(gomock.Any(), id).
					Return(&domain.UserStoredData{ID: id}, nil)

				suite.cryptor.
					EXPECT().
					DecryptBytes(gomock.Any()).
					Return(nil, domain.ErrInternal)

				return id
			},
		},
	}

	for _, testCase := range testCases {
		suite.Run(testCase.name, func() {
			id := testCase.prepare()
			_, err := suite.service.GetByID(context.Background(), id)
			suite.Equal(testCase.err, err)
		})
	}
}

func (suite *userStoredDataTestSuite) TestGetAll() {
	testCases := []struct {
		name    string
		err     error
		prepare func()
	}{
		{
			name: "valid",
			err:  nil,
			prepare: func() {
				data := domain.LogPassData{
					Login:    "Test",
					Password: "test",
				}
				b, _ := json.Marshal(data)

				suite.repository.
					EXPECT().
					GetAll(gomock.Any()).
					Return([]domain.UserStoredData{{ID: 1, DataType: domain.LogPassDataType}}, nil)

				suite.cryptor.
					EXPECT().
					DecryptBytes(gomock.Any()).
					Return(b, nil)
			},
		},
		{
			name: "invalid data type",
			err:  domain.ErrInvalidDataType,
			prepare: func() {
				data := domain.LogPassData{
					Login:    "Test",
					Password: "test",
				}
				b, _ := json.Marshal(data)

				suite.repository.
					EXPECT().
					GetAll(gomock.Any()).
					Return([]domain.UserStoredData{{ID: 1, DataType: "test"}}, nil)

				suite.cryptor.
					EXPECT().
					DecryptBytes(gomock.Any()).
					Return(b, nil)
			},
		},
		{
			name: "internal error GetAll",
			err:  domain.ErrInternal,
			prepare: func() {
				suite.repository.
					EXPECT().
					GetAll(gomock.Any()).
					Return(nil, domain.ErrInternal)
			},
		},
		{
			name: "decryption err",
			err:  domain.ErrInternal,
			prepare: func() {
				suite.repository.
					EXPECT().
					GetAll(gomock.Any()).
					Return([]domain.UserStoredData{{ID: 1}}, nil)

				suite.cryptor.
					EXPECT().
					DecryptBytes(gomock.Any()).
					Return(nil, domain.ErrInternal)
			},
		},
	}

	for _, testCase := range testCases {
		suite.Run(testCase.name, func() {
			testCase.prepare()
			_, err := suite.service.GetAll(context.Background())
			suite.Equal(testCase.err, err)
		})
	}
}

func (suite *userStoredDataTestSuite) TestGetUserData() {
	testCases := []struct {
		name    string
		err     error
		prepare func() (string, *domain.StorageFilters)
	}{
		{
			name: "valid",
			err:  nil,
			prepare: func() (string, *domain.StorageFilters) {
				filters := &domain.StorageFilters{
					IsPaginated:    false,
					IsSortedByDate: false,
				}
				dataType := domain.LogPassDataType
				data := domain.LogPassData{
					Login:    "Test",
					Password: "test",
				}
				b, _ := json.Marshal(data)

				suite.repository.
					EXPECT().
					GetWithType(gomock.Any(), dataType, filters).
					Return([]domain.UserStoredData{{ID: 1, DataType: domain.LogPassDataType}}, nil)

				suite.repository.
					EXPECT().
					CountUserDataOfType(gomock.Any(), dataType).
					Return(1, nil)

				suite.cryptor.
					EXPECT().
					DecryptBytes(gomock.Any()).
					Return(b, nil)

				return dataType, filters
			},
		},
		{
			name: "can't get data with type",
			err:  domain.ErrInternal,
			prepare: func() (string, *domain.StorageFilters) {
				filters := &domain.StorageFilters{
					IsPaginated:    false,
					IsSortedByDate: false,
				}
				dataType := domain.LogPassDataType

				suite.repository.
					EXPECT().
					GetWithType(gomock.Any(), dataType, filters).
					Return(nil, domain.ErrInternal)

				return dataType, filters
			},
		},
		{
			name: "can't get count",
			err:  domain.ErrInternal,
			prepare: func() (string, *domain.StorageFilters) {
				filters := &domain.StorageFilters{
					IsPaginated:    false,
					IsSortedByDate: false,
				}
				dataType := domain.LogPassDataType

				suite.repository.
					EXPECT().
					GetWithType(gomock.Any(), dataType, filters).
					Return([]domain.UserStoredData{{ID: 1, DataType: domain.LogPassDataType}}, nil)

				suite.repository.
					EXPECT().
					CountUserDataOfType(gomock.Any(), dataType).
					Return(0, domain.ErrInternal)

				return dataType, filters
			},
		},
		{
			name: "decryption err",
			err:  domain.ErrInternal,
			prepare: func() (string, *domain.StorageFilters) {
				filters := &domain.StorageFilters{
					IsPaginated:    false,
					IsSortedByDate: false,
				}
				dataType := domain.LogPassDataType

				suite.repository.
					EXPECT().
					GetWithType(gomock.Any(), dataType, filters).
					Return([]domain.UserStoredData{{ID: 1, DataType: domain.LogPassDataType}}, nil)

				suite.repository.
					EXPECT().
					CountUserDataOfType(gomock.Any(), dataType).
					Return(1, nil)

				suite.cryptor.
					EXPECT().
					DecryptBytes(gomock.Any()).
					Return(nil, domain.ErrInternal)

				return dataType, filters
			},
		},
	}

	for _, testCase := range testCases {
		suite.Run(testCase.name, func() {
			dataType, filters := testCase.prepare()
			_, err := suite.service.GetUserData(context.Background(), dataType, filters)
			suite.Equal(testCase.err, err)
		})
	}
}

func (suite *userStoredDataTestSuite) TestAdd() {
	testCases := []struct {
		name    string
		err     error
		prepare func() (string, interface{}, string)
	}{
		{
			name: "valid",
			err:  nil,
			prepare: func() (string, interface{}, string) {
				dataType := domain.LogPassDataType
				data := domain.LogPassData{
					Login:    "Test",
					Password: "test",
				}
				b, _ := json.Marshal(data)
				meta := "meta"
				cryptedBytes := []uint8{1, 2, 3}

				suite.cryptor.
					EXPECT().
					EncryptBytes(b).
					Return(cryptedBytes, nil)

				suite.repository.
					EXPECT().
					AddData(gomock.Any(), dataType, cryptedBytes, meta).
					Return(int64(1), nil)

				return dataType, data, meta
			},
		},
		{
			name: "encrypt error",
			err:  domain.ErrInternal,
			prepare: func() (string, interface{}, string) {
				dataType := domain.LogPassDataType
				data := domain.LogPassData{
					Login:    "Test",
					Password: "test",
				}
				b, _ := json.Marshal(data)
				meta := "meta"

				suite.cryptor.
					EXPECT().
					EncryptBytes(b).
					Return(nil, domain.ErrInternal)

				return dataType, data, meta
			},
		},
		{
			name: "add error",
			err:  domain.ErrInternal,
			prepare: func() (string, interface{}, string) {
				dataType := domain.LogPassDataType
				data := domain.LogPassData{
					Login:    "Test",
					Password: "test",
				}
				b, _ := json.Marshal(data)
				meta := "meta"
				cryptedBytes := []uint8{1, 2, 3}

				suite.cryptor.
					EXPECT().
					EncryptBytes(b).
					Return(cryptedBytes, nil)

				suite.repository.
					EXPECT().
					AddData(gomock.Any(), dataType, cryptedBytes, meta).
					Return(int64(0), domain.ErrInternal)

				return dataType, data, meta
			},
		},
	}

	for _, testCase := range testCases {
		suite.Run(testCase.name, func() {
			dataType, data, meta := testCase.prepare()
			_, err := suite.service.Add(context.Background(), dataType, data, meta)
			suite.Equal(testCase.err, err)
		})
	}
}

func (suite *userStoredDataTestSuite) TestUpdateByID() {
	testCases := []struct {
		name    string
		err     error
		prepare func() (int, interface{}, string)
	}{
		{
			name: "valid",
			err:  nil,
			prepare: func() (int, interface{}, string) {
				id := 1
				data := domain.LogPassData{
					Login:    "Test",
					Password: "test",
				}
				b, _ := json.Marshal(data)
				meta := "meta"
				cryptedBytes := []uint8{1, 2, 3}

				suite.cryptor.
					EXPECT().
					EncryptBytes(b).
					Return(cryptedBytes, nil)

				suite.repository.
					EXPECT().
					UpdateByID(gomock.Any(), id, cryptedBytes, meta).
					Return(&domain.UserStoredData{ID: id}, nil)

				return id, data, meta
			},
		},
		{
			name: "encrypt error",
			err:  domain.ErrInternal,
			prepare: func() (int, interface{}, string) {
				id := 1
				data := domain.LogPassData{
					Login:    "Test",
					Password: "test",
				}
				b, _ := json.Marshal(data)
				meta := "meta"

				suite.cryptor.
					EXPECT().
					EncryptBytes(b).
					Return(nil, domain.ErrInternal)

				return id, data, meta
			},
		},
		{
			name: "update error",
			err:  domain.ErrInternal,
			prepare: func() (int, interface{}, string) {
				id := 1
				data := domain.LogPassData{
					Login:    "Test",
					Password: "test",
				}
				b, _ := json.Marshal(data)
				meta := "meta"
				cryptedBytes := []uint8{1, 2, 3}

				suite.cryptor.
					EXPECT().
					EncryptBytes(b).
					Return(cryptedBytes, nil)

				suite.repository.
					EXPECT().
					UpdateByID(gomock.Any(), id, cryptedBytes, meta).
					Return(nil, domain.ErrInternal)

				return id, data, meta
			},
		},
	}

	for _, testCase := range testCases {
		suite.Run(testCase.name, func() {
			id, data, meta := testCase.prepare()
			_, err := suite.service.UpdateByID(context.Background(), id, data, meta)
			suite.Equal(testCase.err, err)
		})
	}
}

func (suite *userStoredDataTestSuite) TestDeleteBatch() {
	testCases := []struct {
		name    string
		err     error
		prepare func() []int
	}{
		{
			name: "valid",
			err:  nil,
			prepare: func() []int {
				ids := []int{1}
				suite.repository.
					EXPECT().
					DeleteBatch(gomock.Any(), ids).
					Return(nil)
				return ids
			},
		},
		{
			name: "err",
			err:  domain.ErrInternal,
			prepare: func() []int {
				ids := []int{1}
				suite.repository.
					EXPECT().
					DeleteBatch(gomock.Any(), ids).
					Return(domain.ErrInternal)
				return ids
			},
		},
	}

	for _, testCase := range testCases {
		suite.Run(testCase.name, func() {
			ids := testCase.prepare()
			err := suite.service.DeleteBatch(context.Background(), ids)
			suite.Equal(testCase.err, err)
		})
	}
}

func (suite *userStoredDataTestSuite) TestDeleteByID() {
	testCases := []struct {
		name    string
		err     error
		prepare func() int
	}{
		{
			name: "valid",
			err:  nil,
			prepare: func() int {
				id := 1
				suite.repository.
					EXPECT().
					DeleteByID(gomock.Any(), id).
					Return(nil)
				return id
			},
		},
		{
			name: "err",
			err:  domain.ErrInternal,
			prepare: func() int {
				id := 1
				suite.repository.
					EXPECT().
					DeleteByID(gomock.Any(), id).
					Return(domain.ErrInternal)
				return id
			},
		},
	}

	for _, testCase := range testCases {
		suite.Run(testCase.name, func() {
			id := testCase.prepare()
			err := suite.service.DeleteByID(context.Background(), id)
			suite.Equal(testCase.err, err)
		})
	}
}
