package service

import (
	"github.com/esvarez/game-nest-service/internal/dto"
)

type BoardGameReader interface {
	Find(id string) (*model.BoardGame, error)
	FindByUrl(url string) (*model.BoardGame, error)
	GetAll() ([]*model.BoardGame, error)
}

type BoardGameWriter interface {
	Set(item *dto.BoardGame) error
	Update(id string, game *dto.BoardGame) error
	Delete(id string) error
}

type BoardGameRepository interface {
	BoardGameReader
	BoardGameWriter
}

type BoardGameUseCase interface {
	Save(item *dto.BoardGame) error
	GetAll() ([]*model.BoardGame, error)
	Find(key string) (*model.BoardGame, error)
	FindGameUrl(url string) (*model.BoardGame, error)
	Update(key string, value *dto.BoardGame) error
	Delete(key string) error
}
