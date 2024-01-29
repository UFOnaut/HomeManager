package handlers

import (
	"net/http"

	"home_manager/models"
	"home_manager/usecases"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type (
	UserHandler interface {
		Login(c echo.Context) error
	}

	userHttpHandler struct {
		UserUsecase usecases.UserUsecase
	}
)

func NewUserHttpHandler(userUsecase usecases.UserUsecase) UserHandler {
	return &userHttpHandler{
		UserUsecase: userUsecase,
	}
}

func (h *userHttpHandler) Login(c echo.Context) error {
	reqBody := new(models.LoginData)
	var err = c.Bind(reqBody)
	if err != nil {
		log.Errorf("Error binding request body: %v", err)
		return makeBasicResponse(c, http.StatusBadRequest, "Bad request")
	}

	token, err := h.UserUsecase.Login(reqBody)
	if err != nil {
		return makeBasicResponse(c, http.StatusInternalServerError, "Processing data failed")
	}

	return makeLoginResponse(c, http.StatusOK, token)
}
