package repositories

import (
	"home_manager/entities"

	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetUserByEmail(email string) (*entities.User, error)
}

type UserRepositoryPostgress struct {
	db *gorm.DB
}

// GetUserByEmail implements UserRepository.
func (pgDB *UserRepositoryPostgress) GetUserByEmail(email string) (*entities.User, error) {
	var storedUser *entities.User
	//TODO validate at some level and process login
	result := pgDB.db.Find(&entities.User{Email: email}, storedUser)

	if result.Error != nil {
		log.Errorf("GetUserByEmail: %v", result.Error)
		return nil, result.Error
	}

	//TODO return token
	return storedUser, nil
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryPostgress{db: db}
}
