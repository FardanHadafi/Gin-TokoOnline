package handler

import (
	"Toko-Online/dto"
	"Toko-Online/service"
	"Toko-Online/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OrderHandlerImpl struct {
	svc service.OrderService
}

func NewOrderHandler(svc service.OrderService) OrderHandler {
	return &OrderHandlerImpl{svc: svc}
}

func (h *OrderHandlerImpl) Checkout(c *gin.Context) {
	var req dto.CheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": 400,
			"title":  "Validation Error",
			"detail": err.Error(),
		})
		return
	}

	order, err := h.svc.Checkout(c.Request.Context(), req)
	if err != nil {
		utils.NewErrorResponse(c, err)
		return
	}
	utils.NewSuccessResponse(c, http.StatusCreated, "Order created", order)
}

func (h *OrderHandlerImpl) Webhook(c *gin.Context) {
	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		utils.NewErrorResponse(c, err)
		return
	}

	if err := h.svc.HandleMidtransWebhook(c.Request.Context(), payload); err != nil {
		utils.NewErrorResponse(c, err)
		return
	}
	utils.NewSuccessResponse(c, http.StatusOK, "Webhook handled", nil)
}

func (h *OrderHandlerImpl) FindAll(c *gin.Context) {
	orders, err := h.svc.FindAll(c.Request.Context())
	if err != nil {
		utils.NewErrorResponse(c, err)
		return
	}
	utils.NewSuccessResponse(c, http.StatusOK, "Orders found", orders)
}

func (h *OrderHandlerImpl) FindByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.NewErrorResponse(c, err)
		return
	}

	order, err := h.svc.FindByID(c.Request.Context(), id)
	if err != nil {
		utils.NewErrorResponse(c, err)
		return
	}
	utils.NewSuccessResponse(c, http.StatusOK, "Order found", order)
}
