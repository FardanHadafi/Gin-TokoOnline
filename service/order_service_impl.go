package service

import (
	"Toko-Online/config"
	"Toko-Online/dto"
	"Toko-Online/model"
	"Toko-Online/repository"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type OrderServiceImpl struct {
	db          *gorm.DB
	repo        repository.OrderRepository
	productRepo repository.ProductRepository
	logger      *slog.Logger
	snapClient  snap.Client
}

func NewOrderService(
	db *gorm.DB,
	repo repository.OrderRepository,
	productRepo repository.ProductRepository,
	logger *slog.Logger,
) OrderService {
	client := snap.Client{}
	
	return &OrderServiceImpl{
		db:          db,
		repo:        repo,
		productRepo: productRepo,
		logger:      logger,
		snapClient:  client,
	}
}

func (s *OrderServiceImpl) Checkout(ctx context.Context, req dto.CheckoutRequest) (dto.OrderResponse, error) {
	s.logger.InfoContext(ctx, "Starting checkout process", "customer", req.CustomerName)

	var order model.Order
	var totalAmount decimal.Decimal

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		order = model.Order{
			CustomerName:    req.CustomerName,
			CustomerEmail:   req.CustomerEmail,
			CustomerPhone:   req.CustomerPhone,
			ShippingAddress: req.ShippingAddress,
			Status:          "pending",
			OrderNumber:     fmt.Sprintf("ORD-%d", time.Now().UnixNano()),
		}

		for _, itemReq := range req.Items {
			product, err := s.productRepo.FindByID(ctx, itemReq.ProductID)
			if err != nil {
				return &config.ApiError{Status: 400, Title: "Bad Request", Detail: "Product not found: " + itemReq.ProductID.String()}
			}

			if !product.IsActive || product.Stock < itemReq.Quantity {
				return &config.ApiError{Status: 400, Title: "Bad Request", Detail: "Product out of stock or inactive: " + product.Name}
			}

			itemTotal := product.Price.Mul(decimal.NewFromInt(int64(itemReq.Quantity)))
			totalAmount = totalAmount.Add(itemTotal)

			order.Items = append(order.Items, model.OrderItem{
				ProductID: product.ID,
				Quantity:  itemReq.Quantity,
				Price:     product.Price,
			})
		}

		order.TotalAmount = totalAmount

		if err := tx.Create(&order).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		s.logger.ErrorContext(ctx, "Checkout failed in transaction", "error", err)
		if apiErr, ok := err.(*config.ApiError); ok {
			return dto.OrderResponse{}, apiErr
		}
		return dto.OrderResponse{}, &config.ApiError{Status: 500, Title: "Internal Error", Detail: err.Error()}
	}

	s.logger.InfoContext(ctx, "Generating Midtrans SNAP token", "order_id", order.ID)
	snapReq := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  order.OrderNumber,
			GrossAmt: order.TotalAmount.IntPart(),
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: order.CustomerName,
			Email: order.CustomerEmail,
			Phone: order.CustomerPhone,
		},
	}

	snapResp, snapErr := s.snapClient.CreateTransaction(snapReq)
	if snapErr == nil {
		_ = s.repo.UpdatePaymentInfo(ctx, order.ID, snapResp.Token, snapResp.RedirectURL)
		order.SnapToken = snapResp.Token
		order.SnapRedirectURL = snapResp.RedirectURL
	} else {
		s.logger.WarnContext(ctx, "Failed to generate Midtrans SNAP token", "error", snapErr)
	}

	return s.mapToResponse(order), nil
}

func (s *OrderServiceImpl) FindAll(ctx context.Context) ([]dto.OrderResponse, error) {
	orders, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	var res []dto.OrderResponse
	for _, o := range orders {
		res = append(res, s.mapToResponse(o))
	}
	return res, nil
}

func (s *OrderServiceImpl) FindByID(ctx context.Context, id uuid.UUID) (dto.OrderResponse, error) {
	order, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return dto.OrderResponse{}, err
	}
	return s.mapToResponse(order), nil
}

