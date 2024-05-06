package paseto

import (
	"errors"
	"github.com/go-matchmaker/matchmaker-server/internal/adapter/config"
	"github.com/go-matchmaker/matchmaker-server/internal/core/domain/entity"
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
}

func NewPaseto(cfg *config.Container) (token.TokenMaker, error) {
	symmetricKey := cfg.Token.SymmetricKey
	validSymmetricKey := len(symmetricKey) == chacha20poly1305.KeySize
	if !validSymmetricKey {
		return nil, errors.New("invalid token symmetric key")
	}

	return &PasetoToken{
		paseto.NewV2(),
		[]byte(symmetricKey),
	}, nil
}

func (pt *PasetoToken) CreateToken(email, role string, duration time.Duration) (string, *entity.Payload, error) {
	payload, err := entity.NewPayload(email, role, duration)
	if err != nil {
		return "", nil, err
	}

	generatedToken, err := pt.paseto.Encrypt(pt.symmetricKey, payload, nil)
	return generatedToken, payload, nil
}

func (pt *PasetoToken) VerifyToken(pasetoToken string) (*entity.Payload, error) {
	var payload entity.Payload

	err := pt.paseto.Decrypt(pasetoToken, pt.symmetricKey, &payload, nil)
	if err != nil {
		return nil, err
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return &payload, nil
}
