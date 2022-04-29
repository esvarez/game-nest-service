package service

import (
	"github.com/esvarez/game-nest-service/service/boardgame/dto"
)

type Reader interface {
	//Find(key string) (*entity.Game, error)
	///GetAll() ([]*entity.Game, error)
}

type Writer interface {
	Set(item *dto.BoardGame) error
	//Update(item *entity.Game) error
	//Delete(key string) error
}

type Repository interface {
	Reader
	Writer
}

type UseCase interface {
	/*Save(item *dto.Game) error
	GetAll() ([]*entity.Game, error)
	Find(key string) (*entity.Game, error)
	Update(key, value *entity.Game) error
	Delete(key string) error*/
}
