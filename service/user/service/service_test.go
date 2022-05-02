package user

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	errs "github.com/esvarez/game-nest-service/internal/error"
	"github.com/esvarez/game-nest-service/service/user/dto"
	"github.com/esvarez/game-nest-service/service/user/entity"
	"github.com/esvarez/game-nest-service/service/user/service/mocks"
)

//go:generate mockery --name Repository --dir ./service/user/service --outpkg mocks --output ./service/user/service/mocks --case=underscore

func TestService_Create(t *testing.T) {
	tests := map[string]struct {
		data          *dto.User
		expectedError error
		mockSetup     func(repo *mocks.Repository)
	}{
		"success": {
			data: &dto.User{
				Email: "mail",
				User:  "user",
			},
			mockSetup: func(repo *mocks.Repository) {
				repo.On("Create", mock.Anything).Return(nil)
			},
		},
		"return validation error": {
			data: &dto.User{
				Email: "",
				User:  "user",
			},
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

			err := svc.Create(tc.data)
			assert.ErrorIs(t, err, tc.expectedError)
		})
	}
}

func TestService_Get(t *testing.T) {
	tests := map[string]struct {
		usersExpected []*entity.User
		expectedError error
		mockSetup     func(repo *mocks.Repository)
	}{
		"success return users": {
			usersExpected: []*entity.User{
				{
					ID:    "1",
					Email: "1@mail.com",
					User:  "uno",
				},
				{
					ID:    "2",
					Email: "other@mail.com",
					User:  "otherUser",
				},
			},
			mockSetup: func(repo *mocks.Repository) {
				repo.On("Get").Return([]*entity.User{
					{
						ID:    "1",
						Email: "1@mail.com",
						User:  "uno",
					},
					{
						ID:    "2",
						Email: "other@mail.com",
						User:  "otherUser",
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

			users, err := svc.Get()
			assert.ErrorIs(t, err, tc.expectedError)
			assert.Equal(t, users, tc.usersExpected)
		})
	}
}

func TestService_Find(t *testing.T) {
	tests := map[string]struct {
		id            string
		userExpected  *entity.User
		expectedError error
		mockSetup     func(repo *mocks.Repository)
	}{
		"success return user": {
			userExpected: &entity.User{
				ID:    "1",
				Email: "user@mail.com",
				User:  "user",
			},
			mockSetup: func(repo *mocks.Repository) {
				repo.On("Find", mock.Anything).Return(&entity.User{
					ID:    "1",
					Email: "user@mail.com",
					User:  "user",
				}, nil)
			},
		},
		"user not found": {
			expectedError: errs.ErrItemNotFound,
			id:            "123",
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

			user, err := svc.Find(tc.id)
			assert.ErrorIs(t, err, tc.expectedError)
			assert.Equal(t, user, tc.userExpected)
		})
	}
}

func TestService_Update(t *testing.T) {
	tests := map[string]struct {
		data          *dto.User
		id            string
		expectedError error
		mockSetup     func(repo *mocks.Repository)
	}{
		"should return validation error": {
			data: &dto.User{
				User: "user",
			},
			expectedError: errs.ErrInvalidEntity,
		},
		"should update user": {
			id: "1",
			data: &dto.User{
				Email: "user@mail.com",
				User:  "user",
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
		"shoud delete user": {
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
