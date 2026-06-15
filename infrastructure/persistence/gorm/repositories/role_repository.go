package repositories

import (
	"context"

	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/role"
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/infrastructure/persistence/gorm/models"
	"gorm.io/gorm"
)

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) role.Repository {
	return &roleRepository{db: db}
}

func (r *roleRepository) GetByID(ctx context.Context, id string) (*role.Role, error) {
	var model models.RoleModel
	if err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *roleRepository) GetByTenantID(ctx context.Context, tenantID string) ([]role.Role, error) {
	var modelsList []models.RoleModel
	if err := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID).Find(&modelsList).Error; err != nil {
		return nil, err
	}

	roles := make([]role.Role, 0, len(modelsList))
	for _, m := range modelsList {
		roles = append(roles, *m.ToDomain())
	}
	return roles, nil
}

func (r *roleRepository) Save(ctx context.Context, roleEntity *role.Role) error {
	model := models.FromRoleDomain(roleEntity)
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *roleRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.RoleModel{}, "id = ?", id).Error
}
