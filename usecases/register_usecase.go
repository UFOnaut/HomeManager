package usecases

import (
	. "home_manager/entities"
	"home_manager/models"
	"home_manager/repositories"
	"home_manager/utils"
)

type RegisterUseCase interface {
	Register(in *models.RegisterData) Result[string]
}

type RegisterUseCaseImpl struct {
	repository repositories.UserRepository
}

func (u *RegisterUseCaseImpl) Register(in *models.RegisterData) Result[string] {
	getUserResult := u.repository.GetUserByEmail(in.Email)

	if !getUserResult.IsError() {
		return Error[string]("User already exists")
	} else {
		registerResult := u.repository.RegisterNewUserByEmail(in.Email, in.Password)
		if registerResult.IsError() {
			return Error[string](registerResult.Error)
		}
		sendVerificationResult := utils.SendVerificationEmail(in.Email, registerResult.Result)
		if sendVerificationResult != nil {
			return Error[string]("Send verification email error: " + sendVerificationResult.Error())
		}
		return Success(registerResult.Result.Token)
	}

}

func NewRegisterUseCase(
	repository repositories.UserRepository,
) RegisterUseCase {
	return &RegisterUseCaseImpl{
		repository: repository,
	}
}
