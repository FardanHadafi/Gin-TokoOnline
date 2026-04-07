package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Order struct {
	ID              uuid.UUID       `gorm:"type:uuid;primaryKey" json:"id"`
	OrderNumber     string          `gorm:"type:varchar(50);uniqueIndex" json:"order_number"`
	CustomerName    string          `gorm:"type:varchar(255);not null" json:"customer_name" validate:"required,min=3"`
	CustomerPhone   string          `gorm:"type:varchar(20);not null" json:"customer_phone" validate:"required,numeric,min=10"`
	CustomerEmail   string          `gorm:"type:varchar(100);not null" json:"customer_email" validate:"required,email"`
	ShippingAddress string          `gorm:"type:text;not null" json:"shipping_address" validate:"required"`
	TotalAmount     decimal.Decimal `gorm:"type:decimal(12,2)" json:"total_amount" validate:"required,gt=0"`
	Status          string          `gorm:"type:varchar(20);default:'pending'" json:"status"`
	Note            string          `gorm:"type:text" json:"note"`
	SnapToken       string          `gorm:"type:varchar(255)" json:"snap_token"`
	SnapRedirectURL string          `gorm:"type:varchar(255)" json:"snap_redirect_url"`
	PaymentType     string          `gorm:"type:varchar(50)" json:"payment_type"`
	Items           []OrderItem     `json:"items"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

func (o *Order) BeforeCreate(tx *gorm.DB) (err error) {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	return
}
