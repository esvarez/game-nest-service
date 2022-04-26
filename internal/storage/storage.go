package storage

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/esvarez/game-nest-service/entity"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/esvarez/game-nest-service/internal/config"
)

type DynamoConnection interface {
	Query(input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error)
	GetItem(*dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error)
	PutItem(*dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error)
	UpdateItem(*dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error)
	DeleteItem(*dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error)
	// ListTables(*dynamodb.ListTablesInput) (*dynamodb.ListTablesOutput, error)
	// CreateTable(*dynamodb.CreateTableInput) (*dynamodb.CreateTableOutput, error)
}

type DynamoClient struct {
	DB    DynamoConnection
	Table string
}

type DynamoItem interface {
	*entity.Game
}

func CreateDynamoClient(conf *config.Configuration) *DynamoClient {
	return &DynamoClient{
		DB:    connectToDynamo(conf),
		Table: *conf.DynamoDB.Table}
}

func connectToDynamo(conf *config.Configuration) *dynamodb.DynamoDB {
	session, err := session.NewSession(aws.NewConfig().
		WithRegion(*conf.AWS.Region).
		WithEndpoint(*conf.DynamoDB.Endpoint))
	if err != nil {
		log.Fatal("error creating aws session", err)
	}

	db := dynamodb.New(session)
	input := &dynamodb.ListTablesInput{}
	input.SetLimit(5)
	_, err = db.ListTables(input)
	if err != nil {
		log.Fatal("service unavailable", err)
	}
	return db
}

func putItem[V DynamoItem](client *DynamoClient, item V) error {
	av, err := dynamodbattribute.MarshalMap(item)

	if err != nil {
		return fmt.Errorf("%v: error marshalling game entity %w", err, entity.ErrEntityMarshal)
	}

	_, err = client.DB.PutItem(&dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(client.Table),
	})
	if err != nil {
		return fmt.Errorf("%v: error putting item %w", err, entity.ErrDynamoDB)
	}
	return nil
}
