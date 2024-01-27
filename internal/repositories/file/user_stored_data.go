package file

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"sort"
	"time"

	"github.com/MowlCoder/goph-keeper/internal/domain"
)

type UserStoredDataRepository struct {
	file      *os.File
	structure []domain.UserStoredData

	nextID int
}

func NewUserStoredDataRepository(file *os.File) *UserStoredDataRepository {
	repo := &UserStoredDataRepository{
		file:      file,
		structure: make([]domain.UserStoredData, 0, 100),
	}

	if err := json.NewDecoder(file).Decode(&repo.structure); err != nil {
		if !errors.Is(err, io.EOF) {
			panic(err)
		}
	}

	repo.nextID = len(repo.structure)

	return repo
}

func (repo *UserStoredDataRepository) GetAll(ctx context.Context) ([]domain.UserStoredData, error) {
	dataSet := make([]domain.UserStoredData, len(repo.structure))
	copy(dataSet, repo.structure)

	return dataSet, nil
}

func (repo *UserStoredDataRepository) GetWithType(ctx context.Context, dataType string, filters *domain.StorageFilters) ([]domain.UserStoredData, error) {
	dataSet := make([]domain.UserStoredData, 0)

	for _, data := range repo.structure {
		if data.DataType == dataType {
			dataSet = append(dataSet, data)
		}
	}

	if filters.IsSortedByDate {
		sort.Slice(dataSet, func(i, j int) bool {
			if filters.SortDate.IsASC {
				return dataSet[i].CreatedAt.Before(dataSet[j].CreatedAt)
			}

			return dataSet[i].CreatedAt.After(dataSet[j].CreatedAt)
		})
	}

	if filters.IsPaginated {
		startFrom := (filters.Pagination.Page - 1) * filters.Pagination.Count
		endAt := startFrom + filters.Pagination.Count

		if startFrom >= len(dataSet) {
			return []domain.UserStoredData{}, nil
		}

		if endAt > len(dataSet) {
			return dataSet[startFrom:], nil
		}

		return dataSet[startFrom : startFrom+filters.Pagination.Count], nil
	}

	return dataSet, nil
}

func (repo *UserStoredDataRepository) CountUserDataOfType(ctx context.Context, dataType string) (int, error) {
	count := 0

	for _, data := range repo.structure {
		if data.DataType == dataType {
			count += 1
		}
	}

	return count, nil
}

func (repo *UserStoredDataRepository) AddData(ctx context.Context, dataType string, data []byte, meta string) (int64, error) {
	userStoredData := domain.UserStoredData{
		ID:          repo.getNextID() * -1,
		DataType:    dataType,
		CryptedData: data,
		Meta:        meta,
		CreatedAt:   time.Now().UTC(),
		Version:     -1,
	}
	repo.structure = append(repo.structure, userStoredData)

	if err := repo.SaveInFile(); err != nil {
		return 0, err
	}

	return int64(userStoredData.ID), nil
}

func (repo *UserStoredDataRepository) UpdateByID(ctx context.Context, id int, data []byte, meta string) (*domain.UserStoredData, error) {
	foundIdx := 0
	isFound := false

	for idx, pair := range repo.structure {
		if pair.ID == id {
			isFound = true
			foundIdx = idx
			break
		}
	}

	if !isFound {
		return nil, domain.ErrNotFound
	}

	repo.structure[foundIdx].CryptedData = data
	repo.structure[foundIdx].Meta = meta
	if err := repo.SaveInFile(); err != nil {
		return nil, err
	}

	updatedData := repo.structure[foundIdx]

	return &updatedData, nil
}

func (repo *UserStoredDataRepository) DeleteByID(ctx context.Context, id int) error {
	foundIdx := 0
	isFound := false

	for idx, pair := range repo.structure {
		if pair.ID == id {
			isFound = true
			foundIdx = idx
			break
		}
	}

	if !isFound {
		return domain.ErrNotFound
	}

	repo.structure = append(repo.structure[:foundIdx], repo.structure[foundIdx+1:]...)

	return repo.SaveInFile()
}

func (repo *UserStoredDataRepository) DeleteBatch(ctx context.Context, ids []int) error {
	filtered := make([]domain.UserStoredData, 0, len(repo.structure))

	for _, data := range repo.structure {
		isFound := false

		for _, id := range ids {
			if data.ID == id {
				isFound = true
				break
			}
		}

		if !isFound {
			filtered = append(filtered, data)
		}
	}

	repo.structure = filtered

	return repo.SaveInFile()
}

func (repo *UserStoredDataRepository) SyncUpdate(ctx context.Context, oldID int, newID int, version int) error {
	for idx, data := range repo.structure {
		if data.ID == oldID {
			repo.structure[idx].ID = newID
			repo.structure[idx].Version = version
		}
	}

	return repo.SaveInFile()
}

func (repo *UserStoredDataRepository) SaveInFile() error {
	if err := repo.file.Truncate(0); err != nil {
		return err
	}
	if _, err := repo.file.Seek(0, 0); err != nil {
		return err
	}

	writer := json.NewEncoder(repo.file)
	return writer.Encode(repo.structure)
}

func (repo *UserStoredDataRepository) getNextID() int {
	repo.nextID += 1
	return repo.nextID
}
