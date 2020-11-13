package user

import "time"

type User struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TinyUser struct {
	Name     string
	Email    string
	Password string
}

type Claims struct {
	Email string `json:"email"`
}
