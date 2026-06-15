package main

import (
	"log"
	"net/http"

	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/application/auth"
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/application/auth/strategies"
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
	db.AutoMigrate(&models.TenantModel{}, &models.UserModel{}, &models.RoleModel{}, &models.AuditLogModel{})

	// 2. Setup Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.RedisURL,
	})

	// 3. Setup Repositories
	tenantRepo := repositories.NewTenantRepository(db)
	userRepo := repositories.NewUserRepository(db)
	roleRepo := repositories.NewRoleRepository(db)
	auditRepo := repositories.NewAuditRepository(db)
	sessRepo := redisRepo.NewSessionRepository(rdb)

	// 4. Setup Providers
	jwtProvider := infraAuth.NewJWTProvider("my-secret-key")

	// 5. Setup Auth Strategies
	pwdStrategy := strategies.NewPasswordStrategy(userRepo, roleRepo)
	authStrategies := map[string]session.AuthStrategy{
		"password": pwdStrategy,
	}

	// 6. Setup Services
	authService := auth.NewService(tenantRepo, sessRepo, jwtProvider, authStrategies)

	// 7. Setup Router
	r := interfacesHttp.NewRouter(tenantRepo, sessRepo, auditRepo, authService)

	// 8. Start Server
	log.Printf("Server starting on port %s...", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
