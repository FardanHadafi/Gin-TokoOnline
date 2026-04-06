package service

import (
	"Toko-Online/dto"
	"context"

	"github.com/google/uuid"
)

type UserService interface {
	Login(ctx context.Context, username, password string) (dto.LoginResponse, error)
	GetProfile(ctx context.Context, id uuid.UUID) (dto.UserResponse, error)
	UpdateProfile(ctx context.Context, id uuid.UUID, req dto.UpdateUserRequest) (dto.UserResponse, error)
	Logout(ctx context.Context, token string) error
}
