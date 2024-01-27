package postgresql

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/MowlCoder/goph-keeper/internal/domain"
)

type UserStoredDataRepository struct {
	pool *pgxpool.Pool
}

func NewUserStoredDataRepository(pool *pgxpool.Pool) *UserStoredDataRepository {
	return &UserStoredDataRepository{
		pool: pool,
	}
}

func (repo *UserStoredDataRepository) GetByID(ctx context.Context, id int) (*domain.UserStoredData, error) {
	query := `
		SELECT id, user_id, data_type, data, meta, version, created_at FROM user_stored_data
		WHERE id = $1 
	`

	var userData domain.UserStoredData
	if err := repo.pool.QueryRow(ctx, query, id).Scan(
		&userData.ID,
		&userData.UserID,
		&userData.DataType,
		&userData.CryptedData,
		&userData.Meta,
		&userData.Version,
		&userData.CreatedAt,
	); err != nil {
		return nil, err
	}

	return &userData, nil
}

func (repo *UserStoredDataRepository) GetUserAllData(ctx context.Context, userID int) ([]domain.UserStoredData, error) {
	query := `
		SELECT id, user_id, data_type, data, meta, version, created_at FROM user_stored_data
		WHERE user_id = $1
	`

	rows, err := repo.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	dataSet := make([]domain.UserStoredData, 0)
	for rows.Next() {
		var data domain.UserStoredData
		if err := rows.Scan(&data.ID, &data.UserID, &data.DataType, &data.CryptedData, &data.Meta, &data.Version, &data.CreatedAt); err != nil {
			return nil, err
		}

		dataSet = append(dataSet, data)
	}

	return dataSet, nil
}

func (repo *UserStoredDataRepository) GetWithType(ctx context.Context, userID int, dataType string, filters *domain.StorageFilters) ([]domain.UserStoredData, error) {
	baseQuery := `
		SELECT id, user_id, data_type, data, meta, version, created_at FROM user_stored_data
		WHERE user_id = $1 AND data_type = $2
	`

	rows, err := repo.pool.Query(ctx, filters.BuildSQL(baseQuery), userID, dataType)
	if err != nil {
		return nil, err
	}

	dataSet := make([]domain.UserStoredData, 0)
	for rows.Next() {
		var data domain.UserStoredData
		if err := rows.Scan(&data.ID, &data.UserID, &data.DataType, &data.CryptedData, &data.Meta, &data.Version, &data.CreatedAt); err != nil {
			return nil, err
		}

		dataSet = append(dataSet, data)
	}

	return dataSet, nil
}

func (repo *UserStoredDataRepository) CountUserDataOfType(ctx context.Context, userID int, dataType string) (int, error) {
	query := `
		SELECT COUNT(id)
		FROM user_stored_data
		WHERE user_id = $1 AND data_type = $2
	`

	var count int
	err := repo.pool.QueryRow(ctx, query, userID, dataType).Scan(&count)
	if err != nil {
		return -1, err
	}

	return count, nil
}

func (repo *UserStoredDataRepository) AddData(ctx context.Context, userID int, dataType string, data []byte, meta string) (int64, error) {
	query := `
		INSERT INTO user_stored_data (user_id, data_type, data, meta)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	var insertedID int64

	err := repo.pool.QueryRow(
		ctx,
		query,
		userID, dataType, data, meta,
	).Scan(&insertedID)
	if err != nil {
		return 0, err
	}

	return insertedID, nil
}

func (repo *UserStoredDataRepository) UpdateUserData(ctx context.Context, userID int, dataID int, data interface{}, meta string) (*domain.UserStoredData, error) {
	query := `
		UPDATE user_stored_data
		SET data = $1, meta = $2, version = version + 1
		WHERE id = $3 AND user_id = $4
		RETURNING id, user_id, data_type, data, meta, version, created_at
	`

	var userData domain.UserStoredData
	if err := repo.pool.QueryRow(ctx, query, data, meta, dataID, userID).Scan(
		&userData.ID,
		&userData.UserID,
		&userData.DataType,
		&userData.CryptedData,
		&userData.Meta,
		&userData.Version,
		&userData.CreatedAt,
	); err != nil {
		return nil, err
	}

	return &userData, nil
}

func (repo *UserStoredDataRepository) DeleteByID(ctx context.Context, userID int, id int) error {
	query := `
		DELETE FROM user_stored_data
		WHERE id = $1 AND user_id = $2
	`

	result, err := repo.pool.Exec(ctx, query, id, userID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (repo *UserStoredDataRepository) DeleteBatch(ctx context.Context, userID int, id []int) error {
	query := `
		DELETE FROM user_stored_data
		WHERE user_id = $1 AND id = ANY($2)
	`

	_, err := repo.pool.Exec(ctx, query, userID, id)
	if err != nil {
		return err
	}

	return nil
}
