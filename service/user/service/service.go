package user

import (
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

func (s Service) Get() ([]*entity.User, error) {
	//TODO implement me
	panic("implement me")
}

func (s Service) Find(id string) (*entity.User, error) {
	//TODO implement me
	panic("implement me")
}

func (s Service) Create(user *dto.User) (*entity.User, error) {
	//TODO implement me
	panic("implement me")
}

func (s Service) Update(id string, user *dto.User) (*entity.User, error) {
	//TODO implement me
	panic("implement me")
}

func (s Service) Delete(id string) error {
	//TODO implement me
	panic("implement me")
}

func NewService(r Repository, l *logrus.Logger, v *validator.Validate) *Service {
	return &Service{
		repo:      r,
		log:       l,
		validator: v,
	}
}
