package storage

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/sirupsen/logrus"

	"github.com/esvarez/game-nest-service/entity"
)

type GameClient struct {
	db     DynamoConnection
	client *DynamoClient
	log    *logrus.Logger
}

func NewGameClient(client *DynamoClient) *GameClient {
	return &GameClient{db: client.DB}
}

func (g *GameClient) Set(key string, game *entity.Game) error {
	av, err := dynamodbattribute.MarshalMap(game)
	if err != nil {
		g.log.WithError(err).Error("Error marshalling game entity")
		return fmt.Errorf("%v: error marshalling game entity %w", err, entity.ErrEntityMarshal)
	}

	_, err = g.db.PutItem(&dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(g.client.Table),
	})
	if err != nil {
		g.log.WithError(err).Error("Error putting game item")
		return fmt.Errorf("%v: error putting item %w", err, entity.ErrDynamoDB)
	}
	return nil
}
