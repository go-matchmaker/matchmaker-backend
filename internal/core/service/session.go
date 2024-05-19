package service

import (
	"context"
	"github.com/go-matchmaker/matchmaker-server/internal/core/domain/aggregate"
	"github.com/go-matchmaker/matchmaker-server/internal/core/port/auth"
	"github.com/google/uuid"
	"github.com/google/wire"
)

var (
	_                 auth.SessionServicePort = (*SessionService)(nil)
	SessionServiceSet                         = wire.NewSet(NewSessionService)
)

type SessionService struct {
	authRepo auth.SessionRepositoryPort
}

func NewSessionService(authRepo auth.SessionRepositoryPort) auth.SessionServicePort {
	return &SessionService{
		authRepo: authRepo,
	}
}

func (s *SessionService) GetUserSession(ctx context.Context, id uuid.UUID) (*aggregate.Session, error) {
	return s.authRepo.Get(ctx, id)
}

func (s *SessionService) AddSession(ctx context.Context, session *aggregate.Session) (*uuid.UUID, error) {
	return s.authRepo.Set(ctx, session)
}
