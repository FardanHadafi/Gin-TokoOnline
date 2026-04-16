package config

import (
	"context"
	"crypto/tls"
	"log/slog"
	"os"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(log *slog.Logger) *redis.Client {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		log.Warn("REDIS_URL not set")
		return nil
	}

	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Error("Failed to parse REDIS_URL", "error", err)
		return nil
	}

	if opt.TLSConfig == nil {
		opt.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	client := redis.NewClient(opt)
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		log.Warn("Failed to connect to Redis", "error", err)
	} else {
		log.Info("Connected to Upstash Redis successfully")
	}

	return client
}
