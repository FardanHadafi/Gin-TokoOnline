package test

import (
	"Toko-Online/model"
	"Toko-Online/service"
	"Toko-Online/test/mocks"
	"context"
	"fmt"
	"log/slog"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestUserService_Login_Success(t *testing.T) {
	repo := new(mocks.UserRepositoryMock)
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	
	svc := service.NewUserService(repo, logger, nil)

	ctx := context.Background()
	userID := uuid.New()
	user := model.User{
		ID:       userID,
		Username: "admin",
		Name:     "Admin User",
	}
	
	hash, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	user.Password = string(hash)

	repo.On("FindByUsername", ctx, "admin").Return(user, nil)

	res, err := svc.Login(ctx, "admin", "password")

	assert.NoError(t, err)
	assert.Equal(t, userID, res.User.ID)
	assert.NotEmpty(t, res.Token)
	repo.AssertExpectations(t)
}

func TestUserService_Login_InvalidCredentials(t *testing.T) {
	repo := new(mocks.UserRepositoryMock)
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	
	svc := service.NewUserService(repo, logger, nil)

	ctx := context.Background()
	repo.On("FindByUsername", ctx, "wrong").Return(model.User{}, fmt.Errorf("not found"))

	res, err := svc.Login(ctx, "wrong", "password")

	assert.Error(t, err)
	assert.Empty(t, res.Token)
}

func TestUserService_Logout(t *testing.T) {
	repo := new(mocks.UserRepositoryMock)
	svc := service.NewUserService(repo, slog.Default(), nil)
	err := svc.Logout(context.Background(), "test-token")
	assert.NoError(t, err)
}

func TestUserService_GetProfile(t *testing.T) {
	repo := new(mocks.UserRepositoryMock)
	svc := service.NewUserService(repo, slog.Default(), nil)

	id := uuid.New()
	repo.On("FindByID", mock.Anything, id).Return(model.User{ID: id, Username: "test", Name: "Test"}, nil)

	res, err := svc.GetProfile(context.Background(), id)
	assert.NoError(t, err)
	assert.Equal(t, "test", res.Username)
}
