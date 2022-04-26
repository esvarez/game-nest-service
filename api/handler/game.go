package handler

import (
	"github.com/esvarez/game-nest-service/api/presenter"
	"github.com/esvarez/game-nest-service/internal/web"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/esvarez/game-nest-service/service/game"
)

type GameHandler struct {
	GameService *game.Service
	log         *logrus.Logger
}

func (g *GameHandler) getGames() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		games, err := g.GameService.GetAll()
		if err != nil {
			g.log.WithError(err).Error("error trying to retrieve games")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		response := make([]*presenter.GameResponse, len(games))

		for i, game := range games {
			g := &presenter.GameResponse{}
			g.BuildResponse(game)
			response[i] = g
		}

		web.Success(response, http.StatusOK).Send(w)
	})
}

func MakeGameHandler(router *mux.Router, handler *GameHandler) {
	router = router.PathPrefix("/game").Subrouter()

	router.Handle("/", handler.getGames()).
		Methods("GET")
}
