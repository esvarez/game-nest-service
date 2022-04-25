package entity

import (
	"github.com/esvarez/game-nest-service/dto"
	"strings"
	"time"
)

type Game struct {
	PK          string `dynamodbav:"PK"`
	SK          string `dynamodbav:"SK"`
	MinPlayers  int8   `dynamodbav:"MinPlayers"`
	MaxPlayers  int8   `dynamodbav:"MaxPlayers"`
	Description string `dynamodbav:"Description"`
	Duration    uint32 `dynamodbav:"Duration"`
	UpdatedAt   int64  `dynamodbav:"UpdatedAt"`
	CreatedAt   int64  `dynamodbav:"CreatedAt"`
}

func (g *Game) Create(data *dto.Game) {
	g.PK = NewID().String()
	g.SK = data.Name
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
