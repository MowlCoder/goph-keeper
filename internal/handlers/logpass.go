package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/MowlCoder/goph-keeper/internal/domain"
	"github.com/MowlCoder/goph-keeper/internal/dtos"
	"github.com/MowlCoder/goph-keeper/internal/handlers/httperrors"
	"github.com/MowlCoder/goph-keeper/internal/utils/usercontext"
	"github.com/MowlCoder/goph-keeper/pkg/httputils"
	jsonutil "github.com/MowlCoder/goph-keeper/pkg/jsonutils"
)

type logPassService interface {
	GetUserPairs(ctx context.Context, userID int, filters *domain.StorageFilters, needToDecrypt bool) (*domain.PaginatedResult, error)
	AddNewPair(ctx context.Context, needEncrypt bool, userID int, login string, password string, source string) (*domain.LogPass, error)
	DeleteBatchPairs(ctx context.Context, userID int, ids []int) error
}

type LogPassHandler struct {
	logPassService logPassService
}

func NewLogPassHandler(logPassService logPassService) *LogPassHandler {
	return &LogPassHandler{
		logPassService: logPassService,
	}
}

func (h *LogPassHandler) GetMyPairs(w http.ResponseWriter, r *http.Request) {
	userID, err := usercontext.GetUserIDFromContext(r.Context())
	if err != nil {
		httperrors.Handle(w, domain.ErrNotAuth)
		return
	}

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page <= 0 {
		page = 1
	}
	count, err := strconv.Atoi(r.URL.Query().Get("count"))
	if err != nil || count <= 0 || count > 300 {
		count = 50
	}

	paginatedResult, err := h.logPassService.GetUserPairs(r.Context(), userID, &domain.StorageFilters{
		IsPaginated:    true,
		IsSortedByDate: true,
		Pagination: domain.PaginationFilters{
			Page:  page,
			Count: count,
		},
		SortDate: domain.SortDateFilters{
			IsASC: false,
		},
	}, true)
	if err != nil {
		httperrors.Handle(w, err)
		return
	}

	httputils.SendJSONResponse(w, http.StatusOK, paginatedResult)
}

func (h *LogPassHandler) DeleteBatchPairs(w http.ResponseWriter, r *http.Request) {
	userID, err := usercontext.GetUserIDFromContext(r.Context())
	if err != nil {
		httperrors.Handle(w, domain.ErrNotAuth)
		return
	}

	var body dtos.DeleteBatchPairsBody
	if statusCode, err := jsonutil.Unmarshal(w, r, &body); err != nil {
		httputils.SendJSONErrorResponse(w, statusCode, err.Error(), statusCode)
		return
	}

	if !body.Valid() {
		httperrors.Handle(w, domain.ErrInvalidBody)
		return
	}

	err = h.logPassService.DeleteBatchPairs(r.Context(), userID, body.IDs)
	if err != nil {
		httperrors.Handle(w, err)
		return
	}

	httputils.SendStatusCode(w, http.StatusNoContent)
}

func (h *LogPassHandler) AddNewPair(w http.ResponseWriter, r *http.Request) {
	userID, err := usercontext.GetUserIDFromContext(r.Context())
	if err != nil {
		httperrors.Handle(w, domain.ErrNotAuth)
		return
	}

	var body dtos.AddNewLogPassBody
	if statusCode, err := jsonutil.Unmarshal(w, r, &body); err != nil {
		httputils.SendJSONErrorResponse(w, statusCode, err.Error(), statusCode)
		return
	}

	if !body.Valid() {
		httperrors.Handle(w, domain.ErrInvalidBody)
		return
	}

	logPass, err := h.logPassService.AddNewPair(r.Context(), true, userID, body.Login, body.Password, body.Source)
	if err != nil {
		httperrors.Handle(w, err)
		return
	}

	httputils.SendJSONResponse(w, http.StatusCreated, logPass)
}
