package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type UploadServiceImpl struct {
	cloudinary *cloudinary.Cloudinary
	logger     *slog.Logger
}

func NewUploadService(cloudinary *cloudinary.Cloudinary, logger *slog.Logger) UploadService {
	return &UploadServiceImpl{
		cloudinary: cloudinary,
		logger:     logger,
	}
}

func (s *UploadServiceImpl) UploadFile(ctx context.Context, fileHeader *multipart.FileHeader) (string, error) {
	if s.cloudinary == nil {
		return "", errors.New("Cloudinary is not configured")
	}

	file, err := fileHeader.Open()
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to open file for upload", "error", err)
		return "", err
	}
	defer file.Close()

	resp, err := s.cloudinary.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder: "toko-online/uploads",
	})
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to upload to Cloudinary", "error", err)
		return "", fmt.Errorf("failed to upload to Cloudinary: %v", err)
	}

	s.logger.InfoContext(ctx, "File uploaded successfully", "url", resp.SecureURL)
	return resp.SecureURL, nil
}
