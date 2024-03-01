package usecases

import (
	. "home_manager/entities"
	"home_manager/models"
	"home_manager/repositories"
	"home_manager/utils"
)

type LoginUseCase interface {
	Execute(in *models.LoginData) Result[Session]
}

type LoginUseCaseImpl struct {
	repository repositories.UserRepository
}

func (u *LoginUseCaseImpl) Execute(in *models.LoginData) Result[Session] {
	getUserResult := u.repository.GetUserByEmail(in.Email)
	if getUserResult.IsError() {
		return Error[Session](getUserResult.Error)
	}

	user := getUserResult.Result
	if user.IsPasswordCorrect(in.Password) {
		getSessionResult := u.repository.GetSessionByUserId(user.ID)

		session := getSessionResult.Result
		verifyTokenError := utils.VerifyToken(session.AuthToken)
		if verifyTokenError != nil {
			return u.repository.GenerateNewSession(user)
		}

		return Success(session)
	} else {
		return Error[Session]("Incorrect password")
	}
}

func NewLoginUseCase(
	repository repositories.UserRepository,
) LoginUseCase {
	return &LoginUseCaseImpl{
		repository: repository,
	}
}
