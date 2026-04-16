package dto

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ProductResponse struct {
	ID           uuid.UUID       `json:"id"`
	CategoryID   uuid.UUID       `json:"category_id"`
	Name         string          `json:"name"`
	Slug         string          `json:"slug"`
	Description  string          `json:"description"`
	Price        decimal.Decimal `json:"price"`
	Stock        int             `json:"stock"`
	ImageUrl     string          `json:"image_url"`
	Images       []string        `json:"images"`
	CategoryName string          `json:"category_name"`
	CategorySlug string          `json:"category_slug"`
	IsActive     bool            `json:"is_active"`
}

type AddProductRequest struct {
	CategoryID  uuid.UUID       `json:"category_id" binding:"required"`
	Name        string          `json:"name" binding:"required,min=3"`
	Description string          `json:"description"`
	Price       decimal.Decimal `json:"price" binding:"required"`
	Stock       int             `json:"stock" binding:"gte=0"`
	ImageUrl    string          `json:"image_url"`
	Images      []string        `json:"images"`
	IsActive    bool            `json:"is_active"`
}

type UpdateProductRequest struct {
	CategoryID  uuid.UUID       `json:"category_id" binding:"omitempty"`
	Name        string          `json:"name" binding:"omitempty,min=3"`
	Description string          `json:"description"`
	Price       decimal.Decimal `json:"price" binding:"omitempty"`
	Stock       int             `json:"stock" binding:"omitempty,gte=0"`
	ImageUrl    string          `json:"image_url"`
	Images      []string        `json:"images"`
	IsActive    bool            `json:"is_active"`
}
