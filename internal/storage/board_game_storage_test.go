package storage

import (
	"github.com/esvarez/game-nest-service/internal/logger"
	"github.com/esvarez/game-nest-service/service/boardgame/dto"
	"testing"
)

func TestStorage_SetIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	name := createLocalTable(t)
	defer deleteLocalTable(t, name)
	bgs := getBGStorage(name)
	bg := &dto.BoardGame{
		Name:       "Scythe",
		MinPlayers: 1,
		MaxPlayers: 5,
		Duration:   120,
	}

	if err := bgs.Set(bg); err != nil {
		t.Errorf("failed to create boardgame: %v", err)
	}
}

func TestStorage_GetAllIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	name := createLocalTable(t)
	defer deleteLocalTable(t, name)
	bgs := getBGStorage(name)
	bg := &dto.BoardGame{
		Name:       "Scythe",
		MinPlayers: 1,
		MaxPlayers: 5,
		Duration:   120,
	}
	if err := bgs.Set(bg); err != nil {
		t.Errorf("failed to create boardgame: %v", err)
	}
	bg = &dto.BoardGame{
		Name:       "Catan",
		MinPlayers: 3,
		MaxPlayers: 4,
		Duration:   90,
	}
	if err := bgs.Set(bg); err != nil {
		t.Errorf("failed to create boardgame: %v", err)
	}

	bGames, err := bgs.GetAll()
	if err != nil {
		t.Errorf("failed to get all boardgames: %v", err)
	}

	if len(bGames) != 2 {
		t.Errorf("expected 2 boardgame, got %d", len(bGames))
	}
}

func getBGStorage(name string) *BoardGameStorage {
	conf := getConfigFile()
	conf.DynamoDB.Table = &name
	var (
		client = NewDynamoClient(conf)
		l      = logger.CreateLogger(conf)
		stor   = NewStorage(*conf.DynamoDB.Table, client)
	)
	return NewBoardGameStorage(l, stor)
}
