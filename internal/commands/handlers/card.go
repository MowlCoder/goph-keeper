package handlers

import (
	"context"
	"fmt"
	"strconv"

	"github.com/MowlCoder/goph-keeper/internal/domain"
	"github.com/MowlCoder/goph-keeper/internal/session"
	"github.com/MowlCoder/goph-keeper/internal/validators"
)

type CardHandler struct {
	clientSession *session.ClientSession

	userStoredDataService userStoredDataService
}

func NewCardHandler(
	clientSession *session.ClientSession,
	userStoredDataService userStoredDataService,
) *CardHandler {
	return &CardHandler{
		clientSession:         clientSession,
		userStoredDataService: userStoredDataService,
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

	_, err := h.userStoredDataService.Add(
		context.Background(),
		domain.CardDataType,
		domain.CardData{
			Number:    args[0],
			ExpiredAt: args[1],
			CVV:       args[2],
		},
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

	err = h.userStoredDataService.DeleteByID(
		context.Background(),
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

	paginatedResult, err := h.userStoredDataService.GetUserData(
		context.Background(),
		domain.CardDataType,
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

	for _, data := range paginatedResult.Data.([]domain.UserStoredData) {
		cardData := data.Data.(domain.CardData)
		fmt.Println(fmt.Sprintf("ID: %d | %s %s %s | Meta: %s (version %d)", data.ID, cardData.Number, cardData.ExpiredAt, cardData.CVV, data.Meta, data.Version))
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
