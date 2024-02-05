package entities

import (
	"time"

	"gorm.io/gorm"
)

const SessionLifeMinutes = 30

type (
	Session struct {
		gorm.Model
		UserId uint   `json:"user_id"`
		Token  string `json:"token"`
	}
)

func (session *Session) IsValid() bool {
	return (session.CreatedAt.Add(time.Minute*SessionLifeMinutes).After(time.Now()) ||
		session.UpdatedAt.Add(time.Minute*SessionLifeMinutes).After(time.Now())) &&
		session.Token != ""
}
