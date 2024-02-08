package clientsync

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	mock_clientsync "github.com/MowlCoder/goph-keeper/internal/clientsync/mocks"
	"github.com/MowlCoder/goph-keeper/internal/domain"
	"github.com/MowlCoder/goph-keeper/internal/session"
)

type baseSyncerTestSuite struct {
	suite.Suite

	session         *session.ClientSession
	serverApi       *mock_clientsync.MockserverApi
	localService    *mock_clientsync.MocklocalService
	localRepository *mock_clientsync.MocklocalRepository

	file   *os.File
	syncer *BaseSyncer
}

func (suite *baseSyncerTestSuite) SetupSuite() {
	suite.file, _ = os.Create(suite.T().TempDir() + "/temp.json")
}

func (suite *baseSyncerTestSuite) TearDownSuite() {
	suite.file.Close()
}

func (suite *baseSyncerTestSuite) SetupTest() {
	ctrl := gomock.NewController(suite.T())

	suite.serverApi = mock_clientsync.NewMockserverApi(ctrl)
	suite.localService = mock_clientsync.NewMocklocalService(ctrl)
	suite.localRepository = mock_clientsync.NewMocklocalRepository(ctrl)
	suite.session = session.NewClientSession(suite.file)

	suite.syncer = NewBaseSyncer(
		suite.session,
		suite.serverApi,
		suite.localService,
		suite.localRepository,
	)
}

func (suite *baseSyncerTestSuite) TearDownTest() {
}

func TestBaseSyncer(t *testing.T) {
	suite.Run(t, new(baseSyncerTestSuite))
}

