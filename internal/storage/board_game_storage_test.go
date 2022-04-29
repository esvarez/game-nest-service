package storage

import (
	"errors"
	"testing"

	"github.com/esvarez/game-nest-service/config"
	errs "github.com/esvarez/game-nest-service/internal/error"
	"github.com/esvarez/game-nest-service/internal/logger"
	"github.com/esvarez/game-nest-service/service/boardgame/dto"
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

func TestStorage_FindIntegration(t *testing.T) {
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
	bGames, err := bgs.GetAll()
	if err != nil {
		t.Errorf("failed to get all boardgames: %v", err)
	}

	id := bGames[0].ID

	bg = &dto.BoardGame{
		Name:       "Catan",
		MinPlayers: 3,
		MaxPlayers: 4,
		Duration:   90,
	}
	if err = bgs.Set(bg); err != nil {
		t.Errorf("failed to create boardgame: %v", err)
	}

	game, err := bgs.Find(id)
	if err != nil {
		t.Errorf("failed to find boardgame: %v", err)
	}
	if game.ID != id {
		t.Errorf("expected id %s, got %s", id, game.ID)
	}
	if game.Name != "Scythe" {
		t.Errorf("expected name %s, got %s", "Scythe", game.Name)
	}
	if game.MinPlayers != 1 {
		t.Errorf("expected min players %d, got %d", 1, game.MinPlayers)
	}
	if game.MaxPlayers != 5 {
		t.Errorf("expected max players %d, got %d", 5, game.MaxPlayers)
	}
}

func TestStorage_UpdateIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	name := createLocalTable(t)
	defer deleteLocalTable(t, name)
	bgs := getBGStorage(name)
	bg := &dto.BoardGame{
		Name:       "cite",
		MinPlayers: 1,
		MaxPlayers: 1,
		Duration:   1,
	}
	if err := bgs.Set(bg); err != nil {
		t.Errorf("failed to create boardgame: %v", err)
	}
	bGames, err := bgs.GetAll()
	if err != nil {
		t.Errorf("failed to get all boardgames: %v", err)
	}
	id := bGames[0].ID

	game, err := bgs.Find(id)
	if err != nil {
		t.Errorf("failed to find boardgame: %v", err)
	}
	if game.ID != id {
		t.Errorf("expected id %s, got %s", id, game.ID)
	}
	if game.Name != "cite" {
		t.Errorf("expected name %s, got %s", "cite", game.Name)
	}
	if game.MinPlayers != 1 {
		t.Errorf("expected min players %d, got %d", 1, game.MinPlayers)
	}
	if game.MaxPlayers != 1 {
		t.Errorf("expected max players %d, got %d", 5, game.MaxPlayers)
	}

	err = bgs.Update(id, &dto.BoardGame{
		Name:       "Scythe",
		MinPlayers: 1,
		MaxPlayers: 5,
		Duration:   120,
	})
	if err != nil {
		t.Errorf("failed to update boardgame: %v", err)
	}

	game, err = bgs.Find(id)
	if err != nil {
		t.Errorf("failed to find boardgame: %v", err)
	}
	if game.ID != id {
		t.Errorf("expected id %s, got %s", id, game.ID)
	}
	if game.Name != "Scythe" {
		t.Errorf("expected name %s, got %s", "Scythe", game.Name)
	}
	if game.Duration != 120 {
		t.Errorf("expected min players %d, got %d", 1, game.MinPlayers)
	}
	if game.MaxPlayers != 5 {
		t.Errorf("expected max players %d, got %d", 5, game.MaxPlayers)
	}
}

func TestStorage_DeleteIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	name := createLocalTable(t)
	defer deleteLocalTable(t, name)
	bgs := getBGStorage(name)
	bg := &dto.BoardGame{
		Name:       "cite",
		MinPlayers: 1,
		MaxPlayers: 1,
		Duration:   1,
	}
	if err := bgs.Set(bg); err != nil {
		t.Errorf("failed to create boardgame: %v", err)
	}
	bGames, err := bgs.GetAll()
	if err != nil {
		t.Errorf("failed to get all boardgames: %v", err)
	}
	id := bGames[0].ID

	game, err := bgs.Find(id)
	if err != nil {
		t.Errorf("failed to find boardgame: %v", err)
	}
	if game.ID != id {
		t.Errorf("expected id %s, got %s", id, game.ID)
	}
	if game.Name != "cite" {
		t.Errorf("expected name %s, got %s", "cite", game.Name)
	}
	if game.MinPlayers != 1 {
		t.Errorf("expected min players %d, got %d", 1, game.MinPlayers)
	}
	if game.MaxPlayers != 1 {
		t.Errorf("expected max players %d, got %d", 5, game.MaxPlayers)
	}

	err = bgs.Delete(id)
	if err != nil {
		t.Errorf("failed to delete boardgame: %v", err)
	}

	game, err = bgs.Find(id)
	if !errors.Is(err, errs.ErrItemNotFound) {
		t.Errorf("expected nof found, got: %v", err)
	}
}

func TestStorage_FindByUrlIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	name := createLocalTable(t)
	defer deleteLocalTable(t, name)
	bgs := getBGStorage(name)
	bg := &dto.BoardGame{
		Name:       "Scythe game",
		MinPlayers: 1,
		MaxPlayers: 5,
		Duration:   120,
	}
	if err := bgs.Set(bg); err != nil {
		t.Errorf("failed to create boardgame: %v", err)
	}
	bGames, err := bgs.GetAll()
	if err != nil {
		t.Errorf("failed to get all boardgames: %v", err)
	}

	url := bGames[0].Url
	id := bGames[0].ID

	bg = &dto.BoardGame{
		Name:       "Catan",
		MinPlayers: 3,
		MaxPlayers: 4,
		Duration:   90,
	}
	if err = bgs.Set(bg); err != nil {
		t.Errorf("failed to create boardgame: %v", err)
	}

	game, err := bgs.FindByUrl(url)
	if err != nil {
		t.Errorf("failed to find boardgame: %v", err)
	}
	if game.ID != id {
		t.Errorf("expected id %s, got %s", id, game.ID)
	}
	if game.Name != "Scythe game" {
		t.Errorf("expected name %s, got %s", "Scythe", game.Name)
	}
	if game.MinPlayers != 1 {
		t.Errorf("expected min players %d, got %d", 1, game.MinPlayers)
	}
	if game.MaxPlayers != 5 {
		t.Errorf("expected max players %d, got %d", 5, game.MaxPlayers)
	}
}

func getBGStorage(name string) *BoardGameStorage {
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
	return NewBoardGameStorage(l, stor)
}
