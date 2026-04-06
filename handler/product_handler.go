package handler

import "github.com/gin-gonic/gin"

type ProductHandler interface {
	FindAll(c *gin.Context)
	FindByID(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	GetByCategory(c *gin.Context)
}
