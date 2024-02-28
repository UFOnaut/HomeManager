package handlers

import (
	"home_manager/repositories"
	"net/http"
	"strconv"

	"home_manager/handlers/response"
	"home_manager/models"
	"home_manager/usecases"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type (
	RegisterHandler interface {
		Register(c echo.Context) error
		VerifyEmail(c echo.Context) error
	}

	RegisterHttpHandler struct {
		RegisterUseCase    usecases.RegisterUseCase
		VerifyEmailUseCase usecases.VerifyEmailUseCase
	}
)

func NewRegisterHttpHandler(repository repositories.UserRepository) RegisterHandler {
	return &RegisterHttpHandler{
		RegisterUseCase:    usecases.NewRegisterUseCase(repository),
		VerifyEmailUseCase: usecases.NewVerifyEmailUseCase(repository),
	}
}

func (h *RegisterHttpHandler) Register(c echo.Context) error {
	reqBody := new(models.RegisterData)
	var err = c.Bind(reqBody)
	if err != nil {
		log.Errorf("Error binding request body: %v", err)
		return response.MakeBasicResponse(c, http.StatusBadRequest, "Bad request")
	}

	registerResult := h.RegisterUseCase.Execute(reqBody)
	if registerResult.IsError() {
		return response.MakeBasicResponse(c, http.StatusInternalServerError, registerResult.Error)
	}

	return response.MakeBasicResponse(c, http.StatusOK, "Success")
}

func (h *RegisterHttpHandler) VerifyEmail(c echo.Context) error {
	params := c.QueryParams()
	if params == nil {
		log.Errorf("Error binding query params")
		return response.MakeBasicResponse(c, http.StatusBadRequest, "Bad request")
	}

	userIdParam := params.Get("user_id")
	verifyToken := params.Get("verify_token")

	userId, err := strconv.Atoi(userIdParam)

	if err != nil {
		return err
	}

	verifyResult := h.VerifyEmailUseCase.Execute(uint(userId), verifyToken)
	if verifyResult.IsError() {
		return response.MakeBasicResponse(c, http.StatusInternalServerError, verifyResult.Error)
	}

	return response.MakeBasicResponse(c, http.StatusOK, "Success")
}
