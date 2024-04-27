package token

import (
	"github.com/go-matchmaker/matchmaker-server/internal/core/domain/valueobject"
)

type TokenMaker interface {
	CreateToken(email string, role string) (string, *valueobject.TokenPayload, error)
	VerifyToken(token string) (*valueobject.TokenPayload, error)
}
