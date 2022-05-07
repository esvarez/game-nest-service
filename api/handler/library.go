package handler

import (
	"encoding/json"
	"errors"
	"github.com/esvarez/game-nest-service/api/presenter"
	"github.com/esvarez/game-nest-service/internal/dto"
	errs "github.com/esvarez/game-nest-service/pkg/error"
	"github.com/esvarez/game-nest-service/pkg/web"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/esvarez/game-nest-service/internal/service"
)

type LibraryHandler struct {
	libraryService service.LibraryUseCase
	log            *logrus.Logger
}

func NewLibraryHandler(libraryService service.LibraryUseCase, log *logrus.Logger) *LibraryHandler {
	return &LibraryHandler{
		libraryService: libraryService,
		log:            log,
	}
}

func (l *LibraryHandler) getBoardGames() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			params = mux.Vars(r)
			log    = l.log.WithField("user_id", params[userID])
		)
		usrInfo, err := l.libraryService.GetBoardGamesByUser(params[userID])
		if err != nil {
			log.WithError(err).Error("Error getting board games from user")
			web.ErrInternalServer.Send(w)
			return
		}

		info := presenter.BuildUserInfoResponse(usrInfo)

		web.Success(info, http.StatusOK).Send(w)
	})
}

func (l *LibraryHandler) addBoardGame() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			body   = r.Body
			params = mux.Vars(r)
			log    = l.log.WithField("user_id", params[userID])
		)
		defer body.Close()

		lb := &dto.Library{}
		if err := json.NewDecoder(body).Decode(lb); err != nil {
			log.WithError(err).Error("Error decoding body")
			web.ErrInvalidJSON.Send(w)
			return
		}
		lb.UserID = params[userID]

		if err := l.libraryService.AddBoardGameToUser(lb); err != nil {
			var response web.AppError
			switch {
			case errors.Is(err, errs.ErrInvalidEntity):
				response = web.InvalidBody(err)
			default:
				response = web.ErrInternalServer
			}
			log.WithError(err).Error("Error adding board game to user")
			response.Send(w)
			return
		}

		web.Success(nil, http.StatusCreated).Send(w)
	})
}

func MakeLibraryHandler(r *mux.Router, handler *LibraryHandler) *mux.Router {
	r.Handle("/user/{user_id}/boardgames", handler.getBoardGames()).Methods(http.MethodGet)
	r.Handle("/user/{user_id}/boardgames", handler.addBoardGame()).Methods(http.MethodPost)

	return r
}
