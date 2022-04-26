package game

import (
	"github.com/esvarez/game-nest-service/dto"
	"github.com/esvarez/game-nest-service/entity"
)

type Reader interface {
	Find(key string) (*entity.Game, error)
	GetAll() ([]*entity.Game, error)
}

type Writer interface {
	Set(item *entity.Game) error
	Update(item *entity.Game) error
	Delete(key string) error
}

type Repository interface {
	Reader
	Writer
}

type UseCase interface {
	Save(item *dto.Game) error
	GetAll() ([]*entity.Game, error)
	Find(key string) (*entity.Game, error)
	Update(key, value *entity.Game) error
	Delete(key string) error
}
