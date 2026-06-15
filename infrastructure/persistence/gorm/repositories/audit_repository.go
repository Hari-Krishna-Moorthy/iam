package repositories

import (
	"context"

	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/audit"
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/infrastructure/persistence/gorm/models"
	"gorm.io/gorm"
)

type auditRepository struct {
	db *gorm.DB
}

func NewAuditRepository(db *gorm.DB) audit.Repository {
	return &auditRepository{db: db}
}

func (r *auditRepository) Save(ctx context.Context, l *audit.AuditLog) error {
	model := models.FromAuditDomain(l)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *auditRepository) GetByTenantID(ctx context.Context, tenantID string) ([]*audit.AuditLog, error) {
	var modelsList []models.AuditLogModel
	if err := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID).Find(&modelsList).Error; err != nil {
		return nil, err
	}

	logs := make([]*audit.AuditLog, 0, len(modelsList))
	for _, m := range modelsList {
		logs = append(logs, m.ToDomain())
	}
	return logs, nil
}
