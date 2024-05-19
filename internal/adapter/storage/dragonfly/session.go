package dragonfly

import (
	"context"
	"github.com/go-matchmaker/matchmaker-server/internal/core/domain/aggregate"
	"github.com/go-matchmaker/matchmaker-server/internal/core/port/auth"
	"github.com/go-matchmaker/matchmaker-server/internal/core/port/cache"
	"github.com/go-matchmaker/matchmaker-server/internal/core/util"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/google/wire"
)

var (
	_               auth.SessionRepositoryPort = (*SessionCache)(nil)
	SessionCacheSet                            = wire.NewSet(NewSessionCache)
)

type SessionCache struct {
	cache cache.CacheEngine
}

func NewSessionCache(mem cache.CacheEngine) auth.SessionRepositoryPort {
	return &SessionCache{
		cache: mem,
	}
}

func (q *SessionCache) Get(ctx context.Context, id uuid.UUID) (*aggregate.Session, error) {
	cacheKey := util.GenerateCacheKey(id.String(), "session")
	sessionData, err := q.cache.Get(ctx, cacheKey)
	if err != nil {
		return nil, err
	}

	sessionModel := new(aggregate.Session)
	err = json.Unmarshal(sessionData, sessionModel)
	if err != nil {
		return nil, err
	}

	return sessionModel, nil
}

func (q *SessionCache) Set(ctx context.Context, session *aggregate.Session) (*uuid.UUID, error) {
	id := session.ID
	cacheKey := util.GenerateCacheKey(id.String(), "session")
	sessionData, err := json.Marshal(session)
	if err != nil {
		return nil, err
	}

	ttl := session.ExpiredAt.Sub(session.CreatedAt)
	err = q.cache.Set(ctx, cacheKey, sessionData, ttl)
	if err != nil {
		return nil, err
	}

	return &id, nil
}
