package entity

import "github.com/esvarez/game-nest-service/service/boardgame/entity"

type User struct {
	ID    string
	Email string
	User  string
}

type UserDetails struct {
	User
	BoardGames []*entity.BoardGame
}
