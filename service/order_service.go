package service

import (
	"Toko-Online/dto"
	"context"

	"github.com/google/uuid"
)

type OrderService interface {
	Checkout(ctx context.Context, req dto.CheckoutRequest) (dto.OrderResponse, error)
	FindAll(ctx context.Context) ([]dto.OrderResponse, error)
	FindByID(ctx context.Context, id uuid.UUID) (dto.OrderResponse, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) (dto.OrderResponse, error)
	HandleMidtransWebhook(ctx context.Context, payload map[string]interface{}) error
}
