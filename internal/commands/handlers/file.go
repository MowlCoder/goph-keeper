package handlers

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/MowlCoder/goph-keeper/internal/domain"
	"github.com/MowlCoder/goph-keeper/internal/session"
	"github.com/MowlCoder/goph-keeper/pkg/input"
)

type FileHandler struct {
	clientSession *session.ClientSession

	userStoredDataService userStoredDataService
}

func NewFileHandler(
	clientSession *session.ClientSession,
	userStoredDataService userStoredDataService,
) *FileHandler {
	return &FileHandler{
		clientSession:         clientSession,
		userStoredDataService: userStoredDataService,
	}
}

func (h *FileHandler) AddFile(args []string) error {
	filePath := input.GetConsoleInput("Enter file path: ", "")
	if filePath == "" {
		return domain.ErrInvalidInputValue
	}

	meta := input.GetConsoleInput("Enter meta information: ", "")

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	fileData := domain.FileData{
		Content: fileContent,
		Name:    filepath.Base(filePath),
	}

	_, err = h.userStoredDataService.Add(
		context.Background(),
		domain.FileDataType,
		fileData,
		meta,
	)
	if err != nil {
		return err
	}

	fmt.Println("Successfully saved new file!")

	return nil
}

func (h *FileHandler) GetFiles(args []string) error {
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
		domain.FileDataType,
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

	fmt.Println("================== Files ==================")

	for _, data := range paginatedResult.Data.([]domain.UserStoredData) {
		fileData := data.Data.(domain.FileData)
		fmt.Println(fmt.Sprintf("ID: %d | %s | Meta: %s (version %d)", data.ID, fileData.Name, data.Meta, data.Version))
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

func (h *FileHandler) DecryptFile(args []string) error {
	if len(args) != 1 {
		return domain.ErrInvalidCommandUsage
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		return domain.ErrInvalidCommandUsage
	}

	dirPath := input.GetConsoleInput("Enter directory where decrypt file: ", "")
	if dirPath == "" {
		return domain.ErrInvalidInputValue
	}

	if err := os.Mkdir(dirPath, os.ModePerm); err != nil && !errors.Is(err, os.ErrExist) {
		return err
	}

	userData, err := h.userStoredDataService.GetByID(context.Background(), id)
	if err != nil {
		return err
	}

	parsedData := userData.Data.(domain.FileData)
	pathToFile := filepath.Join(dirPath, parsedData.Name)
	if err := os.WriteFile(pathToFile, parsedData.Content, os.ModePerm); err != nil {
		return err
	}

	fmt.Println("Successfully decrypted file")

	return nil
}

func (h *FileHandler) UpdateFile(args []string) error {
	if len(args) != 1 {
		return domain.ErrInvalidCommandUsage
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		return domain.ErrInvalidCommandUsage
	}

	filePath := input.GetConsoleInput("Enter file path: ", "")
	if filePath == "" {
		return domain.ErrInvalidInputValue
	}

	meta := input.GetConsoleInput("Enter meta information: ", "")

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	fileData := domain.FileData{
		Content: fileContent,
		Name:    filepath.Base(filePath),
	}

	_, err = h.userStoredDataService.UpdateByID(
		context.Background(),
		id,
		fileData,
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

	fmt.Println("Successfully updated file data")

	return nil
}

func (h *FileHandler) DeleteFile(args []string) error {
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

	fmt.Printf("Successfully delete file with id %d\n", id)

	return nil
}
