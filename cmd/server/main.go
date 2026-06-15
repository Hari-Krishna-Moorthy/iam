package main

import (
	"log"
	"net/http"

	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/application/auth"
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/application/auth/strategies"
	applicationRateLimit "github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/application/ratelimit"
	applicationRole "github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/application/role"
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/session"
	infraAuth "github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/infrastructure/auth"
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/infrastructure/config"
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/infrastructure/persistence/gorm/models"
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/infrastructure/persistence/gorm/repositories"
	redisRepo "github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/infrastructure/persistence/redis/repositories"
	interfacesHttp "github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/interfaces/http"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()

	// 1. Setup DB (GORM)
	db, err := gorm.Open(postgres.Open(cfg.DBURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate models
	db.AutoMigrate(
		&models.TenantModel{},
		&models.UserModel{},
		&models.RoleModel{},
		&models.AuditLogModel{},
		&models.RateLimitConfigModel{},
	)

	// 2. Setup Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.RedisURL,
	})

	// 3. Setup Repositories
	tenantRepo := repositories.NewTenantRepository(db)
	userRepo := repositories.NewUserRepository(db)
	roleRepo := repositories.NewRoleRepository(db)
	auditRepo := repositories.NewAuditRepository(db)
	ratelimitRepo := repositories.NewRateLimitRepository(db)
	sessRepo := redisRepo.NewSessionRepository(rdb)

	// 4. Setup Providers
	jwtProvider := infraAuth.NewJWTProvider("my-secret-key")

	// 5. Setup Auth Strategies
	pwdStrategy := strategies.NewPasswordStrategy(userRepo, roleRepo)
	authStrategies := map[string]session.AuthStrategy{
		"password": pwdStrategy,
	}

	// 6. Setup Limiters
	limiter := applicationRateLimit.NewRedisLimiter(rdb, ratelimitRepo)

	// 7. Setup Services
	authService := auth.NewService(tenantRepo, sessRepo, jwtProvider, authStrategies)
	roleService := applicationRole.NewService(roleRepo)

	// 8. Setup Router
	r := interfacesHttp.NewRouter(tenantRepo, sessRepo, auditRepo, limiter, authService, roleService)

	// 9. Start Server
	log.Printf("Server starting on port %s...", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
