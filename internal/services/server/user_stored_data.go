package server

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
	GetByID(ctx context.Context, id int) (*domain.UserStoredData, error)
	AddData(ctx context.Context, userID int, dataType string, data []byte, meta string) (int64, error)
	GetUserAllData(ctx context.Context, userID int) ([]domain.UserStoredData, error)
	GetWithType(ctx context.Context, userID int, dataType string, filters *domain.StorageFilters) ([]domain.UserStoredData, error)
	CountUserDataOfType(ctx context.Context, userID int, dataType string) (int, error)
	UpdateUserData(ctx context.Context, userID int, dataID int, data interface{}, meta string) (*domain.UserStoredData, error)
	DeleteBatch(ctx context.Context, userID int, id []int) error
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

func (s *UserStoredDataService) GetAllUserData(ctx context.Context, userID int) ([]domain.UserStoredData, error) {
	dataSet, err := s.repository.GetUserAllData(ctx, userID)
	if err != nil {
		return nil, err
	}

	for idx, data := range dataSet {
		decryptedBytes, err := s.cryptor.DecryptBytes(data.CryptedData)
		if err != nil {
			return nil, err
		}
		parsedData, err := domain.ParseUserStoredData(data.DataType, decryptedBytes)
		if err != nil {
			return nil, err
		}

		dataSet[idx].Data = parsedData
		dataSet[idx].CryptedData = nil
	}

	return dataSet, nil
}

func (s *UserStoredDataService) GetUserData(ctx context.Context, userID int, dataType string, filters *domain.StorageFilters) (*domain.PaginatedResult, error) {
	dataSet, err := s.repository.GetWithType(ctx, userID, dataType, filters)
	if err != nil {
		return nil, err
	}
	pairsCount, err := s.repository.CountUserDataOfType(ctx, userID, dataType)
	if err != nil {
		return nil, err
	}

	for idx, data := range dataSet {
		decryptedBytes, err := s.cryptor.DecryptBytes(data.CryptedData)
		if err != nil {
			return nil, err
		}
		parsedData, err := domain.ParseUserStoredData(data.DataType, decryptedBytes)
		if err != nil {
			return nil, err
		}

		dataSet[idx].Data = parsedData
		dataSet[idx].CryptedData = nil
	}

	return &domain.PaginatedResult{
		Count:       filters.Pagination.Count,
		CurrentPage: filters.Pagination.Page,
		PageCount:   int(math.Ceil(float64(pairsCount) / float64(filters.Pagination.Count))),
		Data:        dataSet,
	}, nil
}

func (s *UserStoredDataService) GetUserDataByID(ctx context.Context, userID int, id int) (*domain.UserStoredData, error) {
	userData, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if userData.UserID != userID {
		return nil, domain.ErrUserStoredDataNotFound
	}

	decryptedBytes, err := s.cryptor.DecryptBytes(userData.CryptedData)
	if err != nil {
		return nil, err
	}
	parsedData, err := domain.ParseUserStoredData(userData.DataType, decryptedBytes)
	if err != nil {
		return nil, err
	}

	userData.Data = parsedData

	return userData, nil
}

func (s *UserStoredDataService) UpdateUserData(ctx context.Context, userID int, dataID int, data interface{}, meta string) (*domain.UserStoredData, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	encrypted, err := s.cryptor.EncryptBytes(jsonData)
	if err != nil {
		return nil, err
	}

	newDate, err := s.repository.UpdateUserData(ctx, userID, dataID, encrypted, meta)
	if err != nil {
		return nil, err
	}

	newDate.Data = data
	newDate.CryptedData = nil

	return newDate, nil
}

func (s *UserStoredDataService) Add(ctx context.Context, userID int, dataType string, data interface{}, meta string) (*domain.UserStoredData, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	encrypted, err := s.cryptor.EncryptBytes(jsonData)
	if err != nil {
		return nil, err
	}

	insertedID, err := s.repository.AddData(ctx, userID, dataType, encrypted, meta)
	if err != nil {
		return nil, err
	}

	return &domain.UserStoredData{
		ID:        int(insertedID),
		UserID:    userID,
		DataType:  dataType,
		Data:      data,
		Meta:      meta,
		Version:   1,
		CreatedAt: time.Now().UTC(),
	}, nil
}

func (s *UserStoredDataService) DeleteBatch(ctx context.Context, userID int, ids []int) error {
	return s.repository.DeleteBatch(ctx, userID, ids)
}
