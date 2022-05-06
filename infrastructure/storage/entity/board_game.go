package storage

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/esvarez/game-nest-service/internal/dto"
	"github.com/esvarez/game-nest-service/internal/model"
	"github.com/esvarez/game-nest-service/pkg/uuid"
)

const (
	BoardGameRecordName = "boardGame"
	boardGameNameField  = "boardGameName"
)

func NewBoardGameFromRecord(bg *BoardGameRecord) *model.BoardGame {
	return &model.BoardGame{
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

func NewBoardGameRecord(b *dto.BoardGame) BoardGameRecord {
	return BoardGameRecord{
		record: record{
			ID:         newBoardGameRecordHashKey(),
			SK:         newBoardGameRecordRangeKey(),
			RecordType: BoardGameRecordName,
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
	return BoardGameRecordName + "#" + uuid.NewID().String()
}

func newBoardGameRecordRangeKey() string {
	return BoardGameRecordName
}

func GetBoardGameKey(id string) map[string]*dynamodb.AttributeValue {
	return map[string]*dynamodb.AttributeValue{
		"PK": {
			S: aws.String(BoardGameRecordName + "#" + id),
		},
		"SK": {
			S: aws.String(BoardGameRecordName),
		},
	}
}

type BoardGameRecord struct {
	record
	boardGameRecordFields
	audit
}

type boardGameRecordFields struct {
	Name             string `json:"Name"`
	Url              string `json:"Url"`
	MinPlayers       int8   `json:"minPlayers"`
	MaxPlayers       int8   `json:"MaxPlayers"`
	Duration         uint32 `json:"Duration"`
	MinDuration      uint32 `json:"MinDuration"`
	MaxDuration      uint32 `json:"MaxDuration"`
	Description      string `json:"Description"`
	ShortDescription string `json:"ShortDescription"`
}

type NameConstraint struct {
	Name string `json:"PK"`
	SK   string `json:"SK"`
}

func NewNameConstraint(b BoardGameRecord) NameConstraint {
	return NameConstraint{
		Name: boardGameNameField + "#" + strings.ToLower(b.Name),
		SK:   boardGameNameField,
	}
}
