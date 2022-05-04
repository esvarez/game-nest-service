package storage

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/sirupsen/logrus"

	errs "github.com/esvarez/game-nest-service/internal/error"
	"github.com/esvarez/game-nest-service/internal/storage/entity"
	"github.com/esvarez/game-nest-service/service/boardgame/dto"
	"github.com/esvarez/game-nest-service/service/boardgame/entity"
)

type BoardGameStorage struct {
	repo      *Storage
	client    *dynamodb.DynamoDB
	tableName string
	log       *logrus.Logger
	now       func() int64
}

func NewBoardGameStorage(l *logrus.Logger, s *Storage) *BoardGameStorage {
	return &BoardGameStorage{
		repo:      s,
		log:       l,
		client:    s.Client,
		tableName: s.TableName,
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

func (g *BoardGameStorage) GetAll() ([]*entity.BoardGame, error) {
	key := expression.Key("SK").Equal(expression.Value(storage.BoardGameRecordName))

	expr, err := expression.NewBuilder().WithKeyCondition(key).Build()
	if err != nil {
		g.log.WithError(err).Error("error building expression")
		return nil, fmt.Errorf("%v: error building expression: %w", err, errs.ErrAWSConfig)
	}

	result, err := g.repo.Query(expr, SKIndex)
	if err != nil {
		g.log.WithError(err).Error("error querying board game records")
		return nil, fmt.Errorf("%v: error querying: %w", err, errs.ErrAWSConfig)
	}
	games := make([]*entity.BoardGame, len(result.Items))
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

func (g *BoardGameStorage) Find(id string) (*entity.BoardGame, error) {
	pk := storage.BoardGameRecordName + "#" + id
	sk := storage.BoardGameRecordName

	record, err := getItem[storage.BoardGameRecord](pk, sk, g.repo.TableName, g.repo.Client)
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

func (g *BoardGameStorage) FindByUrl(url string) (*entity.BoardGame, error) {

	key := expression.Key("Url").Equal(expression.Value(url))

	expr, err := expression.NewBuilder().WithKeyCondition(key).Build()
	if err != nil {
		g.log.WithError(err).Error("error building expression")
		return nil, fmt.Errorf("%v: error building expression: %w", err, errs.ErrAWSConfig)
	}

	result, err := g.repo.Query(expr, UrlIndex)
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
