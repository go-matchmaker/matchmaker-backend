package valueobject

import (
	"time"
)

type TokenPayload struct {
	Email     string
	Role      string
	Duration  time.Duration
	IssuedAt  time.Time
	ExpiredAt time.Time
}
