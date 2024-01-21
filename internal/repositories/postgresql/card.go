package postgresql

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/MowlCoder/goph-keeper/internal/domain"
)

type CardRepository struct {
	pool *pgxpool.Pool
}

func NewCardRepository(pool *pgxpool.Pool) *CardRepository {
	return &CardRepository{
		pool: pool,
	}
}

func (r *CardRepository) GetByUserID(ctx context.Context, userID int, filters *domain.StorageFilters) ([]domain.Card, error) {
	baseQuery := `
		SELECT id, user_id, number, expired_at, cvv, meta, version, created_at
		FROM cards
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

	cards := make([]domain.Card, 0)

	for rows.Next() {
		var card domain.Card

		if err := rows.Scan(&card.ID, &card.UserID, &card.Number, &card.ExpiredAt, &card.CVV, &card.Meta, &card.Version, &card.CreatedAt); err != nil {
			return nil, err
		}

		cards = append(cards, card)
	}

	return cards, nil
}

func (r *CardRepository) CountUserCards(ctx context.Context, userID int) (int, error) {
	query := `
		SELECT COUNT(id)
		FROM cards
		WHERE user_id = $1
	`

	var count int
	err := r.pool.QueryRow(ctx, query, userID).Scan(&count)
	if err != nil {
		return -1, err
	}

	return count, nil
}

func (r *CardRepository) AddCard(
	ctx context.Context,
	userID int,
	number string,
	expiredAt string,
	cvv string,
	meta string,
) (*domain.Card, error) {
	var insertedID int64
	query := `
		INSERT INTO cards (user_id, number, expired_at, cvv, meta)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	err := r.pool.QueryRow(
		ctx,
		query,
		userID, number, expiredAt, cvv, meta,
	).Scan(&insertedID)
	if err != nil {
		return nil, err
	}

	return &domain.Card{
		ID:        int(insertedID),
		Number:    number,
		ExpiredAt: expiredAt,
		UserID:    userID,
		CVV:       cvv,
		Meta:      meta,
		Version:   1,
		CreatedAt: time.Now().UTC(),
	}, nil
}

func (r *CardRepository) DeleteByID(ctx context.Context, userID int, id int) error {
	query := `
		DELETE FROM cards
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

func (r *CardRepository) DeleteBatch(ctx context.Context, userID int, id []int) error {
	query := `
		DELETE FROM cards
		WHERE user_id = $1 AND id = ANY($2)
	`

	_, err := r.pool.Exec(ctx, query, userID, id)
	if err != nil {
		return err
	}

	return nil
}
