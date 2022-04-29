package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/esvarez/game-nest-service/api/handler"
	"github.com/esvarez/game-nest-service/config"
	"github.com/esvarez/game-nest-service/internal/storage"
	"github.com/esvarez/game-nest-service/service/boardgame/service"
)

func Start(conf *config.Configuration, log *logrus.Logger, validate *validator.Validate) {
	var (
		router = mux.NewRouter()
		client = storage.NewDynamoClient(conf)

		stor = storage.NewStorage(*conf.DynamoDB.Table, client)

		boardGameStore = storage.NewBoardGameStorage(log, stor)

		// gameService = game.(gameClient, log, validate)
		boardGameService = service.NewService(boardGameStore, log, validate)

		gameHandler = handler.NewGameHandler(boardGameService, log)
	)

	router = router.PathPrefix("/api/v1").Subrouter()

	handler.MakeGameHandler(router, gameHandler)

	server := newServer(*conf.HTTPPort, router)
	server.Start()
}
