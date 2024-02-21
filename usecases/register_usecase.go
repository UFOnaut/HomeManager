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
		return Result[string]{Error: "User already exists"}
	} else {
		result := u.repository.RegisterNewUserByEmail(in.Email, in.Password)
		if !result.IsError() {
			verificationToken := result.Result
			err := utils.SendVerificationEmail(in.Email, verificationToken)
			if err != nil {
				return Result[string]{Error: "Send verification email error: " + err.Error()}
			}
			//TODO make endpoint to handle this GET
		}
		return result
	}

}

func NewRegisterUseCase(
	repository repositories.UserRepository,
) RegisterUseCase {
	return &RegisterUseCaseImpl{
		repository: repository,
	}
}
