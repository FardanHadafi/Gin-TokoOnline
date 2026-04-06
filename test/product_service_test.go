package test

import (
	"Toko-Online/dto"
	"Toko-Online/model"
	"Toko-Online/service"
	"Toko-Online/test/mocks"
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProductService_Create(t *testing.T) {
	repo := new(mocks.ProductRepositoryMock)
	catRepo := new(mocks.CategoryRepositoryMock)
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	
	svc := service.NewProductService(repo, catRepo, logger, nil)

	ctx := context.Background()
	catID := uuid.New()
	
	req := dto.AddProductRequest{
		Name:        "Test Product",
		Description: "Description",
		Price:       decimal.NewFromInt(1000),
		Stock:       10,
		CategoryID:  catID,
	}

	catRepo.On("FindByID", ctx, catID).Return(model.Category{ID: catID, Name: "Test Category"}, nil)

	repo.On("Create", ctx, mock.MatchedBy(func(p *model.Product) bool {
		return p.Name == req.Name && p.Slug == "test-product"
	})).Return(nil)

	res, err := svc.Create(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, "Test Product", res.Name)
	assert.Equal(t, "test-product", res.Slug)
	repo.AssertExpectations(t)
	catRepo.AssertExpectations(t)
}

func TestProductService_FindAll_CacheMiss(t *testing.T) {
	repo := new(mocks.ProductRepositoryMock)
	catRepo := new(mocks.CategoryRepositoryMock)
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	svc := service.NewProductService(repo, catRepo, logger, nil)

	ctx := context.Background()
	products := []model.Product{
		{ID: uuid.New(), Name: "P1", Price: decimal.NewFromInt(100), Category: model.Category{Name: "C1"}},
	}

	repo.On("FindAll", ctx).Return(products, nil)

	res, err := svc.FindAll(ctx)

	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, "P1", res[0].Name)
	repo.AssertExpectations(t)
}

func TestProductService_Delete_Success(t *testing.T) {
	repo := new(mocks.ProductRepositoryMock)
	catRepo := new(mocks.CategoryRepositoryMock)
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	
	svc := service.NewProductService(repo, catRepo, logger, nil)

	ctx := context.Background()
	id := uuid.New()

	repo.On("FindByID", ctx, id).Return(model.Product{ID: id, CategoryID: uuid.New()}, nil)
	repo.On("Delete", ctx, id).Return(nil)

	err := svc.Delete(ctx, id)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestProductService_Delete_Error(t *testing.T) {
	repo := new(mocks.ProductRepositoryMock)
	catRepo := new(mocks.CategoryRepositoryMock)
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	
	svc := service.NewProductService(repo, catRepo, logger, nil)

	ctx := context.Background()
	id := uuid.New()

	repo.On("FindByID", ctx, id).Return(model.Product{ID: id, CategoryID: uuid.New()}, nil)
	repo.On("Delete", ctx, id).Return(errors.New("db error"))

	err := svc.Delete(ctx, id)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db error")
	repo.AssertExpectations(t)
}
