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
	GetSessionByUser(user User) Result[Session]
	GenerateNewSession(user User, existingSession *Session) Result[Session]
	RegisterNewUserByEmail(email string, password string) Result[VerificationToken]
	VerifyEmail(userId string, verifyToken string) Result[string]
	RefreshToken(refreshToken string) Result[Session]
}

type UserRepositoryPostgress struct {
	db           *gorm.DB
	tokenManager utils.TokenManager
}

func (repository *UserRepositoryPostgress) GetUserByEmail(email string) Result[User] {
	storedUser := User{}
	result := repository.db.Where("email = ?", email).First(&storedUser)

	if result.Error != nil {
		log.Errorf("GetUserByEmail: %v", result.Error)
		return Error[User](result.Error.Error())
	} else if result.RowsAffected == 0 {
		log.Errorf("GetUserByEmail: No user found")
		return Error[User]("No user found")
	}

	return Success(storedUser)
}

func (repository *UserRepositoryPostgress) GetSessionByUser(user User) Result[Session] {
	storedSession := Session{}
	repository.db.Where("user_id = ?", user.ID).First(&storedSession)
	verifyTokenError := repository.tokenManager.VerifyToken(storedSession.AuthToken)
	if verifyTokenError != nil {
		return repository.GenerateNewSession(user, &storedSession)
	}

	return Success(storedSession)
}

func (repository *UserRepositoryPostgress) GenerateNewSession(user User, existingSession *Session) Result[Session] {
	createAuthTokenResult := repository.tokenManager.CreateToken(user.Email, utils.AuthTokenType)
	if createAuthTokenResult.IsError() {
		return Error[Session](createAuthTokenResult.Error)
	}
	createRefreshTokenResult := repository.tokenManager.CreateToken(user.Email, utils.RefreshTokenType)
	if createRefreshTokenResult.IsError() {
		return Error[Session](createRefreshTokenResult.Error)
	}
	newAuthToken := createAuthTokenResult.Result
	newRefreshToken := createRefreshTokenResult.Result

	session := Session{
		UserId:       user.ID,
		AuthToken:    newAuthToken,
		RefreshToken: newRefreshToken,
	}
	if existingSession != nil {
		session.ID = existingSession.ID
	}

	if err := repository.db.Save(session); err.Error != nil {
		return Error[Session]("Save session error")
	}
	return Success(session)
}

func (repository *UserRepositoryPostgress) RegisterNewUserByEmail(email string, password string) Result[VerificationToken] {
	var verificationToken string
	newUser := User{Email: email, Password: password}
	err := repository.db.Transaction(func(tx *gorm.DB) error {
		if createUserErr := tx.Create(&newUser).Error; createUserErr != nil {
			return createUserErr
		}

		createVerifyTokenResult := repository.tokenManager.CreateToken(email, utils.VerifyTokenType)
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

func (repository *UserRepositoryPostgress) VerifyEmail(userId string, verifyToken string) Result[string] {
	err := repository.db.Transaction(func(tx *gorm.DB) error {
		var verificationToken VerificationToken
		result := repository.db.Where(&VerificationToken{UserId: userId, Token: verifyToken}).First(&verificationToken)

		if result.Error != nil || result.RowsAffected == 0 {
			log.Errorf("VerifyEmail: token not found")
			return &errors.GeneralError{
				Message: "VerifyEmail token not found",
			}
		}

		result = repository.db.Model(&User{}).Where("ID = ?", userId).Update("verified", "true")

		if result.Error != nil || result.RowsAffected == 0 {
			log.Errorf("VerifyEmail")
			return &errors.GeneralError{
				Message: "VerifyEmail error",
			}
		}

		result = repository.db.Where("token = ?", verifyToken).Delete(&VerificationToken{})

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

func (repository *UserRepositoryPostgress) getUserById(userId string) Result[User] {
	storedUser := User{}
	result := repository.db.Where("ID = ?", userId).First(&storedUser)

	if result.Error != nil {
		log.Errorf("GetUserById: %v", result.Error)
		return Error[User](result.Error.Error())
	} else if result.RowsAffected == 0 {
		log.Errorf("GetUserById: No user found")
		return Error[User]("No user found")
	}

	return Success(storedUser)
}

func (repository *UserRepositoryPostgress) RefreshToken(refreshToken string) Result[Session] {
	var storedSession Session
	result := repository.db.Take(&Session{RefreshToken: refreshToken}, storedSession)

	if result.Error != nil {
		return Error[Session]("No session available")
	}

	getUserResult := repository.getUserById(storedSession.UserId)

	if getUserResult.IsError() {
		return Error[Session]("No user with this id")
	}

	return repository.GenerateNewSession(getUserResult.Result, &storedSession)
}

func NewUserRepository(db *gorm.DB, manager utils.TokenManager) UserRepository {
	return &UserRepositoryPostgress{db: db, tokenManager: manager}
}
