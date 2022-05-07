package service

import (
	"github.com/esvarez/game-nest-service/internal/dto"
	"github.com/esvarez/game-nest-service/internal/model"
)

type LibraryReader interface {
	GetBoardGames(id string) (*model.UserInfo, error)
}

type LibraryWriter interface {
	AddBoardGame(library *dto.Library) error
}

type LibraryRepository interface {
	LibraryReader
	LibraryWriter
}

type LibraryUseCase interface {
	GetBoardGamesByUser(id string) (*model.UserInfo, error)
	AddBoardGameToUser(library *dto.Library) error
}
