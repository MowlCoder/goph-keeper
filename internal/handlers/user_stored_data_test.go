package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	"github.com/MowlCoder/goph-keeper/internal/domain"
	"github.com/MowlCoder/goph-keeper/internal/dtos"
	mock_handlers "github.com/MowlCoder/goph-keeper/internal/handlers/mocks"
	"github.com/MowlCoder/goph-keeper/internal/utils/usercontext"
)

type userStoredDataTestSuite struct {
	suite.Suite

	service *mock_handlers.MockuserStoredDataService

	handler *UserStoredDataHandler
}

func (suite *userStoredDataTestSuite) SetupSuite() {
}

func (suite *userStoredDataTestSuite) TearDownSuite() {
}

func (suite *userStoredDataTestSuite) SetupTest() {
	ctrl := gomock.NewController(suite.T())

	suite.service = mock_handlers.NewMockuserStoredDataService(ctrl)

	suite.handler = NewUserStoredDataHandler(suite.service)
}

func (suite *userStoredDataTestSuite) TearDownTest() {
}

func TestUserStoredDataSuite(t *testing.T) {
	suite.Run(t, new(userStoredDataTestSuite))
}

func (suite *userStoredDataTestSuite) TestGetUserAll() {
	testCases := []struct {
		name       string
		statusCode int
		prepare    func() int
	}{
		{
			name:       "valid",
			statusCode: http.StatusOK,
			prepare: func() int {
				userID := 1

				suite.service.
					EXPECT().
					GetAllUserData(gomock.Any(), userID).
					Return([]domain.UserStoredData{}, nil)

				return userID
			},
		},
		{
			name:       "invalid",
			statusCode: http.StatusInternalServerError,
			prepare: func() int {
				userID := 1

				suite.service.
					EXPECT().
					GetAllUserData(gomock.Any(), userID).
					Return(nil, domain.ErrInternal)

				return userID
			},
		},
	}

	for _, testCase := range testCases {
		suite.Run(testCase.name, func() {
			userID := testCase.prepare()
			r := httptest.NewRequest(http.MethodGet, "/api/v1/data", nil)
			r = r.WithContext(usercontext.SetUserIDToContext(r.Context(), userID))
			w := httptest.NewRecorder()

			suite.handler.GetUserAll(w, r)
			res := w.Result()
			defer res.Body.Close()

			suite.Equal(testCase.statusCode, res.StatusCode)
		})
	}
}

func (suite *userStoredDataTestSuite) TestGetOfType() {
	testCases := []struct {
		name       string
		statusCode int
		prepare    func() (int, string)
	}{
		{
			name:       "valid",
			statusCode: http.StatusOK,
			prepare: func() (int, string) {
				userID := 1
				dataType := domain.TextDataType

				suite.service.
					EXPECT().
					GetUserData(gomock.Any(), userID, dataType, &domain.StorageFilters{
						IsPaginated:    true,
						IsSortedByDate: true,
						Pagination: domain.PaginationFilters{
							Page:  1,
							Count: 50,
						},
						SortDate: domain.SortDateFilters{
							IsASC: false,
						},
					}).
					Return(&domain.PaginatedResult{}, nil)

				return userID, dataType
			},
		},
		{
			name:       "invalid data type",
			statusCode: http.StatusBadRequest,
			prepare: func() (int, string) {
				userID := 1
				dataType := "invalid"

				suite.service.
					EXPECT().
					GetUserData(gomock.Any(), userID, dataType, &domain.StorageFilters{
						IsPaginated:    true,
						IsSortedByDate: true,
						Pagination: domain.PaginationFilters{
							Page:  1,
							Count: 50,
						},
						SortDate: domain.SortDateFilters{
							IsASC: false,
						},
					}).
					Return(nil, domain.ErrInvalidDataType)

				return userID, dataType
			},
		},
		{
			name:       "invalid",
			statusCode: http.StatusInternalServerError,
			prepare: func() (int, string) {
				userID := 1
				dataType := domain.TextDataType

				suite.service.
					EXPECT().
					GetUserData(gomock.Any(), userID, dataType, &domain.StorageFilters{
						IsPaginated:    true,
						IsSortedByDate: true,
						Pagination: domain.PaginationFilters{
							Page:  1,
							Count: 50,
						},
						SortDate: domain.SortDateFilters{
							IsASC: false,
						},
					}).
					Return(nil, domain.ErrInternal)

				return userID, dataType
			},
		},
	}

	for _, testCase := range testCases {
		suite.Run(testCase.name, func() {
			userID, dataType := testCase.prepare()
			r := httptest.NewRequest(http.MethodGet, "/api/v1/data/"+dataType, nil)
			r = r.WithContext(usercontext.SetUserIDToContext(r.Context(), userID))
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("type", dataType)
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
			w := httptest.NewRecorder()

			suite.handler.GetOfType(w, r)
			res := w.Result()
			defer res.Body.Close()

			suite.Equal(testCase.statusCode, res.StatusCode)
		})
	}
}

