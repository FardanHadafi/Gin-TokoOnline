package service

import (
	"Toko-Online/dto"
	"context"

	"github.com/google/uuid"
)

type CategoryService interface {
	FindAll(ctx context.Context) ([]dto.CategoryResponse, error)
	FindByID(ctx context.Context, id uuid.UUID) (dto.CategoryResponse, error)
	Create(ctx context.Context, req dto.UpdateCategoryRequest) (dto.CategoryResponse, error)
	Update(ctx context.Context, id uuid.UUID, req dto.UpdateCategoryRequest) (dto.CategoryResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
