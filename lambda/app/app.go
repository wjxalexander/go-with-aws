package app

import (
	"lambda/api"
	"lambda/database"
)

type App struct {
	ApiHandler api.ApiHandler
}

func NewApp() App {
	// here init db store
	// pass the db to Api handler
	db := database.NewDynamoDBClient()

	// 这里需要decouple 不一定是DynamoDB 可能是 postgres db 与 NewApiHandler type 兼容
	apiHandler := api.NewApiHandler(db)
	return App{
		ApiHandler: apiHandler,
	}
}
