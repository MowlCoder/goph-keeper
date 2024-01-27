package handlers

import (
	"context"

	"github.com/MowlCoder/goph-keeper/internal/domain"
)

type userStoredDataService interface {
	GetUserData(ctx context.Context, dataType string, filters *domain.StorageFilters) (*domain.PaginatedResult, error)
	Add(ctx context.Context, dataType string, data interface{}, meta string) (*domain.UserStoredData, error)
	GetByID(ctx context.Context, id int) (*domain.UserStoredData, error)
	UpdateByID(ctx context.Context, id int, data interface{}, meta string) (*domain.UserStoredData, error)
	DeleteByID(ctx context.Context, id int) error
}
