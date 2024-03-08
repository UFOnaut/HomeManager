package entities

type (
	Group struct {
		ID            string   `json:"id" gorm:"primaryKey;size:16"`
		Name          string   `json:"name"`
		ParticipantId GroupIds `json:"participant_ids"`
	}
)
