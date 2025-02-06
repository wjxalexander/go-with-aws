package api

import (
	"encoding/json"
	"fmt"
	"lambda/database"
	"lambda/types"
	"net/http"

	// "lambda/types"

	"github.com/aws/aws-lambda-go/events"
)

type ApiHandler struct {
	daStore database.UserStore
}

func NewApiHandler(dbStore database.UserStore) ApiHandler {
	return ApiHandler{
		daStore: dbStore,
	}
}

func (api ApiHandler) RegisterUserHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var registerUser types.RegisterUser
	error := json.Unmarshal([]byte(request.Body), &registerUser)
	if error != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Invalid Request",
			StatusCode: http.StatusBadRequest,
		}, error
	}
	if registerUser.Username == "" || registerUser.Password == "" {
		return events.APIGatewayProxyResponse{
			Body:       "Invalid Request",
			StatusCode: http.StatusBadRequest,
		}, fmt.Errorf("request is invalid")
	}
	// does user already exist
	useExists, err := api.daStore.DoesUserExist(registerUser.Username)
	// fmt.Errorf("there is error %w", err)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Intenal Server Error",
			StatusCode: http.StatusInternalServerError,
		}, err
	}
	if useExists {
		return events.APIGatewayProxyResponse{
			Body:       "User already exists",
			StatusCode: http.StatusConflict,
		}, nil
	}
	// user not exists
	err = api.daStore.InsertUser(registerUser)
	// fmt.Errorf("there is inster user error %w", error)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Intenal Server Error",
			StatusCode: http.StatusInternalServerError,
		}, fmt.Errorf("there is inster user error %w", error)
	}
	return events.APIGatewayProxyResponse{
		Body:       "Success",
		StatusCode: http.StatusOK,
	}, nil
}
