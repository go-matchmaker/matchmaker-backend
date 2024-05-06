package entity

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID           uuid.UUID
	Role         string
	Name         string
	Surname      string
	Email        string
	PhoneNumber  string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
