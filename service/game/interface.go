package game

import "github.com/esvarez/game-nest-service/entity"

type Reader interface {
	Find(key string) (*entity.Game, error)
	Get() ([]*entity.Game, error)
}

type Writer interface {
	Set(key, value *entity.Game) error
	Delete(key string) error
}

type Repository interface {
	Reader
	Writer
}

type UseCase interface {
	Find(key string) (*entity.Game, error)
	Get() ([]*entity.Game, error)
	Save(key, value *entity.Game) error
	Update(key, value *entity.Game) error
	Delete(key string) error
}
