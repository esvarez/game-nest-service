package api

import (
	"github.com/esvarez/game-nest-service/service/game"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/esvarez/game-nest-service/internal/config"
	"github.com/esvarez/game-nest-service/internal/storage"
)

func Start(conf *config.Configuration, log *logrus.Logger, validate *validator.Validate) {
	var (
		router     = mux.NewRouter()
		client     = storage.CreateDynamoClient(conf)
		gameClient = storage.NewGameClient(client)
		_          = game.NewService(gameClient, log, validate)
	)

	router = router.PathPrefix("/api/v1").Subrouter()

	server := newServer(*conf.HTTPPort, router)
	server.Start()
}
