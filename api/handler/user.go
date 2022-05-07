package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/esvarez/game-nest-service/api/presenter"
	"github.com/esvarez/game-nest-service/internal/dto"
	"github.com/esvarez/game-nest-service/internal/service"
	errs "github.com/esvarez/game-nest-service/pkg/error"
	"github.com/esvarez/game-nest-service/pkg/web"
)

type UserHandler struct {
	userService service.UserUseCase
	log         *logrus.Logger
}

const (
	userID = "user_id"
)

func NewUserHandler(u service.UserUseCase, l *logrus.Logger) *UserHandler {
	return &UserHandler{
		userService: u,
		log:         l,
	}
}

func (u *UserHandler) getUsers() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		users, err := u.userService.Get()
		if err != nil {
			u.log.WithError(err).Error("error trying to get users")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		response := make([]*presenter.UserResponse, len(users))

		for i, usr := range users {
			response[i] = presenter.BuildUserResponse(usr)
		}

		web.Success(response, http.StatusOK).Send(w)
	})
}

func (u *UserHandler) findUser() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)

		item, err := u.userService.Find(params[userID])
		if err != nil {
			var response web.AppError
			switch {
			case errors.Is(err, errs.ErrItemNotFound):
				response = web.ErrResourceNotFound
			default:
				response = web.ErrInternalServer
			}
			response.Send(w)
			return
		}

		response := presenter.BuildUserResponse(item)
		web.Success(response, http.StatusOK).Send(w)
	})
}

func (u *UserHandler) createUser() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body := r.Body
		defer body.Close()

		user := &dto.User{}
		if err := json.NewDecoder(body).Decode(user); err != nil {
			u.log.WithError(err).Error("error trying to decode user")
			web.ErrInvalidJSON.Send(w)
			return
		}

		if err := u.userService.Create(user); err != nil {
			var response web.AppError
			switch {
			case errors.Is(err, errs.ErrInvalidEntity):
				response = web.InvalidBody(err)
			default:
				response = web.ErrInternalServer
			}
			u.log.WithError(err).Error("error trying to create user")
			response.Send(w)
			return
		}

		web.Success(nil, http.StatusCreated).Send(w)
	})
}

func (u *UserHandler) updateUser() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		body := r.Body
		defer body.Close()

		id := params[userID]
		usr := &dto.User{}
		if err := json.NewDecoder(body).Decode(usr); err != nil {
			u.log.WithError(err).Error("error trying to decode user")
			web.ErrInvalidJSON.Send(w)
			return
		}

		if err := u.userService.Update(id, usr); err != nil {
			var response web.AppError
			switch {
			case errors.Is(err, errs.ErrInvalidEntity):
				response = web.InvalidBody(err)
			default:
				response = web.ErrInternalServer
			}
			u.log.WithError(err).Error("error trying to update user")
			response.Send(w)
			return
		}

		web.Success(nil, http.StatusOK).Send(w)
	})
}

func MakeUserHandler(r *mux.Router, handler *UserHandler) {
	r.Handle("/user", handler.getUsers()).Methods(http.MethodGet)
	r.Handle("/user", handler.createUser()).Methods(http.MethodPost)
	r.Handle("/user/{user_id}", handler.findUser()).Methods(http.MethodGet)
	r.Handle("/user/{user_id}", handler.updateUser()).Methods(http.MethodPut)
}
