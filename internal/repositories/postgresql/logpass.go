package postgresql

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/MowlCoder/goph-keeper/internal/domain"
)

type LogPassRepository struct {
	pool *pgxpool.Pool
}

func NewLogPassRepository(pool *pgxpool.Pool) *LogPassRepository {
	return &LogPassRepository{
		pool: pool,
	}
}

func (r *LogPassRepository) GetPairsByUserID(ctx context.Context, userID int, filters *domain.StorageFilters) ([]domain.LogPass, error) {
	baseQuery := `
		SELECT id, user_id, login, password, source, version, created_at
		FROM log_pass
		WHERE user_id = $1
	`

	rows, err := r.pool.Query(
		ctx,
		filters.BuildSQL(baseQuery),
		userID,
	)
	if err != nil {
		return nil, err
	}

	pairs := make([]domain.LogPass, 0)

	for rows.Next() {
		var pair domain.LogPass

		if err := rows.Scan(&pair.ID, &pair.UserID, &pair.Login, &pair.Password, &pair.Source, &pair.Version, &pair.CreatedAt); err != nil {
			return nil, err
		}

		pairs = append(pairs, pair)
	}

	return pairs, nil
}

func (r *LogPassRepository) CountUserPairs(ctx context.Context, userID int) (int, error) {
	query := `
		SELECT COUNT(id)
		FROM log_pass
		WHERE user_id = $1
	`

	var count int
	err := r.pool.QueryRow(ctx, query, userID).Scan(&count)
	if err != nil {
		return -1, err
	}

	return count, nil
}

func (r *LogPassRepository) AddPair(
	ctx context.Context,
	userID int,
	login string,
	password string,
	source string,
) (*domain.LogPass, error) {
	var insertedID int64
	query := `
		INSERT INTO log_pass (user_id, login, password, source)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	err := r.pool.QueryRow(
		ctx,
		query,
		userID, login, password, source,
	).Scan(&insertedID)
	if err != nil {
		return nil, err
	}

	return &domain.LogPass{
		ID:        int(insertedID),
		Login:     login,
		Password:  password,
		UserID:    userID,
		Source:    source,
		Version:   1,
		CreatedAt: time.Now().UTC(),
	}, nil
}

func (r *LogPassRepository) DeleteByID(ctx context.Context, userID int, id int) error {
	query := `
		DELETE FROM log_pass
		WHERE id = $1 AND user_id = $2
	`

	result, err := r.pool.Exec(ctx, query, id, userID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *LogPassRepository) DeleteBatch(ctx context.Context, userID int, id []int) error {
	query := `
		DELETE FROM log_pass
		WHERE user_id = $1 AND id = ANY($2)
	`

	_, err := r.pool.Exec(ctx, query, userID, id)
	if err != nil {
		return err
	}

	return nil
}
