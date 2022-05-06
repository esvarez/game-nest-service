package service

import (
	"github.com/esvarez/game-nest-service/internal/dto"
	"github.com/esvarez/game-nest-service/internal/entity"
)

type BoardGameReader interface {
	Find(id string) (*entity.BoardGame, error)
	FindByUrl(url string) (*entity.BoardGame, error)
	GetAll() ([]*entity.BoardGame, error)
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
	GetAll() ([]*entity.BoardGame, error)
	Find(key string) (*entity.BoardGame, error)
	FindGameUrl(url string) (*entity.BoardGame, error)
	Update(key string, value *dto.BoardGame) error
	Delete(key string) error
}
