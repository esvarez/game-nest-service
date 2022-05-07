package service

import (
	"github.com/esvarez/game-nest-service/internal/dto"
	"github.com/esvarez/game-nest-service/internal/model"
	"github.com/esvarez/game-nest-service/internal/service/mocks"
	errs "github.com/esvarez/game-nest-service/pkg/error"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

//go:generate mockery --name LibraryRepository --dir ./internal/service --outpkg mocks --output ./internal/service/mocks --case=underscore

func TestLibraryService_AddBoardGameToUser(t *testing.T) {
	tests := map[string]struct {
		data          *dto.Library
		expectedError error
		mockSetup     func(repo *mocks.LibraryRepository)
	}{
		"success": {
			data: &dto.Library{
				UserID:        "123",
				BoardGameID:   "321",
				BoardGameName: "Scrabble",
			},
			mockSetup: func(repo *mocks.LibraryRepository) {
				repo.On("AddBoardGame", mock.Anything).Return(nil)
			},
		},
		"return validation error": {
			data:          &dto.Library{},
			expectedError: errs.ErrInvalidEntity,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var (
				repo = &mocks.LibraryRepository{}
				l    = logrus.New()
				v    = validator.New()
				svc  = NewLibraryService(repo, l, v)
			)
			l.SetLevel(logrus.DebugLevel)

			if tc.mockSetup != nil {
				tc.mockSetup(repo)
			}

			err := svc.AddBoardGameToUser(tc.data)
			assert.ErrorIs(t, err, tc.expectedError)
		})
	}
}

func TestLibraryService_GetBoardGamesByUser(t *testing.T) {
	tests := map[string]struct {
		userID        string
		usersExpected *model.UserInfo
		expectedError error
		mockSetup     func(repo *mocks.LibraryRepository)
	}{
		"success return boardgames": {
			userID: "123",
			usersExpected: &model.UserInfo{
				User: model.User{
					ID: "123",
				},
				BoardGames: []*model.BoardGame{
					{
						ID:   "123",
						Name: "Uno",
					},
				},
			},
			mockSetup: func(repo *mocks.LibraryRepository) {
				repo.On("GetBoardGames", mock.Anything).Return(&model.UserInfo{
					User: model.User{
						ID: "123",
					},
					BoardGames: []*model.BoardGame{
						{
							ID:   "123",
							Name: "Uno",
						},
					},
				}, nil)
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var (
				repo = &mocks.LibraryRepository{}
				l    = logrus.New()
				v    = validator.New()
				svc  = NewLibraryService(repo, l, v)
			)
			l.SetLevel(logrus.DebugLevel)

			if tc.mockSetup != nil {
				tc.mockSetup(repo)
			}

			users, err := svc.GetBoardGamesByUser(tc.userID)
			assert.ErrorIs(t, err, tc.expectedError)
			assert.Equal(t, users, tc.usersExpected)
		})
	}
}
