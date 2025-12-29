package models

// User struct for stor
type User struct {
	ID           string `json:"id" db:"id"`
	Login        string `json:"login" db:"login"`
	PasswordHash string `json:"-" db:"password_hash"`
}

// Struct for register user
type UserRegister struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// Struct for login
type UserLogin struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