func (suite *userStoredDataTestSuite) TestAdd() {
	testCases := []struct {
		name       string
		statusCode int
		prepare    func() (int, []byte, string)
	}{
		{
			name:       "valid",
			statusCode: http.StatusCreated,
			prepare: func() (int, []byte, string) {
				userID := 1
				body := dtos.AddNewTextBody{
					Data: domain.TextData{
						Text: "text",
					},
					Meta: "meta",
				}
				b, _ := json.Marshal(body)
				dataType := domain.TextDataType

				suite.service.
					EXPECT().
					Add(gomock.Any(), userID, dataType, body.Data, body.Meta).
					Return(&domain.UserStoredData{}, nil)

				return userID, b, dataType
			},
		},
		{
			name:       "invalid body",
			statusCode: http.StatusBadRequest,
			prepare: func() (int, []byte, string) {
				userID := 1
				body := dtos.AddNewTextBody{
					Data: domain.TextData{
						Text: "",
					},
					Meta: "meta",
				}
				b, _ := json.Marshal(body)
				dataType := domain.TextDataType

				return userID, b, dataType
			},
		},
		{
			name:       "invalid data type",
			statusCode: http.StatusBadRequest,
			prepare: func() (int, []byte, string) {
				userID := 1
				body := dtos.AddNewTextBody{
					Data: domain.TextData{
						Text: "test",
					},
					Meta: "meta",
				}
				b, _ := json.Marshal(body)
				dataType := "test"

				return userID, b, dataType
			},
		},
		{
			name:       "internal error",
			statusCode: http.StatusInternalServerError,
			prepare: func() (int, []byte, string) {
				userID := 1
				body := dtos.AddNewTextBody{
					Data: domain.TextData{
						Text: "text",
					},
					Meta: "meta",
				}
				b, _ := json.Marshal(body)
				dataType := domain.TextDataType

				suite.service.
					EXPECT().
					Add(gomock.Any(), userID, dataType, body.Data, body.Meta).
					Return(nil, domain.ErrInternal)

				return userID, b, dataType
			},
		},
	}

	for _, testCase := range testCases {
		suite.Run(testCase.name, func() {
			userID, body, dataType := testCase.prepare()
			r := httptest.NewRequest(http.MethodPost, "/api/v1/data/"+dataType, bytes.NewReader(body))
			r.Header.Set("Content-Type", "application/json")
			r = r.WithContext(usercontext.SetUserIDToContext(r.Context(), userID))
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("type", dataType)
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
			w := httptest.NewRecorder()

			suite.handler.Add(w, r)
			res := w.Result()
			defer res.Body.Close()

			suite.Equal(testCase.statusCode, res.StatusCode)
		})
	}
}

