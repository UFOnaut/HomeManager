package usecases

import (
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"home_manager/entities"
	"home_manager/handlers/errors"
	"home_manager/models"
	"home_manager/repositories"
	"home_manager/utils"
	"testing"
)

const password = "test_password"
const id = "test_id"
const email = "test@test.com"
const userId = "test_user_id"
const empty = ""
const name = "test_name"
const trueString = "true"
const authToken = "test_auth_token"
const refreshToken = "test_refresh_token"

func TestLogin_AuthTokenVerifySuccessAndSessionReturns(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	hashedPassword := entities.Hash(password)
	userRows := sqlmock.NewRows([]string{
		"id", "email", "password", "groups_ids", "name", "verified",
	}).AddRow(userId, email, hashedPassword, empty, name, trueString)
	getUserQuery := `[SELECT * FROM "users" WHERE email = ? ORDER BY "users"."id" LIMIT 1]`
	mock.
		ExpectQuery(getUserQuery).
		WithArgs(email).
		WillReturnRows(userRows)

	sessionRows := sqlmock.NewRows([]string{
		"id", "user_id", "auth_token", "refresh_token",
	}).AddRow(id, userId, authToken, refreshToken)
	getSessionQuery := `[SELECT * FROM "sessions" WHERE "sessions"."user_id" = ? ORDER BY "sessions"."id" LIMIT 1]`
	mock.
		ExpectQuery(getSessionQuery).
		WithArgs(userId).
		WillReturnRows(sessionRows)
	mock.
		ExpectQuery(getSessionQuery).
		WithArgs(userId).
		WillReturnRows(sessionRows)

	mock.ExpectCommit()

	postgresDb := postgres.New(postgres.Config{
		Conn:       mockDB,
		DriverName: "postgres",
	})

	db, _ := gorm.Open(postgresDb, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	tokenManager := TokenManagerMock{}
	repository := repositories.NewUserRepository(db, &tokenManager)
	loginUseCase := NewLoginUseCase(repository)
	data := models.LoginData{
		Email:    email,
		Password: password,
	}
	result := loginUseCase.Execute(&data)
	if result.IsError() {
		t.Errorf("Result was incorrect")
	}
}

func TestLogin_AuthTokenVerifyErrorButNewSessionReturns(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	hashedPassword := entities.Hash(password)
	userRows := sqlmock.NewRows([]string{
		"id", "email", "password", "groups_ids", "name", "verified",
	}).AddRow(userId, email, hashedPassword, empty, name, trueString)
	getUserQuery := `[SELECT * FROM "users" WHERE email = ? ORDER BY "users"."id" LIMIT 1]`
	mock.
		ExpectQuery(getUserQuery).
		WithArgs(email).
		WillReturnRows(userRows)

	sessionRows := sqlmock.NewRows([]string{
		"id", "user_id", "auth_token", "refresh_token",
	}).AddRow(id, userId, authToken, refreshToken)
	mock.
		ExpectQuery(`[SELECT * FROM "sessions" WHERE "sessions"."user_id" = ? ORDER BY "sessions"."id" LIMIT 1]`).
		WithArgs(userId).
		WillReturnRows(sessionRows)

	mock.ExpectBegin()

	mock.
		ExpectExec(`[UPDATE "sessions" SET "id" = ?, "user_id" = ?, "auth_token" = ?, "refresh_token" = ? WHERE "sessions"."id" = ? ]`).
		WithArgs(id, userId, authToken, refreshToken, id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	postgresDb := postgres.New(postgres.Config{
		Conn:       mockDB,
		DriverName: "postgres",
	})

	db, _ := gorm.Open(postgresDb, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	tokenManager := TokenManagerMock{
		AuthToken:    authToken,
		RefreshToken: refreshToken,
		VerifyError:  &errors.GeneralError{Message: ""}}
	repository := repositories.NewUserRepository(db, &tokenManager)
	loginUseCase := NewLoginUseCase(repository)
	data := models.LoginData{
		Email:    email,
		Password: password,
	}
	result := loginUseCase.Execute(&data)
	if result.IsError() {
		t.Errorf("Result was incorrect")
	}
}

func TestLogin_AuthTokenCreateError(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	hashedPassword := entities.Hash(password)
	userRows := sqlmock.NewRows([]string{
		"id", "email", "password", "groups_ids", "name", "verified",
	}).AddRow(userId, email, hashedPassword, empty, name, trueString)
	getUserQuery := `[SELECT * FROM "users" WHERE email = ? ORDER BY "users"."id" LIMIT 1]`
	mock.
		ExpectQuery(getUserQuery).
		WithArgs(email).
		WillReturnRows(userRows)

	sessionRows := sqlmock.NewRows([]string{
		"id", "user_id", "auth_token", "refresh_token",
	}).AddRow(id, userId, authToken, refreshToken)
	mock.
		ExpectQuery(`[SELECT * FROM "sessions" WHERE "sessions"."user_id" = ? ORDER BY "sessions"."id" LIMIT 1]`).
		WithArgs(userId).
		WillReturnRows(sessionRows)

	mock.ExpectBegin()

	mock.
		ExpectExec(`[UPDATE "sessions" SET "id" = ?, "user_id" = ?, "auth_token" = ?, "refresh_token" = ? WHERE "sessions"."id" = ? ]`).
		WithArgs(id, userId, authToken, refreshToken, id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	postgresDb := postgres.New(postgres.Config{
		Conn:       mockDB,
		DriverName: "postgres",
	})

	db, _ := gorm.Open(postgresDb, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	tokenManager := TokenManagerMock{
		AuthToken:    empty,
		RefreshToken: refreshToken,
		VerifyError:  &errors.GeneralError{Message: ""}}
	repository := repositories.NewUserRepository(db, &tokenManager)
	loginUseCase := NewLoginUseCase(repository)
	data := models.LoginData{
		Email:    email,
		Password: password,
	}
	result := loginUseCase.Execute(&data)
	if !result.IsError() {
		t.Errorf("Result was incorrect")
	}
}

type TokenManagerMock struct {
	AuthToken    string
	RefreshToken string
	VerifyError  error
}

func (tokenManager *TokenManagerMock) CreateToken(email string, tokenType string) entities.Result[string] {
	if tokenType == utils.AuthTokenType && len(tokenManager.AuthToken) > 0 {
		return entities.Success(tokenManager.AuthToken)
	} else if tokenType == utils.RefreshTokenType && len(tokenManager.RefreshToken) > 0 {
		return entities.Success(tokenManager.RefreshToken)
	} else {
		return entities.Error[string]("Token creation error")
	}

}

func (tokenManager *TokenManagerMock) VerifyToken(tokenString string) error {
	return tokenManager.VerifyError
}
