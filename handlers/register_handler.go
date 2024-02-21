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
	RegisterHandler interface {
		Register(c echo.Context) error
	}

	RegisterHttpHandler struct {
		RegisterUseCase usecases.RegisterUseCase
	}
)

func NewRegisterHttpHandler(repository repositories.UserRepository) RegisterHandler {
	return &RegisterHttpHandler{
		RegisterUseCase: usecases.NewRegisterUseCase(repository),
	}
}

func (h *RegisterHttpHandler) Register(c echo.Context) error {
	reqBody := new(models.RegisterData)
	var err = c.Bind(reqBody)
	if err != nil {
		log.Errorf("Error binding request body: %v", err)
		return response.MakeBasicResponse(c, http.StatusBadRequest, "Bad request")
	}

	registerResult := h.RegisterUseCase.Register(reqBody)
	if registerResult.IsError() {
		return response.MakeBasicResponse(c, http.StatusInternalServerError, registerResult.Error)
	}

	return response.MakeBasicResponse(c, http.StatusOK, "Success")
}
