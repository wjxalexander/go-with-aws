package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

// lambda: our app backend logic
// serveless: abstract a lot of the server management core principles into a server-based architecture.
// 简单来说就是以前需要考虑的的scale, ngnix, load balance 他给你全做了
// why?: only runs when it get revoked

type MyEvent struct {
	Username string `json:"username"` // tag
}

func HandleRequest(event MyEvent) (string, error) {
	if event.Username == "" {
		return "", fmt.Errorf("username cannot be empty")
	}
	return fmt.Sprintf("Successfully called by - %s", event.Username), nil
}

func main() {
	lambda.Start(HandleRequest)
}
