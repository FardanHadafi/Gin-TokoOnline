package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Category struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(100);not null" json:"name" validate:"required,min=3"`
	Slug      string    `gorm:"type:varchar(100);uniqueIndex" json:"slug"`
	ImageUrl  string    `gorm:"type:varchar(255)" json:"image_url"`
	Products  []Product `json:"products,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (c *Category) BeforeCreate(tx *gorm.DB) (err error) {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return
}
