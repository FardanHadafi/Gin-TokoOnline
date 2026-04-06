package middleware

import (
	"Toko-Online/config"
	"Toko-Online/utils"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func AdminAuthMiddleware(redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.NewErrorResponse(c, &config.ApiError{
				Status: http.StatusUnauthorized,
				Title:  "Unauthorized",
				Detail: "Authorization header is missing",
			})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			utils.NewErrorResponse(c, &config.ApiError{
				Status: http.StatusUnauthorized,
				Title:  "Unauthorized",
				Detail: "Invalid authorization format. Use 'Bearer <token>'",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		if redisClient != nil {
			sessionKey := fmt.Sprintf("session:%s", tokenString)
			_, err := redisClient.Get(context.Background(), sessionKey).Result()
			if err != nil {
				utils.NewErrorResponse(c, &config.ApiError{
					Status: http.StatusUnauthorized,
					Title:  "Unauthorized",
					Detail: "Session has expired or logged out",
				})
				c.Abort()
				return
			}
		}

		userID, err := config.ValidateToken(tokenString)
		if err != nil {
			utils.NewErrorResponse(c, &config.ApiError{
				Status: http.StatusUnauthorized,
				Title:  "Unauthorized",
				Detail: err.Error(),
			})
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Set("token", tokenString)
		c.Next()
	}
}

func GetUserIDFromContext(c *gin.Context) (string, error) {
	userId, exists := c.Get("user_id")
	if !exists {
		return "", errors.New("user ID not found in context")
	}
	return userId.(string), nil
}

func GetTokenFromContext(c *gin.Context) (string, error) {
	token, exists := c.Get("token")
	if !exists {
		return "", errors.New("token not found in context")
	}
	return token.(string), nil
}
