package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/esvarez/game-nest-service/api/handler"
	"github.com/esvarez/game-nest-service/config"
	"github.com/esvarez/game-nest-service/infrastructure/storage"
	"github.com/esvarez/game-nest-service/internal/service"
)

func Start(conf *config.Configuration, log *logrus.Logger, validate *validator.Validate) {
	var (
		router = mux.NewRouter()
		client = storage.NewDynamoClient(conf)

		store = storage.NewStorage(*conf.DynamoDB.Table, client)

		boardGameStore = storage.NewBoardGameStorage(*conf.DynamoDB.Table, log, store, client)
		userStore      = storage.NewUserStorage(*conf.DynamoDB.Table, log, store, client)
		libraryStore   = storage.NewLibraryStorage(*conf.DynamoDB.Table, log, store, client)

		boardGameService = service.NewBoardGameService(boardGameStore, log, validate)
		userService      = service.NewUserService(userStore, log, validate)
		libraryService   = service.NewLibraryService(libraryStore, log, validate)

		gameHandler    = handler.NewBoardGameHandler(boardGameService, log)
		userHandler    = handler.NewUserHandler(userService, log)
		libraryHandler = handler.NewLibraryHandler(libraryService, log)
	)

	router = router.PathPrefix("/api/v1").Subrouter()

	handler.MakeGameHandler(router, gameHandler)
	handler.MakeUserHandler(router, userHandler)
	router = handler.MakeLibraryHandler(router, libraryHandler)

	server := newServer(*conf.HTTPPort, router)
	server.Start()
}
