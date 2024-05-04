package entity

import (
	"github.com/google/uuid"
	"time"
)

type UserRole string

const (
	TypeA = iota
	TypeB
	TypeC
	TypeD
)

const (
	UserRoleAdmin    UserRole = "admin"
	UserRoleCustomer UserRole = "customer"
)

var CompanyTypes = map[int]string{
	TypeA: "Type A",
	TypeB: "Type B",
	TypeC: "Type C",
	TypeD: "Type D",
}

type User struct {
	ID           uuid.UUID
	UserRole     UserRole
	Name         string
	Surname      string
	Email        string
	PhoneNumber  string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
