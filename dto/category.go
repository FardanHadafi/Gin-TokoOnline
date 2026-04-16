package dto

import "github.com/google/uuid"

type CategoryResponse struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Slug         string    `json:"slug"`
	ImageUrl     string    `json:"image_url"`
	ProductCount int       `json:"product_count"`
}

type UpdateCategoryRequest struct {
	Name     string `json:"name" binding:"required,min=3"`
	ImageUrl string `json:"image_url"`
}
