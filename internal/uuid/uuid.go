package uuid

import (
	"encoding/base64"
	"github.com/google/uuid"
	"strings"
)

type ID = uuid.UUID

func NewID() ID {
	return ID(uuid.New())
}

func NewShortID() string {
	u := uuid.New()
	bite, _ := u.MarshalBinary()
	var escaper = strings.NewReplacer("9", "99", "-", "90", "_", "91")
	return escaper.Replace(base64.RawURLEncoding.EncodeToString(bite))
}

//StringToID convert a string to an entity ID
func StringToID(s string) (ID, error) {
	id, err := uuid.Parse(s)
	return ID(id), err
}
