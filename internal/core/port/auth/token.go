package auth

import (
	"context"
	"github.com/go-matchmaker/matchmaker-server/internal/core/domain/aggregate"
	"github.com/go-matchmaker/matchmaker-server/internal/core/domain/valueobject"
	"github.com/google/uuid"
)

type TokenMaker interface {
	CreateToken(id uuid.UUID, email, role string, isBlocked bool) (string, string, *valueobject.Payload, error)
	CreateRefreshToken(payload *valueobject.Payload) (string, string, *valueobject.Payload, error)
	DecodeToken(token, public string) (*valueobject.Payload, error)
}

type SessionRepositoryPort interface {
	Get(ctx context.Context, id uuid.UUID) (*aggregate.Session, error)
	Set(ctx context.Context, session *aggregate.Session) (*uuid.UUID, error)
	//Update(ctx context.Context, session *aggregate.Session) (*aggregate.Session, error)
	//DeleteOne(ctx context.Context, id uuid.UUID) error
}

type SessionServicePort interface {
	GetUserSession(ctx context.Context, id uuid.UUID) (*aggregate.Session, error)
	AddSession(ctx context.Context, session *aggregate.Session) (*uuid.UUID, error)
	//UpdateSession(ctx context.Context, session *aggregate.Session) (*aggregate.Session, error)
	//DeleteSession(ctx context.Context, id uuid.UUID) error
}
