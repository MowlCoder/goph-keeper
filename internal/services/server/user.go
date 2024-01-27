package server

import (
	"context"

	"github.com/MowlCoder/goph-keeper/internal/domain"
)

type userRepository interface {
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Create(ctx context.Context, email string, password string) (*domain.User, error)
}

type passwordHasher interface {
	Hash(original string) (string, error)
	Equal(original string, hash string) bool
}

type UserService struct {
	repository userRepository
	hasher     passwordHasher
}

func NewUserService(
	repository userRepository,
	hasher passwordHasher,
) *UserService {
	return &UserService{
		hasher:     hasher,
		repository: repository,
	}
}

func (s *UserService) Create(ctx context.Context, email string, password string) (*domain.User, error) {
	hash, err := s.hasher.Hash(password)
	if err != nil {
		return nil, err
	}

	return s.repository.Create(ctx, email, hash)
}

func (s *UserService) Authorize(ctx context.Context, email string, password string) (*domain.User, error) {
	user, err := s.repository.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if !s.hasher.Equal(password, user.Password) {
		return nil, domain.ErrWrongCredentials
	}

	return user, nil
}
