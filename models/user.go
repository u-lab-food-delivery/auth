package models

type User struct {
	UserId         string `json:"user_id"`
	Email          string `json:"email"`
	HashedPassword string `json:"hashed_password"`
	Name           string `json:"name"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}
