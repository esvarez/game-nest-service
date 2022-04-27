package handler

import (
	"encoding/json"
	"errors"
	"github.com/esvarez/game-nest-service/api/presenter"
	"github.com/esvarez/game-nest-service/dto"
	"github.com/esvarez/game-nest-service/entity"
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

func NewGameHandler(s *game.Service, l *logrus.Logger) *GameHandler {
	return &GameHandler{
		GameService: s,
		log:         l,
	}
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

func (g *GameHandler) createGame() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body := r.Body
		defer body.Close()

		game := &dto.Game{}
		if err := json.NewDecoder(body).Decode(game); err != nil {
			g.log.WithError(err).Error("error trying to decode body game")
			web.ErrInvalidJSON.Send(w)
			return
		}

		if err := g.GameService.Save(game); err != nil {
			var status web.AppError
			switch {
			case errors.Is(err, entity.ErrInvalidEntity):
				status = web.InvalidBody(err)
			default:
				status = web.ErrInternalServer
			}
			g.log.WithError(err).Error("error trying to create game")
			status.Send(w)
			return
		}

		web.Success(nil, http.StatusCreated).Send(w)
	})
}

func MakeGameHandler(router *mux.Router, handler *GameHandler) {
	router = router.PathPrefix("/game").Subrouter()

	router.Handle("", handler.getGames()).
		Methods("GET")
	router.Handle("", handler.createGame()).
		Methods("POST")
}
