package repositories

import (
	"context"

	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/user"
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/infrastructure/persistence/gorm/models"
	"gorm.io/gorm"
)

type passwordPolicyRepository struct {
	db *gorm.DB
}

func NewPasswordPolicyRepository(db *gorm.DB) user.PasswordPolicyRepository {
	return &passwordPolicyRepository{db: db}
}

func (r *passwordPolicyRepository) GetByTenantID(ctx context.Context, tenantID string) (*user.PasswordPolicy, error) {
	var model models.PasswordPolicyModel
	if err := r.db.WithContext(ctx).First(&model, "tenant_id = ?", tenantID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Default policy
			return &user.PasswordPolicy{
				TenantID:        tenantID,
				MinLength:       8,
				RequireNumber:   true,
				RequireUppercase: true,
				RequireSpecial:   true,
			}, nil
		}
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *passwordPolicyRepository) Save(ctx context.Context, p *user.PasswordPolicy) error {
	model := models.FromPasswordPolicyDomain(p)
	return r.db.WithContext(ctx).Save(model).Error
}
