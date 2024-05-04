package repository

import (
	"context"
	"github.com/go-matchmaker/matchmaker-server/internal/core/domain/entity"
	"github.com/google/uuid"
)

type UserPort interface {
	Insert(ctx context.Context, user *entity.User) (*uuid.UUID, error)
	Update(ctx context.Context, user *entity.User) (*entity.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (entity.User, error)
	GetByEmail(ctx context.Context, email string) (entity.User, error)
	DeleteOne(ctx context.Context, id uuid.UUID) error
	DeleteAll(ctx context.Context) error
}
