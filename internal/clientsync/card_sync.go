package clientsync

import (
	"context"
	"errors"

	"github.com/MowlCoder/goph-keeper/internal/domain"
	"github.com/MowlCoder/goph-keeper/internal/session"
)

type cardApi interface {
	GetUserCards(ctx context.Context) ([]domain.Card, error)
	AddNewCard(ctx context.Context, number, expiredAt, cvv, meta string) (*domain.Card, error)
	DeleteBatchCards(ctx context.Context, ids []int) error
}

type cardService interface {
	GetAllUserCards(ctx context.Context, userID int, needToDecrypt bool) ([]domain.Card, error)
	AddNewCard(ctx context.Context, needEncrypt bool, userID int, number, expiredAt, cvv, meta string) (*domain.Card, error)
	DeleteBatchCards(ctx context.Context, userID int, ids []int) error
}

type cardLocalRepository interface {
	SyncUpdate(ctx context.Context, oldID, newID, version int) error
}

type cardApiAdapter struct {
	cardApi cardApi
}

func (a *cardApiAdapter) GetUserData(ctx context.Context) ([]entity, error) {
	pairs, err := a.cardApi.GetUserCards(ctx)
	if err != nil {
		return nil, err
	}

	d := make([]entity, 0, len(pairs))
	for _, pair := range pairs {
		d = append(d, pair)
	}

	return d, nil
}

func (a *cardApiAdapter) AddNewData(ctx context.Context, data entity) (entity, error) {
	card, ok := data.(domain.Card)
	if !ok {
		return nil, errors.New("invalid type")
	}

	newCard, err := a.cardApi.AddNewCard(ctx, card.Number, card.ExpiredAt, card.CVV, card.Meta)
	if err != nil {
		return nil, err
	}

	return newCard, err
}

func (a *cardApiAdapter) DeleteBatchData(ctx context.Context, ids []int) error {
	return a.cardApi.DeleteBatchCards(ctx, ids)
}

type cardServiceAdapter struct {
	cardService cardService
}

func (a *cardServiceAdapter) GetAllUserData(ctx context.Context, userID int) ([]entity, error) {
	pairs, err := a.cardService.GetAllUserCards(ctx, userID, false)
	if err != nil {
		return nil, err
	}

	d := make([]entity, 0, len(pairs))
	for _, pair := range pairs {
		d = append(d, pair)
	}

	return d, err
}

func (a *cardServiceAdapter) AddNewData(ctx context.Context, userID int, data entity) (entity, error) {
	card, ok := data.(domain.Card)
	if !ok {
		return nil, errors.New("invalid type")
	}

	newCard, err := a.cardService.AddNewCard(ctx, false, userID, card.Number, card.ExpiredAt, card.CVV, card.Meta)
	if err != nil {
		return nil, err
	}

	return newCard, err
}

func (a *cardServiceAdapter) DeleteBatchData(ctx context.Context, userID int, ids []int) error {
	return a.cardService.DeleteBatchCards(ctx, userID, ids)
}

func NewCardSyncer(
	clientSession *session.ClientSession,
	cardApi cardApi,
	cardService cardService,
	cardLocalRepository cardLocalRepository,
) *BaseSyncer {
	return newBaseSyncer(
		clientSession,
		&cardApiAdapter{cardApi: cardApi},
		&cardServiceAdapter{cardService: cardService},
		cardLocalRepository,
	)
}
