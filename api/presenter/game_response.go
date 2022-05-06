package presenter

type BoardGameResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	MinPlayers  int8   `json:"min_players"`
	MaxPlayers  int8   `json:"max_players"`
	Description string `json:"description"`
	Duration    uint32 `json:"duration"`
}

func (g *BoardGameResponse) BuildResponse(game *model.BoardGame) {
	g.ID = game.ID
	g.Name = game.Name
	g.MinPlayers = game.MinPlayers
	g.MaxPlayers = game.MaxPlayers
	g.Description = game.Description
	g.Duration = game.Duration
}
