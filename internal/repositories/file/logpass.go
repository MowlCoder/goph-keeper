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

type LogPassRepository struct {
	file      *os.File
	structure []domain.LogPass
}

func NewLogPassRepository(file *os.File) *LogPassRepository {
	repo := &LogPassRepository{
		file:      file,
		structure: make([]domain.LogPass, 0, 100),
	}

	if err := json.NewDecoder(file).Decode(&repo.structure); err != nil {
		if !errors.Is(err, io.EOF) {
			panic(err)
		}
	}

	return repo
}

func (r *LogPassRepository) GetPairsByUserID(ctx context.Context, userID int, filters *domain.StorageFilters) ([]domain.LogPass, error) {
	pairs := make([]domain.LogPass, 0)

	for _, logPass := range r.structure {
		if logPass.UserID == userID {
			pairs = append(pairs, logPass)
		}
	}

	if filters.IsSortedByDate {
		sort.Slice(pairs, func(i, j int) bool {
			if filters.SortDate.IsASC {
				return pairs[i].CreatedAt.Before(pairs[j].CreatedAt)
			}

			return pairs[i].CreatedAt.After(pairs[j].CreatedAt)
		})
	}

	if filters.IsPaginated {
		startFrom := (filters.Pagination.Page - 1) * filters.Pagination.Count
		endAt := startFrom + filters.Pagination.Count

		if startFrom >= len(pairs) {
			return []domain.LogPass{}, nil
		}

		if endAt > len(pairs) {
			return pairs[startFrom:], nil
		}

		return pairs[startFrom : startFrom+filters.Pagination.Count], nil
	}

	return pairs, nil
}

func (r *LogPassRepository) CountUserPairs(ctx context.Context, userID int) (int, error) {
	count := 0

	for _, logPass := range r.structure {
		if logPass.UserID == userID {
			count += 1
		}
	}

	return count, nil
}

func (r *LogPassRepository) DeleteByID(ctx context.Context, userID int, id int) error {
	foundIdx := 0
	isFound := false

	for idx, pair := range r.structure {
		if pair.ID == id && pair.UserID == userID {
			isFound = true
			foundIdx = idx
			break
		}
	}

	if !isFound {
		return domain.ErrNotFound
	}

	r.structure = append(r.structure[:foundIdx], r.structure[foundIdx+1:]...)

	return r.SaveInFile()
}

func (r *LogPassRepository) DeleteBatch(ctx context.Context, userID int, ids []int) error {
	filtered := make([]domain.LogPass, 0, len(r.structure))

	for _, pair := range r.structure {
		isFound := false

		for _, id := range ids {
			if pair.ID == id && pair.UserID == userID {
				isFound = true
				break
			}
		}

		if !isFound {
			filtered = append(filtered, pair)
		}
	}

	r.structure = filtered

	return r.SaveInFile()
}

func (r *LogPassRepository) AddPair(
	ctx context.Context,
	userID int,
	login string,
	password string,
	source string,
) (*domain.LogPass, error) {
	logPass := domain.LogPass{
		ID:        (len(r.structure) + 1) * -1,
		UserID:    userID,
		Login:     login,
		Password:  password,
		Source:    source,
		CreatedAt: time.Now().UTC(),
		Version:   -1,
	}
	r.structure = append(r.structure, logPass)

	if err := r.SaveInFile(); err != nil {
		return nil, err
	}

	return &logPass, nil
}

func (r *LogPassRepository) SyncUpdate(ctx context.Context, oldID int, newID int, version int) error {
	for idx, pair := range r.structure {
		if pair.ID == oldID {
			r.structure[idx].ID = newID
			r.structure[idx].Version = version
		}
	}

	return r.SaveInFile()
}

func (r *LogPassRepository) SaveInFile() error {
	if err := r.file.Truncate(0); err != nil {
		return err
	}
	if _, err := r.file.Seek(0, 0); err != nil {
		return err
	}

	writer := json.NewEncoder(r.file)
	return writer.Encode(r.structure)
}
