package service

import (
	"Toko-Online/dto"
	"context"

	"github.com/google/uuid"
)

type ProductService interface {
	FindAll(ctx context.Context) ([]dto.ProductResponse, error)
	FindByID(ctx context.Context, id uuid.UUID) (dto.ProductResponse, error)
	Create(ctx context.Context, req dto.AddProductRequest) (dto.ProductResponse, error)
	Update(ctx context.Context, id uuid.UUID, req dto.UpdateProductRequest) (dto.ProductResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
	FindByCategoryID(ctx context.Context, categoryID uuid.UUID) ([]dto.ProductResponse, error)
}
