package storage

import (
	"strings"

	uuid "github.com/esvarez/game-nest-service/internal/entity"
	"github.com/esvarez/game-nest-service/service/boardgame/dto"
	"github.com/esvarez/game-nest-service/service/boardgame/entity"
)

func newBoardGameFromRecord(bg *BoardGameRecord) *entity.BoardGame {
	return &entity.BoardGame{
		ID:          bg.ID[strings.Index(bg.ID, "#")+1:],
		Name:        bg.Name,
		Url:         bg.Url,
		MinPlayers:  bg.MinPlayers,
		MaxPlayers:  bg.MaxPlayers,
		Description: bg.Description,
		Duration:    bg.Duration,
		UpdatedAt:   bg.UpdatedAt,
		CreatedAt:   bg.CreatedAt,
	}
}

func newBoardGameRecord(b *dto.BoardGame) BoardGameRecord {
	return BoardGameRecord{
		record: record{
			ID:         newBoardGameRecordHashKey(),
			SK:         newBoardGameRecordRangeKey(),
			RecordType: boardGameRecordName,
			Version:    0,
		},
		boardGameRecordFields: boardGameRecordFields{
			Name:        b.Name,
			Url:         strings.Replace(strings.ToLower(strings.Trim(b.Name, " ")), " ", "-", -1),
			MinPlayers:  b.MinPlayers,
			MaxPlayers:  b.MaxPlayers,
			Duration:    b.Duration,
			Description: b.Description,
		},
	}
}

func newBoardGameRecordHashKey() string {
	return boardGameRecordName + "#" + uuid.NewID().String()
}

func newBoardGameRecordRangeKey() string {
	return boardGameRecordName
}

type BoardGameRecord struct {
	record
	boardGameRecordFields
	audit
}

type boardGameRecordFields struct {
	Name        string `json:"Name"`
	Url         string `json:"Url"`
	MinPlayers  int8   `json:"minPlayers"`
	MaxPlayers  int8   `json:"MaxPlayers"`
	Duration    uint32 `json:"Duration"`
	Description string `json:"Description"`
}
