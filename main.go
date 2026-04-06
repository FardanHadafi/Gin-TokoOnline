package main

import (
	"Toko-Online/config"
	"Toko-Online/model"
	"Toko-Online/router"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	_ = godotenv.Load()

	logger := config.NewLogger()
	slog.SetDefault(logger)

	db := config.NewDatabase(logger)
	
	// Create default admin if not exists
	var count int64
	db.Model(&model.User{}).Count(&count)
	if count == 0 {
		logger.Info("Seeding default admin user")
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(os.Getenv("ADMIN_PASSWORD")), bcrypt.DefaultCost)
		db.Create(&model.User{
			Username: os.Getenv("ADMIN_USERNAME"),
			Name:     os.Getenv("ADMIN_NAME"),
			Password: string(hashedPassword),
		})
		logger.Info("Default admin created")
	}

	r := router.SetupRouter(db, logger)

	logger.Info("Starting server on port 3000")
	if err := r.Run(":3000"); err != nil {
		logger.Error("Failed to start server", "error", err)
	}
}