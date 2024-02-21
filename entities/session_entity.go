package entities

import (
	"gorm.io/gorm"
)

type (
	Session struct {
		gorm.Model
		UserId uint   `json:"user_id"`
		Token  string `json:"token"`
	}
)
