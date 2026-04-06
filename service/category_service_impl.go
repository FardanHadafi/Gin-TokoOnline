package service

import (
	"Toko-Online/config"
	"Toko-Online/dto"
	"Toko-Online/model"
	"Toko-Online/repository"
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type CategoryServiceImpl struct {
	repo   repository.CategoryRepository
	logger *slog.Logger
	redis  *redis.Client
}

const (
	CategoryCacheKey      = "categories:all"
	CategoryCacheDuration = 24 * time.Hour
)

func NewCategoryService(
	repo repository.CategoryRepository,
	logger *slog.Logger,
	redis *redis.Client,
) CategoryService {
	return &CategoryServiceImpl{
		repo:   repo,
		logger: logger,
		redis:  redis,
	}
}

func (s *CategoryServiceImpl) FindAll(ctx context.Context) ([]dto.CategoryResponse, error) {
	if s.redis != nil {
		cached, err := s.redis.Get(ctx, CategoryCacheKey).Result()
		if err == nil {
			var res []dto.CategoryResponse
			if json.Unmarshal([]byte(cached), &res) == nil {
				s.logger.InfoContext(ctx, "Cache hit: FindAll categories")
				return res, nil
			}
		}
	}

	categories, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, &config.ApiError{Status: 500, Title: "Internal Error", Detail: err.Error()}
	}

	var res []dto.CategoryResponse
	for _, c := range categories {
		res = append(res, dto.CategoryResponse{
			ID:   c.ID,
			Name: c.Name,
		})
	}

	if s.redis != nil {
		data, _ := json.Marshal(res)
		s.redis.Set(ctx, CategoryCacheKey, data, CategoryCacheDuration)
	}

	return res, nil
}

func (s *CategoryServiceImpl) FindByID(ctx context.Context, id uuid.UUID) (dto.CategoryResponse, error) {
	category, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return dto.CategoryResponse{}, &config.ApiError{Status: 404, Title: "Not Found", Detail: "Category not found"}
	}

	return dto.CategoryResponse{
		ID:   category.ID,
		Name: category.Name,
	}, nil
}

func (s *CategoryServiceImpl) Create(ctx context.Context, req dto.UpdateCategoryRequest) (dto.CategoryResponse, error) {
	category := &model.Category{
		Name: req.Name,
	}

	if err := s.repo.Create(ctx, category); err != nil {
		return dto.CategoryResponse{}, err
	}

	if s.redis != nil {
		s.redis.Del(ctx, CategoryCacheKey)
	}

	return dto.CategoryResponse{
		ID:   category.ID,
		Name: category.Name,
	}, nil
}

func (s *CategoryServiceImpl) Update(ctx context.Context, id uuid.UUID, req dto.UpdateCategoryRequest) (dto.CategoryResponse, error) {
	category, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return dto.CategoryResponse{}, &config.ApiError{Status: 404, Title: "Not Found", Detail: "Category not found"}
	}

	if req.Name != "" {
		category.Name = req.Name
	}

	if err := s.repo.Update(ctx, &category); err != nil {
		return dto.CategoryResponse{}, err
	}

	if s.redis != nil {
		s.redis.Del(ctx, CategoryCacheKey)
	}

	return dto.CategoryResponse{
		ID:   category.ID,
		Name: category.Name,
	}, nil
}

func (s *CategoryServiceImpl) Delete(ctx context.Context, id uuid.UUID) error {
	if _, err := s.repo.FindByID(ctx, id); err != nil {
		return err
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	if s.redis != nil {
		s.redis.Del(ctx, CategoryCacheKey)
	}

	return nil
}
