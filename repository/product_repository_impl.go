package repository

import (
	"Toko-Online/model"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductRepositoryImpl struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &ProductRepositoryImpl{db: db}
}

func (r *ProductRepositoryImpl) FindAll(ctx context.Context) ([]model.Product, error) {
	var products []model.Product
	err := r.db.WithContext(ctx).Preload("Category").Find(&products).Error
	return products, err
}

func (r *ProductRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (model.Product, error) {
	var product model.Product
	err := r.db.WithContext(ctx).Preload("Category").First(&product, "id = ?", id).Error
	return product, err
}

func (r *ProductRepositoryImpl) Create(ctx context.Context, product *model.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

func (r *ProductRepositoryImpl) Update(ctx context.Context, product *model.Product) error {
	return r.db.WithContext(ctx).Model(product).Updates(product).Error
}

func (r *ProductRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.Product{}, "id = ?", id).Error
}

func (r *ProductRepositoryImpl) FindByCategoryID(ctx context.Context, categoryID uuid.UUID) ([]model.Product, error) {
	var products []model.Product
	err := r.db.WithContext(ctx).Where("category_id = ?", categoryID).Find(&products).Error
	return products, err
}
