package repositories

import (
	"context"

	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/ratelimit"
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/infrastructure/persistence/gorm/models"
	"gorm.io/gorm"
)

type rateLimitRepository struct {
	db *gorm.DB
}

func NewRateLimitRepository(db *gorm.DB) ratelimit.Repository {
	return &rateLimitRepository{db: db}
}

func (r *rateLimitRepository) GetByTenantID(ctx context.Context, tenantID string) (*ratelimit.Config, error) {
	var model models.RateLimitConfigModel
	if err := r.db.WithContext(ctx).First(&model, "tenant_id = ?", tenantID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Return default config if not found
			return &ratelimit.Config{
				TenantID:  tenantID,
				Algorithm: ratelimit.FixedWindow,
				Limit:     100,
				Window:    60 * 1000 * 1000 * 1000, // 60s
			}, nil
		}
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *rateLimitRepository) Save(ctx context.Context, c *ratelimit.Config) error {
	model := models.FromRateLimitDomain(c)
	return r.db.WithContext(ctx).Save(model).Error
}
