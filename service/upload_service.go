package service

import (
	"context"
	"mime/multipart"
)

type UploadService interface {
	UploadFile(ctx context.Context, fileHeader *multipart.FileHeader) (string, error)
}
