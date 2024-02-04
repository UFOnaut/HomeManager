package response

import "github.com/labstack/echo/v4"

type LoginResponse struct {
	Token string `json:"token"`
}

func MakeLoginResponse(c echo.Context, responseCode int, token string) error {
	return c.JSON(responseCode, &LoginResponse{
		Token: token,
	})
}