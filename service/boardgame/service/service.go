package boardgame

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"

	errs "github.com/esvarez/game-nest-service/internal/error"
	"github.com/esvarez/game-nest-service/service/boardgame/dto"
	"github.com/esvarez/game-nest-service/service/boardgame/entity"
)

type Service struct {
	repo      Repository
	log       *logrus.Logger
	validator *validator.Validate
}

func NewService(r Repository, l *logrus.Logger, v *validator.Validate) *Service {
	return &Service{
		repo:      r,
		log:       l,
		validator: v,
	}
}

func (s *Service) Save(bGame *dto.BoardGame) error {
	log := s.log.WithField("method", "[boardgame] Save")
	if err := s.validator.Struct(bGame); err != nil {
		log.WithError(err).Error("validation error")
		return fmt.Errorf("%v: game item %w", err, errs.ErrInvalidEntity)
	}

	return s.repo.Set(bGame)
}

func (s *Service) GetAll() ([]*entity.BoardGame, error) {
	return s.repo.GetAll()
}

func (s *Service) Find(pk string) (*entity.BoardGame, error) {
	return s.repo.Find(pk)
}

func (s *Service) FindGameUrl(url string) (*entity.BoardGame, error) {
	return s.repo.FindByUrl(url)
}

func (s *Service) Update(id string, data *dto.BoardGame) error {
	log := s.log.WithField("method", "[boardgame] Update")
	if err := s.validator.Struct(data); err != nil {
		log.WithError(err).Error("validation error")
		return fmt.Errorf("%v: game item %w", err, errs.ErrInvalidEntity)
	}
	return s.repo.Update(id, data)
}

func (s *Service) Delete(id string) error {
	return s.repo.Delete(id)
}
