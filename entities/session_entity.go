package entities

import (
	"gorm.io/gorm"
)

type (
	Session struct {
		gorm.Model
		UserId       uint   `json:"user_id"`
		AuthToken    string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}
)
