package entities

import (
	"crypto/sha256"

	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type (
	GroupIds []int64

	User struct {
		gorm.Model
		Id        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
		Email     string    `json:"email"`
		Password  string    `json:"password"`
		GroupIds  GroupIds  `gorm:"type:text" json:"group_ids"`
		Name      string    `json:"name"`
		CreatedAt time.Time `json:"createdAt"`
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

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.Password = hash(u.Password)
	return nil
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
