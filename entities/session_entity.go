package entities

type (
	Session struct {
		ID           string `json:"id" gorm:"primaryKey;size:16"`
		UserId       string `json:"user_id"`
		AuthToken    string `json:"auth_token"`
		RefreshToken string `json:"refresh_token"`
	}
)
