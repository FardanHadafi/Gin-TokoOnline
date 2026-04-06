package test

import (
	"Toko-Online/dto"
	"Toko-Online/model"
	"Toko-Online/service"
	"Toko-Online/test/mocks"
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCategoryService_Create(t *testing.T) {
	repo := new(mocks.CategoryRepositoryMock)
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	
	svc := service.NewCategoryService(repo, logger, nil)

	ctx := context.Background()
	req := dto.UpdateCategoryRequest{
		Name: "New Category",
	}

	repo.On("Create", ctx, mock.MatchedBy(func(c *model.Category) bool {
		return c.Name == "New Category"
	})).Return(nil)

	res, err := svc.Create(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, "New Category", res.Name)
	repo.AssertExpectations(t)
}

func TestCategoryService_FindAll(t *testing.T) {
	repo := new(mocks.CategoryRepositoryMock)
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	
	svc := service.NewCategoryService(repo, logger, nil)

	ctx := context.Background()
	categories := []model.Category{
		{ID: uuid.New(), Name: "Cat 1"},
		{ID: uuid.New(), Name: "Cat 2"},
	}

	repo.On("FindAll", ctx).Return(categories, nil)

	res, err := svc.FindAll(ctx)

	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Equal(t, "Cat 1", res[0].Name)
	repo.AssertExpectations(t)
}

func TestCategoryService_Update(t *testing.T) {
	repo := new(mocks.CategoryRepositoryMock)
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	
	svc := service.NewCategoryService(repo, logger, nil)

	ctx := context.Background()
	id := uuid.New()
	req := dto.UpdateCategoryRequest{
		Name: "Updated Name",
	}

	repo.On("FindByID", ctx, id).Return(model.Category{ID: id, Name: "Old Name"}, nil)
	repo.On("Update", ctx, mock.MatchedBy(func(c *model.Category) bool {
		return c.Name == "Updated Name"
	})).Return(nil)

	res, err := svc.Update(ctx, id, req)

	assert.NoError(t, err)
	assert.Equal(t, "Updated Name", res.Name)
	repo.AssertExpectations(t)
}

func TestCategoryService_Delete(t *testing.T) {
	repo := new(mocks.CategoryRepositoryMock)
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	
	svc := service.NewCategoryService(repo, logger, nil)

	ctx := context.Background()
	id := uuid.New()

	repo.On("FindByID", ctx, id).Return(model.Category{ID: id, Name: "To Delete"}, nil)
	repo.On("Delete", ctx, id).Return(nil)

	err := svc.Delete(ctx, id)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}
