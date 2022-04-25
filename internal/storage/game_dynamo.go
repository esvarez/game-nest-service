package storage

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/sirupsen/logrus"

	"github.com/esvarez/game-nest-service/entity"
)

type GameClient struct {
	client *DynamoClient
	log    *logrus.Logger
}

const (
	pkGame = "game#"
	skGame = "name#"
)

func NewGameClient(client *DynamoClient) *GameClient {
	return &GameClient{client: client}
}

func (g *GameClient) Set(game *entity.Game) error {
	game.PK = pkGame + game.PK
	game.SK = skGame + game.SK
	client := g.client
	av, err := dynamodbattribute.MarshalMap(game)
	if err != nil {
		g.log.WithError(err).Error("Error marshalling game entity")
		return fmt.Errorf("%v: error marshalling game entity %w", err, entity.ErrEntityMarshal)
	}

	_, err = client.DB.PutItem(&dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(client.Table),
	})
	if err != nil {
		g.log.WithError(err).Error("Error putting game item")
		return fmt.Errorf("%v: error putting item %w", err, entity.ErrDynamoDB)
	}
	return nil
}

func (g *GameClient) GetAll() ([]*entity.Game, error) {
	client := g.client
	filt := expression.Name("PK").BeginsWith(pkGame).
		And(expression.Name("SK").BeginsWith(skGame))

	expr, err := expression.NewBuilder().WithFilter(filt).Build()
	if err != nil {
		g.log.WithError(err).Error("error building expression")
		return nil, fmt.Errorf("%v: error building expression %w", err, entity.ErrAWSConfig)
	}

	result, err := client.DB.Query(&dynamodb.QueryInput{
		TableName:                aws.String(client.Table),
		KeyConditionExpression:   expr.KeyCondition(),
		ExpressionAttributeNames: expr.Names(),
	})
	if err != nil {
		g.log.WithError(err).Error("error querying dynamo")
		return nil, fmt.Errorf("%v: error querying dynamo %w", err, entity.ErrDynamoDB)
	}

	games := make([]*entity.Game, len(result.Items))
	if len(games) == 0 {
		g.log.Info("No games found")
		return games, nil
	}

	for i, item := range result.Items {
		game := &entity.Game{}
		if err = dynamodbattribute.UnmarshalMap(item, game); err != nil {
			g.log.WithError(err).Error("error unmarshalling game entity")
			return nil, fmt.Errorf("%v: error unmarshalling game entity %w", err, entity.ErrEntityUnmarshal)
		}
		game.FormatKeys()
		games[i] = game
	}
	return games, nil
}
