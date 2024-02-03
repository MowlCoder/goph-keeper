package handlers

import (
	"context"
	"fmt"
	"strconv"

	"github.com/MowlCoder/goph-keeper/internal/domain"
	"github.com/MowlCoder/goph-keeper/internal/session"
	"github.com/MowlCoder/goph-keeper/internal/validators"
	"github.com/MowlCoder/goph-keeper/pkg/input"
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
	cardNumber, _ := input.GetConsoleInput("Enter card number: ", "")
	if !validators.ValidateCardNumber(cardNumber) {
		return domain.ErrInvalidCardNumber
	}

	expiredAt, _ := input.GetConsoleInput("Enter expired date (e.g. 04/30): ", "")
	if !validators.ValidateExpiredAt(expiredAt) {
		return domain.ErrInvalidCardExpiredAt
	}

	cvv, _ := input.GetConsoleInput("Enter card cvv: ", "")
	if !validators.ValidateCVV(cvv) {
		return domain.ErrInvalidCardCVV
	}

	meta, _ := input.GetConsoleInput("Enter meta information: ", "")

	_, err := h.userStoredDataService.Add(
		context.Background(),
		domain.CardDataType,
		domain.CardData{
			Number:    cardNumber,
			ExpiredAt: expiredAt,
			CVV:       cvv,
		},
		meta,
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
		if err := h.clientSession.AddDeleted(id); err != nil {
			return err
		}
	}

	fmt.Printf("Successfully delete card with id %d\n", id)

	return nil
}

func (h *CardHandler) UpdateCard(args []string) error {
	if len(args) != 1 {
		return domain.ErrInvalidCommandUsage
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		return domain.ErrInvalidCommandUsage
	}

	userStoredData, err := h.userStoredDataService.GetByID(context.Background(), id)
	if err != nil {
		return err
	}

	data := userStoredData.Data.(domain.CardData)

	cardNumber, _ := input.GetConsoleInput(fmt.Sprintf("Enter card number (current - %s): ", data.Number), data.Number)
	if !validators.ValidateCardNumber(cardNumber) {
		return domain.ErrInvalidCardNumber
	}

	expiredAt, _ := input.GetConsoleInput(fmt.Sprintf("Enter expired date (current - %s): ", data.ExpiredAt), data.ExpiredAt)
	if !validators.ValidateExpiredAt(expiredAt) {
		return domain.ErrInvalidCardExpiredAt
	}

	cvv, _ := input.GetConsoleInput(fmt.Sprintf("Enter card cvv (current - %s): ", data.CVV), data.CVV)
	if !validators.ValidateCVV(cvv) {
		return domain.ErrInvalidCardCVV
	}

	meta, _ := input.GetConsoleInput(fmt.Sprintf("Enter meta information (current - %s): ", userStoredData.Meta), userStoredData.Meta)

	_, err = h.userStoredDataService.UpdateByID(
		context.Background(),
		id,
		domain.CardData{
			Number:    cardNumber,
			ExpiredAt: expiredAt,
			CVV:       cvv,
		},
		meta,
	)
	if err != nil {
		return err
	}

	if id >= 0 {
		if err := h.clientSession.AddEdited(id); err != nil {
			return err
		}
	}

	fmt.Println("Successfully update card data")

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
