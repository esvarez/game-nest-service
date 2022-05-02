package user

import (
	"fmt"
	errs "github.com/esvarez/game-nest-service/internal/error"
	"github.com/esvarez/game-nest-service/service/user/dto"
	"github.com/esvarez/game-nest-service/service/user/entity"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
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

func (s Service) Get() ([]*entity.User, error) {
	return s.repo.Get()
}

func (s Service) Find(id string) (*entity.User, error) {
	return s.repo.Find(id)
}

func (s Service) Create(user *dto.User) error {
	log := s.log.WithField("method", "[User]Create")
	if err := s.validator.Struct(user); err != nil {
		log.WithError(err).Error("Validation error")
		return fmt.Errorf("%v: user item %w", err, errs.ErrInvalidEntity)
	}

	return s.repo.Create(user)
}

func (s Service) Update(id string, user *dto.User) error {
	log := s.log.WithField("method", "[User]Update")
	if err := s.validator.Struct(user); err != nil {
		log.WithError(err).Error("Validation error")
		return fmt.Errorf("%v: user item %w", err, errs.ErrInvalidEntity)
	}
	return s.repo.Update(id, user)
}

func (s Service) Delete(id string) error {
	return s.repo.Delete(id)
}
