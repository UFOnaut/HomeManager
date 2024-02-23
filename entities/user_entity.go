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
	GroupIds []uint

	User struct {
		gorm.Model
		Email    string   `json:"email"`
		Password string   `json:"password"`
		GroupIds GroupIds `json:"group_ids" gorm:"type:text" `
		Name     string   `json:"name"`
		Verified bool     `json:"verified"`
	}

	VerificationToken struct {
		gorm.Model
		UserId uint   `json:"user_id"`
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
	u.Password = hash(u.Password)
	return nil
}

func (u *User) IsPasswordCorrect(password string) bool {
	return hash(password) == u.Password
}

func hash(in string) string {
	h := sha256.Sum256([]byte(in))
	return fmt.Sprintf("%x", h[:])
}
