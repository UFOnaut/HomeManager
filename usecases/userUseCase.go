package usecases

import (
	"home_manager/models"
	"home_manager/repositories"
)

type UserUsecase interface {
	Login(in *models.LoginData) (string, error)
}

type UserUsecaseImpl struct {
	repository repositories.UserRepository
}

func (u *UserUsecaseImpl) Login(in *models.LoginData) (string, error) {
	user, err := u.repository.GetUserByEmail(in.Email)
	if err != nil {
		return "", err
	}
	if user.IsPasswordCorrect(in.Password) {
		//TODO get session by user id and return token
		//If no teken yet or expired - generate new one and save to session table
		token := ""
		return token, nil
	} else {
		return "", err
	}
}

func NewUserUsecase(
	repository repositories.UserRepository,
) UserUsecase {
	return &UserUsecaseImpl{
		repository: repository,
	}
}
