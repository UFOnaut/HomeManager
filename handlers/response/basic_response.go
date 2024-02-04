package response

import "github.com/labstack/echo/v4"

type BasicResponse struct {
	Message string `json:"message"`
}

func MakeBasicResponse(c echo.Context, responseCode int, message string) error {
	return c.JSON(responseCode, &BasicResponse{
		Message: message,
	})
}
