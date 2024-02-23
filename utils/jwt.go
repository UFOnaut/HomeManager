package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"home_manager/config"
	"home_manager/entities"
	"time"
)

func getSignedString() string {
	cfg := config.GetConfig()
	return cfg.Jwt.SecretKey
}

func CreateToken(username string) entities.Result[string] {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username + time.Now().String(),
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})
	tokenString, err := token.SignedString([]byte(getSignedString()))
	if err != nil {
		return entities.Error[string](err.Error())
	}

	return entities.Success(tokenString)
}

func VerifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return getSignedString(), nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}
