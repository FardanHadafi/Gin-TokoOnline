package config

import (
	"log/slog"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
)

func NewCloudinaryClient(log *slog.Logger) *cloudinary.Cloudinary {
	cloudinaryURL := os.Getenv("CLOUDINARY_URL")
	if cloudinaryURL == "" {
		log.Warn("CLOUDINARY_URL not set, image upload will fail")
		return nil
	}

	cld, err := cloudinary.NewFromURL(cloudinaryURL)
	if err != nil {
		log.Error("Failed to initialize Cloudinary", "error", err)
		return nil
	}

	log.Info("Cloudinary initialized successfully")
	return cld
}
