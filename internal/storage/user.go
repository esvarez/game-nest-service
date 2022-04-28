package storage

import "github.com/esvarez/game-nest-service/entity"

func newUserRecord(user *entity.User) userRecord {
	return userRecord{
		record: record{
			ID:         newUserRecordHashKey(""),
			SK:         newUSerRecordRangeKey(),
			RecordType: userRecordName,
			Version:    0,
		},
		userRecordFields: userRecordFields{
			User:  "",
			Email: "",
		},
	}
}

func newUserRecordHashKey(id string) string {
	return pkGame + id
}

func newUSerRecordRangeKey() string {
	return userRecordName
}

type userRecord struct {
	record
	userRecordFields
}

type userRecordFields struct {
	User  string `json:"user"`
	Email string `json:"email"`
}
