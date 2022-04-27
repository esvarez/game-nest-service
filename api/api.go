package api

import (
	"github.com/esvarez/game-nest-service/api/handler"
	"github.com/esvarez/game-nest-service/service/game"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/esvarez/game-nest-service/internal/config"
	"github.com/esvarez/game-nest-service/internal/storage"
)

func Start(conf *config.Configuration, log *logrus.Logger, validate *validator.Validate) {
	var (
		router = mux.NewRouter()
		client = storage.CreateDynamoClient(conf)

		gameClient = storage.NewGameClient(client, log)

		gameService = game.NewService(gameClient, log, validate)

		gameHandler = handler.NewGameHandler(gameService, log)
	)

	router = router.PathPrefix("/api/v1").Subrouter()

	handler.MakeGameHandler(router, gameHandler)

	server := newServer(*conf.HTTPPort, router)
	server.Start()
}
