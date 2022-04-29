package boardgame

import (
	"github.com/esvarez/game-nest-service/service/boardgame/dto"
	"github.com/esvarez/game-nest-service/service/boardgame/entity"
)

type Reader interface {
	Find(id string) (*entity.BoardGame, error)
	FindByUrl(url string) (*entity.BoardGame, error)
	GetAll() ([]*entity.BoardGame, error)
}

type Writer interface {
	Set(item *dto.BoardGame) error
	Update(id string, game *dto.BoardGame) error
	Delete(id string) error
}

type Repository interface {
	Reader
	Writer
}

type UseCase interface {
	Save(item *dto.BoardGame) error
	GetAll() ([]*entity.BoardGame, error)
	Find(key string) (*entity.BoardGame, error)
	FindGameUrl(url string) (*entity.BoardGame, error)
	Update(key string, value *dto.BoardGame) error
	Delete(key string) error
}
