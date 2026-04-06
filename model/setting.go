package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Setting struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	StoreName      string    `gorm:"type:varchar(100);not null" json:"store_name" validate:"required,min=3"`
	WhatsAppNumber string    `gorm:"type:varchar(20);not null" json:"whatsapp_number" validate:"required,numeric,min=10"`
	AddressInfo    string    `gorm:"type:text" json:"address_info"`
	WelcomeMessage string    `gorm:"type:text" json:"welcome_message"`
}

func (s *Setting) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return
}
