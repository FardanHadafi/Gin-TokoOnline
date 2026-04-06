package repository

import (
	"Toko-Online/model"
	"context"
)

type SettingRepository interface {
	Get(ctx context.Context) (model.Setting, error)
	Update(ctx context.Context, setting *model.Setting) error
}
