package clientsync

import (
	"context"
	"fmt"
	"time"

	"github.com/MowlCoder/goph-keeper/internal/domain"
	"github.com/MowlCoder/goph-keeper/internal/session"
)

type entity interface {
	GetID() int
	GetVersion() int
	IsLocal() bool
}

type serverApi interface {
	GetUserData(ctx context.Context) ([]entity, error)
	AddNewData(ctx context.Context, entity entity) (entity, error)
	DeleteBatchData(ctx context.Context, ids []int) error
}

type localService interface {
	GetAllUserData(ctx context.Context, userID int) ([]entity, error)
	AddNewData(ctx context.Context, userID int, entity entity) (entity, error)
	DeleteBatchData(ctx context.Context, userID int, ids []int) error
}

type preparedData struct {
	DelFromServer []int
	DelFromClient []int

	AddToServer []entity
	AddToClient []entity
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

func newBaseSyncer(
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

	serverPairsMap, err := s.getServerPairs(ctx)
	if err != nil {
		return err
	}
	clientPairsMap, err := s.getClientPairs(ctx)
	if err != nil {
		return err
	}

	data := s.prepareData(serverPairsMap, clientPairsMap)

	if len(data.DelFromClient) > 0 {
		if err := s.localService.DeleteBatchData(ctx, domain.LocalUserID, data.DelFromClient); err != nil {
			return err
		}
	}

	if len(data.DelFromServer) > 0 {
		err := s.serverApi.DeleteBatchData(ctx, data.DelFromServer)
		if err != nil {
			return err
		}
	}

	if len(data.AddToServer) > 0 {
		for _, pair := range data.AddToServer {
			newLogPass, err := s.serverApi.AddNewData(
				ctx,
				pair,
			)
			if err != nil {
				fmt.Println(err)
				continue
			}

			s.localRepository.SyncUpdate(
				context.Background(),
				pair.GetID(),
				newLogPass.GetID(),
				newLogPass.GetVersion(),
			)
		}
	}

	if len(data.AddToClient) > 0 {
		for _, pair := range data.AddToClient {
			newPair, _ := s.localService.AddNewData(
				ctx,
				domain.LocalUserID,
				pair,
			)

			s.localRepository.SyncUpdate(
				ctx,
				newPair.GetID(),
				pair.GetID(),
				pair.GetVersion(),
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

func (s *BaseSyncer) InfiniteSync(interval time.Duration) {
	ticker := time.NewTicker(interval)

	for range ticker.C {
		ctx := context.Background()

		err := s.Sync(ctx)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (s *BaseSyncer) prepareData(serverData map[int]entity, clientData map[int]entity) *preparedData {
	pd := &preparedData{
		DelFromServer: make([]int, 0),
		DelFromClient: make([]int, 0),

		AddToServer: make([]entity, 0),
		AddToClient: make([]entity, 0),
	}

	for _, serverPair := range serverData {
		if s.clientSession.IsLogPassDeleted(serverPair.GetID()) {
			pd.DelFromServer = append(pd.DelFromServer, serverPair.GetID())
			continue
		}

		if _, ok := clientData[serverPair.GetID()]; !ok {
			pd.AddToClient = append(pd.AddToClient, serverPair)
		}
	}

	for _, clientPair := range clientData {
		_, ok := serverData[clientPair.GetID()]

		if !ok {
			if clientPair.IsLocal() {
				pd.AddToServer = append(pd.AddToServer, clientPair)
			} else {
				pd.DelFromClient = append(pd.DelFromClient, clientPair.GetID())
			}
		}
	}

	return pd
}

func (s *BaseSyncer) getServerPairs(ctx context.Context) (map[int]entity, error) {
	serverData, err := s.serverApi.GetUserData(ctx)
	if err != nil {
		return nil, err
	}

	serverPairsMap := make(map[int]entity)
	for _, serverPair := range serverData {
		serverPairsMap[serverPair.GetID()] = serverPair
	}

	return serverPairsMap, nil
}

func (s *BaseSyncer) getClientPairs(ctx context.Context) (map[int]entity, error) {
	clientPairs, err := s.localService.GetAllUserData(ctx, domain.LocalUserID)
	if err != nil {
		return nil, err
	}

	clientPairsMap := make(map[int]entity)

	for _, clientPair := range clientPairs {
		clientPairsMap[clientPair.GetID()] = clientPair
	}

	return clientPairsMap, nil
}
