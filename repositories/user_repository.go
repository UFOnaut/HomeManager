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
	RegisterNewUserByEmail(email string, password string) Result[VerificationToken]
	VerifyEmail(userId uint, verifyToken string) Result[string]
}

type UserRepositoryPostgress struct {
	db *gorm.DB
}

func (pgDB *UserRepositoryPostgress) GetUserByEmail(email string) Result[User] {
	storedUser := User{}
	result := pgDB.db.Where("email = ?", email).First(&storedUser)

	if result.Error != nil {
		log.Errorf("GetUserByEmail: %v", result.Error)
		return Error[User](result.Error.Error())
	} else if result.RowsAffected == 0 {
		log.Errorf("GetUserByEmail: No user found")
		return Error[User]("No user found")
	}

	return Success(storedUser)
}

func (pgDB *UserRepositoryPostgress) GetSessionByUserId(id uint) Result[Session] {
	var storedSession Session
	result := pgDB.db.Take(&Session{UserId: id}, storedSession)

	if result.Error != nil {
		return Success(Session{}) //new session
	}

	return Success(storedSession)
}

func (pgDB *UserRepositoryPostgress) SaveSession(token string, userId uint) bool {
	if err := pgDB.db.Where(Session{UserId: userId}).
		Assign(Session{UserId: userId, Token: token}).
		FirstOrCreate(&Session{UserId: userId, Token: token}); err != nil {
		return false
	}
	return true
}

func (pgDB *UserRepositoryPostgress) RegisterNewUserByEmail(email string, password string) Result[VerificationToken] {
	var verificationToken string
	newUser := User{Email: email, Password: password}
	err := pgDB.db.Transaction(func(tx *gorm.DB) error {
		if createUserErr := tx.Create(&newUser).Error; createUserErr != nil {
			return createUserErr
		}

		createTokenResult := utils.CreateToken(email)
		if createTokenResult.IsError() {
			return &errors.GeneralError{
				Message: "Verification token cannot be generated",
			}
		}

		if err := tx.Create(&VerificationToken{UserId: newUser.ID, Token: createTokenResult.Result}).Error; err != nil {
			return err
		}

		verificationToken = createTokenResult.Result
		return nil
	})

	if err != nil {
		return Error[VerificationToken]("Verification token error: " + err.Error())
	}
	return Success(VerificationToken{
		Token:  verificationToken,
		UserId: newUser.ID,
	})
}

func (pgDB *UserRepositoryPostgress) VerifyEmail(userId uint, verifyToken string) Result[string] {
	err := pgDB.db.Transaction(func(tx *gorm.DB) error {
		var verificationToken VerificationToken
		result := pgDB.db.Where(&VerificationToken{UserId: userId, Token: verifyToken}).First(&verificationToken)

		if result.Error != nil || result.RowsAffected == 0 {
			log.Errorf("VerifyEmail: token not found")
			return &errors.GeneralError{
				Message: "VerifyEmail token not found",
			}
		}

		result = pgDB.db.Model(&User{}).Where("ID = ?", userId).Update("verified", "true")

		if result.Error != nil || result.RowsAffected == 0 {
			log.Errorf("VerifyEmail")
			return &errors.GeneralError{
				Message: "VerifyEmail error",
			}
		}

		result = pgDB.db.Where("token = ?", verifyToken).Delete(&VerificationToken{})

		if result.Error != nil || result.RowsAffected == 0 {
			log.Errorf("VerifyEmail: token not found")
			return &errors.GeneralError{
				Message: "VerifyEmail token not found",
			}
		}

		return nil
	})

	if err != nil {
		return Error[string](err.Error())
	}

	return Success("Verified successfully")
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryPostgress{db: db}
}
