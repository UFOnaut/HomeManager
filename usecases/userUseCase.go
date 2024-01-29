package usecases

import (
	entities "home_manager/entities"
	models "home_manager/models"
	repositories "home_manager/repositories"
)

type UserUsecase interface {
	Login(in *models.LoginData) (string, error)
}

type UserUsecaseImpl struct {
	UserUsecase
	repository repositories.UserRepository
}

func NewUserUsecase(
	repository repositories.UserRepository,
) UserUsecase {
	return &UserUsecaseImpl{
		repository: repository,
	}
}

func (u *UserUsecaseImpl) Login(in *models.LoginData) (string, error) {
	loginUserData := &entities.LoginUserDto{
		Email: in.Email, Password: in.Password,
	}
	return u.repository.Login(loginUserData)
}
