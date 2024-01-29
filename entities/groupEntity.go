package entities

import (
	"time"
)

type (
	Group struct {
		Id            int64     `gorm:"primaryKey;autoIncrement" json:"id"`
		Name          string    `json:"name"`
		ParticipantId []int64   `json:"participant_ids"`
		CreatedAt     time.Time `json:"createdAt"`
	}
)
