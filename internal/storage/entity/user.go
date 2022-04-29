package storage

import (
	"github.com/esvarez/game-nest-service/internal/uuid"
	"github.com/esvarez/game-nest-service/service/user/dto"
	"github.com/esvarez/game-nest-service/service/user/entity"
	"strings"
)

func NewUserFromRecord(r *UserRecord) *entity.User {
	return &entity.User{
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

func newUserRecordHashKey() string {
	return UserRecordName + "#" + uuid.NewID().String()
}

func newUSerRecordRangeKey() string {
	return UserRecordName
}

type UserRecord struct {
	record
	userRecordFields
	audit
}

type userRecordFields struct {
	User  string `json:"user"`
	Email string `json:"email"`
}
