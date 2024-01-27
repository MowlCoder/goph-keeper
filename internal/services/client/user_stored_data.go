package client

import (
	"context"
	"encoding/json"
	"math"
	"time"

	"github.com/MowlCoder/goph-keeper/internal/domain"
)

type cryptorForUserStoredDataService interface {
	EncryptBytes(raw []byte) ([]byte, error)
	DecryptBytes(crypted []byte) ([]byte, error)
}

type userStoredDataRepository interface {
	GetAll(ctx context.Context) ([]domain.UserStoredData, error)
	AddData(ctx context.Context, dataType string, data []byte, meta string) (int64, error)
	GetWithType(ctx context.Context, dataType string, filters *domain.StorageFilters) ([]domain.UserStoredData, error)
	CountUserDataOfType(ctx context.Context, dataType string) (int, error)
	DeleteByID(ctx context.Context, id int) error
	DeleteBatch(ctx context.Context, ids []int) error
}

type UserStoredDataService struct {
	repository userStoredDataRepository
	cryptor    cryptorForUserStoredDataService
}

func NewUserStoredDataService(repository userStoredDataRepository, cryptor cryptorForUserStoredDataService) *UserStoredDataService {
	return &UserStoredDataService{
		repository: repository,
		cryptor:    cryptor,
	}
}

func (s *UserStoredDataService) GetAll(ctx context.Context) ([]domain.UserStoredData, error) {
	dataSet, err := s.repository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	for idx, data := range dataSet {
		decryptedBytes, err := s.cryptor.DecryptBytes(data.CryptedData)
		if err != nil {
			return nil, err
		}

		switch data.DataType {
		case domain.LogPassDataType:
			var parsedData domain.LogPassData
			if err := json.Unmarshal(decryptedBytes, &parsedData); err != nil {
				return nil, err
			}
			dataSet[idx].Data = parsedData
		case domain.CardDataType:
			var parsedData domain.CardData
			if err := json.Unmarshal(decryptedBytes, &parsedData); err != nil {
				return nil, err
			}
			dataSet[idx].Data = parsedData
		default:
			return nil, domain.ErrInvalidDataType
		}
	}

	return dataSet, err
}

func (s *UserStoredDataService) GetUserData(ctx context.Context, dataType string, filters *domain.StorageFilters) (*domain.PaginatedResult, error) {
	dataSet, err := s.repository.GetWithType(ctx, dataType, filters)
	if err != nil {
		return nil, err
	}
	pairsCount, err := s.repository.CountUserDataOfType(ctx, dataType)
	if err != nil {
		return nil, err
	}

	for idx, data := range dataSet {
		decryptedBytes, err := s.cryptor.DecryptBytes(data.CryptedData)
		if err != nil {
			return nil, err
		}

		switch dataType {
		case domain.LogPassDataType:
			var parsedData domain.LogPassData
			if err := json.Unmarshal(decryptedBytes, &parsedData); err != nil {
				return nil, err
			}
			dataSet[idx].Data = parsedData
		case domain.CardDataType:
			var parsedData domain.CardData
			if err := json.Unmarshal(decryptedBytes, &parsedData); err != nil {
				return nil, err
			}
			dataSet[idx].Data = parsedData
		default:
			return nil, domain.ErrInvalidDataType
		}
	}

	return &domain.PaginatedResult{
		Count:       filters.Pagination.Count,
		CurrentPage: filters.Pagination.Page,
		PageCount:   int(math.Ceil(float64(pairsCount) / float64(filters.Pagination.Count))),
		Data:        dataSet,
	}, nil
}

func (s *UserStoredDataService) Add(ctx context.Context, dataType string, data interface{}, meta string) (*domain.UserStoredData, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	encrypted, err := s.cryptor.EncryptBytes(jsonData)
	if err != nil {
		return nil, err
	}

	insertedID, err := s.repository.AddData(ctx, dataType, encrypted, meta)
	if err != nil {
		return nil, err
	}

	return &domain.UserStoredData{
		ID:        int(insertedID),
		DataType:  dataType,
		Data:      data,
		Meta:      meta,
		Version:   1,
		CreatedAt: time.Now().UTC(),
	}, nil
}

func (s *UserStoredDataService) DeleteBatch(ctx context.Context, ids []int) error {
	return s.repository.DeleteBatch(ctx, ids)
}

func (s *UserStoredDataService) DeleteByID(ctx context.Context, id int) error {
	return s.repository.DeleteByID(ctx, id)
}
