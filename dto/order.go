package dto

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type OrderItemResponse struct {
	ID        uuid.UUID       `json:"id"`
	OrderID   uuid.UUID       `json:"order_id"`
	ProductID uuid.UUID       `json:"product_id"`
	Product   ProductResponse `json:"product"`
	Quantity  int             `json:"quantity"`
	Price     decimal.Decimal `json:"price"`
}

type OrderResponse struct {
	ID              uuid.UUID           `json:"id"`
	OrderNumber     string              `json:"order_number"`
	CustomerName    string              `json:"customer_name"`
	CustomerEmail   string              `json:"customer_email"`
	CustomerPhone   string              `json:"customer_phone"`
	ShippingAddress string              `json:"shipping_address"`
	TotalAmount     decimal.Decimal     `json:"total_amount"`
	Status          string              `json:"status"`
	Note            string              `json:"note"`
	SnapToken       string              `json:"snap_token"`
	SnapRedirectURL string              `json:"snap_redirect_url"`
	PaymentType     string              `json:"payment_type"`
	Items           []OrderItemResponse `json:"items"`
	WhatsAppURL     string              `json:"whatsapp_url"`
}

type CheckoutItemRequest struct {
	ProductID uuid.UUID `json:"product_id" binding:"required"`
	Quantity  int       `json:"quantity" binding:"required,gt=0"`
}

type CheckoutRequest struct {
	CustomerName    string                `json:"customer_name" binding:"required,min=3"`
	CustomerEmail   string                `json:"customer_email" binding:"required,email"`
	CustomerPhone   string                `json:"customer_phone" binding:"required,numeric,min=10"`
	ShippingAddress string                `json:"shipping_address" binding:"required"`
	Note            string                `json:"note"`
	Items           []CheckoutItemRequest `json:"items" binding:"required,dive"`
}
