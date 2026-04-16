package repository

import (
	"Toko-Online/model"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryRepositoryImpl struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &CategoryRepositoryImpl{db: db}
}

func (r *CategoryRepositoryImpl) FindAll(ctx context.Context) ([]model.Category, error) {
	var categories []model.Category
	err := r.db.WithContext(ctx).Find(&categories).Error
	return categories, err
}

func (r *CategoryRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (model.Category, error) {
	var category model.Category
	err := r.db.WithContext(ctx).First(&category, "id = ?", id).Error
	return category, err
}

func (r *CategoryRepositoryImpl) Create(ctx context.Context, category *model.Category) error {
	return r.db.WithContext(ctx).Create(category).Error
}

func (r *CategoryRepositoryImpl) Update(ctx context.Context, category *model.Category) error {
	return r.db.WithContext(ctx).Model(category).Updates(category).Error
}

func (r *CategoryRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.Category{}, "id = ?", id).Error
}

func (r *CategoryRepositoryImpl) CountProductsByCategory(ctx context.Context) (map[uuid.UUID]int, error) {
	type countResult struct {
		CategoryID uuid.UUID `gorm:"column:category_id"`
		Count      int       `gorm:"column:count"`
	}
	var results []countResult
	err := r.db.WithContext(ctx).
		Table("products").
		Select("category_id, COUNT(*) as count").
		Where("deleted_at IS NULL").
		Group("category_id").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	counts := make(map[uuid.UUID]int)
	for _, r := range results {
		counts[r.CategoryID] = r.Count
	}
	return counts, nil
}
