package app

import (
	"lambda-func/api"
	"lambda-func/database"
)

type App struct {
	ApiHandler api.ApiHandler
}

func NewApp() App {
	// here init db store
	// pass the db to Api handler
	db := database.NewDynamoDBClient()
	apiHandler := api.NewApiHandler(db)
	return App{
		ApiHandler: apiHandler,
	}
}
