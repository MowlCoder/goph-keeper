package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/MowlCoder/goph-keeper/internal/domain"
	"github.com/MowlCoder/goph-keeper/internal/dtos"
	"github.com/MowlCoder/goph-keeper/internal/handlers/httperrors"
	"github.com/MowlCoder/goph-keeper/internal/utils/usercontext"
	"github.com/MowlCoder/goph-keeper/pkg/httputils"
	jsonutil "github.com/MowlCoder/goph-keeper/pkg/jsonutils"
)

type userStoredDataService interface {
	GetAllUserData(ctx context.Context, userID int) ([]domain.UserStoredData, error)
	Add(ctx context.Context, userID int, dataType string, data interface{}, meta string) (*domain.UserStoredData, error)
	GetUserDataByID(ctx context.Context, userID int, id int) (*domain.UserStoredData, error)
	GetUserData(ctx context.Context, userID int, dataType string, filters *domain.StorageFilters) (*domain.PaginatedResult, error)
	UpdateUserData(ctx context.Context, userID int, dataID int, data interface{}, meta string) (*domain.UserStoredData, error)
	DeleteBatch(ctx context.Context, userID int, ids []int) error
}

type UserStoredDataHandler struct {
	service userStoredDataService
}

func NewUserStoredDataHandler(service userStoredDataService) *UserStoredDataHandler {
	return &UserStoredDataHandler{
		service: service,
	}
}

// GetUserAll godoc
// @Summary Get all user saved data
// @Produce json
// @Tags data
// @Security Bearer
// @Success 200 {array} domain.UserStoredData
// @Failure 400 {object} httputils.HTTPError
// @Failure 401
// @Failure 500 {object} httputils.HTTPError
// @Router /api/v1/data [get]
func (h *UserStoredDataHandler) GetUserAll(w http.ResponseWriter, r *http.Request) {
	userID, err := usercontext.GetUserIDFromContext(r.Context())
	if err != nil {
		httperrors.Handle(w, domain.ErrNotAuth)
		return
	}

	dataSet, err := h.service.GetAllUserData(r.Context(), userID)
	if err != nil {
		httperrors.Handle(w, err)
		return
	}

	httputils.SendJSONResponse(w, http.StatusOK, dataSet)
}

// GetOfType godoc
// @Summary Get all user saved data with type
// @Produce json
// @Tags data
// @Security Bearer
// @Param type path string true "Data Type"
// @Success 200 {object} domain.UserStoredData
// @Failure 400 {object} httputils.HTTPError
// @Failure 401
// @Failure 500 {object} httputils.HTTPError
// @Router /api/v1/data/{type} [get]
func (h *UserStoredDataHandler) GetOfType(w http.ResponseWriter, r *http.Request) {
	userID, err := usercontext.GetUserIDFromContext(r.Context())
	if err != nil {
		httperrors.Handle(w, domain.ErrNotAuth)
		return
	}
	dataType := chi.URLParam(r, "type")

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page <= 0 {
		page = 1
	}
	count, err := strconv.Atoi(r.URL.Query().Get("count"))
	if err != nil || count <= 0 || count > 300 {
		count = 50
	}

	paginatedResult, err := h.service.GetUserData(r.Context(), userID, dataType, &domain.StorageFilters{
		IsPaginated:    true,
		IsSortedByDate: true,
		Pagination: domain.PaginationFilters{
			Page:  page,
			Count: count,
		},
		SortDate: domain.SortDateFilters{
			IsASC: false,
		},
	})
	if err != nil {
		httperrors.Handle(w, err)
		return
	}

	httputils.SendJSONResponse(w, http.StatusOK, paginatedResult)
}

// Add godoc
// @Summary Save user data
// @Accept json
// @Produce json
// @Tags data
// @Security Bearer
// @Param type path string true "Data Type"
// @Param dto body dtos.AddNewCardBody true "body"
// @Success 200 {object} domain.UserStoredData
// @Failure 400 {object} httputils.HTTPError
// @Failure 401
// @Failure 500 {object} httputils.HTTPError
// @Router /api/v1/data/{type} [post]
func (h *UserStoredDataHandler) Add(w http.ResponseWriter, r *http.Request) {
	userID, err := usercontext.GetUserIDFromContext(r.Context())
	if err != nil {
		httperrors.Handle(w, domain.ErrNotAuth)
		return
	}

	dataType := chi.URLParam(r, "type")
	dataBody, err := h.parseUserDataBody(w, r, dataType)
	if err != nil {
		httperrors.Handle(w, err)
		return
	}

	if !dataBody.Valid() {
		httperrors.Handle(w, domain.ErrInvalidBody)
		return
	}

	data, err := h.service.Add(r.Context(), userID, dataType, dataBody.GetData(), dataBody.GetMeta())
	if err != nil {
		httperrors.Handle(w, err)
		return
	}

	httputils.SendJSONResponse(w, http.StatusCreated, data)
}

