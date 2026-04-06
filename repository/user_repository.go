package repository

import (
	"Toko-Online/model"
	"context"

	"github.com/google/uuid"
)

type UserRepository interface {
	FindByUsername(ctx context.Context, username string) (model.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (model.User, error)
	Update(ctx context.Context, user *model.User) error
}
