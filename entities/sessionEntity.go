package entities

import (
	"time"
)

type (
	SessionDto struct {
		Token string `json:"token"`
	}

	Session struct {
		Id        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
		USerId    string    `json:"user_id"`
		Token     string    `json:"token"`
		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
	}
)
