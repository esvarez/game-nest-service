package handler

import (
	"encoding/json"
	"errors"
	"github.com/esvarez/game-nest-service/api/presenter"
	"github.com/esvarez/game-nest-service/dto"
	"github.com/esvarez/game-nest-service/internal/entity"
	"github.com/esvarez/game-nest-service/internal/web"
	"github.com/esvarez/game-nest-service/src/game"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
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

const (
	gameName = "game_name"
)

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

func (g *GameHandler) findGameByName() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)

		item, err := g.GameService.Find(params[gameName])
		if err != nil {
			var status web.AppError
			switch {
			case errors.Is(err, entity.ErrItemNotFound):
				status = web.ErrResourceNotFound
			default:
				status = web.ErrInternalServer
			}
			status.Send(w)
			return
		}

		response := &presenter.GameResponse{}
		response.BuildResponse(item)
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
	router.Handle("/{game_name}", handler.findGameByName()).
		Methods("GET")
	router.Handle("", handler.createGame()).
		Methods("POST")
}
