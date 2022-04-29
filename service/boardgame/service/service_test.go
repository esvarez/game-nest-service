package boardgame

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	errs "github.com/esvarez/game-nest-service/internal/error"
	"github.com/esvarez/game-nest-service/service/boardgame/dto"
	"github.com/esvarez/game-nest-service/service/boardgame/entity"
	"github.com/esvarez/game-nest-service/service/boardgame/service/mocks"
)

//go:generate mockery --name Repository --dir ./service/boardgame/service --outpkg mocks --output ./service/boardgame/service/mocks --case=underscore

func TestService_Save(t *testing.T) {
	tests := map[string]struct {
		data          *dto.BoardGame
		expectedError error
		mockSetup     func(repo *mocks.Repository)
	}{
		"should save board game": {
			data: &dto.BoardGame{
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
			data:          &dto.BoardGame{},
			expectedError: errs.ErrInvalidEntity,
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
		gamesExpected []*entity.BoardGame
		expectedError error
		mockSetup     func(repo *mocks.Repository)
	}{
		"returns all games": {
			gamesExpected: []*entity.BoardGame{
				{
					ID:          "1",
					Name:        "Catan",
					MinPlayers:  3,
					MaxPlayers:  4,
					Description: "",
					Duration:    120,
				},
				{
					ID:          "2",
					Name:        "Monopoly",
					MinPlayers:  2,
					MaxPlayers:  4,
					Description: "",
					Duration:    120,
				},
			},
			mockSetup: func(repo *mocks.Repository) {
				repo.On("GetAll").Return([]*entity.BoardGame{
					{
						ID:          "1",
						Name:        "Catan",
						MinPlayers:  3,
						MaxPlayers:  4,
						Description: "",
						Duration:    120,
					},
					{
						ID:          "2",
						Name:        "Monopoly",
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
		gameExpected  *entity.BoardGame
		expectedError error
		partitionKey  string
		mockSetup     func(repo *mocks.Repository)
	}{
		"returns game": {
			gameExpected: &entity.BoardGame{
				ID:          "1",
				Name:        "Catan",
				MinPlayers:  3,
				MaxPlayers:  4,
				Description: "",
				Duration:    120,
			},
			partitionKey: "Catan",
			mockSetup: func(repo *mocks.Repository) {
				repo.On("Find", mock.Anything).Return(&entity.BoardGame{
					ID:          "1",
					Name:        "Catan",
					MinPlayers:  3,
					MaxPlayers:  4,
					Description: "",
					Duration:    120,
				}, nil)
			},
		},
		"game not found": {
			expectedError: errs.ErrItemNotFound,
			partitionKey:  "not a game",
			mockSetup: func(repo *mocks.Repository) {
				repo.On("Find", mock.Anything).Return(nil, errs.ErrItemNotFound)
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

func TestService_Update(t *testing.T) {
	tests := map[string]struct {
		data          *dto.BoardGame
		id            string
		expectedError error
		mockSetup     func(repo *mocks.Repository)
	}{
		"should return invalid validation": {
			data: &dto.BoardGame{
				Name:       "Root",
				MinPlayers: -1,
				MaxPlayers: 4,
			},
			expectedError: errs.ErrInvalidEntity,
		},
		"should update board game": {
			id: "1",
			data: &dto.BoardGame{
				Name:        "Root",
				MinPlayers:  1,
				MaxPlayers:  4,
				Description: "",
				Duration:    120,
			},
			mockSetup: func(repo *mocks.Repository) {
				repo.On("Update", mock.Anything, mock.Anything).Return(nil)
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

			err := svc.Update(tc.id, tc.data)
			assert.ErrorIs(t, err, tc.expectedError)
		})
	}
}

func TestService_Delete(t *testing.T) {
	tests := map[string]struct {
		id            string
		expectedError error
		mockSetup     func(repo *mocks.Repository)
	}{
		"should delete board game": {
			id: "1",
			mockSetup: func(repo *mocks.Repository) {
				repo.On("Delete", mock.Anything).Return(nil)
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

			err := svc.Delete(tc.id)
			assert.ErrorIs(t, err, tc.expectedError)
		})
	}
}
