package token

import (
	"github.com/go-matchmaker/matchmaker-server/internal/core/domain/entity"
	"time"
)

type TokenMaker interface {
	CreateToken(email string, role string, duration time.Duration) (string, *entity.Payload, error)
	VerifyToken(token string) (*entity.Payload, error)
}
