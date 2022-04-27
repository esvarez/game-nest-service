package presenter

import "github.com/esvarez/game-nest-service/entity"

type GameResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	MinPlayers  int8   `json:"min_players"`
	MaxPlayers  int8   `json:"max_players"`
	Description string `json:"description"`
	Duration    uint32 `json:"duration"`
}

func (g *GameResponse) BuildResponse(game *entity.Game) {
	g.ID = game.PK
	g.Name = game.Name
	g.MinPlayers = game.MinPlayers
	g.MaxPlayers = game.MaxPlayers
	g.Description = game.Description
	g.Duration = game.Duration
}
