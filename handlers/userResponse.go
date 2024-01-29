package handlers

import "github.com/labstack/echo/v4"

type BaseResponse struct {
	Message string `json:"message"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func makeBasicResponse(c echo.Context, responseCode int, message string) error {
	return c.JSON(responseCode, &BaseResponse{
		Message: message,
	})
}

func makeLoginResponse(c echo.Context, responseCode int, token string) error {
	return c.JSON(responseCode, &LoginResponse{
		Token: token,
	})
}
