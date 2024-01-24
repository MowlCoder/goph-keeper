package postgresql

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgxpool"
)

type cryptorForUserStoredDataRepo interface {
	EncryptBytes(raw []byte) ([]byte, error)
	DecryptBytes(crypted []byte) ([]byte, error)
}

type UserStoredDataRepository struct {
	pool    *pgxpool.Pool
	cryptor cryptorForUserStoredDataRepo
}

func NewUserStoredDataRepository(pool *pgxpool.Pool, cryptor cryptorForUserStoredDataRepo) *UserStoredDataRepository {
	return &UserStoredDataRepository{
		pool:    pool,
		cryptor: cryptor,
	}
}

func (repo *UserStoredDataRepository) AddData(ctx context.Context, userID int, dataType string, data map[string]interface{}, meta string) (int64, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return 0, err
	}

	encryptedData, err := repo.cryptor.EncryptBytes(jsonData)
	if err != nil {
		return 0, err
	}

	query := `
		INSERT INTO user_stored_data (user_id, data_type, data, meta)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	var insertedID int64

	err = repo.pool.QueryRow(
		ctx,
		query,
		userID, dataType, encryptedData, meta,
	).Scan(&insertedID)
	if err != nil {
		return 0, err
	}

	return insertedID, nil
}
