package dto

type Game struct {
	Name        string `json:"name" validate:"required"`
	MinPlayers  int8   `json:"min_players" validate:"required"`
	MaxPlayers  int8   `json:"max_players" validate:"required"`
	Description string `json:"description"`
	Duration    uint32 `json:"duration" validate:"required"`
}
