package storage

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/esvarez/game-nest-service/entity"
	"github.com/sirupsen/logrus"
	"time"
)

type GameClient struct {
	client *DynamoClient
	log    *logrus.Logger
}

const (
	pkGame = "game#"
	skGame = "data"

	userRecordName = "user"
)

func NewGameClient(client *DynamoClient, l *logrus.Logger) *GameClient {
	return &GameClient{client: client, log: l}
}

func (g *GameClient) Set(game *entity.Game) error {
	game.PK = pkGame + game.PK
	game.SK = skGame + game.SK
	return putItem(g.client, game)
}

func (g *GameClient) GetAll() ([]*entity.Game, error) {
	client := g.client
	keys := expression.Key("SK").Equal(expression.Value(skGame))
	// TODO: add pagination
	// TODO Move to a common function

	expr, err := expression.NewBuilder().WithKeyCondition(keys).Build()
	if err != nil {
		g.log.WithError(err).Error("error building expression")
		return nil, fmt.Errorf("%v: error building expression: %w", err, entity.ErrAWSConfig)
	}
	result, err := client.DB.Query(&dynamodb.QueryInput{
		TableName:                 aws.String(client.Table),
		IndexName:                 aws.String(SKIndex),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
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

func (g *GameClient) Find(key string) (*entity.Game, error) {
	key = pkGame + key
	sk := skGame + key
	client := g.client

	// TODO move to a common function
	result, err := client.DB.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(client.Table),
		Key: map[string]*dynamodb.AttributeValue{
			"PK": {
				S: aws.String(key),
			},
			"SK": {
				S: aws.String(sk),
			},
		},
	})
	if err != nil {
		g.log.WithError(err).Error("error getting game item")
		return nil, fmt.Errorf("%v: error getting item %w", err, entity.ErrDynamoDB)
	}
	if len(result.Item) == 0 {
		g.log.Info("No game found")
		return nil, fmt.Errorf("no game found %w", entity.ErrItemNotFound)
	}

	game := &entity.Game{}
	if err = dynamodbattribute.UnmarshalMap(result.Item, game); err != nil {
		g.log.WithError(err).Error("error unmarshalling game entity")
		return nil, fmt.Errorf("%v: error unmarshalling game entity %w", err, entity.ErrEntityUnmarshal)
	}
	return game, nil
}

func (g *GameClient) Update(game *entity.Game) error {
	game.PK = pkGame + game.PK
	game.SK = skGame + game.SK

	update := expression.Set(expression.Name("SK"), expression.Value(game.SK)).
		Set(expression.Name("MinPlayers"), expression.Value(game.MinPlayers)).
		Set(expression.Name("MaxPlayers"), expression.Value(game.MaxPlayers)).
		Set(expression.Name("Description"), expression.Value(game.Description)).
		Set(expression.Name("Duration"), expression.Value(game.Duration)).
		Set(expression.Name("UpdatedAt"), expression.Value(time.Now().Unix()))
	// TODO move to method
	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		g.log.WithError(err).Error("error building expression")
		return fmt.Errorf("%v: error building expression %w", err, entity.ErrAWSConfig)
	}
	_, err = g.client.DB.UpdateItem(&dynamodb.UpdateItemInput{
		TableName: aws.String(g.client.Table),
		Key: map[string]*dynamodb.AttributeValue{
			"PK": {
				S: aws.String(game.PK),
			},
			"SK": {
				S: aws.String(game.SK),
			},
		},
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		UpdateExpression:          expr.Update(),
	})
	if err != nil {
		g.log.WithError(err).Error("error updating game item")
		return fmt.Errorf("%v: error updating item %w", err, entity.ErrDynamoDB)
	}
	return nil
}

func (g *GameClient) Delete(key string) error {
	filt := expression.Name("PK").BeginsWith(pkGame).
		And(expression.Name("SK").BeginsWith(skGame))
	// TODO move to a common function
	expr, err := expression.NewBuilder().WithFilter(filt).Build()
	if err != nil {
		g.log.WithError(err).Error("error building expression")
		return fmt.Errorf("%v: error building expression %w", err, entity.ErrAWSConfig)
	}

	_, err = g.client.DB.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: aws.String(g.client.Table),
		Key: map[string]*dynamodb.AttributeValue{
			"PK": {
				S: aws.String(pkGame + key),
			},
			"SK": {
				S: aws.String(skGame + key),
			},
		},
		ConditionExpression: expr.Condition(),
	})
	if err != nil {
		g.log.WithError(err).Error("error deleting game item")
		return fmt.Errorf("%v: error deleting item %w", err, entity.ErrDynamoDB)
	}
	return nil
}
