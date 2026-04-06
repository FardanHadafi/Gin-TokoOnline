package mocks

import (
	"Toko-Online/model"
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type ProductRepositoryMock struct {
	mock.Mock
}

func (m *ProductRepositoryMock) FindAll(ctx context.Context) ([]model.Product, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.Product), args.Error(1)
}

func (m *ProductRepositoryMock) FindByID(ctx context.Context, id uuid.UUID) (model.Product, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(model.Product), args.Error(1)
}

func (m *ProductRepositoryMock) FindByCategoryID(ctx context.Context, categoryID uuid.UUID) ([]model.Product, error) {
	args := m.Called(ctx, categoryID)
	return args.Get(0).([]model.Product), args.Error(1)
}

func (m *ProductRepositoryMock) Create(ctx context.Context, product *model.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *ProductRepositoryMock) Update(ctx context.Context, product *model.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *ProductRepositoryMock) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type CategoryRepositoryMock struct {
	mock.Mock
}

func (m *CategoryRepositoryMock) FindAll(ctx context.Context) ([]model.Category, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.Category), args.Error(1)
}

func (m *CategoryRepositoryMock) FindByID(ctx context.Context, id uuid.UUID) (model.Category, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(model.Category), args.Error(1)
}

func (m *CategoryRepositoryMock) Create(ctx context.Context, category *model.Category) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *CategoryRepositoryMock) Update(ctx context.Context, category *model.Category) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *CategoryRepositoryMock) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type UserRepositoryMock struct {
	mock.Mock
}

func (m *UserRepositoryMock) FindByUsername(ctx context.Context, username string) (model.User, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *UserRepositoryMock) FindByID(ctx context.Context, id uuid.UUID) (model.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *UserRepositoryMock) Update(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}
func (m *UserRepositoryMock) Create(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}
func (m *UserRepositoryMock) IsEmpty(ctx context.Context) (bool, error) {
	args := m.Called(ctx)
	return args.Bool(0), args.Error(1)
}

type OrderRepositoryMock struct {
	mock.Mock
}

func (m *OrderRepositoryMock) Create(ctx context.Context, order *model.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *OrderRepositoryMock) FindAll(ctx context.Context) ([]model.Order, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.Order), args.Error(1)
}

func (m *OrderRepositoryMock) FindByID(ctx context.Context, id uuid.UUID) (model.Order, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(model.Order), args.Error(1)
}

func (m *OrderRepositoryMock) FindByOrderNumber(ctx context.Context, orderNumber string) (model.Order, error) {
	args := m.Called(ctx, orderNumber)
	return args.Get(0).(model.Order), args.Error(1)
}

func (m *OrderRepositoryMock) Update(ctx context.Context, order *model.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *OrderRepositoryMock) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *OrderRepositoryMock) UpdatePaymentInfo(ctx context.Context, id uuid.UUID, snapToken string, redirectURL string) error {
	args := m.Called(ctx, id, snapToken, redirectURL)
	return args.Error(0)
}
