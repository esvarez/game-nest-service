package storage

import (
	"time"

	"github.com/sirupsen/logrus"

	"github.com/esvarez/game-nest-service/service/user/dto"
	"github.com/esvarez/game-nest-service/service/user/entity"
)

type UserStorage struct {
	repo *Storage
	log  *logrus.Logger
	now  func() int64
}

func NewUserStorage(l *logrus.Logger, r *Storage) *UserStorage {
	return &UserStorage{
		repo: r,
		log:  l,
		now:  func() int64 { return time.Now().Unix() },
	}
}

func (u UserStorage) Get() ([]*entity.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u UserStorage) Find(id string) (*entity.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u UserStorage) Create(user *dto.User) error {
	//TODO implement me
	panic("implement me")
}

func (u UserStorage) Update(id string, user *dto.User) error {
	//TODO implement me
	panic("implement me")
}

func (u UserStorage) Delete(id string) error {
	//TODO implement me
	panic("implement me")
}
