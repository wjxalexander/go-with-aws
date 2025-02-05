package api

import (
	"fmt"
	"lambda/database"
	"lambda/types"
)

type ApiHandler struct {
	daStore database.UserStore
}

func NewApiHandler(dbStore database.UserStore) ApiHandler {
	return ApiHandler{
		daStore: dbStore,
	}
}

func (api ApiHandler) RegisterUserHandler(event types.RegisterUser) error {
	if event.Username == "" || event.Password == "" {
		return fmt.Errorf("request is invalid")
	}
	// does user already exist
	useExists, err := api.daStore.DoesUserExist(event.Username)
	if err != nil {
		return fmt.Errorf("there is error %w", err)
	}
	if useExists {
		return fmt.Errorf("user exists")
	}
	// user not exists
	error := api.daStore.InsertUser(event)
	if err != nil {
		return fmt.Errorf("there is inster user error %w", error)
	}
	return nil
}
