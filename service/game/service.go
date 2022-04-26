package game

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"

	"github.com/esvarez/game-nest-service/dto"
	"github.com/esvarez/game-nest-service/entity"
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

func (s *Service) Save(data *dto.Game) error {
	log := s.log.WithField("method", "Save")
	if err := s.validator.Struct(data); err != nil {
		log.WithError(err).Error("validation error")
		return fmt.Errorf("%v: game item %w", err, entity.ErrInvalidEntity)
	}

	game := &entity.Game{}
	game.Create(data)
	return s.repo.Set(game)
}

func (s *Service) GetAll() ([]*entity.Game, error) {
	return s.repo.GetAll()
}

func (s *Service) Find(pk string) (*entity.Game, error) {
	return s.repo.Find(pk)
}

func (s *Service) Update(data *dto.Game) error {
	return nil
}

func (s *Service) Delete(key string) error {
	return nil
}
