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
	GenerateNewSession(user User) Result[Session]
	RegisterNewUserByEmail(email string, password string) Result[VerificationToken]
	VerifyEmail(userId uint, verifyToken string) Result[string]
	RefreshToken(refreshToken string) Result[Session]
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

func (pgDB *UserRepositoryPostgress) GenerateNewSession(user User) Result[Session] {
	createAuthTokenResult := utils.CreateToken(user.Email)
	if createAuthTokenResult.IsError() {
		return Error[Session](createAuthTokenResult.Error)
	}
	createRefreshTokenResult := utils.CreateToken(user.Email)
	if createRefreshTokenResult.IsError() {
		return Error[Session](createRefreshTokenResult.Error)
	}
	newToken := createAuthTokenResult.Result
	newRefreshToken := createAuthTokenResult.Result

	session := Session{UserId: user.ID, AuthToken: newToken, RefreshToken: newRefreshToken}
	if err := pgDB.db.Where(Session{UserId: user.ID}).
		Assign(session).
		FirstOrCreate(&session); err != nil {
		return Error[Session]("Save session error")
	}
	return Success(session)
}

func (pgDB *UserRepositoryPostgress) RegisterNewUserByEmail(email string, password string) Result[VerificationToken] {
	var verificationToken string
	newUser := User{Email: email, Password: password}
	err := pgDB.db.Transaction(func(tx *gorm.DB) error {
		if createUserErr := tx.Create(&newUser).Error; createUserErr != nil {
			return createUserErr
		}

		createVerifyTokenResult := utils.CreateToken(email)
		if createVerifyTokenResult.IsError() {
			return &errors.GeneralError{
				Message: "Verification token cannot be generated",
			}
		}

		if err := tx.Create(&VerificationToken{UserId: newUser.ID, Token: createVerifyTokenResult.Result}).Error; err != nil {
			return err
		}

		verificationToken = createVerifyTokenResult.Result
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

func (pgDB *UserRepositoryPostgress) getUserById(userId uint) Result[User] {
	storedUser := User{}
	result := pgDB.db.Where("ID = ?", userId).First(&storedUser)

	if result.Error != nil {
		log.Errorf("GetUserById: %v", result.Error)
		return Error[User](result.Error.Error())
	} else if result.RowsAffected == 0 {
		log.Errorf("GetUserById: No user found")
		return Error[User]("No user found")
	}

	return Success(storedUser)
}

func (pgDB *UserRepositoryPostgress) RefreshToken(refreshToken string) Result[Session] {
	var storedSession Session
	result := pgDB.db.Take(&Session{RefreshToken: refreshToken}, storedSession)

	if result.Error != nil {
		return Error[Session]("No session available")
	}

	getUserResult := pgDB.getUserById(storedSession.UserId)

	if getUserResult.IsError() {
		return Error[Session]("No user with this id")
	}

	return pgDB.GenerateNewSession(getUserResult.Result)
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryPostgress{db: db}
}
