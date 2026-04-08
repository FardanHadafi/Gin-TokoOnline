package utils

import (
	"Toko-Online/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func NewSuccessResponse(c *gin.Context, status int, message string, data interface{}) {
	c.JSON(status, Response{
		Status:  status,
		Message: message,
		Data:    data,
	})
}

func NewErrorResponse(c *gin.Context, err error) {
	if apiErr, ok := err.(*config.ApiError); ok {
		c.JSON(int(apiErr.Status), apiErr)
		return
	}

	c.JSON(http.StatusBadRequest, config.ApiError{
		Status: http.StatusBadRequest,
		Title:  "Bad Request",
		Detail: err.Error(),
	})
}
