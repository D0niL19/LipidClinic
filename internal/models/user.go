package models

import "time"

type User struct {
	Id             int64     `json:"id"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"hashed_password"`
	Role           string    `json:"role"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
