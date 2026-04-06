package handler

import "github.com/gin-gonic/gin"

type SettingHandler interface {
	Get(c *gin.Context)
	Update(c *gin.Context)
}
