package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/golang-jwt/jwt/v5"
)

// extract request header
// extract pur claims
// validating everyting

func ValidateJWTMiddleware(next func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)) func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		// extract headers
		tokenStr := extractTokenFromHeaders(request.Headers)
		if tokenStr == "" {
			return events.APIGatewayProxyResponse{
				Body:       "missing Auth token",
				StatusCode: http.StatusUnauthorized,
			}, nil
		}
		claims, error := parseToken(tokenStr)

		if error != nil {
			return events.APIGatewayProxyResponse{
				Body:       "user unauthorized",
				StatusCode: http.StatusUnauthorized,
			}, error
		}
		exp, ok := claims["expires"].(float64) // JWT timestamps come as float64
		if !ok {
			log.Printf("[AUTH] [%s] Invalid expiration format in claims", claims["expires"])
			return events.APIGatewayProxyResponse{
				Body:       "invalid token expiration",
				StatusCode: http.StatusUnauthorized,
			}, fmt.Errorf("invalid token expiration format")
		}
		expInt := int64(exp)
		if time.Now().Unix() > expInt {
			return events.APIGatewayProxyResponse{
				Body:       "token expired",
				StatusCode: http.StatusUnauthorized,
			}, nil
		}
		return next(request)
	}
}

func extractTokenFromHeaders(header map[string]string) string {
	authHeader, ok := header["Authorization"]
	if !ok {
		return ""
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")
	return strings.TrimSpace(token)
}

var secretKey = []byte("secret")

func parseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("unauthorized")
	}
	if !token.Valid {
		return nil, fmt.Errorf("token is not valid - unauthorized")
	}
	claims, ok := token.Claims.(jwt.MapClaims) // same as create token claim
	if !ok {
		return nil, fmt.Errorf("token of unrecognized type - unauthorized")
	}

	return claims, nil
}
