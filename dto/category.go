package dto

import "github.com/google/uuid"

type CategoryResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type UpdateCategoryRequest struct {
	Name string `json:"name" binding:"required,min=3"`
}
