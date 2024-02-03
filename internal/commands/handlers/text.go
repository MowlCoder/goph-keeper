package handlers

import (
	"context"
	"fmt"
	"strconv"

	"github.com/MowlCoder/goph-keeper/internal/domain"
	"github.com/MowlCoder/goph-keeper/internal/session"
	"github.com/MowlCoder/goph-keeper/pkg/input"
)

type TextHandler struct {
	clientSession *session.ClientSession

	userStoredDataService userStoredDataService
}

func NewTextHandler(
	clientSession *session.ClientSession,
	userStoredDataService userStoredDataService,
) *TextHandler {
	return &TextHandler{
		clientSession:         clientSession,
		userStoredDataService: userStoredDataService,
	}
}

func (h *TextHandler) AddText(args []string) error {
	title, _ := input.GetConsoleInput("Enter title: ", "")
	if title == "" {
		return domain.ErrInvalidInputValue
	}

	text, _ := input.GetConsoleInput("Enter text: ", "")
	if text == "" {
		return domain.ErrInvalidInputValue
	}

	_, err := h.userStoredDataService.Add(
		context.Background(),
		domain.TextDataType,
		domain.TextData{
			Text: text,
		},
		title,
	)
	if err != nil {
		return err
	}

	fmt.Println("Successfully saved new text!")

	return nil
}

func (h *TextHandler) UpdateText(args []string) error {
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

	data := userStoredData.Data.(domain.TextData)

	title, _ := input.GetConsoleInput(fmt.Sprintf("Enter title (current - %s): ", userStoredData.Meta), userStoredData.Meta)
	if title == "" {
		return domain.ErrInvalidInputValue
	}

	text, _ := input.GetConsoleInput(fmt.Sprintf("Enter text (current - %s): ", data.Text), data.Text)
	if text == "" {
		return domain.ErrInvalidInputValue
	}

	_, err = h.userStoredDataService.UpdateByID(
		context.Background(),
		id,
		domain.TextData{
			Text: text,
		},
		title,
	)
	if err != nil {
		return err
	}

	if id >= 0 {
		if err := h.clientSession.AddEdited(id); err != nil {
			return err
		}
	}

	fmt.Println("Successfully update text data")

	return nil
}

func (h *TextHandler) DeleteText(args []string) error {
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

func (h *TextHandler) GetTexts(args []string) error {
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
		domain.TextDataType,
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

	fmt.Println("================== Text ==================")

	for _, data := range paginatedResult.Data.([]domain.UserStoredData) {
		textData := data.Data.(domain.TextData)
		fmt.Println(fmt.Sprintf("ID: %d | %s | Meta: %s (version %d)", data.ID, textData.Text, data.Meta, data.Version))
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
