package paseto

import (
	"errors"
	"github.com/go-matchmaker/matchmaker-server/internal/adapter/config"
	"github.com/go-matchmaker/matchmaker-server/internal/core/domain/valueobject"
	"github.com/go-matchmaker/matchmaker-server/internal/core/port/token"
	"github.com/google/wire"
	"time"

	"github.com/o1egl/paseto"
	"golang.org/x/crypto/chacha20poly1305"
)

var (
	_         token.TokenMaker = (*PasetoToken)(nil)
	PasetoSet                  = wire.NewSet(NewPaseto)
)

type PasetoToken struct {
	paseto       *paseto.V2
	symmetricKey []byte
	TTL          time.Duration
}

func NewPaseto(cfg *config.Container) (token.TokenMaker, error) {
	symmetricKey := cfg.Token.SymmetricKey
	durationStr := cfg.Token.TTL

	validSymmetricKey := len(symmetricKey) == chacha20poly1305.KeySize
	if !validSymmetricKey {
		return nil, errors.New("invalid token symmetric key")
	}

	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return nil, err
	}

	return &PasetoToken{
		paseto.NewV2(),
		[]byte(symmetricKey),
		duration,
	}, nil
}

func (pt *PasetoToken) CreateToken(email, role string) (string, *valueobject.TokenPayload, error) {
	payload := valueobject.TokenPayload{
		Email:     email,
		Role:      role,
		Duration:  pt.TTL,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(pt.TTL),
	}

	generatedToken, err := pt.paseto.Encrypt(pt.symmetricKey, payload, nil)
	if err != nil {
		return "", nil, err
	}

	return generatedToken, &payload, nil

}

func (pt *PasetoToken) VerifyToken(pasetoToken string) (*valueobject.TokenPayload, error) {
	var payload valueobject.TokenPayload

	err := pt.paseto.Decrypt(pasetoToken, pt.symmetricKey, &payload, nil)
	if err != nil {
		return nil, err
	}

	isExpired := time.Now().After(payload.ExpiredAt)
	if isExpired {
		return nil, errors.New("token is expired")
	}

	return &payload, nil
}
