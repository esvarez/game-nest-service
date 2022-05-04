package storage

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/sirupsen/logrus"

	errs "github.com/esvarez/game-nest-service/internal/error"
	storage "github.com/esvarez/game-nest-service/internal/storage/entity"
	"github.com/esvarez/game-nest-service/service/user/dto"
	"github.com/esvarez/game-nest-service/service/user/entity"
)

type UserStorage struct {
	repo      *Storage
	client    *dynamodb.DynamoDB
	tableName string
	log       *logrus.Logger
	now       func() int64
}

func NewUserStorage(l *logrus.Logger, r *Storage) *UserStorage {
	return &UserStorage{
		repo:      r,
		log:       l,
		client:    r.Client,
		tableName: r.TableName,
		now: func() int64 {
			return time.Now().Unix()
		},
	}
}

func (u UserStorage) Get() ([]*entity.User, error) {
	key := expression.Key("SK").Equal(expression.Value(storage.UserRecordName))
	expr, err := expression.NewBuilder().WithKeyCondition(key).Build()
	if err != nil {
		return nil, fmt.Errorf("%v: error building expression: %w", err, errs.ErrAWSConfig)
	}

	result, err := u.repo.Client.Query(&dynamodb.QueryInput{
		TableName:                 aws.String(u.repo.TableName),
		IndexName:                 aws.String(SKIndex),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})
	if err != nil {
		return nil, fmt.Errorf("%v: error querying dynamo %w", err, errs.ErrDynamoDB)
	}

	games := make([]*entity.User, len(result.Items))
	if len(games) == 0 {
		u.log.Warn("No games found")
		return games, nil
	}

	for i, item := range result.Items {
		us := &storage.UserRecord{}
		if err = dynamodbattribute.UnmarshalMap(item, us); err != nil {
			u.log.WithError(err).Error("error unmarshalling game entity")
			return nil, fmt.Errorf("%v: error unmarshalling game entity %w", err, errs.ErrEntityUnmarshal)
		}
		games[i] = storage.NewUserFromRecord(us)
	}
	return games, nil
}

func (u UserStorage) Find(id string) (*entity.User, error) {
	pk := storage.UserRecordName + "#" + id
	sk := storage.UserRecordName
	rec, err := getItem[storage.UserRecord](pk, sk, u.repo.TableName, u.repo.Client)
	if err != nil {
		u.log.WithError(err).Error("error getting user")
		return nil, fmt.Errorf("error getting user %w", err)
	}
	return storage.NewUserFromRecord(rec), nil
}

func (u UserStorage) Create(user *dto.User) error {
	us := storage.NewUserRecord(user)
	us.CreatedAt = u.now()
	us.UpdatedAt = u.now()

	usrName := storage.NewUsernameConstraint(us.User)
	usrMail := storage.NewEmailConstraint(us.Email)

	avUrs, err := dynamodbattribute.MarshalMap(us)
	if err != nil {
		u.log.WithError(err).Error("error marshalling user")
		return fmt.Errorf("error marshalling user %w", err)
	}
	avUrsName, err := dynamodbattribute.MarshalMap(usrName)
	if err != nil {
		u.log.WithError(err).Error("error marshalling user name")
		return fmt.Errorf("error marshalling user name %w", err)
	}
	avUsrMail, err := dynamodbattribute.MarshalMap(usrMail)
	if err != nil {
		u.log.WithError(err).Error("error marshalling user email")
		return fmt.Errorf("error marshalling user email %w", err)
	}

	_, err = u.client.TransactWriteItems(&dynamodb.TransactWriteItemsInput{
		TransactItems: []*dynamodb.TransactWriteItem{
			{
				Put: &dynamodb.Put{
					TableName:           aws.String(u.tableName),
					Item:                avUrs,
					ConditionExpression: aws.String(conditionPKNotExist),
				},
			},
			{
				Put: &dynamodb.Put{
					TableName:           aws.String(u.tableName),
					Item:                avUrsName,
					ConditionExpression: aws.String(conditionPKNotExist),
				},
			},
			{
				Put: &dynamodb.Put{
					TableName:           aws.String(u.tableName),
					Item:                avUsrMail,
					ConditionExpression: aws.String(conditionPKNotExist),
				},
			},
		},
	})
	if err != nil {
		u.log.WithError(err).Error("error creating user")
		return fmt.Errorf("error creating user: %w", errs.ErrFailTransaction)
	}
	return nil
}

func (u UserStorage) Update(id string, user *dto.User) error {
	key := storage.GetUserKey(id)

	update := expression.Set(expression.Name("User"), expression.Value(user.User)).
		Set(expression.Name("Email"), expression.Value(user.Email)).
		Set(expression.Name("UpdatedAt"), expression.Value(u.now()))

	return u.repo.UpdateItem(key, update)
}

func (u UserStorage) Delete(id string) error {
	pk := storage.UserRecordName + "#" + id
	sk := storage.UserRecordName
	f := expression.Name("PK").Equal(expression.Value(pk)).
		And(expression.Name("SK").Equal(expression.Value(sk)))
	expr, err := expression.NewBuilder().WithCondition(f).Build()
	if err != nil {
		u.log.WithError(err).Error("error building expression")
		return fmt.Errorf("%v: error building expression %w", err, errs.ErrAWSConfig)
	}
	return u.repo.DeleteItem(pk, sk, expr)
}
