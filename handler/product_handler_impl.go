package handler

import (
	"Toko-Online/dto"
	"Toko-Online/service"
	"Toko-Online/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProductHandlerImpl struct {
	svc service.ProductService
}

func NewProductHandler(svc service.ProductService) ProductHandler {
	return &ProductHandlerImpl{svc: svc}
}

func (h *ProductHandlerImpl) FindAll(c *gin.Context) {
	products, err := h.svc.FindAll(c.Request.Context())
	if err != nil {
		utils.NewErrorResponse(c, err)
		return
	}
	utils.NewSuccessResponse(c, http.StatusOK, "Products found", products)
}

func (h *ProductHandlerImpl) FindByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.NewErrorResponse(c, err)
		return
	}

	product, err := h.svc.FindByID(c.Request.Context(), id)
	if err != nil {
		utils.NewErrorResponse(c, err)
		return
	}
	utils.NewSuccessResponse(c, http.StatusOK, "Product found", product)
}

func (h *ProductHandlerImpl) Create(c *gin.Context) {
	var req dto.AddProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.NewErrorResponse(c, err)
		return
	}

	product, err := h.svc.Create(c.Request.Context(), req)
	if err != nil {
		utils.NewErrorResponse(c, err)
		return
	}
	utils.NewSuccessResponse(c, http.StatusCreated, "Product created", product)
}

func (h *ProductHandlerImpl) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.NewErrorResponse(c, err)
		return
	}

	var req dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.NewErrorResponse(c, err)
		return
	}

	product, err := h.svc.Update(c.Request.Context(), id, req)
	if err != nil {
		utils.NewErrorResponse(c, err)
		return
	}
	utils.NewSuccessResponse(c, http.StatusOK, "Product updated", product)
}

func (h *ProductHandlerImpl) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.NewErrorResponse(c, err)
		return
	}

	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		utils.NewErrorResponse(c, err)
		return
	}
	utils.NewSuccessResponse(c, http.StatusOK, "Product deleted", nil)
}

func (h *ProductHandlerImpl) GetByCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.NewErrorResponse(c, err)
		return
	}

	products, err := h.svc.FindByCategoryID(c.Request.Context(), id)
	if err != nil {
		utils.NewErrorResponse(c, err)
		return
	}
	utils.NewSuccessResponse(c, http.StatusOK, "Products found by category", products)
}
