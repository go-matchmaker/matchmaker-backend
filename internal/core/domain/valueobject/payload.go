package valueobject

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

const (
	AccessToken  = "access"
	RefreshToken = "refresh"
)

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

type (
	Payload struct {
		ID        uuid.UUID `json:"id"`
		Email     string    `json:"email"`
		Role      string    `json:"role"`
		IsBlocked bool      `json:"is_blocked"`
		IssuedAt  time.Time `json:"issued_at"`
		ExpiredAt time.Time `json:"expired_at"`
	}
)

func NewPayload(userID uuid.UUID, email, role string, isBlocked bool, duration time.Duration) (*Payload, error) {
	payload := &Payload{
		ID:        userID,
		Role:      role,
		Email:     email,
		IsBlocked: isBlocked,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	return payload, nil
}

func (payload *Payload) Valid() error {
	if !time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}
