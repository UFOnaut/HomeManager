package entities

import (
	"gorm.io/gorm"
)

type (
	Group struct {
		gorm.Model
		Name          string `json:"name"`
		ParticipantId []uint `json:"participant_ids"`
	}
)
