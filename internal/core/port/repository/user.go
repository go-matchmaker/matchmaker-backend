package repository

import (
	"context"
	"github.com/go-matchmaker/matchmaker-server/internal/core/domain/entity"
	"github.com/google/uuid"
)

type UserMaker interface {
	Insert(ctx context.Context, user *entity.User) (*uuid.UUID, error)
	Update(ctx context.Context, user *entity.User) (*entity.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (entity.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
}
