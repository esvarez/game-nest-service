package storage

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	storage "github.com/esvarez/game-nest-service/infrastructure/storage/entity"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/sirupsen/logrus"

	"github.com/esvarez/game-nest-service/internal/dto"
	"github.com/esvarez/game-nest-service/internal/model"
	errs "github.com/esvarez/game-nest-service/pkg/error"
)

type LibraryStorage struct {
	repo      *Storage
	client    *dynamodb.DynamoDB
	tableName string
	log       *logrus.Logger
	now       func() int64
}

func NewLibraryStorage(tableName string, log *logrus.Logger, repo *Storage, client *dynamodb.DynamoDB) *LibraryStorage {
	return &LibraryStorage{
		repo:      repo,
		client:    client,
		tableName: tableName,
		log:       log,
		now: func() int64 {
			return time.Now().Unix()
		},
	}
}

func (l *LibraryStorage) AddBoardGame(userBoardGame *dto.Library) error {
	usrBoardGame := storage.NewLibraryRecord(userBoardGame)
	usrBoardGame.CreatedAt = l.now()
	usrBoardGame.UpdatedAt = l.now()

	return l.repo.PutItem(usrBoardGame)
}

func (l *LibraryStorage) GetBoardGames(id string) (*model.UserInfo, error) {
	key := storage.GetLibraryKey(id)
	expr, err := expression.NewBuilder().WithKeyCondition(key).Build()
	if err != nil {
		l.log.WithError(err).Error("error building expression")
		return nil, fmt.Errorf("%v: error building expression: %w", err, errs.ErrAWSConfig)
	}

	result, err := l.client.Query(&dynamodb.QueryInput{
		TableName:                 aws.String(l.tableName),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})
	if err != nil {
		return nil, fmt.Errorf("%v: error querying dynamo %w", err, errs.ErrDynamoDB)
	}

	usrInfo := &model.UserInfo{
		BoardGames: make([]*model.BoardGame, len(result.Items)),
	}

	for i, item := range result.Items {
		lib := &storage.LibraryRecord{}
		if err = dynamodbattribute.UnmarshalMap(item, lib); err != nil {
			l.log.WithError(err).Error("error unmarshalling item")
			return nil, fmt.Errorf("%v: error unmarshalling item: %w", err, errs.ErrDynamoDB)
		}
		usrInfo.BoardGames[i] = storage.NewBoardGameFromLibraryRecord(lib)
	}

	return usrInfo, nil
}
