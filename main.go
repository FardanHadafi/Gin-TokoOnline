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
	
	var admin model.User
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(os.Getenv("ADMIN_PASSWORD")), bcrypt.DefaultCost)
	
	err := db.Where("username = ?", os.Getenv("ADMIN_USERNAME")).First(&admin).Error
	if err != nil {
		logger.Info("Seeding / Updating default admin user")
		db.Create(&model.User{
			Username: os.Getenv("ADMIN_USERNAME"),
			Name:     os.Getenv("ADMIN_NAME"),
			Password: string(hashedPassword),
		})
		logger.Info("Default admin created")
	} else {
		admin.Password = string(hashedPassword)
		admin.Name = os.Getenv("ADMIN_NAME")
		db.Save(&admin)
		logger.Info("Default admin credentials synced with .env")
	}

	r := router.SetupRouter(db, logger)

	logger.Info("Starting server on port 3000")
	if err := r.Run(":3000"); err != nil {
		logger.Error("Failed to start server", "error", err)
	}
}