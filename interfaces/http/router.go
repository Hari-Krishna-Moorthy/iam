package http

import (
	"net/http"

	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/application/auth"
	applicationRole "github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/application/role"
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/audit"
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/ratelimit"
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/session"
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/tenant"
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/interfaces/http/handlers"
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/interfaces/http/middleware"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

func NewRouter(
	tenantRepo tenant.Repository,
	sessionRepo session.Repository,
	auditRepo audit.Repository,
	limiter ratelimit.Limiter,
	authService auth.Service,
	roleService applicationRole.Service,
) *chi.Mux {
	r := chi.NewRouter()

	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(middleware.AuditMiddleware(auditRepo))

	authHandler := handlers.NewAuthHandler(authService)
	roleHandler := handlers.NewRoleHandler(roleService)

	// Public routes
	r.Group(func(r chi.Router) {
		r.Post("/login", authHandler.Login)
	})

	// Protected routes (Feature D: Hydration)
	r.Group(func(r chi.Router) {
		r.Use(middleware.TenantMiddleware(tenantRepo))
		r.Use(middleware.RateLimitMiddleware(limiter))
		r.Use(middleware.AuthMiddleware(sessionRepo))

		r.Get("/me", func(w http.ResponseWriter, r *http.Request) {
			// Example of accessing hydrated headers
			userID := r.Header.Get("X-User-ID")
			w.Write([]byte("Hello, user " + userID))
		})

		// Role Management
		r.Route("/roles", func(r chi.Router) {
			r.Post("/", roleHandler.CreateRole)
			r.Get("/", roleHandler.ListRoles)
			r.Get("/{id}", roleHandler.GetRole)
			r.Put("/{id}", roleHandler.UpdateRole)
			r.Delete("/{id}", roleHandler.DeleteRole)
		})
	})

	return r
}
