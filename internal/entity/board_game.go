package entity

type BoardGame struct {
	ID          string
	Name        string
	Url         string
	MinPlayers  int8
	MaxPlayers  int8
	Description string
	Duration    uint32
	UpdatedAt   int64
	CreatedAt   int64
}

/*
func (g *Game) Create(data *dto.Game) {
	g.PK = NewID().String()
	g.Name = data.Name
	g.Url = strings.Replace(strings.Trim(strings.ToLower(data.Name), " "), " ", "-", -1)
	g.MinPlayers = data.MinPlayers
	g.MaxPlayers = data.MaxPlayers
	g.Description = data.Description
	g.Duration = data.Duration
	g.CreatedAt = time.Now().Unix()
	g.UpdatedAt = time.Now().Unix()
}

func (g *Game) FormatKeys() {
	g.PK = g.PK[strings.Index(g.PK, "#")+1:]
	g.SK = g.SK[strings.Index(g.SK, "#")+1:]
}
*/
