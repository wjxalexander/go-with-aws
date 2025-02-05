package types

import "golang.org/x/crypto/bcrypt"

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
