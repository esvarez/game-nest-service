package storage

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/esvarez/game-nest-service/internal/entity"
)

type Record interface {
	BoardGameRecord
}

const (
	SKIndex = "SKIndex"
)

type Storage struct {
	TableName string
	Client    *dynamodb.DynamoDB
}

func NewStorage(tableName string, client *dynamodb.DynamoDB) *Storage {
	return &Storage{
		TableName: tableName,
		Client:    client,
	}
}

func (r *Storage) PutItem(item any) error {
	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("%v: error marshalling game entity %w", err, entity.ErrEntityMarshal)
	}

	_, err = r.Client.PutItem(&dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(r.TableName),
	})
	if err != nil {
		return fmt.Errorf("%v: error putting item %w", err, entity.ErrDynamoDB)
	}
	return nil
}

func (r *Storage) Query() {

}