func (s *OrderServiceImpl) UpdateStatus(ctx context.Context, id uuid.UUID, status string) (dto.OrderResponse, error) {
	if err := s.repo.UpdateStatus(ctx, id, status); err != nil {
		return dto.OrderResponse{}, err
	}
	return s.FindByID(ctx, id)
}

func (s *OrderServiceImpl) mapToResponse(o model.Order) dto.OrderResponse {
	var items []dto.OrderItemResponse
	var itemDetails string
	for _, i := range o.Items {
		items = append(items, dto.OrderItemResponse{
			ID:        i.ID,
			OrderID:   i.OrderID,
			ProductID: i.ProductID,
			Product: dto.ProductResponse{
				ID:   i.Product.ID,
				Name: i.Product.Name,
			},
			Quantity: i.Quantity,
			Price:    i.Price,
		})
		itemDetails += fmt.Sprintf("- %dx %s (Rp %s)\n", i.Quantity, i.Product.Name, i.Price.String())
	}

	adminWA := os.Getenv("ADMIN_WHATSAPP")
	// Message Template
	msg := fmt.Sprintf("Halo Admin,\n\nSaya ingin mengonfirmasi pesanan saya:\n\n*No. Pesanan*: %s\n*Nama*: %s\n*Total*: Rp %s\n*Alamat*: %s\n\n*Detail Barang*:\n%s\n*Status*: %s\n\nTerima kasih.",
		o.OrderNumber, o.CustomerName, o.TotalAmount.String(), o.ShippingAddress, itemDetails, o.Status)

	whatsappURL := ""
	if adminWA != "" {
		whatsappURL = fmt.Sprintf("https://wa.me/%s?text=%s", adminWA, url.QueryEscape(msg))
	}

	return dto.OrderResponse{
		ID:              o.ID,
		OrderNumber:     o.OrderNumber,
		CustomerName:    o.CustomerName,
		CustomerEmail:   o.CustomerEmail,
		CustomerPhone:   o.CustomerPhone,
		ShippingAddress: o.ShippingAddress,
		TotalAmount:     o.TotalAmount,
		Status:          o.Status,
		Note:            o.Note,
		SnapToken:       o.SnapToken,
		SnapRedirectURL: o.SnapRedirectURL,
		PaymentType:     o.PaymentType,
		Items:           items,
		WhatsAppURL:     whatsappURL,
	}
}

func (s *OrderServiceImpl) HandleMidtransWebhook(ctx context.Context, payload map[string]interface{}) error {
	orderID, ok := payload["order_id"].(string)
	if !ok {
		return &config.ApiError{Status: http.StatusBadRequest, Title: "Bad Request", Detail: "Invalid order_id in payload"}
	}

	transactionStatus, _ := payload["transaction_status"].(string)
	fraudStatus, _ := payload["fraud_status"].(string)
	paymentType, _ := payload["payment_type"].(string)

	s.logger.InfoContext(ctx, "Received Midtrans webhook", "order_id", orderID, "status", transactionStatus)

	order, err := s.repo.FindByOrderNumber(ctx, orderID)
	if err != nil {
		s.logger.ErrorContext(ctx, "Order not found for webhook", "order_id", orderID)
		return &config.ApiError{Status: http.StatusNotFound, Title: "Not Found", Detail: "Order not found"}
	}

	var status string
	switch transactionStatus {
	case "capture":
		if fraudStatus == "challenge" {
			status = "challenge"
		} else if fraudStatus == "accept" {
			status = "success"
		}
	case "settlement":
		status = "success"
	case "cancel", "deny", "expire":
		status = "failed"
	case "pending":
		status = "pending"
	default:
		status = "unknown"
	}

	err = s.db.WithContext(ctx).Model(&model.Order{}).Where("id = ?", order.ID).Updates(map[string]interface{}{
		"status":       status,
		"payment_type": paymentType,
	}).Error

	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to update order status from webhook", "error", err)
		return &config.ApiError{Status: http.StatusInternalServerError, Title: "Internal Error", Detail: "Failed to update order status"}
	}

	s.logger.InfoContext(ctx, "Successfully updated order from webhook", "order_id", order.ID, "status", status)
	return nil
}
