package entity

import (
	"github.com/google/uuid"
	"time"
)

type UserRole string

const (
	UserRoleAdmin    UserRole = "admin"
	UserRoleCustomer UserRole = "customer"
)

type User struct {
	ID             uuid.UUID
	UserRole       UserRole
	Name           string
	Surname        string
	Email          string
	PhoneNumber    string
	CompanyName    string
	CompanyType    int
	CompanyWebSite string
	PasswordHash   string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
