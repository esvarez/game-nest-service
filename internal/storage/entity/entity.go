package storage

import (
	"github.com/esvarez/game-nest-service/service/boardgame/entity"
)

func newUserRecord(user *entity.User) userRecord {
	return userRecord{
		record: record{
			ID:         newUserRecordHashKey(""),
			SK:         newUSerRecordRangeKey(),
			RecordType: UserRecordName,
			Version:    0,
		},
		userRecordFields: userRecordFields{
			User:  "",
			Email: "",
		},
	}
}

func newUserRecordHashKey(id string) string {
	return id
}

func newUSerRecordRangeKey() string {
	return UserRecordName
}

type userRecord struct {
	record
	userRecordFields
	audit
}

type userRecordFields struct {
	User  string `json:"user"`
	Email string `json:"email"`
}