// UpdateOne godoc
// @Summary Update one record with given id
// @Accept json
// @Produce json
// @Tags data
// @Security Bearer
// @Param id path string true "Data Record ID"
// @Param dto body dtos.AddNewCardBody true "body"
// @Success 200 {object} domain.UserStoredData
// @Failure 400 {object} httputils.HTTPError
// @Failure 401
// @Failure 500 {object} httputils.HTTPError
// @Router /api/v1/data/update/{id} [put]
func (h *UserStoredDataHandler) UpdateOne(w http.ResponseWriter, r *http.Request) {
	userID, err := usercontext.GetUserIDFromContext(r.Context())
	if err != nil {
		httperrors.Handle(w, domain.ErrNotAuth)
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		httperrors.Handle(w, domain.ErrUserStoredDataNotFound)
		return
	}

	oldData, err := h.service.GetUserDataByID(r.Context(), userID, id)
	if err != nil {
		httperrors.Handle(w, domain.ErrUserStoredDataNotFound)
		return
	}

	dataBody, err := h.parseUserDataBody(w, r, oldData.DataType)
	if err != nil {
		httperrors.Handle(w, err)
		return
	}

	if !dataBody.Valid() {
		httperrors.Handle(w, domain.ErrInvalidBody)
		return
	}

	updatedUserData, err := h.service.UpdateUserData(r.Context(), userID, id, dataBody.GetData(), dataBody.GetMeta())
	if err != nil {
		httperrors.Handle(w, err)
		return
	}

	httputils.SendJSONResponse(w, http.StatusOK, updatedUserData)
}

// DeleteBatch godoc
// @Summary Save user data
// @Accept json
// @Produce json
// @Tags data
// @Security Bearer
// @Param dto body dtos.DeleteBatchBody true "body"
// @Success 204
// @Failure 400 {object} httputils.HTTPError
// @Failure 401
// @Failure 500 {object} httputils.HTTPError
// @Router /api/v1/data [delete]
func (h *UserStoredDataHandler) DeleteBatch(w http.ResponseWriter, r *http.Request) {
	userID, err := usercontext.GetUserIDFromContext(r.Context())
	if err != nil {
		httperrors.Handle(w, domain.ErrNotAuth)
		return
	}

	var body dtos.DeleteBatchBody
	if statusCode, err := jsonutil.Unmarshal(w, r, &body); err != nil {
		httputils.SendJSONErrorResponse(w, statusCode, err.Error(), statusCode)
		return
	}

	if !body.Valid() {
		httperrors.Handle(w, domain.ErrInvalidBody)
		return
	}

	err = h.service.DeleteBatch(r.Context(), userID, body.IDs)
	if err != nil {
		httperrors.Handle(w, err)
		return
	}

	httputils.SendStatusCode(w, http.StatusNoContent)
}

func (h *UserStoredDataHandler) parseUserDataBody(w http.ResponseWriter, r *http.Request, dataType string) (domain.AddUserStoredDataBody, error) {
	switch dataType {
	case domain.LogPassDataType:
		var body dtos.AddNewLogPassBody
		if _, err := jsonutil.Unmarshal(w, r, &body); err != nil {
			return nil, err
		}

		return &body, nil
	case domain.CardDataType:
		var body dtos.AddNewCardBody
		if _, err := jsonutil.Unmarshal(w, r, &body); err != nil {
			return nil, err
		}

		return &body, nil
	case domain.TextDataType:
		var body dtos.AddNewTextBody
		if _, err := jsonutil.Unmarshal(w, r, &body); err != nil {
			return nil, err
		}

		return &body, nil
	case domain.FileDataType:
		var body dtos.AddNewFileBody
		if _, err := jsonutil.Unmarshal(w, r, &body); err != nil {
			return nil, err
		}

		return &body, nil
	default:
		return nil, domain.ErrInvalidDataType
	}
}
