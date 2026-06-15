package migrations

import (
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/infrastructure/persistence/gorm/models"
	"gorm.io/gorm"
)

// Run executes the GORM auto-migrations for all domain models.
// In a production scenario, you might replace this with a tool like golang-migrate or goose.
func Run(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.TenantModel{},
		&models.UserModel{},
		&models.RoleModel{},
		&models.AuditLogModel{},
		&models.RateLimitConfigModel{},
		&models.PasswordPolicyModel{},
		&models.GroupModel{},
	)
}
