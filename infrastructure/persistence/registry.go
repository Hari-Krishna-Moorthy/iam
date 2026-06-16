package persistence

import (
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/audit"
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/job"
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/ratelimit"
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/role"
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/session"
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/tenant"
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/user"
	gormRepos "github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/infrastructure/persistence/gorm/repositories"
	redisRepos "github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/infrastructure/persistence/redis/repositories"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// Repositories acts as a registry/container for all domain repositories.
type Repositories struct {
	Tenant         tenant.Repository
	User           user.Repository
	Role           role.Repository
	Group          user.GroupRepository
	Audit          audit.Repository
	RateLimit      ratelimit.Repository
	PasswordPolicy user.PasswordPolicyRepository
	Session        session.Repository
	Job            job.Repository
}

// NewRepositories initializes and returns a complete registry of repositories.
func NewRepositories(db *gorm.DB, rdb *redis.Client) *Repositories {
	return &Repositories{
		Tenant:         gormRepos.NewTenantRepository(db),
		User:           gormRepos.NewUserRepository(db),
		Role:           gormRepos.NewRoleRepository(db),
		Group:          gormRepos.NewGroupRepository(db),
		Audit:          gormRepos.NewAuditRepository(db),
		RateLimit:      gormRepos.NewRateLimitRepository(db),
		PasswordPolicy: gormRepos.NewPasswordPolicyRepository(db),
		Session:        redisRepos.NewSessionRepository(rdb),
		Job:            redisRepos.NewJobRepository(rdb),
	}
}
