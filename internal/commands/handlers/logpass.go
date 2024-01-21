package handlers

import (
	"context"
	"fmt"
	"strconv"

	"github.com/MowlCoder/goph-keeper/internal/domain"
	"github.com/MowlCoder/goph-keeper/internal/session"
)

type logPassService interface {
	GetUserPairs(ctx context.Context, userID int, filters *domain.StorageFilters, needToDecrypt bool) (*domain.PaginatedResult, error)
	AddNewPair(ctx context.Context, needEncrypt bool, userID int, login string, password string, source string) (*domain.LogPass, error)
	DeletePairByID(ctx context.Context, userID int, id int) error
}

type LogPassHandler struct {
	clientSession *session.ClientSession

	logPassService logPassService
}

func NewLogPassHandler(
	clientSession *session.ClientSession,
	logPassService logPassService,
) *LogPassHandler {
	return &LogPassHandler{
		clientSession:  clientSession,
		logPassService: logPassService,
	}
}

func (h *LogPassHandler) AddPair(args []string) error {
	if len(args) != 3 {
		return domain.ErrInvalidCommandUsage
	}

	_, err := h.logPassService.AddNewPair(
		context.Background(),
		true,
		domain.LocalUserID,
		args[0],
		args[1],
		args[2],
	)
	if err != nil {
		return err
	}

	fmt.Println("Successfully saved new log:pass pair!")

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

	err = h.logPassService.DeletePairByID(
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

	paginatedResult, err := h.logPassService.GetUserPairs(
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
		true,
	)
	if err != nil {
		return err
	}

	fmt.Println("================== Log Pass ==================")

	for _, pair := range paginatedResult.Data.([]domain.LogPass) {
		fmt.Println(fmt.Sprintf("ID: %d | %s:%s | Source: %s (version %d)", pair.ID, pair.Login, pair.Password, pair.Source, pair.Version))
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
