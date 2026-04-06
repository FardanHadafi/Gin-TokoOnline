package repository

import (
	"Toko-Online/model"
	"context"

	"github.com/google/uuid"
)

type OrderRepository interface {
	Create(ctx context.Context, order *model.Order) error
	FindAll(ctx context.Context) ([]model.Order, error)
	FindByID(ctx context.Context, id uuid.UUID) (model.Order, error)
	FindByOrderNumber(ctx context.Context, orderNumber string) (model.Order, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
	UpdatePaymentInfo(ctx context.Context, id uuid.UUID, snapToken string, redirectURL string) error
}
