package repositories

import (
	"home_manager/entities"

	"gorm.io/gorm"
)

type UserRepository interface {
	Login(in *entities.LoginUserDto) (string, error)
}

type userRepositoryPostgress struct {
	UserRepository
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepositoryPostgress{db: db}
}

func (r *userRepositoryPostgress) Login(in *entities.LoginUserDto) (string, error) {

	//TODO validate at some level and process login
	println("Email: " + in.Email + " Password: " + in.Password)
	// result := r.db.Create(data)

	// if result.Error != nil {
	// 	log.Errorf("InsertUser: %v", result.Error)
	// 	return result.Error
	// }

	// log.Debugf("InsertUser: %v", result.RowsAffected)

	//TODO return token
	return "test_token", nil
}
