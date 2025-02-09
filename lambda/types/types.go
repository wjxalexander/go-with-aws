package types

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type RegisterUser struct {
	Username string `json:"username"` // tag
	Password string `json:"password"`
}

type User struct {
	Username     string `json:"username"`
	PasswordHash string `json:"password"`
}

func NewUser(registerUser RegisterUser) (User, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(registerUser.Password), 10)
	if err != nil {
		return User{}, err
	}
	return User{
		Username:     registerUser.Username,
		PasswordHash: string(hashPassword),
	}, nil
}

func ValidatePassword(hashPassword string, plainTextPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(plainTextPassword))
	return err == nil
}

func CreateToken(user User) string {
	now := time.Now()
	validUntil := now.Add(time.Hour * 1).Unix()
	claims := jwt.MapClaims{
		"user":    user.Username,
		"expires": validUntil,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims, nil)
	secret := "secret" // NOT A PROD IMPLEMENTATION!!
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		return ""
	}
	return tokenStr
}
