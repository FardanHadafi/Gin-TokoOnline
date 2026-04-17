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
	err := r.db.WithContext(ctx).Preload("Category").Preload("Images").Find(&products).Error
	return products, err
}

func (r *ProductRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (model.Product, error) {
	var product model.Product
	err := r.db.WithContext(ctx).Preload("Category").Preload("Images").First(&product, "id = ?", id).Error
	return product, err
}

func (r *ProductRepositoryImpl) Create(ctx context.Context, product *model.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

func (r *ProductRepositoryImpl) Update(ctx context.Context, product *model.Product) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(product).Association("Images").Replace(product.Images); err != nil {
			return err
		}
		return tx.Session(&gorm.Session{FullSaveAssociations: false}).Save(product).Error
	})
}

func (r *ProductRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.Product{}, "id = ?", id).Error
}

func (r *ProductRepositoryImpl) FindByCategoryID(ctx context.Context, categoryID uuid.UUID) ([]model.Product, error) {
	var products []model.Product
	err := r.db.WithContext(ctx).Preload("Category").Preload("Images").Where("category_id = ?", categoryID).Find(&products).Error
	return products, err
}
func (r *ProductRepositoryImpl) FindByCategoryIDUnscoped(ctx context.Context, categoryID uuid.UUID) ([]model.Product, error) {
	var products []model.Product
	err := r.db.WithContext(ctx).Unscoped().Where("category_id = ?", categoryID).Find(&products).Error
	return products, err
}

func (r *ProductRepositoryImpl) HardDeleteByCategoryID(ctx context.Context, categoryID uuid.UUID) error {
	var productIDs []uuid.UUID
	if err := r.db.WithContext(ctx).Unscoped().Model(&model.Product{}).
		Where("category_id = ?", categoryID).
		Pluck("id", &productIDs).Error; err != nil {
		return err
	}

	if len(productIDs) == 0 {
		return nil
	}

	if err := r.db.WithContext(ctx).Where("product_id IN ?", productIDs).
		Delete(&model.ProductImage{}).Error; err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Exec(
		"UPDATE order_items SET product_id = NULL WHERE product_id IN ?", productIDs,
	).Error; err != nil {}

	return r.db.WithContext(ctx).Unscoped().Where("category_id = ?", categoryID).Delete(&model.Product{}).Error
}
