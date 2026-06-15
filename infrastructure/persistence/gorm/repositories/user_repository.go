package repositories

import (
	"context"

	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/user"
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/infrastructure/persistence/gorm/models"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) user.Repository {
	return &userRepository{db: db}
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*user.User, error) {
	var model models.UserModel
	if err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *userRepository) GetByEmail(ctx context.Context, tenantID, email string) (*user.User, error) {
	var model models.UserModel
	if err := r.db.WithContext(ctx).Where("tenant_id = ? AND email = ?", tenantID, email).First(&model).Error; err != nil {
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *userRepository) GetByUsername(ctx context.Context, tenantID, username string) (*user.User, error) {
	var model models.UserModel
	if err := r.db.WithContext(ctx).Where("tenant_id = ? AND username = ?", tenantID, username).First(&model).Error; err != nil {
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *userRepository) Save(ctx context.Context, u *user.User) error {
	model := models.FromUserDomain(u)
	return r.db.WithContext(ctx).Save(model).Error
}
