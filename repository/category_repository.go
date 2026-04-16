package repository

import (
	"Toko-Online/model"
	"context"

	"github.com/google/uuid"
)

type CategoryRepository interface {
	FindAll(ctx context.Context) ([]model.Category, error)
	FindByID(ctx context.Context, id uuid.UUID) (model.Category, error)
	Create(ctx context.Context, category *model.Category) error
	Update(ctx context.Context, category *model.Category) error
	Delete(ctx context.Context, id uuid.UUID) error
	CountProductsByCategory(ctx context.Context) (map[uuid.UUID]int, error)
}