func (suite *baseSyncerTestSuite) TestBaseSyncer_Sync() {
	testCases := []struct {
		name    string
		prepare func()
		err     error
	}{
		{
			name: "success",
			prepare: func() {
				localData := domain.UserStoredData{
					ID:       -2,
					Version:  -1,
					DataType: domain.TextDataType,
					Data: domain.TextData{
						Text: "text",
					},
				}

				serverData := domain.UserStoredData{
					ID:       1,
					Version:  1,
					DataType: domain.CardDataType,
					Data: domain.CardData{
						Number:    "1111111122222222",
						ExpiredAt: "12/30",
						CVV:       "123",
					},
				}

				suite.session.SetToken("some-token")

				suite.serverApi.
					EXPECT().
					GetAll(gomock.Any()).
					Return([]domain.UserStoredData{serverData}, nil)

				suite.localService.
					EXPECT().
					GetAll(gomock.Any()).
					Return([]domain.UserStoredData{localData, {ID: 3, Version: 1}}, nil)

				suite.localService.
					EXPECT().
					DeleteBatch(gomock.Any(), []int{3}).
					Return(nil)

				suite.serverApi.
					EXPECT().
					Add(gomock.Any(), localData).
					Return(&domain.UserStoredData{ID: 10, Version: 1}, nil)

				suite.localRepository.
					EXPECT().
					SyncUpdate(gomock.Any(), localData.ID, 10, 1).
					Return(nil)

				suite.localService.
					EXPECT().
					Add(gomock.Any(), serverData.DataType, serverData.Data, "").
					Return(&domain.UserStoredData{ID: 12}, nil)

				suite.localRepository.
					EXPECT().
					SyncUpdate(gomock.Any(), 12, serverData.ID, 1).
					Return(nil)
			},
			err: nil,
		},
		{
			name: "no session",
			prepare: func() {
				suite.session.SetToken("")
			},
			err: nil,
		},
		{
			name: "err getting data from server",
			prepare: func() {
				suite.session.SetToken("some-token")
				suite.serverApi.
					EXPECT().
					GetAll(gomock.Any()).
					Return(nil, domain.ErrInternal)
			},
			err: domain.ErrInternal,
		},
		{
			name: "err getting data from client",
			prepare: func() {
				suite.session.SetToken("some-token")
				suite.serverApi.
					EXPECT().
					GetAll(gomock.Any()).
					Return([]domain.UserStoredData{}, nil)

				suite.localService.
					EXPECT().
					GetAll(gomock.Any()).
					Return(nil, domain.ErrInternal)
			},
			err: domain.ErrInternal,
		},
		{
			name: "delete from client",
			prepare: func() {
				suite.session.SetToken("some-token")

				suite.serverApi.
					EXPECT().
					GetAll(gomock.Any()).
					Return([]domain.UserStoredData{{ID: 5, Version: 1}}, nil)

				suite.localService.
					EXPECT().
					GetAll(gomock.Any()).
					Return([]domain.UserStoredData{{ID: 3, Version: 1}, {ID: 4, Version: 1}, {ID: 5, Version: 1}}, nil)

				suite.localService.
					EXPECT().
					DeleteBatch(gomock.Any(), []int{3, 4}).
					Return(nil)
			},
			err: nil,
		},
		{
			name: "delete from server",
			prepare: func() {
				suite.session.SetToken("some-token")
				suite.session.AddDeleted(5)
				suite.session.AddDeleted(6)

				suite.serverApi.
					EXPECT().
					GetAll(gomock.Any()).
					Return([]domain.UserStoredData{{ID: 5, Version: 1}, {ID: 6, Version: 1}, {ID: 7, Version: 1}}, nil)

				suite.localService.
					EXPECT().
					GetAll(gomock.Any()).
					Return([]domain.UserStoredData{{ID: 7, Version: 1}}, nil)

				suite.serverApi.
					EXPECT().
					DeleteBatch(gomock.Any(), []int{5, 6}).
					Return(nil)
			},
			err: nil,
		},
		{
			name: "add to client",
			prepare: func() {
				suite.session.SetToken("some-token")
				addToClientData := domain.UserStoredData{
					ID:       5,
					Version:  1,
					DataType: domain.TextDataType,
					Data: domain.TextData{
						Text: "Text",
					},
				}

				suite.serverApi.
					EXPECT().
					GetAll(gomock.Any()).
					Return([]domain.UserStoredData{addToClientData, {ID: 6, Version: 1}}, nil)

				suite.localService.
					EXPECT().
					GetAll(gomock.Any()).
					Return([]domain.UserStoredData{{ID: 6, Version: 1}}, nil)

				suite.localService.
					EXPECT().
					Add(gomock.Any(), addToClientData.DataType, addToClientData.Data, "").
					Return(&addToClientData, nil)

				suite.localRepository.
					EXPECT().
					SyncUpdate(gomock.Any(), addToClientData.ID, addToClientData.ID, addToClientData.Version).
					Return(nil)
			},
			err: nil,
		},
		{
			name: "add to server",
			prepare: func() {
				suite.session.SetToken("some-token")
				addToServerData := domain.UserStoredData{
					ID:       -5,
					Version:  1,
					DataType: domain.TextDataType,
					Data: domain.TextData{
						Text: "Text",
					},
				}

				suite.serverApi.
					EXPECT().
					GetAll(gomock.Any()).
					Return([]domain.UserStoredData{}, nil)

				suite.localService.
					EXPECT().
					GetAll(gomock.Any()).
					Return([]domain.UserStoredData{addToServerData}, nil)

				suite.serverApi.
					EXPECT().
					Add(gomock.Any(), addToServerData).
					Return(&domain.UserStoredData{ID: 10, Version: 1}, nil)

				suite.localRepository.
					EXPECT().
					SyncUpdate(gomock.Any(), addToServerData.ID, 10, 1).
					Return(nil)
			},
			err: nil,
		},
		{
			name: "edit on client",
			prepare: func() {
				suite.session.SetToken("some-token")
				editOnClientData := domain.UserStoredData{
					ID:      1,
					Version: 2,
					Data: domain.TextData{
						Text: "123",
					},
					Meta: "123",
				}

				suite.serverApi.
					EXPECT().
					GetAll(gomock.Any()).
					Return([]domain.UserStoredData{editOnClientData}, nil)

				suite.localService.
					EXPECT().
					GetAll(gomock.Any()).
					Return([]domain.UserStoredData{{ID: 1, Version: 1}}, nil)

				suite.localService.
					EXPECT().
					UpdateByID(gomock.Any(), editOnClientData.ID, editOnClientData.Data, editOnClientData.Meta).
					Return(nil, nil)

				suite.localRepository.
					EXPECT().
					SyncUpdate(gomock.Any(), editOnClientData.ID, editOnClientData.ID, editOnClientData.Version).
					Return(nil)
			},
			err: nil,
		},
		{
			name: "edit on server",
			prepare: func() {
				suite.session.SetToken("some-token")
				suite.session.AddEdited(1)
				editOnServer := domain.UserStoredData{
					ID:      1,
					Version: 1,
					Data: domain.TextData{
						Text: "123",
					},
					Meta: "123",
				}

				suite.serverApi.
					EXPECT().
					GetAll(gomock.Any()).
					Return([]domain.UserStoredData{{ID: 1, Version: 1}}, nil)

				suite.localService.
					EXPECT().
					GetAll(gomock.Any()).
					Return([]domain.UserStoredData{editOnServer}, nil)

				suite.serverApi.
					EXPECT().
					UpdateByID(gomock.Any(), editOnServer.ID, editOnServer.Data, editOnServer.Meta).
					Return(&domain.UserStoredData{ID: 1, Version: 2}, nil)

				suite.localRepository.
					EXPECT().
					SyncUpdate(gomock.Any(), editOnServer.ID, 1, 2).
					Return(nil)
			},
			err: nil,
		},
	}

	for _, testCase := range testCases {
		suite.Run(testCase.name, func() {
			testCase.prepare()
			err := suite.syncer.Sync(context.Background())
			suite.Equal(testCase.err, err)
		})
	}
}
