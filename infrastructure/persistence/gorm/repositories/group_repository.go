package repositories

import (
	"context"

	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/user"
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/infrastructure/persistence/gorm/models"
	"gorm.io/gorm"
)

type groupRepository struct {
	db *gorm.DB
}

func NewGroupRepository(db *gorm.DB) user.GroupRepository {
	return &groupRepository{db: db}
}

func (r *groupRepository) GetByID(ctx context.Context, id string) (*user.Group, error) {
	var model models.GroupModel
	if err := r.db.WithContext(ctx).Preload("Users").Preload("Roles").First(&model, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *groupRepository) GetByTenantID(ctx context.Context, tenantID string) ([]user.Group, error) {
	var modelsList []models.GroupModel
	if err := r.db.WithContext(ctx).Preload("Users").Preload("Roles").Where("tenant_id = ?", tenantID).Find(&modelsList).Error; err != nil {
		return nil, err
	}

	groups := make([]user.Group, 0, len(modelsList))
	for _, m := range modelsList {
		groups = append(groups, *m.ToDomain())
	}
	return groups, nil
}

func (r *groupRepository) Save(ctx context.Context, g *user.Group) error {
	model := models.FromGroupDomain(g)
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *groupRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.GroupModel{}, "id = ?", id).Error
}

func (r *groupRepository) AddUser(ctx context.Context, groupID, userID string) error {
	return r.db.WithContext(ctx).Model(&models.GroupModel{ID: groupID}).Association("Users").Append(&models.UserModel{ID: userID})
}

func (r *groupRepository) RemoveUser(ctx context.Context, groupID, userID string) error {
	return r.db.WithContext(ctx).Model(&models.GroupModel{ID: groupID}).Association("Users").Delete(&models.UserModel{ID: userID})
}

func (r *groupRepository) AddRole(ctx context.Context, groupID, roleID string) error {
	return r.db.WithContext(ctx).Model(&models.GroupModel{ID: groupID}).Association("Roles").Append(&models.RoleModel{ID: roleID})
}

func (r *groupRepository) RemoveRole(ctx context.Context, groupID, roleID string) error {
	return r.db.WithContext(ctx).Model(&models.GroupModel{ID: groupID}).Association("Roles").Delete(&models.RoleModel{ID: roleID})
}
