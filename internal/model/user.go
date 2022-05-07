package model

type User struct {
	ID    string
	Email string
	User  string
}

type UserInfo struct {
	User
	BoardGames []*BoardGame
}
