package repository

import (
	"context"
	"github.com/bulutcan99/company-matcher/internal/core/domain/entity"
	"github.com/google/uuid"
)

type UserMaker interface {
	Insert(ctx context.Context, user *entity.User) (*uuid.UUID, error)
	//Update(ctx context.Context, user *domain.User) error
	//GetUserByID(ctx context.Context, id string) (domain.User, error)
	//GetUserByEmail(ctx context.Context, email string) (domain.User, error)
}
