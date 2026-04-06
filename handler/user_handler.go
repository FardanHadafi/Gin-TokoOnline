package handler

import "github.com/gin-gonic/gin"

type UserHandler interface {
	Login(c *gin.Context)
	GetProfile(c *gin.Context)
	UpdateProfile(c *gin.Context)
	Logout(c *gin.Context)
}
