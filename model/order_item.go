package model

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type OrderItem struct {
	ID        uuid.UUID       `gorm:"type:uuid;primaryKey" json:"id"`
	OrderID   uuid.UUID       `gorm:"type:uuid" json:"order_id"`
	ProductID uuid.UUID       `gorm:"type:uuid" json:"product_id" validate:"required"`
	Product   Product         `json:"product"`
	Quantity  int             `json:"quantity" validate:"required,gt=0"`
	Price     decimal.Decimal `gorm:"type:decimal(12,2)" json:"price" validate:"required,gt=0"`
}


func (oi *OrderItem) BeforeCreate(tx *gorm.DB) (err error) {
	if oi.ID == uuid.Nil {
		oi.ID = uuid.New()
	}
	return
}