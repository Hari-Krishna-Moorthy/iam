package repositories

import (
	"context"

	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/tenant"
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/infrastructure/persistence/gorm/models"
	"gorm.io/gorm"
)

type tenantRepository struct {
	db *gorm.DB
}

func NewTenantRepository(db *gorm.DB) tenant.Repository {
	return &tenantRepository{db: db}
}

func (r *tenantRepository) GetByID(ctx context.Context, id string) (*tenant.Tenant, error) {
	var model models.TenantModel
	if err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *tenantRepository) GetByDomain(ctx context.Context, domain string) (*tenant.Tenant, error) {
	var model models.TenantModel
	// PostgreSQL array overlap or inclusion check
	if err := r.db.WithContext(ctx).Where("? = ANY(domains)", domain).First(&model).Error; err != nil {
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *tenantRepository) GetAll(ctx context.Context) ([]tenant.Tenant, error) {
	var models []models.TenantModel
	if err := r.db.WithContext(ctx).Find(&models).Error; err != nil {
		return nil, err
	}
	var tenants []tenant.Tenant
	for _, m := range models {
		tenants = append(tenants, *m.ToDomain())
	}
	return tenants, nil
}

func (r *tenantRepository) Save(ctx context.Context, t *tenant.Tenant) error {
	model := models.FromTenantDomain(t)
	return r.db.WithContext(ctx).Save(model).Error
}
