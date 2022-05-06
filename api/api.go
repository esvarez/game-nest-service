package api

import (
	storage2 "github.com/esvarez/game-nest-service/infrastructure/storage"
	"github.com/esvarez/game-nest-service/internal/service"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/esvarez/game-nest-service/api/handler"
	"github.com/esvarez/game-nest-service/config"
)

func Start(conf *config.Configuration, log *logrus.Logger, validate *validator.Validate) {
	var (
		router = mux.NewRouter()
		client = storage2.NewDynamoClient(conf)

		store = storage2.NewStorage(*conf.DynamoDB.Table, client)

		boardGameStore = storage2.NewBoardGameStorage(*conf.DynamoDB.Table, log, store, client)
		userStore      = storage2.NewUserStorage(*conf.DynamoDB.Table, log, store, client)

		boardGameService = service.NewBoardGameService(boardGameStore, log, validate)
		userService      = service.NewUserService(userStore, log, validate)

		gameHandler = handler.NewBoardGameHandler(boardGameService, log)
		userHandler = handler.NewUserHandler(userService, log)
	)

	router = router.PathPrefix("/api/v1").Subrouter()

	handler.MakeGameHandler(router, gameHandler)
	handler.MakeUserHandler(router, userHandler)

	server := newServer(*conf.HTTPPort, router)
	server.Start()
}
