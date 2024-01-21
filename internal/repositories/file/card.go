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

type CardRepository struct {
	file      *os.File
	structure []domain.Card
}

func NewCardRepository(file *os.File) *CardRepository {
	repo := &CardRepository{
		file:      file,
		structure: make([]domain.Card, 0, 100),
	}

	if err := json.NewDecoder(file).Decode(&repo.structure); err != nil {
		if !errors.Is(err, io.EOF) {
			panic(err)
		}
	}

	return repo
}

func (r *CardRepository) GetByUserID(ctx context.Context, userID int, filters *domain.StorageFilters) ([]domain.Card, error) {
	cards := make([]domain.Card, 0)

	for _, card := range r.structure {
		if card.UserID == userID {
			cards = append(cards, card)
		}
	}

	if filters.IsSortedByDate {
		sort.Slice(cards, func(i, j int) bool {
			if filters.SortDate.IsASC {
				return cards[i].CreatedAt.Before(cards[j].CreatedAt)
			}

			return cards[i].CreatedAt.After(cards[j].CreatedAt)
		})
	}

	if filters.IsPaginated {
		startFrom := (filters.Pagination.Page - 1) * filters.Pagination.Count
		endAt := startFrom + filters.Pagination.Count

		if startFrom >= len(cards) {
			return []domain.Card{}, nil
		}

		if endAt > len(cards) {
			return cards[startFrom:], nil
		}

		return cards[startFrom : startFrom+filters.Pagination.Count], nil
	}

	return cards, nil
}

func (r *CardRepository) CountUserCards(ctx context.Context, userID int) (int, error) {
	count := 0

	for _, card := range r.structure {
		if card.UserID == userID {
			count += 1
		}
	}

	return count, nil
}

func (r *CardRepository) DeleteByID(ctx context.Context, userID int, id int) error {
	foundIdx := 0
	isFound := false

	for idx, card := range r.structure {
		if card.ID == id && card.UserID == userID {
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

func (r *CardRepository) DeleteBatch(ctx context.Context, userID int, ids []int) error {
	filtered := make([]domain.Card, 0, len(r.structure))

	for _, card := range r.structure {
		isFound := false

		for _, id := range ids {
			if card.ID == id && card.UserID == userID {
				isFound = true
				break
			}
		}

		if !isFound {
			filtered = append(filtered, card)
		}
	}

	r.structure = filtered

	return r.SaveInFile()
}

func (r *CardRepository) AddCard(
	ctx context.Context,
	userID int,
	number string,
	expiredAt string,
	cvv string,
	meta string,
) (*domain.Card, error) {
	card := domain.Card{
		ID:        (len(r.structure) + 1) * -1,
		UserID:    userID,
		Number:    number,
		ExpiredAt: expiredAt,
		CVV:       cvv,
		Meta:      meta,
		CreatedAt: time.Now().UTC(),
		Version:   -1,
	}
	r.structure = append(r.structure, card)

	if err := r.SaveInFile(); err != nil {
		return nil, err
	}

	return &card, nil
}

func (r *CardRepository) SyncUpdate(ctx context.Context, oldID int, newID int, version int) error {
	for idx, pair := range r.structure {
		if pair.ID == oldID {
			r.structure[idx].ID = newID
			r.structure[idx].Version = version
		}
	}

	return r.SaveInFile()
}

func (r *CardRepository) SaveInFile() error {
	if err := r.file.Truncate(0); err != nil {
		return err
	}
	if _, err := r.file.Seek(0, 0); err != nil {
		return err
	}

	writer := json.NewEncoder(r.file)
	return writer.Encode(r.structure)
}
