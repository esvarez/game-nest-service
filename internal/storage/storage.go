package storage

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"

	errs "github.com/esvarez/game-nest-service/internal/error"
	storage "github.com/esvarez/game-nest-service/internal/storage/entity"
)

type Record interface {
	storage.BoardGameRecord | storage.UserRecord
}

const (
	SKIndex  = "SKIndex"
	UrlIndex = "UrlIndex"
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
		return fmt.Errorf("%v: error marshalling game entity %w", err, errs.ErrEntityMarshal)
	}

	_, err = r.Client.PutItem(&dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(r.TableName),
	})
	if err != nil {
		return fmt.Errorf("%v: error putting item %w", err, errs.ErrDynamoDB)
	}
	return nil
}

func (r *Storage) Query(expr expression.Expression, index string) (*dynamodb.QueryOutput, error) {
	result, err := r.Client.Query(&dynamodb.QueryInput{
		TableName:                 aws.String(r.TableName),
		IndexName:                 aws.String(index),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})
	if err != nil {
		return nil, fmt.Errorf("%v: error querying dynamo %w", err, errs.ErrDynamoDB)
	}
	return result, nil
}

func (r *Storage) UpdateItem(pk, sk string, expr expression.Expression) error {
	_, err := r.Client.UpdateItem(&dynamodb.UpdateItemInput{
		TableName: aws.String(r.TableName),
		Key: map[string]*dynamodb.AttributeValue{
			"PK": {
				S: aws.String(pk),
			},
			"SK": {
				S: aws.String(sk),
			},
		},
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		UpdateExpression:          expr.Update(),
	})
	if err != nil {
		return fmt.Errorf("%v: error updating item %w", err, errs.ErrDynamoDB)
	}
	return nil
}

func (r *Storage) DeleteItem(pk, sk string, expr expression.Expression) error {
	_, err := r.Client.DeleteItem(&dynamodb.DeleteItemInput{
		TableName:                 aws.String(r.TableName),
		ConditionExpression:       expr.Condition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
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
		return fmt.Errorf("%v: error deleting item %w", err, errs.ErrDynamoDB)
	}
	return nil
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
		return nil, fmt.Errorf("%v: error getting item %w", err, errs.ErrDynamoDB)
	}
	if len(result.Item) == 0 {
		return nil, fmt.Errorf("item not found %w", errs.ErrItemNotFound)
	}

	record := new(T)
	if err = dynamodbattribute.UnmarshalMap(result.Item, record); err != nil {
		return nil, fmt.Errorf("%v: error unmarshalling item entity %w", err, errs.ErrEntityUnmarshal)
	}
	return record, nil
}
