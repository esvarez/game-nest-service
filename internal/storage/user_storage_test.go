package storage

import (
	"errors"
	"testing"

	"github.com/esvarez/game-nest-service/config"
	errs "github.com/esvarez/game-nest-service/internal/error"
	"github.com/esvarez/game-nest-service/internal/logger"
	"github.com/esvarez/game-nest-service/service/user/dto"
)

func TestUser_GetAllIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	name := createLocalTable(t)
	defer deleteLocalTable(t, name)
	us := getUserStorage(name)
	u := &dto.User{
		Email: "test@email.com",
		User:  "test",
	}
	us.Create(u)
	u = &dto.User{
		Email: "test2@email.com",
		User:  "test2",
	}
	us.Create(u)

	users, err := us.Get()
	if err != nil {
		t.Errorf("failed to get users: %v", err)
	}

	if len(users) != 2 {
		t.Errorf("expected 2 users, got %d", len(users))
	}
}

func TestUser_FindIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	name := createLocalTable(t)
	defer deleteLocalTable(t, name)
	us := getUserStorage(name)
	u := &dto.User{
		Email: "test@email.com",
		User:  "test",
	}

	us.Create(u)

	users, _ := us.Get()
	id := users[0].ID

	user, err := us.Find(id)
	if err != nil {
		t.Errorf("failed to find user: %v", err)
	}
	if user.ID != id {
		t.Errorf("expected id %s, got %s", id, user.ID)
	}
	if user.Email != "test@email.com" {
		t.Errorf("expected email %s, got %s", "test@email.com", user.Email)
	}
	if user.User != "test" {
		t.Errorf("expected user %s, got %s", "test", user.User)
	}
}

func TestUser_CreateIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	name := createLocalTable(t)
	defer deleteLocalTable(t, name)
	us := getUserStorage(name)
	u := &dto.User{
		Email: "test@email.com",
		User:  "test",
	}

	if err := us.Create(u); err != nil {
		t.Errorf("faied to create user: %v", err)
	}
}

func TestUser_UpdateIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	name := createLocalTable(t)
	defer deleteLocalTable(t, name)
	us := getUserStorage(name)
	u := &dto.User{
		Email: "test@email.com",
		User:  "test",
	}

	us.Create(u)

	users, _ := us.Get()
	id := users[0].ID

	u = &dto.User{
		Email: "nuevo@mail.com",
		User:  "nuevo",
	}
	if err := us.Update(id, u); err != nil {
		t.Errorf("failed to update user: %v", err)
	}
	user, _ := us.Find(id)
	if user.Email != "nuevo@mail.com" {
		t.Errorf("expected email %s, got %s", "nuevo@mail.com", user.Email)
	}
	if user.User != "nuevo" {
		t.Errorf("expected user %s, got %s", "nuevo", user.User)
	}

}

func TestUser_DeleteIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	name := createLocalTable(t)
	defer deleteLocalTable(t, name)
	us := getUserStorage(name)
	u := &dto.User{
		Email: "test@email.com",
		User:  "test",
	}

	us.Create(u)

	users, _ := us.Get()
	id := users[0].ID

	if err := us.Delete(id); err != nil {
		t.Errorf("failed to delete user: %v", err)
	}
	_, err := us.Find(id)
	if !errors.Is(err, errs.ErrItemNotFound) {
		t.Errorf("expected nof found, got: %v", err)
	}
}

func getUserStorage(name string) *UserStorage {
	var (
		ep   = "http://localhost:4566"
		r    = "us-east-1"
		ll   = "DEBUG"
		conf = &config.Configuration{
			AWS: &config.AWS{Region: &r},
			DynamoDB: &config.DynamoDB{
				Endpoint: &ep,
				Table:    &name,
			},
			LogLevel: &ll,
		}
		client = NewDynamoClient(conf)
		l      = logger.CreateLogger(conf)
		stor   = NewStorage(*conf.DynamoDB.Table, client)
	)
	return NewUserStorage(l, stor)
}
