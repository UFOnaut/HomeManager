package usecases

import (
	. "home_manager/entities"
	"home_manager/repositories"
)

type VerifyEmailUseCase interface {
	Execute(userId uint, verifyToken string) Result[string]
}

type VerifyEmailUseCaseImpl struct {
	repository repositories.UserRepository
}

func (u *VerifyEmailUseCaseImpl) Execute(userId uint, verifyToken string) Result[string] {
	return u.repository.VerifyEmail(userId, verifyToken)
}

func NewVerifyEmailUseCase(
	repository repositories.UserRepository,
) VerifyEmailUseCase {
	return &VerifyEmailUseCaseImpl{
		repository: repository,
	}
}
