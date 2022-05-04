package error

import "errors"

var ErrEntityMarshal = errors.New("entity marshal error")
var ErrEntityUnmarshal = errors.New("entity unmarshal error")
var ErrDynamoDB = errors.New("dynamodb error")
var ErrInvalidEntity = errors.New("invalid entity")
var ErrAWSConfig = errors.New("aws config error")
var ErrItemNotFound = errors.New("item not found")
var ErrFailTransaction = errors.New("fail transaction")
