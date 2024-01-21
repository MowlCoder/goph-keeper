package clientsync

import (
	"context"
	"errors"

	"github.com/MowlCoder/goph-keeper/internal/domain"
	"github.com/MowlCoder/goph-keeper/internal/session"
)

type logPassApi interface {
	GetUserPairs(ctx context.Context) ([]domain.LogPass, error)
	AddNewPair(ctx context.Context, login, password, source string) (*domain.LogPass, error)
	DeleteBatchPairs(ctx context.Context, ids []int) error
}

type logPassService interface {
	GetAllUserPairs(ctx context.Context, userID int, needToDecrypt bool) ([]domain.LogPass, error)
	AddNewPair(ctx context.Context, needEncrypt bool, userID int, login, password, source string) (*domain.LogPass, error)
	DeleteBatchPairs(ctx context.Context, userID int, ids []int) error
}

type logPassLocalRepository interface {
	SyncUpdate(ctx context.Context, oldID, newID, version int) error
}

type logPassApiAdapter struct {
	logPassApi logPassApi
}

func (a *logPassApiAdapter) GetUserData(ctx context.Context) ([]entity, error) {
	pairs, err := a.logPassApi.GetUserPairs(ctx)
	if err != nil {
		return nil, err
	}

	d := make([]entity, 0, len(pairs))
	for _, pair := range pairs {
		d = append(d, pair)
	}

	return d, nil
}

func (a *logPassApiAdapter) AddNewData(ctx context.Context, data entity) (entity, error) {
	pair, ok := data.(domain.LogPass)
	if !ok {
		return nil, errors.New("invalid type")
	}

	newPair, err := a.logPassApi.AddNewPair(ctx, pair.Login, pair.Password, pair.Source)
	if err != nil {
		return nil, err
	}

	return newPair, err
}

func (a *logPassApiAdapter) DeleteBatchData(ctx context.Context, ids []int) error {
	return a.logPassApi.DeleteBatchPairs(ctx, ids)
}

type logPassServiceAdapter struct {
	logPassService logPassService
}

func (a *logPassServiceAdapter) GetAllUserData(ctx context.Context, userID int) ([]entity, error) {
	pairs, err := a.logPassService.GetAllUserPairs(ctx, userID, false)
	if err != nil {
		return nil, err
	}

	d := make([]entity, 0, len(pairs))
	for _, pair := range pairs {
		d = append(d, pair)
	}

	return d, err
}

func (a *logPassServiceAdapter) AddNewData(ctx context.Context, userID int, data entity) (entity, error) {
	pair, ok := data.(domain.LogPass)
	if !ok {
		return nil, errors.New("invalid type")
	}

	newPair, err := a.logPassService.AddNewPair(ctx, false, userID, pair.Login, pair.Password, pair.Source)
	if err != nil {
		return nil, err
	}

	return newPair, err
}

func (a *logPassServiceAdapter) DeleteBatchData(ctx context.Context, userID int, ids []int) error {
	return a.logPassService.DeleteBatchPairs(ctx, userID, ids)
}

func NewLogPassSyncer(
	clientSession *session.ClientSession,
	logPassApi logPassApi,
	logPassService logPassService,
	logPassLocalRepository logPassLocalRepository,
) *BaseSyncer {
	return newBaseSyncer(
		clientSession,
		&logPassApiAdapter{logPassApi: logPassApi},
		&logPassServiceAdapter{logPassService: logPassService},
		logPassLocalRepository,
	)
}
