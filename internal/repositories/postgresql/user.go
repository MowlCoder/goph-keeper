package postgresql

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/MowlCoder/goph-keeper/internal/domain"
	"github.com/MowlCoder/goph-keeper/internal/storage/postgresql"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		pool: pool,
	}
}

func (r *UserRepository) Create(ctx context.Context, email string, hashedPassword string) (*domain.User, error) {
	var insertedID int64

	query := `
		INSERT INTO users (email, password)
		VALUES ($1, $2)
		RETURNING id
	`

	err := r.pool.QueryRow(
		ctx,
		query,
		email, hashedPassword,
	).Scan(&insertedID)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == postgresql.PgUniqueIndexErrorCode {
			return nil, domain.ErrEmailAlreadyTaken
		}

		return nil, err
	}

	return &domain.User{
		ID:        int(insertedID),
		Email:     email,
		Password:  hashedPassword,
		CreatedAt: time.Now().UTC(),
	}, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User

	query := `
		SELECT id, email, password, created_at
		FROM users
		WHERE email = $1
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		email,
	)

	if err := row.Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt); err != nil {
		return nil, domain.ErrUserNotFound
	}

	return &user, nil
}
