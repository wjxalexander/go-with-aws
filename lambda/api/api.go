package api

import (
	"encoding/json"
	"fmt"
	"lambda/database"
	"lambda/types"
	"net/http"
	"strings"

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
	user, err := types.NewUser(registerUser)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Intenal Server Error",
			StatusCode: http.StatusInternalServerError,
		}, err
	}
	// user not exists
	err = api.daStore.InsertUser(user)
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

func (api ApiHandler) LoginUser(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	type LoginRequest struct {
		Username string `json:"username"` // tag
		Password string `json:"password"`
	}
	var loginUser LoginRequest
	err := json.Unmarshal([]byte(request.Body), &loginUser)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Invalid Request",
			StatusCode: http.StatusBadRequest,
		}, err
	}
	// Validate input
	if strings.TrimSpace(loginUser.Username) == "" {
		return events.APIGatewayProxyResponse{
			Body:       "Username cannot be empty",
			StatusCode: http.StatusBadRequest,
		}, fmt.Errorf("empty username")
	}

	if strings.TrimSpace(loginUser.Password) == "" {
		return events.APIGatewayProxyResponse{
			Body:       "Password cannot be empty",
			StatusCode: http.StatusBadRequest,
		}, fmt.Errorf("empty password")
	}

	serverData, err := api.daStore.GetUser(loginUser.Username)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Intenal Server Error",
			StatusCode: http.StatusInternalServerError,
		}, err
	}
	if !types.ValidatePassword(serverData.PasswordHash, loginUser.Password) {
		return events.APIGatewayProxyResponse{
			Body:       "Wrong username or password",
			StatusCode: http.StatusBadRequest,
		}, err
	}

	accessToken := types.CreateToken(serverData)
	successMsg := fmt.Sprintf(`{"access_token": "%s"}`, accessToken)
	return events.APIGatewayProxyResponse{
		Body:       successMsg,
		StatusCode: http.StatusOK,
	}, nil
}
