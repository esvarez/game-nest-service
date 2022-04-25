package storage

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/esvarez/game-nest-service/internal/config"
)

type DynamoConnection interface {
	// ListTables(*dynamodb.ListTablesInput) (*dynamodb.ListTablesOutput, error)
	PutItem(*dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error)
	// GetItem(*dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error)
	// UpdateItem(*dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error)
	// DeleteItem(*dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error)
	// Query(input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error)
	// CreateTable(*dynamodb.CreateTableInput) (*dynamodb.CreateTableOutput, error)
}

type DynamoClient struct {
	DB    DynamoConnection
	Table string
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
