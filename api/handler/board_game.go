package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/esvarez/game-nest-service/api/presenter"
	errs "github.com/esvarez/game-nest-service/internal/error"
	"github.com/esvarez/game-nest-service/internal/web"
	"github.com/esvarez/game-nest-service/service/boardgame/dto"
	"github.com/esvarez/game-nest-service/service/boardgame/service"
)

type BoardGameHandler struct {
	BoardGameService boardgame.UseCase
	log              *logrus.Logger
}

func NewBoardGameHandler(s boardgame.UseCase, l *logrus.Logger) *BoardGameHandler {
	return &BoardGameHandler{
		BoardGameService: s,
		log:              l,
	}
}

const (
	boardGameUrl = "board_game_url"
	boardGameID  = "board_game_id"
)

func (g *BoardGameHandler) getBoardGames() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		games, err := g.BoardGameService.GetAll()
		if err != nil {
			g.log.WithError(err).Error("error trying to retrieve games")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		response := make([]*presenter.BoardGameResponse, len(games))

		for i, game := range games {
			g := &presenter.BoardGameResponse{}
			g.BuildResponse(game)
			response[i] = g
		}

		web.Success(response, http.StatusOK).Send(w)
	})
}

func (g *BoardGameHandler) findBoardGameByUrl() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)

		item, err := g.BoardGameService.FindGameUrl(params[boardGameUrl])
		if err != nil {
			var status web.AppError
			switch {
			case errors.Is(err, errs.ErrItemNotFound):
				status = web.ErrResourceNotFound
			default:
				status = web.ErrInternalServer
			}
			status.Send(w)
			return
		}

		response := &presenter.BoardGameResponse{}
		response.BuildResponse(item)
		web.Success(response, http.StatusOK).Send(w)
	})
}

func (g *BoardGameHandler) createBoardGame() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body := r.Body
		defer body.Close()

		game := &dto.BoardGame{}
		if err := json.NewDecoder(body).Decode(game); err != nil {
			g.log.WithError(err).Error("error trying to decode body game")
			web.ErrInvalidJSON.Send(w)
			return
		}

		if err := g.BoardGameService.Save(game); err != nil {
			var status web.AppError
			switch {
			case errors.Is(err, errs.ErrInvalidEntity):
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

func (g *BoardGameHandler) updateBoardGame() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		body := r.Body
		defer body.Close()

		id := params[boardGameID]
		bGame := &dto.BoardGame{}
		if err := json.NewDecoder(body).Decode(bGame); err != nil {
			g.log.WithError(err).Error("error trying to decode body game")
			web.ErrInvalidJSON.Send(w)
			return
		}

		if err := g.BoardGameService.Update(id, bGame); err != nil {
			var status web.AppError
			switch {
			case errors.Is(err, errs.ErrInvalidEntity):
				status = web.InvalidBody(err)
			default:
				status = web.ErrInternalServer
			}
			g.log.WithError(err).Error("error trying to create game")
			status.Send(w)
			return
		}

		web.Success(nil, http.StatusOK).Send(w)
	})
}

func MakeGameHandler(router *mux.Router, handler *BoardGameHandler) {
	router = router.PathPrefix("/game").Subrouter()

	router.Handle("", handler.getBoardGames()).
		Methods("GET")
	router.Handle("/{board_game_url}", handler.findBoardGameByUrl()).
		Methods("GET")
	router.Handle("", handler.createBoardGame()).
		Methods("POST")
	router.Handle("/{board_game_id}", handler.updateBoardGame()).
		Methods("PUT")
}
