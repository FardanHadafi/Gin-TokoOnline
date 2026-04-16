package config

import (
	"Toko-Online/model"
	"log/slog"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDatabase(log *slog.Logger) *gorm.DB {
	dsn := os.Getenv("DATABASE_URL")
		if dsn == "" {
			log.Error("DATABASE_URL is not set")
			os.Exit(1)
		}

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
		if err != nil {
			log.Error("Failed to connect to database", "error", err)
			os.Exit(1)
		}
	log.Info("Database connected successfully")

	log.Info("Starting database migration")
	err = db.AutoMigrate(
		&model.Category{},
		&model.Order{},
		&model.OrderItem{},
		&model.Product{},
		&model.ProductImage{},
		&model.User{},
		&model.Setting{},
	)
		if err != nil {
			log.Error("Failed to migrate database", "error", err)
			os.Exit(1)
		}
	log.Info("Database migration completed successfully")

	return db
}