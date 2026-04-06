package handler

import (
	"Toko-Online/config"
	"Toko-Online/service"
	"Toko-Online/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	service service.UploadService
}

func NewUploadHandler(service service.UploadService) *UploadHandler {
	return &UploadHandler{
		service: service,
	}
}

func (h *UploadHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		utils.NewErrorResponse(c, err)
		return
	}

	// Limit file size (5MB)
	if file.Size > 5*1024*1024 {
		utils.NewErrorResponse(c, &config.ApiError{
			Status: http.StatusBadRequest,
			Title:  "Bad Request",
			Detail: "File too large (max 5MB)",
		})
		return
	}

	url, err := h.service.UploadFile(c.Request.Context(), file)
	if err != nil {
		utils.NewErrorResponse(c, err)
		return
	}

	utils.NewSuccessResponse(c, http.StatusOK, "File uploaded successfully", gin.H{
		"url": url,
	})
}
