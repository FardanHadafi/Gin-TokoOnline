package dto

import "github.com/google/uuid"

type SettingResponse struct {
	ID             uuid.UUID `json:"id"`
	StoreName      string    `json:"store_name"`
	WhatsAppNumber string    `json:"whatsapp_number"`
	AddressInfo    string    `json:"address_info"`
	WelcomeMessage string    `json:"welcome_message"`
}

type UpdateSettingRequest struct {
	StoreName      string `json:"store_name" binding:"omitempty,min=3"`
	WhatsAppNumber string `json:"whatsapp_number" binding:"omitempty,numeric,min=10"`
	AddressInfo    string `json:"address_info" binding:"omitempty"`
	WelcomeMessage string `json:"welcome_message" binding:"omitempty"`
}
