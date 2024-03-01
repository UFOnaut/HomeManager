package handlers

import (
	"home_manager/repositories"
	"net/http"

	"home_manager/handlers/response"
	"home_manager/models"
	"home_manager/usecases"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type (
	LoginHandler interface {
		Login(c echo.Context) error
		RefreshToken(c echo.Context) error
	}

	LoginHttpHandler struct {
		LoginUsecase        usecases.LoginUseCase
		RefreshTokenUseCase usecases.RefreshTokenUseCase
	}
)

func NewLoginHttpHandler(repository repositories.UserRepository) LoginHandler {
	return &LoginHttpHandler{
		LoginUsecase:        usecases.NewLoginUseCase(repository),
		RefreshTokenUseCase: usecases.NewRefreshTokenUseCase(repository),
	}
}

func (h *LoginHttpHandler) Login(c echo.Context) error {
	reqBody := new(models.LoginData)
	var err = c.Bind(reqBody)
	if err != nil {
		log.Errorf("Error binding request body: %v", err)
		return response.MakeBasicResponse(c, http.StatusBadRequest, "Bad request")
	}

	loginResult := h.LoginUsecase.Execute(reqBody)
	if loginResult.IsError() {
		return response.MakeBasicResponse(c, http.StatusInternalServerError, loginResult.Error)
	}

	return response.MakeLoginResponse(c, http.StatusOK, loginResult.Result)
}

func (h *LoginHttpHandler) RefreshToken(c echo.Context) error {
	reqBody := new(models.RefreshToken)
	var err = c.Bind(reqBody)
	if err != nil {
		log.Errorf("Error binding request body: %v", err)
		return response.MakeBasicResponse(c, http.StatusBadRequest, "Bad request")
	}

	refreshTokenResult := h.RefreshTokenUseCase.Execute(reqBody.RefreshToken)
	if refreshTokenResult.IsError() {
		return response.MakeBasicResponse(c, http.StatusInternalServerError, refreshTokenResult.Error)
	}

	return response.MakeLoginResponse(c, http.StatusOK, refreshTokenResult.Result)
}
