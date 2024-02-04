package usecases

import (
	"github.com/golang-jwt/jwt"
	"home_manager/config"
	. "home_manager/entities"
	"home_manager/models"
	"home_manager/repositories"
	"time"
)

type LoginUseCase interface {
	Login(in *models.LoginData) Result[string]
}

type LoginUseCaseImpl struct {
	repository repositories.UserRepository
}

func (u *LoginUseCaseImpl) Login(in *models.LoginData) Result[string] {
	getUserResult := u.repository.GetUserByEmail(in.Email)
	if getUserResult.IsError() {
		return Result[string]{Error: getUserResult.Error}
	}

	user := getUserResult.Result
	if user.IsPasswordCorrect(in.Password) {
		getSessionResult := u.repository.GetSessionByUserId(user.ID)

		if getSessionResult.IsError() {
			return Result[string]{Error: getSessionResult.Error}
		}

		session := getSessionResult.Result
		if !session.IsValid() {
			createTokenResult := createToken(user.Email)
			if createTokenResult.IsError() {
				return Result[string]{Error: createTokenResult.Error}
			}
			newToken := createTokenResult.Result
			u.repository.SaveSession(newToken, user.ID)
			return Result[string]{Result: newToken}
		}

		return Result[string]{Result: session.Token}
	} else {
		return Result[string]{Error: "Incorrect password"}
	}
}

func NewLoginUseCase(
	repository repositories.UserRepository,
) LoginUseCase {
	return &LoginUseCaseImpl{
		repository: repository,
	}
}

func createToken(username string) Result[string] {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})
	cfg := config.GetConfig()
	tokenString, err := token.SignedString(cfg.Jwt.SecretKey)
	if err != nil {
		return Result[string]{Error: err.Error()}
	}

	return Result[string]{Result: tokenString}
}
