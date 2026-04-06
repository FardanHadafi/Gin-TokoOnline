package handler

import (
	"Toko-Online/dto"
	"Toko-Online/service"
	"Toko-Online/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CategoryHandlerImpl struct {
	svc service.CategoryService
}

func NewCategoryHandler(svc service.CategoryService) CategoryHandler {
	return &CategoryHandlerImpl{svc: svc}
}

func (h *CategoryHandlerImpl) FindAll(c *gin.Context) {
	categories, err := h.svc.FindAll(c.Request.Context())
	if err != nil {
		utils.NewErrorResponse(c, err)
		return
	}
	utils.NewSuccessResponse(c, http.StatusOK, "Categories found", categories)
}

func (h *CategoryHandlerImpl) FindByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.NewErrorResponse(c, err)
		return
	}

	category, err := h.svc.FindByID(c.Request.Context(), id)
	if err != nil {
		utils.NewErrorResponse(c, err)
		return
	}
	utils.NewSuccessResponse(c, http.StatusOK, "Category found", category)
}

func (h *CategoryHandlerImpl) Create(c *gin.Context) {
	var req dto.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.NewErrorResponse(c, err)
		return
	}

	category, err := h.svc.Create(c.Request.Context(), req)
	if err != nil {
		utils.NewErrorResponse(c, err)
		return
	}
	utils.NewSuccessResponse(c, http.StatusCreated, "Category created", category)
}

func (h *CategoryHandlerImpl) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.NewErrorResponse(c, err)
		return
	}

	var req dto.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.NewErrorResponse(c, err)
		return
	}

	category, err := h.svc.Update(c.Request.Context(), id, req)
	if err != nil {
		utils.NewErrorResponse(c, err)
		return
	}
	utils.NewSuccessResponse(c, http.StatusOK, "Category updated", category)
}

func (h *CategoryHandlerImpl) Delete(c *gin.Context) {
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
	utils.NewSuccessResponse(c, http.StatusOK, "Category deleted", nil)
}
