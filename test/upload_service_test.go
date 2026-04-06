package test

import (
	"Toko-Online/service"
	"context"
	"log/slog"
	"mime/multipart"
	"os"
	"testing"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/stretchr/testify/assert"
)

func TestUploadService_UploadFile_ConfigError(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	svc := service.NewUploadService(nil, logger)

	ctx := context.Background()
	header := &multipart.FileHeader{
		Filename: "test.jpg",
	}

	url, err := svc.UploadFile(ctx, header)

	assert.Error(t, err)
	assert.Empty(t, url)
	assert.Contains(t, err.Error(), "Cloudinary is not configured")
}

func TestUploadService_Initialization(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	cld := &cloudinary.Cloudinary{}
	
	svc := service.NewUploadService(cld, logger)
	assert.NotNil(t, svc)
}
