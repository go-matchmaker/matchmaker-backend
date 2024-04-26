package token

import (
	"github.com/bulutcan99/company-matcher/internal/core/domain"
)

type TokenMaker interface {
	CreateToken(email string, role string) (string, *domain.TokenPayload, error)
	VerifyToken(token string) (*domain.TokenPayload, error)
}
