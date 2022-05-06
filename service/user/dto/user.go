package dto

type User struct {
	Email string `json:"email" validate:"required"`
	User  string `json:"user" validate:"required"`
}

type UserBoardGame struct {
	UserID        string `json:"user_id" validate:"required"`
	BoardGameID   string `json:"board_game_id" validate:"required"`
	BoardGameName string `json:"board_game_name" validate:"required"`
}
