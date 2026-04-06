package service

import (
	"Toko-Online/config"
	"Toko-Online/dto"
	"Toko-Online/repository"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceImpl struct {
	repo   repository.UserRepository
	logger *slog.Logger
	redis  *redis.Client
}

const (
	UserSessionKey      = "session:%s"
	UserSessionDuration = 24 * time.Hour
)

func NewUserService(repo repository.UserRepository, logger *slog.Logger, redis *redis.Client) UserService {
	return &UserServiceImpl{
		repo:   repo,
		logger: logger,
		redis:  redis,
	}
}

func (s *UserServiceImpl) Login(ctx context.Context, username, password string) (dto.LoginResponse, error) {
	s.logger.InfoContext(ctx, "Login attempt", "username", username)
	user, err := s.repo.FindByUsername(ctx, username)
	if err != nil {
		s.logger.WarnContext(ctx, "Login failed: user not found", "username", username)
		return dto.LoginResponse{}, &config.ApiError{
			Status: http.StatusUnauthorized,
			Title:  "Unauthorized",
			Detail: "Invalid username or password",
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		s.logger.WarnContext(ctx, "Login failed: invalid password", "username", username)
		return dto.LoginResponse{}, &config.ApiError{
			Status: http.StatusUnauthorized,
			Title:  "Unauthorized",
			Detail: "Invalid username or password",
		}
	}

	token, err := config.GenerateToken(user.ID)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to generate token", "error", err)
		return dto.LoginResponse{}, err
	}

	res := dto.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Name:     user.Name,
	}

	if s.redis != nil {
		sessionKey := fmt.Sprintf(UserSessionKey, token)
		data, _ := json.Marshal(res)
		s.redis.Set(ctx, sessionKey, data, UserSessionDuration)
	}

	return dto.LoginResponse{
		User:  res,
		Token: token,
	}, nil
}

func (s *UserServiceImpl) GetProfile(ctx context.Context, id uuid.UUID) (dto.UserResponse, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return dto.UserResponse{}, &config.ApiError{Status: 404, Title: "Not Found", Detail: "User not found"}
	}

	return dto.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Name:     user.Name,
	}, nil
}

func (s *UserServiceImpl) UpdateProfile(ctx context.Context, id uuid.UUID, req dto.UpdateUserRequest) (dto.UserResponse, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return dto.UserResponse{}, &config.ApiError{Status: 404, Title: "Not Found", Detail: "User not found"}
	}

	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Password != "" {
		hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		user.Password = string(hashed)
	}

	if err := s.repo.Update(ctx, &user); err != nil {
		return dto.UserResponse{}, err
	}

	return dto.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Name:     user.Name,
	}, nil
}

func (s *UserServiceImpl) Logout(ctx context.Context, token string) error {
	if s.redis == nil {
		return nil
	}
	sessionKey := fmt.Sprintf(UserSessionKey, token)
	return s.redis.Del(ctx, sessionKey).Err()
}
