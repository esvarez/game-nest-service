package storage

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/sirupsen/logrus"

	"github.com/esvarez/game-nest-service/infrastructure/storage/entity"
	"github.com/esvarez/game-nest-service/internal/dto"
	"github.com/esvarez/game-nest-service/internal/model"
	errs "github.com/esvarez/game-nest-service/pkg/error"
)

type BoardGameStorage struct {
	repo      *Storage
	client    *dynamodb.DynamoDB
	tableName string
	log       *logrus.Logger
	now       func() int64
}

func NewBoardGameStorage(t string, l *logrus.Logger, s *Storage, c *dynamodb.DynamoDB) *BoardGameStorage {
	return &BoardGameStorage{
		repo:      s,
		log:       l,
		client:    c,
		tableName: t,
		now: func() int64 {
			return time.Now().Unix()
		},
	}
}

func (g *BoardGameStorage) Set(item *dto.BoardGame) error {
	bg := storage.NewBoardGameRecord(item)
	bg.CreatedAt = g.now()
	bg.UpdatedAt = g.now()
	nc := storage.NewNameConstraint(bg)

	avBG, err := dynamodbattribute.MarshalMap(bg)
	if err != nil {
		return fmt.Errorf("failed to marshal board game: %w", errs.ErrEntityMarshal)
	}
	avNC, err := dynamodbattribute.MarshalMap(nc)
	if err != nil {
		return fmt.Errorf("failed to marshal name constraint: %v", errs.ErrEntityMarshal)
	}

	_, err = g.client.TransactWriteItems(&dynamodb.TransactWriteItemsInput{
		TransactItems: []*dynamodb.TransactWriteItem{
			{
				Put: &dynamodb.Put{
					TableName:           aws.String(g.tableName),
					Item:                avBG,
					ConditionExpression: aws.String(conditionPKNotExist),
				},
			},
			{
				Put: &dynamodb.Put{
					TableName:           aws.String(g.tableName),
					Item:                avNC,
					ConditionExpression: aws.String(conditionPKNotExist),
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to put board game and name constraint: %w", errs.ErrFailTransaction)
	}
	return nil
}

func (g *BoardGameStorage) GetAll() ([]*model.BoardGame, error) {
	key := expression.Key("SK").Equal(expression.Value(storage.BoardGameRecordName))

	result, err := g.repo.Query(key, SKIndex)
	if err != nil {
		g.log.WithError(err).Error("error querying board game records")
		return nil, fmt.Errorf("%v: error querying: %w", err, errs.ErrAWSConfig)
	}
	games := make([]*model.BoardGame, len(result.Items))
	if len(games) == 0 {
		g.log.Warn("No games found")
		return games, nil
	}

	for i, item := range result.Items {
		bg := &storage.BoardGameRecord{}
		if err = dynamodbattribute.UnmarshalMap(item, bg); err != nil {
			g.log.WithError(err).Error("error unmarshalling game entity")
			return nil, fmt.Errorf("%v: error unmarshalling game entity %w", err, errs.ErrEntityUnmarshal)
		}
		games[i] = storage.NewBoardGameFromRecord(bg)
	}
	return games, nil
}

func (g *BoardGameStorage) Find(id string) (*model.BoardGame, error) {
	pk := storage.BoardGameRecordName + "#" + id
	sk := storage.BoardGameRecordName

	record, err := getItem[storage.BoardGameRecord](pk, sk, g.tableName, g.client)
	if err != nil {
		g.log.WithError(err).Error("error board game")
		return nil, fmt.Errorf("error getting board game: %w", err)
	}
	return storage.NewBoardGameFromRecord(record), nil
}

func (g *BoardGameStorage) Update(id string, game *dto.BoardGame) error {
	key := storage.GetBoardGameKey(id)

	update := expression.Set(expression.Name("MinPlayers"), expression.Value(game.MinPlayers)).
		Set(expression.Name("MaxPlayers"), expression.Value(game.MaxPlayers)).
		Set(expression.Name("Description"), expression.Value(game.Description)).
		Set(expression.Name("Duration"), expression.Value(game.Duration)).
		Set(expression.Name("UpdatedAt"), expression.Value(g.now()))

	return g.repo.UpdateItem(key, update)
}

func (g *BoardGameStorage) Delete(id string) error {
	pk := storage.BoardGameRecordName + "#" + id
	sk := storage.BoardGameRecordName
	f := expression.Name("PK").Equal(expression.Value(pk)).
		And(expression.Name("SK").Equal(expression.Value(sk)))
	expr, err := expression.NewBuilder().WithCondition(f).Build()
	if err != nil {
		g.log.WithError(err).Error("error building expression")
		return fmt.Errorf("%v: error building expression %w", err, errs.ErrAWSConfig)
	}

	return g.repo.DeleteItem(pk, sk, expr)
}

func (g *BoardGameStorage) FindByUrl(url string) (*model.BoardGame, error) {
	key := expression.Key("Url").Equal(expression.Value(url))

	result, err := g.repo.Query(key, UrlIndex)
	if err != nil {
		g.log.WithError(err).Error("error getting board game")
		return nil, fmt.Errorf("%v: error getting board game %w", err, errs.ErrAWSConfig)
	}
	if len(result.Items) == 0 {
		g.log.Warn("No games found")
		return nil, fmt.Errorf("no baord game found: %w", errs.ErrItemNotFound)
	}

	bg := &storage.BoardGameRecord{}
	if err = dynamodbattribute.UnmarshalMap(result.Items[0], bg); err != nil {
		g.log.WithError(err).Error("error unmarshalling game entity")
		return nil, fmt.Errorf("%v: error unmarshalling game entity %w", err, errs.ErrEntityUnmarshal)
	}
	return storage.NewBoardGameFromRecord(bg), nil
}
