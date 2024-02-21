package repositories

import (
	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
	. "home_manager/entities"
	"home_manager/handlers/errors"
	"home_manager/utils"
)

type UserRepository interface {
	GetUserByEmail(email string) Result[User]
	GetSessionByUserId(id uint) Result[Session]
	SaveSession(token string, userId uint) bool
	RegisterNewUserByEmail(email string, password string) Result[string]
}

type UserRepositoryPostgress struct {
	db *gorm.DB
}

func (pgDB *UserRepositoryPostgress) GetUserByEmail(email string) Result[User] {
	storedUser := User{}
	result := pgDB.db.Where("email = ?", email).First(&storedUser)

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

func (pgDB *UserRepositoryPostgress) RegisterNewUserByEmail(email string, password string) Result[string] {
	var verificationToken string
	err := pgDB.db.Transaction(func(tx *gorm.DB) error {
		newUser := User{Email: email, Password: password}
		if err := tx.Create(&newUser).Error; err != nil {
			return err
		}

		verificationToken := utils.CreateToken(email)
		if verificationToken.IsError() {
			return &errors.GeneralError{
				Message: "Verification token cannot be generated",
			}
		}

		if err := tx.Create(&VerificationToken{UserId: newUser.ID, VerificationToken: verificationToken.Result}).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return Result[string]{Error: "Verification token error: " + err.Error()}
	}
	return Result[string]{Result: verificationToken}
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryPostgress{db: db}
}
