package service

import (
	"github.com/esvarez/game-nest-service/internal/entity"
	entity2 "github.com/esvarez/game-nest-service/service/boardgame/entity"
	"github.com/esvarez/game-nest-service/src/game/mocks"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/esvarez/game-nest-service/dto"
)

//go:generate mockery --name Repository --dir ./service/game --outpkg mocks --output ./service/game/mocks --case=underscore

func TestService_Save(t *testing.T) {
	tests := map[string]struct {
		data          *dto.Game
		expectedError error
		mockSetup     func(repo *mocks.Repository)
	}{
		"should update satellite": {
			data: &dto.Game{
				Name:        "Catan",
				MinPlayers:  3,
				MaxPlayers:  4,
				Description: "",
				Duration:    120,
			},
			mockSetup: func(repo *mocks.Repository) {
				repo.On("Set", mock.Anything).Return(nil)
			},
		},
		"should return error validation": {
			data:          &dto.Game{},
			expectedError: entity.ErrInvalidEntity,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var (
				repo = &mocks.Repository{}
				l    = logrus.New()
				v    = validator.New()
				svc  = NewService(repo, l, v)
			)
			l.SetLevel(logrus.DebugLevel)

			if tc.mockSetup != nil {
				tc.mockSetup(repo)
			}

			err := svc.Save(tc.data)
			assert.ErrorIs(t, err, tc.expectedError)
		})
	}
}

func TestService_Get(t *testing.T) {
	tests := map[string]struct {
		gamesExpected []*entity2.Game
		expectedError error
		mockSetup     func(repo *mocks.Repository)
	}{
		"returns all games": {
			gamesExpected: []*entity2.Game{
				{
					PK:          "1",
					SK:          "Catan",
					MinPlayers:  3,
					MaxPlayers:  4,
					Description: "",
					Duration:    120,
				},
				{
					PK:          "2",
					SK:          "Monopoly",
					MinPlayers:  2,
					MaxPlayers:  4,
					Description: "",
					Duration:    120,
				},
			},
			mockSetup: func(repo *mocks.Repository) {
				repo.On("GetAll").Return([]*entity2.Game{
					{
						PK:          "1",
						SK:          "Catan",
						MinPlayers:  3,
						MaxPlayers:  4,
						Description: "",
						Duration:    120,
					},
					{
						PK:          "2",
						SK:          "Monopoly",
						MinPlayers:  2,
						MaxPlayers:  4,
						Description: "",
						Duration:    120,
					},
				}, nil)
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var (
				repo = &mocks.Repository{}
				l    = logrus.New()
				v    = validator.New()
				svc  = NewService(repo, l, v)
			)
			l.SetLevel(logrus.DebugLevel)

			if tc.mockSetup != nil {
				tc.mockSetup(repo)
			}

			games, err := svc.GetAll()
			assert.ErrorIs(t, err, tc.expectedError)
			assert.Equal(t, games, tc.gamesExpected)
		})
	}
}

func TestService_Find(t *testing.T) {
	tests := map[string]struct {
		gameExpected  *entity2.Game
		expectedError error
		partitionKey  string
		mockSetup     func(repo *mocks.Repository)
	}{
		"returns game": {
			gameExpected: &entity2.Game{
				PK:          "1",
				SK:          "Catan",
				MinPlayers:  3,
				MaxPlayers:  4,
				Description: "",
				Duration:    120,
			},
			partitionKey: "Catan",
			mockSetup: func(repo *mocks.Repository) {
				repo.On("Find", mock.Anything).Return(&entity2.Game{
					PK:          "1",
					SK:          "Catan",
					MinPlayers:  3,
					MaxPlayers:  4,
					Description: "",
					Duration:    120,
				}, nil)
			},
		},
		"game not found": {
			expectedError: entity.ErrItemNotFound,
			partitionKey:  "not a game",
			mockSetup: func(repo *mocks.Repository) {
				repo.On("Find", mock.Anything).Return(nil, entity.ErrItemNotFound)
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var (
				repo = &mocks.Repository{}
				l    = logrus.New()
				v    = validator.New()
				svc  = NewService(repo, l, v)
			)
			l.SetLevel(logrus.DebugLevel)

			if tc.mockSetup != nil {
				tc.mockSetup(repo)
			}

			game, err := svc.Find(tc.partitionKey)
			assert.ErrorIs(t, err, tc.expectedError)
			assert.Equal(t, game, tc.gameExpected)
		})
	}
}
