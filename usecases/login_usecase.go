package usecases

import (
	. "home_manager/entities"
	"home_manager/models"
	"home_manager/repositories"
	"home_manager/utils"
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
		return Error[string](getUserResult.Error)
	}

	user := getUserResult.Result
	if user.IsPasswordCorrect(in.Password) {
		getSessionResult := u.repository.GetSessionByUserId(user.ID)

		token := getSessionResult.Result.Token
		verifyTokenError := utils.VerifyToken(token)
		if verifyTokenError != nil {
			createTokenResult := utils.CreateToken(user.Email)
			if createTokenResult.IsError() {
				return createTokenResult
			}
			newToken := createTokenResult.Result
			u.repository.SaveSession(newToken, user.ID)
			return Success(newToken)
		}

		return Success(token)
	} else {
		return Error[string]("Incorrect password")
	}
}

func NewLoginUseCase(
	repository repositories.UserRepository,
) LoginUseCase {
	return &LoginUseCaseImpl{
		repository: repository,
	}
}
