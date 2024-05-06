// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package user

import (
	"time"

	"github.com/google/uuid"
)

type Users struct {
	ID           uuid.UUID `json:"id"`
	Role         string    `json:"role"`
	Name         string    `json:"name"`
	Surname      string    `json:"surname"`
	Email        string    `json:"email"`
	PhoneNumber  string    `json:"phone_number"`
	PasswordHash string    `json:"password_hash"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
