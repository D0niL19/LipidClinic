package models

import "time"

type TempUser struct {
	Id             int64     `json:"id"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"hashed_password"`
	Token          string    `json:"token"`
	CreatedAt      time.Time `json:"created_at"`
}
