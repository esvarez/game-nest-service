package storage

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"

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

func (r *Storage) QueryIndex(expr expression.Expression, index string) (*dynamodb.QueryOutput, error) {
	result, err := r.Client.Query(&dynamodb.QueryInput{
		TableName:                 aws.String(r.TableName),
		IndexName:                 aws.String(index),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})
	if err != nil {
		return nil, fmt.Errorf("%v: error querying dynamo %w", err, entity.ErrDynamoDB)
	}
	return result, nil
}

func getItem[T Record](pk, sk, table string, client *dynamodb.DynamoDB) (*T, error) {
	result, err := client.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(table),
		Key: map[string]*dynamodb.AttributeValue{
			"PK": {
				S: aws.String(pk),
			},
			"SK": {
				S: aws.String(sk),
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("%v: error getting item %w", err, entity.ErrDynamoDB)
	}
	if len(result.Item) == 0 {
		return nil, fmt.Errorf("item not found %w", entity.ErrItemNotFound)
	}

	game := &T{}
	if err = dynamodbattribute.UnmarshalMap(result.Item, game); err != nil {
		return nil, fmt.Errorf("%v: error unmarshalling item entity %w", err, entity.ErrEntityUnmarshal)
	}
	return game, nil
}
