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

	conditionPKNotExist = "attribute_not_exists(PK)"
)

type Storage struct {
	tableName string
	client    *dynamodb.DynamoDB
}

func NewStorage(tableName string, client *dynamodb.DynamoDB) *Storage {
	return &Storage{
		tableName: tableName,
		client:    client,
	}
}

func (r *Storage) PutItem(item any) error {
	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("%v: error marshalling game entity %w", err, errs.ErrEntityMarshal)
	}

	_, err = r.client.PutItem(&dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(r.tableName),
	})
	if err != nil {
		return fmt.Errorf("%v: error putting item %w", err, errs.ErrDynamoDB)
	}
	return nil
}

func (r *Storage) Query(key expression.KeyConditionBuilder, index string) (*dynamodb.QueryOutput, error) {
	expr, err := expression.NewBuilder().WithKeyCondition(key).Build()
	if err != nil {
		return nil, fmt.Errorf("%v: error building expression: %w", err, errs.ErrAWSConfig)
	}
	result, err := r.client.Query(&dynamodb.QueryInput{
		TableName:                 aws.String(r.tableName),
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

func (r *Storage) UpdateItem(key map[string]*dynamodb.AttributeValue, update expression.UpdateBuilder) error {
	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		return fmt.Errorf("%v: error building expression %w", err, errs.ErrAWSConfig)
	}

	_, err = r.client.UpdateItem(&dynamodb.UpdateItemInput{
		Key:                       key,
		TableName:                 aws.String(r.tableName),
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
	_, err := r.client.DeleteItem(&dynamodb.DeleteItemInput{
		TableName:                 aws.String(r.tableName),
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
