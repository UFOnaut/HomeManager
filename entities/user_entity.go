package entities

import (
	"crypto/sha256"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"gorm.io/gorm"
)

type (
	GroupIds []string

	User struct {
		ID       string   `json:"id" gorm:"primaryKey;size:16"`
		Email    string   `json:"email" validate:"required,email"`
		Password string   `json:"password"`
		GroupIds GroupIds `json:"group_ids" gorm:"type:text" `
		Name     string   `json:"name"`
		Verified bool     `json:"verified"`
	}

	VerificationToken struct {
		ID     string `json:"id" gorm:"primaryKey;size:16"`
		UserId string `json:"user_id"`
		Token  string `json:"token"`
	}
)

func (groupIds GroupIds) Value() (driver.Value, error) {
	if len(groupIds) == 0 {
		return sql.NullString{}, nil
	}
	val, err := json.Marshal(groupIds)
	return string(val), err
}

func (groupIds GroupIds) Scan(value interface{}) error {
	return json.Unmarshal([]byte(value.(string)), &groupIds)
}

func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	u.Password = Hash(u.Password)
	return nil
}

func (u *User) IsPasswordCorrect(password string) bool {
	return Hash(password) == u.Password
}

func Hash(in string) string {
	h := sha256.Sum256([]byte(in))
	return fmt.Sprintf("%x", h[:])
}
