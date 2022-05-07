package dto

type Library struct {
	UserID        string `json:"user_id" validate:"required"`
	BoardGameID   string `json:"board_game_id" validate:"required"`
	BoardGameName string `json:"board_game_name" validate:"required"`
}
