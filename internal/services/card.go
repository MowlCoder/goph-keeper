package services

import (
	"context"
	"math"

	"github.com/MowlCoder/goph-keeper/internal/domain"
)

type cryptorForCard interface {
	Encrypt(raw string) (string, error)
	Decrypt(crypted string) (string, error)
}

type cardRepository interface {
	GetByUserID(ctx context.Context, userID int, filters *domain.StorageFilters) ([]domain.Card, error)
	CountUserCards(ctx context.Context, userID int) (int, error)
	AddCard(ctx context.Context, userID int, number string, expiredAt string, cvv string, meta string) (*domain.Card, error)
	DeleteByID(ctx context.Context, userID int, id int) error
	DeleteBatch(ctx context.Context, userID int, id []int) error
}

type CardService struct {
	cardRepository cardRepository
	cryptor        cryptorForCard
}

func NewCardService(
	cardRepository cardRepository,
	cryptor cryptorForCard,
) *CardService {
	return &CardService{
		cardRepository: cardRepository,
		cryptor:        cryptor,
	}
}

func (s *CardService) GetUserCards(ctx context.Context, userID int, filters *domain.StorageFilters) (*domain.PaginatedResult, error) {
	cards, err := s.cardRepository.GetByUserID(ctx, userID, filters)
	if err != nil {
		return nil, err
	}
	cardsCount, err := s.cardRepository.CountUserCards(ctx, userID)
	if err != nil {
		return nil, err
	}

	for i := range cards {
		cards[i].CVV, err = s.cryptor.Decrypt(cards[i].CVV)
		if err != nil {
			return nil, err
		}

		cards[i].Number, err = s.cryptor.Decrypt(cards[i].Number)
		if err != nil {
			return nil, err
		}
	}

	return &domain.PaginatedResult{
		Count:       filters.Pagination.Count,
		CurrentPage: filters.Pagination.Page,
		PageCount:   int(math.Ceil(float64(cardsCount) / float64(filters.Pagination.Count))),
		Data:        cards,
	}, nil
}

func (s *CardService) GetAllUserCards(ctx context.Context, userID int, needToDecrypt bool) ([]domain.Card, error) {
	cards, err := s.cardRepository.GetByUserID(ctx, userID, &domain.StorageFilters{
		IsPaginated:    false,
		IsSortedByDate: true,

		SortDate: domain.SortDateFilters{
			IsASC: false,
		},
	})
	if err != nil {
		return nil, err
	}

	if needToDecrypt {
		for i := range cards {
			cards[i].CVV, err = s.cryptor.Decrypt(cards[i].CVV)
			if err != nil {
				return nil, err
			}

			cards[i].Number, err = s.cryptor.Decrypt(cards[i].Number)
			if err != nil {
				return nil, err
			}
		}
	}

	return cards, nil
}

func (s *CardService) AddNewCard(
	ctx context.Context,
	needEncrypt bool,
	userID int,
	number string,
	expiredAt string,
	cvv string,
	meta string,
) (*domain.Card, error) {
	savedNumber := number
	savedCVV := cvv
	var err error

	if needEncrypt {
		savedNumber, err = s.cryptor.Encrypt(number)
		if err != nil {
			return nil, err
		}

		savedCVV, err = s.cryptor.Encrypt(cvv)
		if err != nil {
			return nil, err
		}
	}

	card, err := s.cardRepository.AddCard(ctx, userID, savedNumber, expiredAt, savedCVV, meta)
	if err != nil {
		return nil, err
	}

	card.Number = number
	card.CVV = cvv

	return card, nil
}

func (s *CardService) DeleteBatchCards(ctx context.Context, userID int, ids []int) error {
	return s.cardRepository.DeleteBatch(ctx, userID, ids)
}

func (s *CardService) DeleteCardByID(
	ctx context.Context,
	userID int,
	id int,
) error {
	return s.cardRepository.DeleteByID(ctx, userID, id)
}
