package migrations

import (
	"log"

	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/infrastructure/persistence/gorm/models"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Run executes the GORM auto-migrations for all domain models.
func Run(db *gorm.DB) error {
	err := db.AutoMigrate(
		&models.TenantModel{},
		&models.UserModel{},
		&models.RoleModel{},
		&models.AuditLogModel{},
		&models.RateLimitConfigModel{},
		&models.PasswordPolicyModel{},
		&models.GroupModel{},
	)
	if err != nil {
		return err
	}

	return Seed(db)
}

// Seed populates the database with initial required data.
func Seed(db *gorm.DB) error {
	log.Println("Seeding initial data...")

	// 1. Create System Tenant
	systemTenant := models.TenantModel{
		ID:       "00000000-0000-0000-0000-000000000000",
		Name:     "System",
		IsSystem: true,
		IsActive: true,
		Domains:  pq.StringArray{"system.local", "localhost"},
	}
	db.FirstOrCreate(&systemTenant, models.TenantModel{ID: systemTenant.ID})

	// 2. Create a default Standard Tenant
	acmeTenant := models.TenantModel{
		ID:       "11111111-1111-1111-1111-111111111111",
		Name:     "Acme Corp",
		IsSystem: false,
		IsActive: true,
		Domains:  pq.StringArray{"acme.com"},
	}
	db.FirstOrCreate(&acmeTenant, models.TenantModel{ID: acmeTenant.ID})

	// 3. Create Super Admin Role for System Tenant
	superAdminRole := models.RoleModel{
		ID:          "22222222-2222-2222-2222-222222222222",
		TenantID:    systemTenant.ID,
		Name:        "SuperAdmin",
		Permissions: pq.StringArray{"*:*:*"}, // Ultimate access
	}
	db.FirstOrCreate(&superAdminRole, models.RoleModel{ID: superAdminRole.ID})

	// 4. Create Admin Role for Acme Tenant
	acmeAdminRole := models.RoleModel{
		ID:          "33333333-3333-3333-3333-333333333333",
		TenantID:    acmeTenant.ID,
		Name:        "Admin",
		Permissions: pq.StringArray{"tenant:iam:manage", "user:iam:manage", "role:iam:manage"},
	}
	db.FirstOrCreate(&acmeAdminRole, models.RoleModel{ID: acmeAdminRole.ID})

	// 5. Create Password Policies
	policy := models.PasswordPolicyModel{
		TenantID:        systemTenant.ID,
		MinLength:       8,
		RequireNumber:   true,
		RequireUppercase: true,
		RequireSpecial:   true,
	}
	db.FirstOrCreate(&policy, models.PasswordPolicyModel{TenantID: policy.TenantID})

	acmePolicy := models.PasswordPolicyModel{
		TenantID:        acmeTenant.ID,
		MinLength:       8,
		RequireNumber:   true,
		RequireUppercase: true,
		RequireSpecial:   true,
	}
	db.FirstOrCreate(&acmePolicy, models.PasswordPolicyModel{TenantID: acmePolicy.TenantID})

	// 6. Create Users
	hash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)

	superAdminUser := models.UserModel{
		ID:           "44444444-4444-4444-4444-444444444444",
		TenantID:     systemTenant.ID,
		Username:     "superadmin",
		Email:        "superadmin@system.local",
		PasswordHash: string(hash),
		RoleID:       superAdminRole.ID,
		IsActive:     true,
	}
	db.FirstOrCreate(&superAdminUser, models.UserModel{ID: superAdminUser.ID})

	acmeAdminUser := models.UserModel{
		ID:           "55555555-5555-5555-5555-555555555555",
		TenantID:     acmeTenant.ID,
		Username:     "admin",
		Email:        "admin@acme.com",
		PasswordHash: string(hash),
		RoleID:       acmeAdminRole.ID,
		IsActive:     true,
	}
	db.FirstOrCreate(&acmeAdminUser, models.UserModel{ID: acmeAdminUser.ID})

	log.Println("Seeding completed successfully")
	return nil
}
