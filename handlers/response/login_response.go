package response

import (
	"github.com/labstack/echo/v4"
	"home_manager/entities"
)

type LoginResponse struct {
	AuthToken    string `json:"auth_token"`
	RefreshToken string `json:"refresh_token"`
}

func MakeLoginResponse(c echo.Context, responseCode int, session entities.Session) error {
	return c.JSON(responseCode, &LoginResponse{
		AuthToken:    session.AuthToken,
		RefreshToken: session.RefreshToken,
	})
}
