package clientsync

import (
	"context"
	"fmt"

	"github.com/MowlCoder/goph-keeper/internal/domain"
	"github.com/MowlCoder/goph-keeper/internal/session"
)

type serverApi interface {
	GetAll(ctx context.Context) ([]domain.UserStoredData, error)
	Add(ctx context.Context, entity domain.UserStoredData) (*domain.UserStoredData, error)
	DeleteBatch(ctx context.Context, ids []int) error
}

type localService interface {
	GetAll(ctx context.Context) ([]domain.UserStoredData, error)
	Add(ctx context.Context, dataType string, data interface{}, meta string) (*domain.UserStoredData, error)
	DeleteBatch(ctx context.Context, ids []int) error
}

type preparedData struct {
	DelFromServer []int
	DelFromClient []int

	AddToServer []domain.UserStoredData
	AddToClient []domain.UserStoredData
}

type localRepository interface {
	SyncUpdate(ctx context.Context, oldID, newID, version int) error
}

type BaseSyncer struct {
	clientSession *session.ClientSession

	serverApi       serverApi
	localService    localService
	localRepository localRepository
}

func NewBaseSyncer(
	clientSession *session.ClientSession,
	serverApi serverApi,
	localService localService,
	localRepository localRepository,
) *BaseSyncer {
	return &BaseSyncer{
		clientSession: clientSession,

		serverApi:       serverApi,
		localService:    localService,
		localRepository: localRepository,
	}
}

func (s *BaseSyncer) Sync(ctx context.Context) error {
	if !s.clientSession.IsAuth() {
		return nil
	}

	serverDataMap, err := s.getServerData(ctx)
	if err != nil {
		return err
	}
	clientDataMap, err := s.getClientData(ctx)
	if err != nil {
		return err
	}

	data := s.prepareData(serverDataMap, clientDataMap)

	if len(data.DelFromClient) > 0 {
		if err := s.localService.DeleteBatch(ctx, data.DelFromClient); err != nil {
			return err
		}
	}

	if len(data.DelFromServer) > 0 {
		err := s.serverApi.DeleteBatch(ctx, data.DelFromServer)
		if err != nil {
			return err
		}
	}

	if len(data.AddToServer) > 0 {
		for _, data := range data.AddToServer {
			newData, err := s.serverApi.Add(
				ctx,
				data,
			)
			if err != nil {
				fmt.Println("AddToServer:", err)
				continue
			}

			s.localRepository.SyncUpdate(
				context.Background(),
				data.ID,
				newData.ID,
				newData.Version,
			)
		}
	}

	if len(data.AddToClient) > 0 {
		for _, d := range data.AddToClient {
			newData, _ := s.localService.Add(
				ctx,
				d.DataType,
				d.Data,
				d.Meta,
			)

			s.localRepository.SyncUpdate(
				ctx,
				newData.ID,
				d.ID,
				d.Version,
			)
		}
	}

	s.clientSession.ClearLogPassDeleted()

	return nil
}

func (s *BaseSyncer) SyncCommandHandler(args []string) error {
	if !s.clientSession.IsAuth() {
		return domain.ErrInvalidCommandUsage
	}

	return s.Sync(context.Background())
}

func (s *BaseSyncer) prepareData(serverData map[int]domain.UserStoredData, clientData map[int]domain.UserStoredData) *preparedData {
	pd := &preparedData{
		DelFromServer: make([]int, 0),
		DelFromClient: make([]int, 0),

		AddToServer: make([]domain.UserStoredData, 0),
		AddToClient: make([]domain.UserStoredData, 0),
	}

	for _, data := range serverData {
		if s.clientSession.IsLogPassDeleted(data.ID) {
			pd.DelFromServer = append(pd.DelFromServer, data.ID)
			continue
		}

		if _, ok := clientData[data.ID]; !ok {
			pd.AddToClient = append(pd.AddToClient, data)
		}
	}

	for _, data := range clientData {
		_, ok := serverData[data.ID]

		if !ok {
			if data.IsLocal() {
				pd.AddToServer = append(pd.AddToServer, data)
			} else {
				pd.DelFromClient = append(pd.DelFromClient, data.ID)
			}
		}
	}

	return pd
}

func (s *BaseSyncer) getServerData(ctx context.Context) (map[int]domain.UserStoredData, error) {
	serverData, err := s.serverApi.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	serverDataMap := make(map[int]domain.UserStoredData)
	for _, data := range serverData {
		serverDataMap[data.ID] = data
	}

	return serverDataMap, nil
}

func (s *BaseSyncer) getClientData(ctx context.Context) (map[int]domain.UserStoredData, error) {
	dataSet, err := s.localService.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	clientData := make(map[int]domain.UserStoredData)

	for _, data := range dataSet {
		clientData[data.ID] = data
	}

	return clientData, nil
}
