package services

import (
	"context"
	"math"

	"github.com/MowlCoder/goph-keeper/internal/domain"
)

type cryptorForLogPass interface {
	Encrypt(raw string) (string, error)
	Decrypt(crypted string) (string, error)
}

type logPassRepository interface {
	GetPairsByUserID(ctx context.Context, userID int, filters *domain.StorageFilters) ([]domain.LogPass, error)
	CountUserPairs(ctx context.Context, userID int) (int, error)
	AddPair(ctx context.Context, userID int, login string, password string, source string) (*domain.LogPass, error)
	DeleteByID(ctx context.Context, userID int, id int) error
	DeleteBatch(ctx context.Context, userID int, id []int) error
}

type LogPassService struct {
	logPassRepository logPassRepository
	cryptor           cryptorForLogPass
}

func NewLogPassService(
	logPassRepository logPassRepository,
	cryptor cryptorForLogPass,
) *LogPassService {
	return &LogPassService{
		logPassRepository: logPassRepository,
		cryptor:           cryptor,
	}
}

func (s *LogPassService) GetUserPairs(ctx context.Context, userID int, filters *domain.StorageFilters, needToDecrypt bool) (*domain.PaginatedResult, error) {
	pairs, err := s.logPassRepository.GetPairsByUserID(ctx, userID, filters)
	if err != nil {
		return nil, err
	}
	pairsCount, err := s.logPassRepository.CountUserPairs(ctx, userID)
	if err != nil {
		return nil, err
	}

	if needToDecrypt {
		for i := range pairs {
			pairs[i].Password, err = s.cryptor.Decrypt(pairs[i].Password)
			if err != nil {
				return nil, err
			}
		}
	}

	return &domain.PaginatedResult{
		Count:       filters.Pagination.Count,
		CurrentPage: filters.Pagination.Page,
		PageCount:   int(math.Ceil(float64(pairsCount) / float64(filters.Pagination.Count))),
		Data:        pairs,
	}, nil
}

func (s *LogPassService) GetAllUserPairs(ctx context.Context, userID int, needToDecrypt bool) ([]domain.LogPass, error) {
	pairs, err := s.logPassRepository.GetPairsByUserID(ctx, userID, &domain.StorageFilters{
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
		for i := range pairs {
			pairs[i].Password, err = s.cryptor.Decrypt(pairs[i].Password)
			if err != nil {
				return nil, err
			}
		}
	}

	return pairs, nil
}

func (s *LogPassService) AddNewPair(
	ctx context.Context,
	needEncrypt bool,
	userID int,
	login string,
	password string,
	source string,
) (*domain.LogPass, error) {
	var savedPassword = password
	var err error

	if needEncrypt {
		savedPassword, err = s.cryptor.Encrypt(password)
		if err != nil {
			return nil, err
		}
	}

	logPass, err := s.logPassRepository.AddPair(ctx, userID, login, savedPassword, source)
	if err != nil {
		return nil, err
	}

	logPass.Password = password

	return logPass, nil
}

func (s *LogPassService) DeleteBatchPairs(ctx context.Context, userID int, ids []int) error {
	return s.logPassRepository.DeleteBatch(ctx, userID, ids)
}

func (s *LogPassService) DeletePairByID(
	ctx context.Context,
	userID int,
	id int,
) error {
	return s.logPassRepository.DeleteByID(ctx, userID, id)
}
