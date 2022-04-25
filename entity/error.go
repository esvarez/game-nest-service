package entity

import "errors"

var ErrAWSSession = errors.New("AWS session error")
var ErrEntityMarshal = errors.New("entity marshal error")
var ErrDynamoDB = errors.New("dynamodb error")
