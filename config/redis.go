package config

import (
	"context"
	"log/slog"
	"os"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(log *slog.Logger) *redis.Client {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "localhost:6379"
	}

	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		opt = &redis.Options{
			Addr: redisURL,
		}
	}

	client := redis.NewClient(opt)
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		log.Warn("Failed to connect to Redis", "error", err)
	} else {
		log.Info("Connected to Redis successfully")
	}

	return client
}
