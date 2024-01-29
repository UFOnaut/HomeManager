package entities

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type (
	GroupIds []int64

	LoginUserDto struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	InsertUserDto struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}

	User struct {
		Id        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
		Email     string    `json:"email"`
		GroupIds  GroupIds  `gorm:"type:text" json:"group_ids"`
		Name      string    `json:"name"`
		CreatedAt time.Time `json:"createdAt"`
	}
)

func (groupIds GroupIds) Value() (driver.Value, error) {
	fmt.Println("Value fired ", len(groupIds))
	if len(groupIds) == 0 {
		return sql.NullString{}, nil
	}
	val, err := json.Marshal(groupIds)
	return string(val), err
}

func (groupIds GroupIds) Scan(value interface{}) error {
	fmt.Println("Scan fired ", value.(string))
	return json.Unmarshal([]byte(value.(string)), &groupIds)
}