func (suite *userStoredDataTestSuite) TestUpdateOne() {
	testCases := []struct {
		name       string
		statusCode int
		prepare    func() (int, []byte, string)
	}{
		{
			name:       "valid",
			statusCode: http.StatusOK,
			prepare: func() (int, []byte, string) {
				id := "1"
				userID := 1
				body := dtos.AddNewTextBody{
					Data: domain.TextData{
						Text: "text",
					},
					Meta: "meta",
				}
				b, _ := json.Marshal(body)

				suite.service.
					EXPECT().
					GetUserDataByID(gomock.Any(), userID, 1).
					Return(&domain.UserStoredData{DataType: domain.TextDataType}, nil)

				suite.service.
					EXPECT().
					UpdateUserData(gomock.Any(), userID, 1, body.Data, body.Meta).
					Return(&domain.UserStoredData{}, nil)

				return userID, b, id
			},
		},
		{
			name:       "invalid id",
			statusCode: http.StatusNotFound,
			prepare: func() (int, []byte, string) {
				id := "test"
				userID := 1
				body := dtos.AddNewTextBody{
					Data: domain.TextData{
						Text: "Test",
					},
					Meta: "meta",
				}
				b, _ := json.Marshal(body)

				return userID, b, id
			},
		},
		{
			name:       "data not found by id",
			statusCode: http.StatusNotFound,
			prepare: func() (int, []byte, string) {
				id := "1"
				userID := 1
				body := dtos.AddNewTextBody{
					Data: domain.TextData{
						Text: "Test",
					},
					Meta: "meta",
				}
				b, _ := json.Marshal(body)

				suite.service.
					EXPECT().
					GetUserDataByID(gomock.Any(), userID, 1).
					Return(nil, domain.ErrUserStoredDataNotFound)

				return userID, b, id
			},
		},
		{
			name:       "invalid body",
			statusCode: http.StatusBadRequest,
			prepare: func() (int, []byte, string) {
				id := "1"
				userID := 1
				body := dtos.AddNewTextBody{
					Data: domain.TextData{
						Text: "",
					},
					Meta: "meta",
				}
				b, _ := json.Marshal(body)

				suite.service.
					EXPECT().
					GetUserDataByID(gomock.Any(), userID, 1).
					Return(&domain.UserStoredData{DataType: domain.TextDataType}, nil)

				return userID, b, id
			},
		},
		{
			name:       "internal error",
			statusCode: http.StatusInternalServerError,
			prepare: func() (int, []byte, string) {
				id := "1"
				userID := 1
				body := dtos.AddNewTextBody{
					Data: domain.TextData{
						Text: "text",
					},
					Meta: "meta",
				}
				b, _ := json.Marshal(body)

				suite.service.
					EXPECT().
					GetUserDataByID(gomock.Any(), userID, 1).
					Return(&domain.UserStoredData{DataType: domain.TextDataType}, nil)

				suite.service.
					EXPECT().
					UpdateUserData(gomock.Any(), userID, 1, body.Data, body.Meta).
					Return(nil, domain.ErrInternal)

				return userID, b, id
			},
		},
	}

	for _, testCase := range testCases {
		suite.Run(testCase.name, func() {
			userID, body, id := testCase.prepare()
			r := httptest.NewRequest(http.MethodPut, "/api/v1/data/update"+id, bytes.NewReader(body))
			r.Header.Set("Content-Type", "application/json")
			r = r.WithContext(usercontext.SetUserIDToContext(r.Context(), userID))
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", id)
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
			w := httptest.NewRecorder()

			suite.handler.UpdateOne(w, r)
			res := w.Result()
			defer res.Body.Close()

			suite.Equal(testCase.statusCode, res.StatusCode)
		})
	}
}

func (suite *userStoredDataTestSuite) TestDeleteBatch() {
	testCases := []struct {
		name       string
		statusCode int
		prepare    func() (int, []byte)
	}{
		{
			name:       "valid",
			statusCode: http.StatusNoContent,
			prepare: func() (int, []byte) {
				userID := 1
				body := dtos.DeleteBatchBody{
					IDs: []int{1, 2, 3},
				}
				b, _ := json.Marshal(body)

				suite.service.
					EXPECT().
					DeleteBatch(gomock.Any(), userID, body.IDs).
					Return(nil)

				return userID, b
			},
		},
		{
			name:       "invalid body",
			statusCode: http.StatusBadRequest,
			prepare: func() (int, []byte) {
				userID := 1
				body := dtos.DeleteBatchBody{
					IDs: []int{},
				}
				b, _ := json.Marshal(body)

				return userID, b
			},
		},
		{
			name:       "internal error",
			statusCode: http.StatusInternalServerError,
			prepare: func() (int, []byte) {
				userID := 1
				body := dtos.DeleteBatchBody{
					IDs: []int{1, 2, 3},
				}
				b, _ := json.Marshal(body)

				suite.service.
					EXPECT().
					DeleteBatch(gomock.Any(), userID, body.IDs).
					Return(domain.ErrInternal)

				return userID, b
			},
		},
	}

	for _, testCase := range testCases {
		suite.Run(testCase.name, func() {
			userID, body := testCase.prepare()
			r := httptest.NewRequest(http.MethodDelete, "/api/v1/data", bytes.NewReader(body))
			r.Header.Set("Content-Type", "application/json")
			r = r.WithContext(usercontext.SetUserIDToContext(r.Context(), userID))
			w := httptest.NewRecorder()

			suite.handler.DeleteBatch(w, r)
			res := w.Result()
			defer res.Body.Close()

			suite.Equal(testCase.statusCode, res.StatusCode)
		})
	}
}
