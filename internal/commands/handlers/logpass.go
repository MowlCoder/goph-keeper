package handlers

import (
	"context"
	"fmt"
	"strconv"

	"github.com/MowlCoder/goph-keeper/internal/domain"
	"github.com/MowlCoder/goph-keeper/internal/session"
	"github.com/MowlCoder/goph-keeper/pkg/input"
)

type LogPassHandler struct {
	clientSession         *session.ClientSession
	userStoredDataService userStoredDataService
}

func NewLogPassHandler(
	clientSession *session.ClientSession,
	userStoredDataService userStoredDataService,
) *LogPassHandler {
	return &LogPassHandler{
		clientSession:         clientSession,
		userStoredDataService: userStoredDataService,
	}
}

func (h *LogPassHandler) AddPair(args []string) error {
	login, _ := input.GetConsoleInput("Enter login: ", "")
	if login == "" {
		return domain.ErrInvalidInputValue
	}

	password, _ := input.GetConsoleInput("Enter password: ", "")
	if password == "" {
		return domain.ErrInvalidInputValue
	}

	meta, _ := input.GetConsoleInput("Enter source: ", "")
	if meta == "" {
		return domain.ErrInvalidInputValue
	}

	_, err := h.userStoredDataService.Add(
		context.Background(),
		domain.LogPassDataType,
		&domain.LogPassData{
			Login:    login,
			Password: password,
		},
		meta,
	)
	if err != nil {
		return err
	}

	fmt.Println("Successfully saved new log:pass pair!")

	return nil
}

func (h *LogPassHandler) UpdatePair(args []string) error {
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

	data := userStoredData.Data.(domain.LogPassData)

	login, _ := input.GetConsoleInput(fmt.Sprintf("Enter login (current - %s): ", data.Login), data.Login)
	if login == "" {
		return domain.ErrInvalidInputValue
	}

	password, _ := input.GetConsoleInput(fmt.Sprintf("Enter password (current - %s): ", data.Password), data.Password)
	if password == "" {
		return domain.ErrInvalidInputValue
	}

	meta, _ := input.GetConsoleInput(fmt.Sprintf("Enter source (current - %s): ", userStoredData.Meta), userStoredData.Meta)
	if meta == "" {
		return domain.ErrInvalidInputValue
	}

	_, err = h.userStoredDataService.UpdateByID(
		context.Background(),
		id,
		domain.LogPassData{
			Login:    login,
			Password: password,
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

	fmt.Println("Successfully update logpass data")

	return nil
}

func (h *LogPassHandler) DeletePair(args []string) error {
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

	fmt.Printf("Successfully delete log:pass pair with id %d\n", id)

	return nil
}

func (h *LogPassHandler) GetPairs(args []string) error {
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
		domain.LogPassDataType,
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

	fmt.Println("================== Log Pass ==================")

	for _, data := range paginatedResult.Data.([]domain.UserStoredData) {
		logPassData := data.Data.(domain.LogPassData)
		fmt.Println(fmt.Sprintf("ID: %d | %s:%s | Source: %s (version %d)", data.ID, logPassData.Login, logPassData.Password, data.Meta, data.Version))
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
