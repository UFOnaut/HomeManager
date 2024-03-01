package usecases

import (
	. "home_manager/entities"
	"home_manager/repositories"
)

type RefreshTokenUseCase interface {
	Execute(refreshToken string) Result[Session]
}

type RefreshTokenUseCaseImpl struct {
	repository repositories.UserRepository
}

func (u *RefreshTokenUseCaseImpl) Execute(refreshToken string) Result[Session] {
	return u.repository.RefreshToken(refreshToken)
}

func NewRefreshTokenUseCase(
	repository repositories.UserRepository,
) RefreshTokenUseCase {
	return &RefreshTokenUseCaseImpl{
		repository: repository,
	}
}
