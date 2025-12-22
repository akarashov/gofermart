package models

type User struct {
	ID           string `json:"id" db:"id"`
	Login        string `json:"login" db:"login"`
	PasswordHash string `json:"-" db:"password_hash"`
}

type UserRegister struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserLogin struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
