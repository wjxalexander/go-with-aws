package api

import (
	"lambda-func/database"
)

type ApiHandler struct {
	daStore database.DynamoDBClient
}

func NewApiHandler(dbStore database.DynamoDBClient) ApiHandler {
	return ApiHandler{
		daStore: dbStore,
	}
}
