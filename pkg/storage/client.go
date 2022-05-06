package storage

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/esvarez/game-nest-service/config"
)

func NewDynamoClient(conf *config.Configuration) *dynamodb.DynamoDB {
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
