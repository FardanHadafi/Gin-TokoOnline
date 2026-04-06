package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Product struct {
	ID          uuid.UUID       `gorm:"type:uuid;primaryKey" json:"id"`
	CategoryID  uuid.UUID       `gorm:"type:uuid" json:"category_id" validate:"required"`
	Category    Category        `json:"category"`
	Name        string          `gorm:"type:varchar(200);not null" json:"name" validate:"required,min=3"`
	Slug        string          `gorm:"type:varchar(200);uniqueIndex" json:"slug"`
	Description string          `gorm:"type:text" json:"description"`
	Price       decimal.Decimal `gorm:"type:decimal(12,2);not null" json:"price" validate:"required,gt=0"`
	Stock       int             `gorm:"default:0" json:"stock" validate:"gte=0"`
	IsActive    bool            `gorm:"default:true" json:"is_active"`
	ImageUrl    string          `json:"image_url"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	DeletedAt   gorm.DeletedAt  `gorm:"index" json:"-"`
}


func (p *Product) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return
}