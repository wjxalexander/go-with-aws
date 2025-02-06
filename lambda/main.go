package main

import (
	"lambda/app"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// lambda: our app backend logic
// serveless: abstract a lot of the server management core principles into a server-based architecture.
// 简单来说就是以前需要考虑的的scale, ngnix, load balance 他给你全做了
// why?: only runs when it get revoked

type MyEvent struct {
	Username string `json:"username"` // tag
}

// func HandleRequest(event MyEvent) (string, error) {
// 	if event.Username == "" {
// 		return "", fmt.Errorf("username cannot be empty")
// 	}
// 	return fmt.Sprintf("Successfully called by - %s", event.Username), nil
// }

func main() {
	myApp := app.NewApp()
	// lambda.Start(myApp.ApiHandler.RegisterUserHandler)
	lambda.Start(func(requset events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		switch requset.Path {
		case "/register":
			return myApp.ApiHandler.RegisterUserHandler(requset)
		default:
			return events.APIGatewayProxyResponse{
				Body:       "not found",
				StatusCode: http.StatusNotFound,
			}, nil
		}
	})
}
