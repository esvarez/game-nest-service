package dto

type BoardGame struct {
	Name        string `json:"name" validate:"required"`
	MinPlayers  int8   `json:"min_players" validate:"required" min:"1"`
	MaxPlayers  int8   `json:"max_players" validate:"required" min:"1"`
	Description string `json:"description"`
	MinDuration uint32 `json:"min_duration"`
	MaxDuration uint32 `json:"max_duration"`
	Duration    uint32 `json:"duration" validate:"required" min:"0"`
}
