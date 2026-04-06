package service

import (
	"Toko-Online/config"
	"Toko-Online/dto"
	"Toko-Online/model"
	"Toko-Online/repository"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"
)

type ProductServiceImpl struct {
	repo    repository.ProductRepository
	catRepo repository.CategoryRepository
	logger  *slog.Logger
	redis   *redis.Client
}

const (
	ProductCacheKey      = "products:all"
	ProductDetailKey     = "products:detail:%s"
	ProductCategoryKey   = "products:category:%s"
	ProductCacheDuration = 1 * time.Hour
)

func NewProductService(
	repo repository.ProductRepository,
	catRepo repository.CategoryRepository,
	logger *slog.Logger,
	redis *redis.Client,
) ProductService {
	return &ProductServiceImpl{
		repo:    repo,
		catRepo: catRepo,
		logger:  logger,
		redis:   redis,
	}
}

func (s *ProductServiceImpl) FindAll(ctx context.Context) ([]dto.ProductResponse, error) {
	if s.redis != nil {
		cached, err := s.redis.Get(ctx, ProductCacheKey).Result()
		if err == nil {
			var res []dto.ProductResponse
			if json.Unmarshal([]byte(cached), &res) == nil {
				s.logger.InfoContext(ctx, "Cache hit: FindAll products")
				return res, nil
			}
		}
	}

	s.logger.InfoContext(ctx, "Cache miss: Finding all products from DB")
	products, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, &config.ApiError{Status: 500, Title: "Internal Error", Detail: err.Error()}
	}

	var res []dto.ProductResponse
	for _, p := range products {
		res = append(res, s.mapToResponse(p))
	}

	if s.redis != nil {
		data, _ := json.Marshal(res)
		s.redis.Set(ctx, ProductCacheKey, data, ProductCacheDuration)
	}

	return res, nil
}

func (s *ProductServiceImpl) FindByID(ctx context.Context, id uuid.UUID) (dto.ProductResponse, error) {
	cacheKey := fmt.Sprintf(ProductDetailKey, id.String())

	if s.redis != nil {
		cached, err := s.redis.Get(ctx, cacheKey).Result()
		if err == nil {
			var res dto.ProductResponse
			if json.Unmarshal([]byte(cached), &res) == nil {
				s.logger.InfoContext(ctx, "Cache hit: FindByID product", "id", id)
				return res, nil
			}
		}
	}

	product, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return dto.ProductResponse{}, &config.ApiError{Status: 404, Title: "Not Found", Detail: "Product not found"}
	}

	res := s.mapToResponse(product)

	if s.redis != nil {
		data, _ := json.Marshal(res)
		s.redis.Set(ctx, cacheKey, data, ProductCacheDuration)
	}

	return res, nil
}

func (s *ProductServiceImpl) Create(ctx context.Context, req dto.AddProductRequest) (dto.ProductResponse, error) {
	if req.Price.LessThanOrEqual(decimal.Zero) {
		return dto.ProductResponse{}, &config.ApiError{Status: 400, Title: "Bad Request", Detail: "Price must be > 0"}
	}

	if _, err := s.catRepo.FindByID(ctx, req.CategoryID); err != nil {
		return dto.ProductResponse{}, &config.ApiError{Status: 400, Title: "Bad Request", Detail: "Invalid category"}
	}

	product := &model.Product{
		Name:        req.Name,
		Slug:        generateSlug(req.Name),
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		CategoryID:  req.CategoryID,
		ImageUrl:    req.ImageUrl,
		IsActive:    req.IsActive,
	}

	if err := s.repo.Create(ctx, product); err != nil {
		return dto.ProductResponse{}, err
	}

	s.clearCache(ctx, product.CategoryID.String())

	return s.mapToResponse(*product), nil
}

func (s *ProductServiceImpl) Update(ctx context.Context, id uuid.UUID, req dto.UpdateProductRequest) (dto.ProductResponse, error) {
	product, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return dto.ProductResponse{}, &config.ApiError{Status: 404, Title: "Not Found", Detail: "Product not found"}
	}

	oldCategoryID := product.CategoryID.String()

	if req.CategoryID != uuid.Nil {
		product.CategoryID = req.CategoryID
	}
	if req.Name != "" {
		product.Name = req.Name
		product.Slug = generateSlug(req.Name)
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if !req.Price.IsZero() {
		product.Price = req.Price
	}
	if req.Stock != 0 {
		product.Stock = req.Stock
	}
	if req.ImageUrl != "" {
		product.ImageUrl = req.ImageUrl
	}
	product.IsActive = req.IsActive

	if err := s.repo.Update(ctx, &product); err != nil {
		return dto.ProductResponse{}, err
	}

	s.clearCache(ctx, oldCategoryID)
	if product.CategoryID.String() != oldCategoryID {
		s.clearCache(ctx, product.CategoryID.String())
	}
	if s.redis != nil {
		s.redis.Del(ctx, fmt.Sprintf(ProductDetailKey, id.String()))
	}

	return s.mapToResponse(product), nil
}

func (s *ProductServiceImpl) Delete(ctx context.Context, id uuid.UUID) error {
	product, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	s.clearCache(ctx, product.CategoryID.String())
	if s.redis != nil {
		s.redis.Del(ctx, fmt.Sprintf(ProductDetailKey, id.String()))
	}

	return nil
}

func (s *ProductServiceImpl) FindByCategoryID(ctx context.Context, categoryID uuid.UUID) ([]dto.ProductResponse, error) {
	cacheKey := fmt.Sprintf(ProductCategoryKey, categoryID.String())

	if s.redis != nil {
		cached, err := s.redis.Get(ctx, cacheKey).Result()
		if err == nil {
			var res []dto.ProductResponse
			if json.Unmarshal([]byte(cached), &res) == nil {
				s.logger.InfoContext(ctx, "Cache hit: FindByCategoryID", "category_id", categoryID)
				return res, nil
			}
		}
	}

	products, err := s.repo.FindByCategoryID(ctx, categoryID)
	if err != nil {
		return nil, err
	}

	var res []dto.ProductResponse
	for _, p := range products {
		res = append(res, s.mapToResponse(p))
	}

	if s.redis != nil {
		data, _ := json.Marshal(res)
		s.redis.Set(ctx, cacheKey, data, ProductCacheDuration)
	}

	return res, nil
}

func (s *ProductServiceImpl) clearCache(ctx context.Context, categoryID string) {
	if s.redis == nil {
		return
	}
	s.redis.Del(ctx, ProductCacheKey)
	if categoryID != "" {
		s.redis.Del(ctx, fmt.Sprintf(ProductCategoryKey, categoryID))
	}
}

func (s *ProductServiceImpl) mapToResponse(p model.Product) dto.ProductResponse {
	return dto.ProductResponse{
		ID:          p.ID,
		CategoryID:  p.CategoryID,
		Name:        p.Name,
		Slug:        p.Slug,
		Description: p.Description,
		Price:       p.Price,
		Stock:       p.Stock,
		ImageUrl:    p.ImageUrl,
		IsActive:    p.IsActive,
	}
}

func generateSlug(name string) string {
	return strings.ToLower(strings.Join(strings.Fields(name), "-"))
}
