package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"home_manager/config"
	"home_manager/entities"
)

const AuthTokenType = "auth"
const RefreshTokenType = "refresh"
const VerifyTokenType = "verify"

type (
	TokenManager interface {
		VerifyToken(tokenString string) error
		CreateToken(email string, tokenType string) entities.Result[string]
	}
	TokenManagerImpl struct {
		Config config.Config
	}
)

func (tokenManager *TokenManagerImpl) CreateToken(email string, tokenType string) entities.Result[string] {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, GetUniqueClaims(email, tokenType))

	tokenString, err := token.SignedString([]byte(config.GetConfig().Jwt.SecretKey))
	if err != nil {
		return entities.Error[string](err.Error())
	}

	return entities.Success(tokenString)
}

func (tokenManager *TokenManagerImpl) VerifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GetConfig().Jwt.SecretKey), nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}
