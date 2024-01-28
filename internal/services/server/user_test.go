package server

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	"github.com/MowlCoder/goph-keeper/internal/domain"
	mock_server "github.com/MowlCoder/goph-keeper/internal/services/server/mocks"
)

type userTestSuite struct {
	suite.Suite

	repository *mock_server.MockuserRepository
	hasher     *mock_server.MockpasswordHasher

	service *UserService
}

func (suite *userTestSuite) SetupSuite() {
}

func (suite *userTestSuite) TearDownSuite() {
}

func (suite *userTestSuite) SetupTest() {
	ctrl := gomock.NewController(suite.T())

	suite.repository = mock_server.NewMockuserRepository(ctrl)
	suite.hasher = mock_server.NewMockpasswordHasher(ctrl)

	suite.service = NewUserService(suite.repository, suite.hasher)
}

func (suite *userTestSuite) TearDownTest() {
}

func TestUserSuite(t *testing.T) {
	suite.Run(t, new(userTestSuite))
}

func (suite *userTestSuite) TestCreate() {
	testCases := []struct {
		name    string
		err     error
		prepare func() (string, string)
	}{
		{
			name: "valid",
			err:  nil,
			prepare: func() (string, string) {
				email, password := "test@gmail.com", "test"
				hash := "hashed-password"

				suite.hasher.
					EXPECT().
					Hash(password).
					Return(hash, nil)

				suite.repository.
					EXPECT().
					Create(gomock.Any(), email, hash).
					Return(&domain.User{}, nil)

				return email, password
			},
		},
		{
			name: "hash error",
			err:  domain.ErrInternal,
			prepare: func() (string, string) {
				email, password := "test@gmail.com", "test"

				suite.hasher.
					EXPECT().
					Hash(password).
					Return("", domain.ErrInternal)

				return email, password
			},
		},
		{
			name: "email already taken",
			err:  domain.ErrEmailAlreadyTaken,
			prepare: func() (string, string) {
				email, password := "test@gmail.com", "test"
				hash := "hashed-password"

				suite.hasher.
					EXPECT().
					Hash(password).
					Return(hash, nil)

				suite.repository.
					EXPECT().
					Create(gomock.Any(), email, hash).
					Return(nil, domain.ErrEmailAlreadyTaken)

				return email, password
			},
		},
		{
			name: "create error",
			err:  domain.ErrInternal,
			prepare: func() (string, string) {
				email, password := "test@gmail.com", "test"
				hash := "hashed-password"

				suite.hasher.
					EXPECT().
					Hash(password).
					Return(hash, nil)

				suite.repository.
					EXPECT().
					Create(gomock.Any(), email, hash).
					Return(nil, domain.ErrInternal)

				return email, password
			},
		},
	}

	for _, testCase := range testCases {
		suite.Run(testCase.name, func() {
			email, password := testCase.prepare()
			_, err := suite.service.Create(context.Background(), email, password)
			suite.Equal(testCase.err, err)
		})
	}
}

func (suite *userTestSuite) TestAuthorize() {
	testCases := []struct {
		name    string
		err     error
		prepare func() (string, string)
	}{
		{
			name: "valid",
			err:  nil,
			prepare: func() (string, string) {
				email, password := "test@gmail.com", "test"
				hash := "hashed-password"

				suite.repository.
					EXPECT().
					GetByEmail(gomock.Any(), email).
					Return(&domain.User{Password: hash}, nil)

				suite.hasher.
					EXPECT().
					Equal(password, hash).
					Return(true)

				return email, password
			},
		},
		{
			name: "user not found",
			err:  domain.ErrUserNotFound,
			prepare: func() (string, string) {
				email, password := "test@gmail.com", "test"

				suite.repository.
					EXPECT().
					GetByEmail(gomock.Any(), email).
					Return(nil, domain.ErrUserNotFound)

				return email, password
			},
		},
		{
			name: "hash not equals",
			err:  domain.ErrWrongCredentials,
			prepare: func() (string, string) {
				email, password := "test@gmail.com", "test"
				hash := "hashed-password"

				suite.repository.
					EXPECT().
					GetByEmail(gomock.Any(), email).
					Return(&domain.User{Password: hash}, nil)

				suite.hasher.
					EXPECT().
					Equal(password, hash).
					Return(false)

				return email, password
			},
		},
	}

	for _, testCase := range testCases {
		suite.Run(testCase.name, func() {
			email, password := testCase.prepare()
			_, err := suite.service.Authorize(context.Background(), email, password)
			suite.Equal(testCase.err, err)
		})
	}
}
