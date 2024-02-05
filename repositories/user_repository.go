package repositories

import (
	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
	. "home_manager/entities"
)

type UserRepository interface {
	GetUserByEmail(email string) Result[User]
	GetSessionByUserId(id uint) Result[Session]
	SaveSession(token string, userId uint) bool
}

type UserRepositoryPostgress struct {
	db *gorm.DB
}

func (pgDB *UserRepositoryPostgress) GetUserByEmail(email string) Result[User] {
	storedUser := User{Email: email}
	result := pgDB.db.Take(&storedUser)

	if result.Error != nil {
		log.Errorf("GetUserByEmail: %v", result.Error)
		return Result[User]{Error: result.Error.Error()}
	} else if result.RowsAffected == 0 {
		log.Errorf("GetUserByEmail: No user found")
		return Result[User]{Error: "No user found"}
	}

	return Result[User]{Result: storedUser}
}

func (pgDB *UserRepositoryPostgress) GetSessionByUserId(id uint) Result[Session] {
	var storedSession Session
	result := pgDB.db.Take(&Session{UserId: id}, storedSession)

	if result.Error != nil {
		return Result[Session]{Result: Session{}}
	}

	return Result[Session]{Result: storedSession}
}

func (pgDB *UserRepositoryPostgress) SaveSession(token string, userId uint) bool {
	if err := pgDB.db.Where(Session{UserId: userId}).
		Assign(Session{UserId: userId, Token: token}).
		FirstOrCreate(&Session{UserId: userId, Token: token}); err != nil {
		return false
	}
	return true
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryPostgress{db: db}
}
