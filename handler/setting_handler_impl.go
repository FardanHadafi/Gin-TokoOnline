package handler

import (
	"Toko-Online/dto"
	"Toko-Online/service"
	"Toko-Online/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SettingHandlerImpl struct {
	svc service.SettingService
}

func NewSettingHandler(svc service.SettingService) SettingHandler {
	return &SettingHandlerImpl{svc: svc}
}

func (h *SettingHandlerImpl) Get(c *gin.Context) {
	settings, err := h.svc.Get(c.Request.Context())
	if err != nil {
		utils.NewErrorResponse(c, err)
		return
	}
	utils.NewSuccessResponse(c, http.StatusOK, "Settings found", settings)
}

func (h *SettingHandlerImpl) Update(c *gin.Context) {
	var req dto.UpdateSettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.NewErrorResponse(c, err)
		return
	}

	settings, err := h.svc.Update(c.Request.Context(), req)
	if err != nil {
		utils.NewErrorResponse(c, err)
		return
	}
	utils.NewSuccessResponse(c, http.StatusOK, "Settings updated", settings)
}
