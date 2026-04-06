package handler

import (
	"Toko-Online/dto"
	"Toko-Online/middleware"
	"Toko-Online/service"
	"Toko-Online/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandlerImpl struct {
	svc service.UserService
}

func NewUserHandler(svc service.UserService) UserHandler {
	return &UserHandlerImpl{svc: svc}
}

func (h *UserHandlerImpl) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.NewErrorResponse(c, err)
		return
	}

	user, err := h.svc.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		utils.NewErrorResponse(c, err)
		return
	}

	utils.NewSuccessResponse(c, http.StatusOK, "Login successful", user)
}

func (h *UserHandlerImpl) GetProfile(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.NewErrorResponse(c, err)
		return
	}

	user, err := h.svc.GetProfile(c.Request.Context(), id)
	if err != nil {
		utils.NewErrorResponse(c, err)
		return
	}
	utils.NewSuccessResponse(c, http.StatusOK, "Profile found", user)
}

func (h *UserHandlerImpl) UpdateProfile(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.NewErrorResponse(c, err)
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.NewErrorResponse(c, err)
		return
	}

	user, err := h.svc.UpdateProfile(c.Request.Context(), id, req)
	if err != nil {
		utils.NewErrorResponse(c, err)
		return
	}
	utils.NewSuccessResponse(c, http.StatusOK, "Profile updated", user)
}

func (h *UserHandlerImpl) Logout(c *gin.Context) {
	token, err := middleware.GetTokenFromContext(c)
	if err != nil {
		utils.NewErrorResponse(c, err)
		return
	}

	if err := h.svc.Logout(c.Request.Context(), token); err != nil {
		utils.NewErrorResponse(c, err)
		return
	}

	utils.NewSuccessResponse(c, http.StatusOK, "Logout successful", nil)
}
