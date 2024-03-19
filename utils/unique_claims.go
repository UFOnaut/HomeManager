package utils

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/golang-jwt/jwt"
	"time"
)

type UniqueClaims struct {
	jwt.MapClaims
	TokenId string `json:"jti,omitempty"`
}

func GetUniqueClaims(email string, tokenType string) UniqueClaims {
	bits := make([]byte, 12)
	_, err := rand.Read(bits)
	if err != nil {
		panic(err)
	}
	return UniqueClaims{
		MapClaims: jwt.MapClaims{
			"type":     tokenType,
			"username": email,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		},
		TokenId: base64.StdEncoding.EncodeToString(bits),
	}
}
