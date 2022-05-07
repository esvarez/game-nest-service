package storage

import (
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/esvarez/game-nest-service/internal/dto"
	"github.com/esvarez/game-nest-service/internal/model"
	"strings"
)

type LibraryRecord struct {
	record
	libraryRecordFields
	audit
}

type libraryRecordFields struct {
	BoardGameID string `json:"BoardGameID"`
}

func NewLibraryRecord(userBoardGame *dto.Library) *LibraryRecord {
	return &LibraryRecord{
		record: record{
			ID:         newLibraryRecordHashKey(userBoardGame.UserID),
			SK:         newLibraryRecordRangeKey(userBoardGame.BoardGameName),
			RecordType: UserBoardGameRecordName,
			Version:    0,
		},
		libraryRecordFields: libraryRecordFields{
			BoardGameID: userBoardGame.BoardGameID,
		},
	}
}

func newLibraryRecordHashKey(userID string) string {
	return UserRecordName + "#" + userID
}

func newLibraryRecordRangeKey(boardGameName string) string {
	return boardGameNameField + "#" + boardGameName
}

func GetLibraryKey(id string) expression.KeyConditionBuilder {
	return expression.Key("PK").Equal(expression.Value(UserRecordName + "#" + id)).
		And(expression.Key("SK").BeginsWith(boardGameNameField))
}

func NewBoardGameFromLibraryRecord(l *LibraryRecord) *model.BoardGame {
	return &model.BoardGame{
		ID:   l.BoardGameID[strings.Index(l.BoardGameID, "#")+1:],
		Name: l.SK[strings.Index(l.SK, "#")+1:],
	}
}
