// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0

package db

import (
	"time"
)

type User struct {
	Username       string    `json:"username"`
	HashedPassword string    `json:"hashed_password"`
	FullName       string    `json:"full_name"`
	Email          string    `json:"email"`
	CreatedAt      time.Time `json:"created_at"`
}
