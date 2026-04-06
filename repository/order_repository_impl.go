package repository

import (
	"Toko-Online/model"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderRepositoryImpl struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &OrderRepositoryImpl{db: db}
}

func (r *OrderRepositoryImpl) Create(ctx context.Context, order *model.Order) error {
	return r.db.WithContext(ctx).Create(order).Error
}

func (r *OrderRepositoryImpl) FindAll(ctx context.Context) ([]model.Order, error) {
	var orders []model.Order
	err := r.db.WithContext(ctx).Preload("Items.Product").Find(&orders).Error
	return orders, err
}

func (r *OrderRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (model.Order, error) {
	var order model.Order
	err := r.db.WithContext(ctx).Preload("Items.Product").First(&order, "id = ?", id).Error
	return order, err
}

func (r *OrderRepositoryImpl) FindByOrderNumber(ctx context.Context, orderNumber string) (model.Order, error) {
	var order model.Order
	err := r.db.WithContext(ctx).First(&order, "order_number = ?", orderNumber).Error
	return order, err
}

func (r *OrderRepositoryImpl) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	return r.db.WithContext(ctx).Model(&model.Order{}).Where("id = ?", id).Update("status", status).Error
}

func (r *OrderRepositoryImpl) UpdatePaymentInfo(ctx context.Context, id uuid.UUID, snapToken string, redirectURL string) error {
	return r.db.WithContext(ctx).Model(&model.Order{}).Where("id = ?", id).Updates(map[string]interface{}{
		"snap_token":        snapToken,
		"snap_redirect_url": redirectURL,
	}).Error
}
