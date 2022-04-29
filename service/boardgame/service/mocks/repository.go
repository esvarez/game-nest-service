// Code generated by mockery v2.12.0. DO NOT EDIT.

package mocks

import (
	dto "github.com/esvarez/game-nest-service/service/boardgame/dto"
	entity "github.com/esvarez/game-nest-service/service/boardgame/entity"

	mock "github.com/stretchr/testify/mock"

	testing "testing"
)

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

// Delete provides a mock function with given fields: id
func (_m *Repository) Delete(id string) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Find provides a mock function with given fields: id
func (_m *Repository) Find(id string) (*entity.BoardGame, error) {
	ret := _m.Called(id)

	var r0 *entity.BoardGame
	if rf, ok := ret.Get(0).(func(string) *entity.BoardGame); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.BoardGame)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByUrl provides a mock function with given fields: url
func (_m *Repository) FindByUrl(url string) (*entity.BoardGame, error) {
	ret := _m.Called(url)

	var r0 *entity.BoardGame
	if rf, ok := ret.Get(0).(func(string) *entity.BoardGame); ok {
		r0 = rf(url)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.BoardGame)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(url)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAll provides a mock function with given fields:
func (_m *Repository) GetAll() ([]*entity.BoardGame, error) {
	ret := _m.Called()

	var r0 []*entity.BoardGame
	if rf, ok := ret.Get(0).(func() []*entity.BoardGame); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*entity.BoardGame)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Set provides a mock function with given fields: item
func (_m *Repository) Set(item *dto.BoardGame) error {
	ret := _m.Called(item)

	var r0 error
	if rf, ok := ret.Get(0).(func(*dto.BoardGame) error); ok {
		r0 = rf(item)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Update provides a mock function with given fields: id, game
func (_m *Repository) Update(id string, game *dto.BoardGame) error {
	ret := _m.Called(id, game)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, *dto.BoardGame) error); ok {
		r0 = rf(id, game)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewRepository creates a new instance of Repository. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewRepository(t testing.TB) *Repository {
	mock := &Repository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
