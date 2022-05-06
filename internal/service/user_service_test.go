package service

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/esvarez/game-nest-service/internal/dto"
	"github.com/esvarez/game-nest-service/internal/service/mocks"
	errs "github.com/esvarez/game-nest-service/pkg/error"
)

//go:generate mockery --name UserRepository --dir ./internal/service --outpkg mocks --output ./internal/service/mocks --case=underscore

func TestService_Create(t *testing.T) {
	tests := map[string]struct {
		data          *dto.User
		expectedError error
		mockSetup     func(repo *mocks.UserRepository)
	}{
		"success": {
			data: &dto.User{
				Email: "mail",
				User:  "user",
			},
			mockSetup: func(repo *mocks.UserRepository) {
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
				repo = &mocks.UserRepository{}
				l    = logrus.New()
				v    = validator.New()
				svc  = NewUserService(repo, l, v)
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
		usersExpected []*model.User
		expectedError error
		mockSetup     func(repo *mocks.UserRepository)
	}{
		"success return users": {
			usersExpected: []*model.User{
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
			mockSetup: func(repo *mocks.UserRepository) {
				repo.On("Get").Return([]*model.User{
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
				repo = &mocks.UserRepository{}
				l    = logrus.New()
				v    = validator.New()
				svc  = NewUserService(repo, l, v)
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
		userExpected  *model.User
		expectedError error
		mockSetup     func(repo *mocks.UserRepository)
	}{
		"success return user": {
			userExpected: &model.User{
				ID:    "1",
				Email: "user@mail.com",
				User:  "user",
			},
			mockSetup: func(repo *mocks.UserRepository) {
				repo.On("Find", mock.Anything).Return(&model.User{
					ID:    "1",
					Email: "user@mail.com",
					User:  "user",
				}, nil)
			},
		},
		"user not found": {
			expectedError: errs.ErrItemNotFound,
			id:            "123",
			mockSetup: func(repo *mocks.UserRepository) {
				repo.On("Find", mock.Anything).Return(nil, errs.ErrItemNotFound)
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var (
				repo = &mocks.UserRepository{}
				l    = logrus.New()
				v    = validator.New()
				svc  = NewUserService(repo, l, v)
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
		mockSetup     func(repo *mocks.UserRepository)
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
			mockSetup: func(repo *mocks.UserRepository) {
				repo.On("Update", mock.Anything, mock.Anything).Return(nil)
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var (
				repo = &mocks.UserRepository{}
				l    = logrus.New()
				v    = validator.New()
				svc  = NewUserService(repo, l, v)
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
		mockSetup     func(repo *mocks.UserRepository)
	}{
		"shoud delete user": {
			id: "1",
			mockSetup: func(repo *mocks.UserRepository) {
				repo.On("Delete", mock.Anything).Return(nil)
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var (
				repo = &mocks.UserRepository{}
				l    = logrus.New()
				v    = validator.New()
				svc  = NewUserService(repo, l, v)
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
