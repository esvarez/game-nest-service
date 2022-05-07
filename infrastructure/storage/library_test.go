package storage

import (
	"github.com/esvarez/game-nest-service/config"
	"github.com/esvarez/game-nest-service/internal/dto"
	"github.com/esvarez/game-nest-service/pkg/logger"
	"testing"
)

func TestUser_AddBoardGame(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	name := createLocalTable(t)
	defer deleteLocalTable(t, name)
	us := getLibraryStorage(name)

	if err := us.AddBoardGame(&dto.Library{
		UserID:        "123",
		BoardGameID:   "321",
		BoardGameName: "Catan",
	}); err != nil {
		t.Errorf("failed to create user board game record: %v", err)
	}
}

func TestUser_GetAllBoardgames(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration tests")
	}

	name := createLocalTable(t)
	defer deleteLocalTable(t, name)
	us := getLibraryStorage(name)

	if err := us.AddBoardGame(&dto.Library{
		UserID:        "123",
		BoardGameID:   "321",
		BoardGameName: "Catan",
	}); err != nil {
		t.Errorf("failed to create user board game record: %v", err)
	}

	if err := us.AddBoardGame(&dto.Library{
		UserID:        "123",
		BoardGameID:   "234",
		BoardGameName: "Cubitos",
	}); err != nil {
		t.Errorf("failed to create user board game record: %v", err)
	}

	usrInfo, err := us.GetBoardGames("123")
	if err != nil {
		t.Errorf("failed to get user board games: %v", err)
	}
	if len(usrInfo.BoardGames) != 2 {
		t.Errorf("expected 2 boardgames got %d", len(usrInfo.BoardGames))
	}
}

func getLibraryStorage(name string) *LibraryStorage {
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
		client  = NewDynamoClient(conf)
		l       = logger.CreateLogger(conf)
		storage = NewStorage(*conf.DynamoDB.Table, client)
	)
	return NewLibraryStorage(name, l, storage, client)
}
