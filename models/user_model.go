package models

type LoginData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}
