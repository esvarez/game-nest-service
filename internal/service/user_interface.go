package service

import (
	"github.com/esvarez/game-nest-service/internal/dto"
	"github.com/esvarez/game-nest-service/internal/entity"
)

type UserReader interface {
	Get() ([]*entity.User, error)
	Find(id string) (*entity.User, error)
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

type UseCase interface {
	Get() ([]*entity.User, error)
	Find(id string) (*entity.User, error)
	Create(user *dto.User) error
	Update(id string, user *dto.User) error
	Delete(id string) error
}
