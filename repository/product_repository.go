package repository

import (
	"Toko-Online/model"
	"context"

	"github.com/google/uuid"
)

type ProductRepository interface {
	FindAll(ctx context.Context) ([]model.Product, error)
	FindByID(ctx context.Context, id uuid.UUID) (model.Product, error)
	Create(ctx context.Context, product *model.Product) error
	Update(ctx context.Context, product *model.Product) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindByCategoryID(ctx context.Context, categoryID uuid.UUID) ([]model.Product, error)
}
