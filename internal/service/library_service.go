package service

import (
	"fmt"
	"github.com/esvarez/game-nest-service/internal/dto"
	"github.com/esvarez/game-nest-service/internal/model"
	errs "github.com/esvarez/game-nest-service/pkg/error"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type LibraryService struct {
	repo      LibraryRepository
	log       *logrus.Logger
	validator *validator.Validate
}

func NewLibraryService(repo LibraryRepository, log *logrus.Logger, validator *validator.Validate) *LibraryService {
	return &LibraryService{
		repo:      repo,
		log:       log,
		validator: validator,
	}
}

func (l *LibraryService) GetBoardGamesByUser(id string) (*model.UserInfo, error) {
	return l.repo.GetBoardGames(id)
}

func (l *LibraryService) AddBoardGameToUser(library *dto.Library) error {
	log := l.log.WithField("method", "AddBoardGameToUser")
	if err := l.validator.Struct(library); err != nil {
		log.WithError(err).Error("validation error")
		return fmt.Errorf("%v: validation error: %w", err, errs.ErrInvalidEntity)
	}

	return l.repo.AddBoardGame(library)
}
