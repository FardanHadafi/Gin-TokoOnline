package test

import (
	"Toko-Online/model"
	"Toko-Online/service"
	"Toko-Online/test/mocks"
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestOrderService_FindAll(t *testing.T) {
	repo := new(mocks.OrderRepositoryMock)
	prodRepo := new(mocks.ProductRepositoryMock)
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	
	svc := service.NewOrderService(nil, repo, prodRepo, logger, nil)

	ctx := context.Background()
	orders := []model.Order{
		{ID: uuid.New(), OrderNumber: "ORD-1", CustomerName: "User 1", TotalAmount: decimal.NewFromInt(1000)},
	}

	repo.On("FindAll", ctx).Return(orders, nil)

	res, err := svc.FindAll(ctx)

	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, "ORD-1", res[0].OrderNumber)
	repo.AssertExpectations(t)
}

func TestOrderService_FindByID(t *testing.T) {
	repo := new(mocks.OrderRepositoryMock)
	prodRepo := new(mocks.ProductRepositoryMock)
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	
	svc := service.NewOrderService(nil, repo, prodRepo, logger, nil)

	ctx := context.Background()
	id := uuid.New()
	order := model.Order{ID: id, OrderNumber: "ORD-1", CustomerName: "User 1", TotalAmount: decimal.NewFromInt(1000)}

	repo.On("FindByID", ctx, id).Return(order, nil)

	res, err := svc.FindByID(ctx, id)

	assert.NoError(t, err)
	assert.Equal(t, "ORD-1", res.OrderNumber)
	repo.AssertExpectations(t)
}

func TestOrderService_MapToResponse_WhatsAppURL(t *testing.T) {
	repo := new(mocks.OrderRepositoryMock)
	prodRepo := new(mocks.ProductRepositoryMock)
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	
	// Set dummy env
	os.Setenv("ADMIN_WHATSAPP", "628123456789")
	defer os.Unsetenv("ADMIN_WHATSAPP")

	svc := service.NewOrderService(nil, repo, prodRepo, logger, nil).(*service.OrderServiceImpl)

	id := uuid.New()
	order := model.Order{
		ID:           id,
		OrderNumber:  "ORD-999",
		CustomerName: "John Doe",
		TotalAmount:  decimal.NewFromInt(50000),
		Status:       "pending",
		Items: []model.OrderItem{
			{Quantity: 1, Price: decimal.NewFromInt(50000), Product: model.Product{Name: "Keyboard"}},
		},
	}

	repo.On("FindByID", mock.Anything, id).Return(order, nil)
	
	res, err := svc.FindByID(context.Background(), id)
	
	assert.NoError(t, err)
	assert.Contains(t, res.WhatsAppURL, "https://wa.me/628123456789")
	assert.Contains(t, res.WhatsAppURL, "ORD-999")
	assert.Contains(t, res.WhatsAppURL, "John+Doe")
}
