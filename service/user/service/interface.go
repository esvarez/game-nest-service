package user

import (
	"github.com/esvarez/game-nest-service/service/user/dto"
	"github.com/esvarez/game-nest-service/service/user/entity"
)

type Reader interface {
	Get() ([]*entity.User, error)
	Find(id string) (*entity.User, error)
}

type Writer interface {
	Create(user *dto.User) error
	Update(id string, user *dto.User) error
	Delete(id string) error
}

type Repository interface {
	Reader
	Writer
}

type UseCase interface {
	Get() ([]*entity.User, error)
	Find(id string) (*entity.User, error)
	Create(user *dto.User) (*entity.User, error)
	Update(id string, user *dto.User) (*entity.User, error)
	Delete(id string) error
}
