package dto

import "github.com/google/uuid"

type UserResponse struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Name     string    `json:"name"`
}

type LoginResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}

type UpdateUserRequest struct {
	Username string `json:"username" binding:"omitempty,min=3"`
	Name     string `json:"name" binding:"omitempty,min=3"`
	Password string `json:"password" binding:"omitempty,min=8"`
}
