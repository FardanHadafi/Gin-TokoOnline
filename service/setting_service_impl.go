package service

import (
	"Toko-Online/config"
	"Toko-Online/dto"
	"Toko-Online/repository"
	"context"
	"log/slog"
	"net/http"
)

type SettingServiceImpl struct {
	repo   repository.SettingRepository
	logger *slog.Logger
}

func NewSettingService(repo repository.SettingRepository, logger *slog.Logger) SettingService {
	return &SettingServiceImpl{
		repo:   repo,
		logger: logger,
	}
}

func (s *SettingServiceImpl) Get(ctx context.Context) (dto.SettingResponse, error) {
	s.logger.InfoContext(ctx, "Getting store settings")
	setting, err := s.repo.Get(ctx)
	if err != nil {
		return dto.SettingResponse{}, &config.ApiError{
			Status: http.StatusInternalServerError,
			Title:  "Internal Error",
			Detail: "Failed to fetch settings",
		}
	}

	return dto.SettingResponse{
		ID:             setting.ID,
		StoreName:      setting.StoreName,
		WhatsAppNumber: setting.WhatsAppNumber,
		AddressInfo:    setting.AddressInfo,
		WelcomeMessage: setting.WelcomeMessage,
	}, nil
}

func (s *SettingServiceImpl) Update(ctx context.Context, req dto.UpdateSettingRequest) (dto.SettingResponse, error) {
	s.logger.InfoContext(ctx, "Updating store settings")
	setting, err := s.repo.Get(ctx)
	if err != nil {
		return dto.SettingResponse{}, &config.ApiError{
			Status: http.StatusInternalServerError,
			Title:  "Internal Error",
			Detail: "Failed to fetch setting",
		}
	}

	if req.StoreName != "" {
		setting.StoreName = req.StoreName
	}
	if req.WhatsAppNumber != "" {
		setting.WhatsAppNumber = req.WhatsAppNumber
	}
	if req.AddressInfo != "" {
		setting.AddressInfo = req.AddressInfo
	}
	if req.WelcomeMessage != "" {
		setting.WelcomeMessage = req.WelcomeMessage
	}

	if err := s.repo.Update(ctx, &setting); err != nil {
		s.logger.ErrorContext(ctx, "Failed to update settings", "error", err)
		return dto.SettingResponse{}, &config.ApiError{
			Status: http.StatusInternalServerError,
			Title:  "Internal Error",
			Detail: "Failed to update settings",
		}
	}

	return dto.SettingResponse{
		ID:             setting.ID,
		StoreName:      setting.StoreName,
		WhatsAppNumber: setting.WhatsAppNumber,
		AddressInfo:    setting.AddressInfo,
		WelcomeMessage: setting.WelcomeMessage,
	}, nil
}
