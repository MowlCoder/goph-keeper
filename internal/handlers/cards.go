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

type cardService interface {
	GetUserCards(ctx context.Context, userID int, filters *domain.StorageFilters) (*domain.PaginatedResult, error)
	AddNewCard(ctx context.Context, needEncrypt bool, userID int, number string, expiredAt string, cvv string, meta string) (*domain.Card, error)
	DeleteBatchCards(ctx context.Context, userID int, ids []int) error
}

type CardHandler struct {
	cardService cardService
}

func NewCardHandler(cardService cardService) *CardHandler {
	return &CardHandler{
		cardService: cardService,
	}
}

func (h *CardHandler) GetMyCards(w http.ResponseWriter, r *http.Request) {
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

	paginatedResult, err := h.cardService.GetUserCards(r.Context(), userID, &domain.StorageFilters{
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

func (h *CardHandler) DeleteBatchCards(w http.ResponseWriter, r *http.Request) {
	userID, err := usercontext.GetUserIDFromContext(r.Context())
	if err != nil {
		httperrors.Handle(w, domain.ErrNotAuth)
		return
	}

	var body dtos.DeleteBatchCardsBody
	if statusCode, err := jsonutil.Unmarshal(w, r, &body); err != nil {
		httputils.SendJSONErrorResponse(w, statusCode, err.Error(), statusCode)
		return
	}

	if !body.Valid() {
		httperrors.Handle(w, domain.ErrInvalidBody)
		return
	}

	err = h.cardService.DeleteBatchCards(r.Context(), userID, body.IDs)
	if err != nil {
		httperrors.Handle(w, err)
		return
	}

	httputils.SendStatusCode(w, http.StatusNoContent)
}

func (h *CardHandler) AddNewCard(w http.ResponseWriter, r *http.Request) {
	userID, err := usercontext.GetUserIDFromContext(r.Context())
	if err != nil {
		httperrors.Handle(w, domain.ErrNotAuth)
		return
	}

	var body dtos.AddNewCardBody
	if statusCode, err := jsonutil.Unmarshal(w, r, &body); err != nil {
		httputils.SendJSONErrorResponse(w, statusCode, err.Error(), statusCode)
		return
	}

	if !body.Valid() {
		httperrors.Handle(w, domain.ErrInvalidBody)
		return
	}

	card, err := h.cardService.AddNewCard(r.Context(), true, userID, body.Number, body.ExpiredAt, body.CVV, body.Meta)
	if err != nil {
		httperrors.Handle(w, err)
		return
	}

	httputils.SendJSONResponse(w, http.StatusCreated, card)
}
