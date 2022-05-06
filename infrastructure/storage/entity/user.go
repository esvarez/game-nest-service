package storage

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"

	"github.com/esvarez/game-nest-service/internal/dto"
	"github.com/esvarez/game-nest-service/internal/model"
	"github.com/esvarez/game-nest-service/pkg/uuid"
)

const (
	UserRecordName          = "user"
	UserBoardGameRecordName = "userBoardGame"
	userNameField           = "userUsername"
	userEmailField          = "userEmail"
)

type UserRecord struct {
	record
	userRecordFields
	audit
}

type userRecordFields struct {
	User  string `json:"User"`
	Email string `json:"Email"`
}

type UsernameConstraint struct {
	Username string `json:"PK"`
	SK       string `json:"SK"`
}

type EmailConstraint struct {
	Email string `json:"PK"`
	SK    string `json:"SK"`
}

type UserBoardGameRecord struct {
	record
	userBoardGameRecordFields
	audit
}

type userBoardGameRecordFields struct {
	BoardGameID string `json:"BoardGameID"`
}

func NewUserFromRecord(r *UserRecord) *model.User {
	return &model.User{
		ID:    r.ID[strings.Index(r.ID, "#")+1:],
		Email: r.Email,
		User:  r.User,
	}
}

func NewUserRecord(user *dto.User) *UserRecord {
	return &UserRecord{
		record: record{
			ID:         newUserRecordHashKey(),
			SK:         newUSerRecordRangeKey(),
			RecordType: UserRecordName,
			Version:    0,
		},
		userRecordFields: userRecordFields{
			User:  user.User,
			Email: user.Email,
		},
	}
}

func GetUserKey(id string) map[string]*dynamodb.AttributeValue {
	return map[string]*dynamodb.AttributeValue{
		"PK": {
			S: aws.String(UserRecordName + "#" + id),
		},
		"SK": {
			S: aws.String(UserRecordName),
		},
	}
}

func GetUserGamesKey(id string) expression.KeyConditionBuilder {
	return expression.Key("PK").Equal(expression.Value(UserRecordName + "#" + id)).
		And(expression.Key("SK").BeginsWith(boardGameNameField))
}

func newUserRecordHashKey() string {
	return UserRecordName + "#" + uuid.NewID().String()
}

func newUSerRecordRangeKey() string {
	return UserRecordName
}

func NewUserBoardGameRecord(userBoardGame *dto.UserBoardGame) *UserBoardGameRecord {
	return &UserBoardGameRecord{
		record: record{
			ID:         UserBoardGameRecordName + "#" + userBoardGame.UserID,
			SK:         boardGameNameField + "#" + userBoardGame.BoardGameName,
			RecordType: UserBoardGameRecordName,
			Version:    0,
		},
		userBoardGameRecordFields: userBoardGameRecordFields{
			BoardGameID: userBoardGame.BoardGameID,
		},
	}
}

func NewUsernameConstraint(username string) *UsernameConstraint {
	return &UsernameConstraint{
		Username: userNameField + "#" + username,
		SK:       userNameField,
	}
}

func NewEmailConstraint(email string) *EmailConstraint {
	return &EmailConstraint{
		Email: userEmailField + "#" + email,
		SK:    userEmailField,
	}
}
