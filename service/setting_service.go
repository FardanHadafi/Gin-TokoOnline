package service

import (
	"Toko-Online/dto"
	"context"
)

type SettingService interface {
	Get(ctx context.Context) (dto.SettingResponse, error)
	Update(ctx context.Context, req dto.UpdateSettingRequest) (dto.SettingResponse, error)
}
