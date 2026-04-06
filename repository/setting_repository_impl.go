package repository

import (
	"Toko-Online/model"
	"context"

	"gorm.io/gorm"
)

type SettingRepositoryImpl struct {
	db *gorm.DB
}

func NewSettingRepository(db *gorm.DB) SettingRepository {
	return &SettingRepositoryImpl{db: db}
}

func (r *SettingRepositoryImpl) Get(ctx context.Context) (model.Setting, error) {
	var setting model.Setting
	err := r.db.WithContext(ctx).First(&setting).Error
	return setting, err
}

func (r *SettingRepositoryImpl) Update(ctx context.Context, setting *model.Setting) error {
	return r.db.WithContext(ctx).Model(setting).Updates(setting).Error
}
