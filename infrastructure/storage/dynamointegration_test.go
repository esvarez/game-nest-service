package storage

import (
	"flag"
	"github.com/esvarez/game-nest-service/config"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

const region = "eu-west-1"
const url = "http://localhost:4566"

type dynamoStorage interface {
	BoardGameStorage | UserStorage | LibraryStorage
}

func createLocalTable(t *testing.T) (name string) {
	sess, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	if err != nil {
		t.Fatalf("failed to create test db session: %v", err)
		return
	}
	name = uuid.New().String()
	client := dynamodb.New(sess)
	client.Endpoint = url
	_, err = client.CreateTable(&dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("PK"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("SK"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("Url"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("PK"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("SK"),
				KeyType:       aws.String("RANGE"),
			},
		},
		GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{
			{
				IndexName: aws.String("SKIndex"),
				KeySchema: []*dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String("SK"),
						KeyType:       aws.String("HASH"),
					},
				},
				Projection: &dynamodb.Projection{
					ProjectionType: aws.String("ALL"),
				},
			},
			{
				IndexName: aws.String("UrlIndex"),
				KeySchema: []*dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String("Url"),
						KeyType:       aws.String("HASH"),
					},
				},
				Projection: &dynamodb.Projection{
					ProjectionType: aws.String("ALL"),
				},
			},
		},
		BillingMode: aws.String(dynamodb.BillingModePayPerRequest),
		TableName:   aws.String(name),
	})
	if err != nil {
		t.Fatalf("failed to create local table: %v", err)
	}
	return
}

func deleteLocalTable(t *testing.T, name string) {
	sess, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	if err != nil {
		return
	}
	client := dynamodb.New(sess)
	client.Endpoint = url
	_, err = client.DeleteTable(&dynamodb.DeleteTableInput{
		TableName: aws.String(name),
	})
	if err != nil {
		t.Fatalf("failed to delete table: %v", err)
	}
}

func getConfigFile() *config.Configuration {
	var pathFile string
	flag.StringVar(&pathFile, "public-config-file",
		"./test_file/config.yml", "Path to public config file")
	v := validator.New()
	return config.LoadConfiguration(pathFile, v)
}
