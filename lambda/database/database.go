package database

import (
	"fmt"
	"lambda/types"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type DynamoDBClient struct {
	databaseStore *dynamodb.DynamoDB
}

const (
	TABLE_NAME = "user_table"
)

// define functions
type UserStore interface {
	DoesUserExist(username string) (bool, error)
	InsertUser(user types.User) error
	GetUser(username string) (types.User, error)
}

// aws: session to connect db NewDynamoDBClient will generate a session
func NewDynamoDBClient() DynamoDBClient {
	dbSession := session.Must(session.NewSession())
	db := dynamodb.New(dbSession) //*dynamodb.DynamoDB
	return DynamoDBClient{
		databaseStore: db,
	}
}

// DOC https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/example_dynamodb_GetItem_section.html
// 1. Does user exists

// 2. insert user to dynamodb

func (u DynamoDBClient) DoesUserExist(username string) (bool, error) {
	itemInput := dynamodb.GetItemInput{
		TableName: aws.String(TABLE_NAME),
		Key: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(username),
			},
		},
	}
	result, error := u.databaseStore.GetItem(&itemInput)
	if error != nil {
		return true, error
	}
	if result.Item == nil {
		return false, nil
	}
	return true, nil
}

func (u DynamoDBClient) InsertUser(user types.User) error {
	// assemble item
	item := &dynamodb.PutItemInput{
		TableName: aws.String(TABLE_NAME),
		// https://go.dev/blog/maps map[KeyType]ValueType
		Item: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(user.Username),
			},
			"password": {
				S: aws.String(user.PasswordHash),
			},
		},
	}
	_, err := u.databaseStore.PutItem(item)
	return err
}

func (u DynamoDBClient) GetUser(username string) (types.User, error) {
	var user types.User
	result, err := u.databaseStore.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(TABLE_NAME),
		Key: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(username),
			},
		},
	})
	if err != nil {
		return user, err
	}
	if result.Item == nil {
		return user, fmt.Errorf("user not found")
	}
	marshalErr := dynamodbattribute.UnmarshalMap(result.Item, &user)
	if marshalErr != nil {
		return user, marshalErr
	}
	return user, nil
}

// curl -X POST https://kkqx5qn448.execute-api.eu-west-1.amazonaws.com/prod/register -H "Content-Type: application/json" -d '{"username":"Jingxin.Wang", "password":"Aa123456"}'
