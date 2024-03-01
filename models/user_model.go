package models

type LoginData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RefreshToken struct {
	RefreshToken string `json:"refresh_token"`
}

type RegisterData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}
