package service

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"

	"github.com/esvarez/game-nest-service/internal/dto"
	"github.com/esvarez/game-nest-service/internal/model"
	errs "github.com/esvarez/game-nest-service/pkg/error"
)

type UserService struct {
	repo      UserRepository
	log       *logrus.Logger
	validator *validator.Validate
}

func NewUserService(r UserRepository, l *logrus.Logger, v *validator.Validate) *UserService {
	return &UserService{
		repo:      r,
		log:       l,
		validator: v,
	}
}

func (s *UserService) Get() ([]*model.User, error) {
	return s.repo.Get()
}

func (s *UserService) Find(id string) (*model.User, error) {
	return s.repo.Find(id)
}

func (s *UserService) Create(user *dto.User) error {
	log := s.log.WithField("method", "[User]Create")
	if err := s.validator.Struct(user); err != nil {
		log.WithError(err).Error("Validation error")
		return fmt.Errorf("%v: user item %w", err, errs.ErrInvalidEntity)
	}

	return s.repo.Create(user)
}

func (s *UserService) Update(id string, user *dto.User) error {
	log := s.log.WithField("method", "[User]Update")
	if err := s.validator.Struct(user); err != nil {
		log.WithError(err).Error("Validation error")
		return fmt.Errorf("%v: user item %w", err, errs.ErrInvalidEntity)
	}
	return s.repo.Update(id, user)
}

func (s *UserService) Delete(id string) error {
	return s.repo.Delete(id)
}
