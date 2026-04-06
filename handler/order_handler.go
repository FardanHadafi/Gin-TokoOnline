package handler

import "github.com/gin-gonic/gin"

type OrderHandler interface {
	Checkout(c *gin.Context)
	FindAll(c *gin.Context)
	FindByID(c *gin.Context)
	Webhook(c *gin.Context)
}
