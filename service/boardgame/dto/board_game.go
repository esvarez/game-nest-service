package dto

type BoardGame struct {
	Name        string `json:"name" validate:"required"`
	MinPlayers  int8   `json:"min_players" validate:"required" min:"1"`
	MaxPlayers  int8   `json:"max_players" validate:"required" min:"1"`
	Description string `json:"description"`
	Duration    uint32 `json:"duration" validate:"required" min:"0"`
}