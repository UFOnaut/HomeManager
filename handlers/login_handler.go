package handlers

import (
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
	}

	LoginHttpHandler struct {
		LoginUsecase usecases.LoginUseCase
	}
)

func NewLoginHttpHandler(loginUseCase usecases.LoginUseCase) LoginHandler {
	return &LoginHttpHandler{
		LoginUsecase: loginUseCase,
	}
}

func (h *LoginHttpHandler) Login(c echo.Context) error {
	reqBody := new(models.LoginData)
	var err = c.Bind(reqBody)
	if err != nil {
		log.Errorf("Error binding request body: %v", err)
		return response.MakeBasicResponse(c, http.StatusBadRequest, "Bad request")
	}

	loginResult := h.LoginUsecase.Login(reqBody)
	if loginResult.IsError() {
		return response.MakeBasicResponse(c, http.StatusInternalServerError, loginResult.Error)
	}

	return response.MakeLoginResponse(c, http.StatusOK, loginResult.Result)
}