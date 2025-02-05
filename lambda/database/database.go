package database

import (
	"lambda/types"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type DynamoDBClient struct {
	databaseStore *dynamodb.DynamoDB
}

const (
	TABLE_NAME = "user_table"
)

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

func (u DynamoDBClient) InsertUser(user types.RegisterUser) error {
	// assemble item
	item := &dynamodb.PutItemInput{
		TableName: aws.String(TABLE_NAME),
		// https://go.dev/blog/maps map[KeyType]ValueType
		Item: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(user.Username),
			},
			"password": {
				S: aws.String(user.Password),
			},
		},
	}
	_, err := u.databaseStore.PutItem(item)
	return err
}
