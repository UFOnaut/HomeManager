package handlers

import "github.com/labstack/echo/v4"

type BaseResponse struct {
	Message string `json:"message"`
}

func makeBasicResponse(c echo.Context, responseCode int, message string) error {
	return c.JSON(responseCode, &BaseResponse{
		Message: message,
	})
}
