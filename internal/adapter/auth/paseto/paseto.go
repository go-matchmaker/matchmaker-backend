package paseto

import (
	"aidanwoods.dev/go-paseto"
	"github.com/go-matchmaker/matchmaker-server/internal/adapter/config"
	"github.com/go-matchmaker/matchmaker-server/internal/core/domain/valueobject"
	"github.com/go-matchmaker/matchmaker-server/internal/core/port/auth"
	"github.com/google/uuid"
	"github.com/google/wire"
	"time"
)

const (
	SymmetricKeySize = 128
)

var (
	_         auth.TokenMaker = (*PasetoToken)(nil)
	PasetoSet                 = wire.NewSet(NewPaseto)
)

type PasetoToken struct {
	tokenTTL   time.Duration
	refreshTTL time.Duration
}

func NewPaseto(cfg *config.Container) (auth.TokenMaker, error) {
	tokenDuration := cfg.Token.TokenTTL
	refreshDuration := cfg.Token.RefreshTTL

	return &PasetoToken{
		tokenTTL:   tokenDuration,
		refreshTTL: refreshDuration,
	}, nil
}

func (pt *PasetoToken) CreateToken(id uuid.UUID, email, role string, isBlocked bool) (string, string, *valueobject.Payload, error) {
	duration := pt.tokenTTL
	payload, err := valueobject.NewPayload(id, email, role, isBlocked, duration)
	if err != nil {
		return "", "", nil, err
	}

	tokenPaseto := paseto.NewToken()
	tokenPaseto.SetExpiration(payload.ExpiredAt)
	tokenPaseto.SetIssuedAt(payload.IssuedAt)
	tokenPaseto.SetString("id", payload.ID.String())
	tokenPaseto.SetString("role", payload.Role)
	tokenPaseto.SetString("email", payload.Email)
	tokenPaseto.SetString("role", payload.Role)
	tokenPaseto.Set("is_blocked", payload.IsBlocked)
	secretKey := paseto.NewV4AsymmetricSecretKey()
	publicKey := secretKey.Public().ExportHex()
	encrypted := tokenPaseto.V4Sign(secretKey, nil)

	return encrypted, publicKey, payload, nil
}

func (pt *PasetoToken) CreateRefreshToken(payload *valueobject.Payload) (string, string, *valueobject.Payload, error) {
	duration := pt.refreshTTL
	tokenPaseto := paseto.NewToken()
	payload.ExpiredAt = payload.ExpiredAt.Add(duration)
	tokenPaseto.SetExpiration(payload.ExpiredAt)
	tokenPaseto.SetIssuedAt(payload.IssuedAt)
	tokenPaseto.SetString("id", payload.ID.String())
	tokenPaseto.SetString("role", payload.Role)
	tokenPaseto.SetString("email", payload.Email)
	tokenPaseto.SetString("role", payload.Role)

	secretKey := paseto.NewV4AsymmetricSecretKey()
	publicKey := secretKey.Public().ExportHex()
	encrypted := tokenPaseto.V4Sign(secretKey, nil)
	return encrypted, publicKey, payload, nil
}

func (pt *PasetoToken) DecodeToken(pasetoToken, publicKeyHex string) (*valueobject.Payload, error) {
	publicKey, err := paseto.NewV4AsymmetricPublicKeyFromHex(publicKeyHex)
	if err != nil {
		return nil, err
	}

	parser := paseto.NewParser()
	parsedToken, err := parser.ParseV4Public(publicKey, pasetoToken, nil)
	if err != nil {
		return nil, err
	}

	payload := new(valueobject.Payload)
	expiredAt, err := parsedToken.GetExpiration()
	if err != nil {
		return nil, err
	}
	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	issuedAt, err := parsedToken.GetIssuedAt()
	if err != nil {
		return nil, err
	}

	idHex, err := parsedToken.GetString("id")
	if err != nil {
		return nil, err
	}
	id := uuid.MustParse(idHex)
	email, err := parsedToken.GetString("email")
	if err != nil {
		return nil, err
	}
	role, err := parsedToken.GetString("role")
	if err != nil {
		return nil, err
	}
	var isBlocked bool
	err = parsedToken.Get("is_blocked", isBlocked)

	if err != nil {
		return nil, err

	}
	payload = &valueobject.Payload{
		ID:        id,
		Role:      role,
		Email:     email,
		IssuedAt:  issuedAt,
		ExpiredAt: expiredAt,
	}
	return payload, nil

}
