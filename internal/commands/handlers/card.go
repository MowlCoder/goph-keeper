package handlers

import (
	"context"
	"fmt"
	"strconv"

	"github.com/MowlCoder/goph-keeper/internal/domain"
	"github.com/MowlCoder/goph-keeper/internal/session"
	"github.com/MowlCoder/goph-keeper/internal/validators"
)

type cardService interface {
	GetUserCards(ctx context.Context, userID int, filters *domain.StorageFilters) (*domain.PaginatedResult, error)
	AddNewCard(ctx context.Context, needEncrypt bool, userID int, number string, expiredAt string, cvv string, meta string) (*domain.Card, error)
	DeleteCardByID(ctx context.Context, userID int, id int) error
}

type CardHandler struct {
	clientSession *session.ClientSession

	cardService cardService
}

func NewCardHandler(
	clientSession *session.ClientSession,
	cardService cardService,
) *CardHandler {
	return &CardHandler{
		clientSession: clientSession,
		cardService:   cardService,
	}
}

func (h *CardHandler) AddCard(args []string) error {
	if len(args) < 4 {
		return domain.ErrInvalidCommandUsage
	}

	if !validators.ValidateCardNumber(args[0]) {
		return domain.ErrInvalidCardNumber
	}

	if !validators.ValidateExpiredAt(args[1]) {
		return domain.ErrInvalidCardExpiredAt
	}

	if !validators.ValidateCVV(args[2]) {
		return domain.ErrInvalidCardCVV
	}

	_, err := h.cardService.AddNewCard(
		context.Background(),
		true,
		domain.LocalUserID,
		args[0],
		args[1],
		args[2],
		args[3],
	)
	if err != nil {
		return err
	}

	fmt.Println("Successfully saved new card!")

	return nil
}

func (h *CardHandler) DeleteCard(args []string) error {
	if len(args) != 1 {
		return domain.ErrInvalidCommandUsage
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		return domain.ErrInvalidCommandUsage
	}

	err = h.cardService.DeleteCardByID(
		context.Background(),
		domain.LocalUserID,
		id,
	)
	if err != nil {
		return err
	}

	if id >= 0 {
		h.clientSession.AddDeletedLogPassID(id)
	}

	fmt.Printf("Successfully delete card with id %d\n", id)

	return nil
}

func (h *CardHandler) GetCards(args []string) error {
	count := 15
	var page int
	var err error

	if len(args) == 2 {
		page, err = strconv.Atoi(args[0])
		if err != nil || page <= 0 {
			page = 1
		}
	} else {
		page = 1
	}

	paginatedResult, err := h.cardService.GetUserCards(
		context.Background(),
		domain.LocalUserID,
		&domain.StorageFilters{
			IsPaginated:    true,
			IsSortedByDate: true,
			Pagination: domain.PaginationFilters{
				Page:  page,
				Count: count,
			},
			SortDate: domain.SortDateFilters{
				IsASC: false,
			},
		},
	)
	if err != nil {
		return err
	}

	fmt.Println("================== Cards ==================")

	for _, card := range paginatedResult.Data.([]domain.Card) {
		fmt.Println(fmt.Sprintf("ID: %d | %s %s %s | Meta: %s (version %d)", card.ID, card.Number, card.ExpiredAt, card.CVV, card.Meta, card.Version))
	}

	fmt.Println(
		fmt.Sprintf(
			"================== (%d/%d) ==================",
			paginatedResult.CurrentPage,
			paginatedResult.PageCount,
		),
	)

	return nil
}
