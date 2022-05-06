package service

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"

	"github.com/esvarez/game-nest-service/internal/dto"
	errs "github.com/esvarez/game-nest-service/pkg/error"
)

type BoardGameService struct {
	repo      BoardGameRepository
	log       *logrus.Logger
	validator *validator.Validate
}

func NewBoardGameService(r BoardGameRepository, l *logrus.Logger, v *validator.Validate) *BoardGameService {
	return &BoardGameService{
		repo:      r,
		log:       l,
		validator: v,
	}
}

func (s *BoardGameService) Save(bGame *dto.BoardGame) error {
	log := s.log.WithField("method", "[boardgame] Save")
	if err := s.validator.Struct(bGame); err != nil {
		log.WithError(err).Error("validation error")
		return fmt.Errorf("%v: game item %w", err, errs.ErrInvalidEntity)
	}

	return s.repo.Set(bGame)
}

func (s *BoardGameService) GetAll() ([]*model.BoardGame, error) {
	return s.repo.GetAll()
}

func (s *BoardGameService) Find(pk string) (*model.BoardGame, error) {
	return s.repo.Find(pk)
}

func (s *BoardGameService) FindGameUrl(url string) (*model.BoardGame, error) {
	return s.repo.FindByUrl(url)
}

func (s *BoardGameService) Update(id string, data *dto.BoardGame) error {
	log := s.log.WithField("method", "[boardgame] Update")
	if err := s.validator.Struct(data); err != nil {
		log.WithError(err).Error("validation error")
		return fmt.Errorf("%v: game item %w", err, errs.ErrInvalidEntity)
	}
	return s.repo.Update(id, data)
}

func (s *BoardGameService) Delete(id string) error {
	return s.repo.Delete(id)
}
