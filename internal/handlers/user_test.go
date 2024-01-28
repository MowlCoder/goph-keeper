package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	"github.com/MowlCoder/goph-keeper/internal/domain"
	"github.com/MowlCoder/goph-keeper/internal/dtos"
	mock_handlers "github.com/MowlCoder/goph-keeper/internal/handlers/mocks"
)

type userTestSuite struct {
	suite.Suite

	service        *mock_handlers.MockuserService
	tokenGenerator *mock_handlers.MocktokenGenerator

	handler *UserHandler
}

func (suite *userTestSuite) SetupSuite() {
}

func (suite *userTestSuite) TearDownSuite() {
}

func (suite *userTestSuite) SetupTest() {
	ctrl := gomock.NewController(suite.T())

	suite.service = mock_handlers.NewMockuserService(ctrl)
	suite.tokenGenerator = mock_handlers.NewMocktokenGenerator(ctrl)

	suite.handler = NewUserHandler(suite.service, suite.tokenGenerator)
}

func (suite *userTestSuite) TearDownTest() {
}

func TestUserSuite(t *testing.T) {
	suite.Run(t, new(userTestSuite))
}

func (suite *userTestSuite) TestAuthorize() {
	testCases := []struct {
		name       string
		statusCode int
		prepare    func() []byte
	}{
		{
			name:       "valid",
			statusCode: http.StatusOK,
			prepare: func() []byte {
				body := dtos.AuthorizeBody{
					Email:    "test@gmail.com",
					Password: "test123",
				}
				b, _ := json.Marshal(body)

				suite.service.
					EXPECT().
					Authorize(gomock.Any(), body.Email, body.Password).
					Return(&domain.User{ID: 1}, nil)

				suite.tokenGenerator.
					EXPECT().
					Generate(gomock.Any(), domain.User{ID: 1}).
					Return("token", nil)

				return b
			},
		},
		{
			name:       "invalid body",
			statusCode: http.StatusBadRequest,
			prepare: func() []byte {
				body := dtos.AuthorizeBody{
					Email:    "",
					Password: "",
				}
				b, _ := json.Marshal(body)

				return b
			},
		},
		{
			name:       "user not found",
			statusCode: http.StatusNotFound,
			prepare: func() []byte {
				body := dtos.AuthorizeBody{
					Email:    "test@gmail.com",
					Password: "test123",
				}
				b, _ := json.Marshal(body)

				suite.service.
					EXPECT().
					Authorize(gomock.Any(), body.Email, body.Password).
					Return(nil, domain.ErrUserNotFound)

				return b
			},
		},
		{
			name:       "error when generate token",
			statusCode: http.StatusInternalServerError,
			prepare: func() []byte {
				body := dtos.AuthorizeBody{
					Email:    "test@gmail.com",
					Password: "test123",
				}
				b, _ := json.Marshal(body)

				suite.service.
					EXPECT().
					Authorize(gomock.Any(), body.Email, body.Password).
					Return(&domain.User{ID: 1}, nil)

				suite.tokenGenerator.
					EXPECT().
					Generate(gomock.Any(), domain.User{ID: 1}).
					Return("", domain.ErrInternal)

				return b
			},
		},
	}

	for _, testCase := range testCases {
		suite.Run(testCase.name, func() {
			body := testCase.prepare()
			r := httptest.NewRequest(http.MethodPost, "/api/v1/user/authorize", bytes.NewReader(body))
			r.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			suite.handler.Authorize(w, r)

			res := w.Result()
			defer res.Body.Close()

			suite.Equal(testCase.statusCode, res.StatusCode)
		})
	}
}

func (suite *userTestSuite) TestRegister() {
	testCases := []struct {
		name       string
		statusCode int
		prepare    func() []byte
	}{
		{
			name:       "valid",
			statusCode: http.StatusCreated,
			prepare: func() []byte {
				body := dtos.RegisterBody{
					Email:    "test@gmail.com",
					Password: "Test123!",
				}
				b, _ := json.Marshal(body)

				suite.service.
					EXPECT().
					Create(gomock.Any(), body.Email, body.Password).
					Return(&domain.User{ID: 1}, nil)

				suite.tokenGenerator.
					EXPECT().
					Generate(gomock.Any(), domain.User{ID: 1}).
					Return("token", nil)

				return b
			},
		},
		{
			name:       "invalid body",
			statusCode: http.StatusBadRequest,
			prepare: func() []byte {
				body := dtos.RegisterBody{
					Email:    "",
					Password: "",
				}
				b, _ := json.Marshal(body)

				return b
			},
		},
		{
			name:       "user already existed",
			statusCode: http.StatusConflict,
			prepare: func() []byte {
				body := dtos.RegisterBody{
					Email:    "test@gmail.com",
					Password: "Test123!",
				}
				b, _ := json.Marshal(body)

				suite.service.
					EXPECT().
					Create(gomock.Any(), body.Email, body.Password).
					Return(nil, domain.ErrEmailAlreadyTaken)

				return b
			},
		},
		{
			name:       "error when generate token",
			statusCode: http.StatusInternalServerError,
			prepare: func() []byte {
				body := dtos.RegisterBody{
					Email:    "test@gmail.com",
					Password: "Test123!",
				}
				b, _ := json.Marshal(body)

				suite.service.
					EXPECT().
					Create(gomock.Any(), body.Email, body.Password).
					Return(&domain.User{ID: 1}, nil)

				suite.tokenGenerator.
					EXPECT().
					Generate(gomock.Any(), domain.User{ID: 1}).
					Return("", domain.ErrInternal)

				return b
			},
		},
	}

	for _, testCase := range testCases {
		suite.Run(testCase.name, func() {
			body := testCase.prepare()
			r := httptest.NewRequest(http.MethodPost, "/api/v1/user/register", bytes.NewReader(body))
			r.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			suite.handler.Register(w, r)

			res := w.Result()
			defer res.Body.Close()

			suite.Equal(testCase.statusCode, res.StatusCode)
		})
	}
}
