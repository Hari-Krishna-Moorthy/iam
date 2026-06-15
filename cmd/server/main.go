package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/application/auth"
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/application/auth/strategies"
	applicationRateLimit "github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/application/ratelimit"
	applicationRole "github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/application/role"
	applicationTenant "github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/application/tenant"
	applicationUser "github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/application/user"
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/session"
	infraAuth "github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/infrastructure/auth"
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/infrastructure/config"
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/infrastructure/persistence"
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/infrastructure/persistence/gorm/migrations"
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

	// 2. Run Database Migrations
	if err := migrations.Run(db); err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}
	log.Println("Database migrations applied successfully")

	// 3. Setup Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.RedisURL,
	})

	// 4. Setup Repositories (Registry pattern)
	repos := persistence.NewRepositories(db, rdb)

	// 4. Setup Providers
	jwtProvider := infraAuth.NewJWTProvider("my-secret-key")
	oauthProvider := infraAuth.NewDummyOAuth2Provider()

	// 5. Setup Auth Strategies
	pwdStrategy := strategies.NewPasswordStrategy(repos.User, repos.Role)
	oauthStrategy := strategies.NewOAuth2Strategy(repos.User, repos.Role, oauthProvider)
	authStrategies := map[string]session.AuthStrategy{
		"password": pwdStrategy,
		"oauth2":   oauthStrategy,
	}

	// 6. Setup Limiters

	limiter := applicationRateLimit.NewRedisLimiter(rdb, repos.RateLimit)

	// 8. Setup Application Services
	authService := auth.NewService(repos.Tenant, repos.Session, jwtProvider, authStrategies)
	roleService := applicationRole.NewService(repos.Role)
	groupService := applicationUser.NewGroupService(repos.Group)
	tenantService := applicationTenant.NewService(repos.Tenant)
	userService := applicationUser.NewService(repos.User, repos.PasswordPolicy)

	// 9. Setup HTTP Router
	r := interfacesHttp.NewRouter(
		repos.Tenant,
		repos.Session,
		repos.Audit,
		limiter,
		authService,
		roleService,
		groupService,
		tenantService,
		userService,
	)

	// 10. Setup Server for Graceful Shutdown
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	// Start server in a goroutine so it doesn't block
	go func() {
		log.Printf("Server starting on port %s...", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// 11. Listen for OS signals to initiate graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Block until a signal is received
	<-quit
	log.Println("Shutdown signal received, shutting down gracefully...")

	// Create a context with a timeout for the shutdown process
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown the server (waits for active connections to finish)
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited cleanly")
}
