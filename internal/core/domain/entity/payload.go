package entity

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

type Payload struct {
	ID        uuid.UUID
	Email     string
	Role      string
	Duration  time.Duration
	IssuedAt  time.Time
	ExpiredAt time.Time
}

func NewPayload(email string, role string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:        tokenID,
		Email:     email,
		Role:      role,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	return payload, nil
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}
