package service

import (
	"github.com/esvarez/game-nest-service/internal/dto"
	"github.com/esvarez/game-nest-service/internal/model"
)

type UserReader interface {
	Get() ([]*model.User, error)
	Find(id string) (*model.User, error)
}

type UserWriter interface {
	Create(user *dto.User) error
	Update(id string, user *dto.User) error
	Delete(id string) error
}

type UserRepository interface {
	UserReader
	UserWriter
}

type UserUseCase interface {
	Get() ([]*model.User, error)
	Find(id string) (*model.User, error)
	Create(user *dto.User) error
	Update(id string, user *dto.User) error
	Delete(id string) error
}
