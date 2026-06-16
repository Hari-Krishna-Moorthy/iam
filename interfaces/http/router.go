package http

import (
	"net/http"

	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/application/auth"
	applicationRole "github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/application/role"
	applicationTenant "github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/application/tenant"
	applicationUser "github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/application/user"
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
	groupService applicationUser.GroupService,
	tenantService applicationTenant.Service,
	userService applicationUser.Service,
) *chi.Mux {
	r := chi.NewRouter()

	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(middleware.AuditMiddleware(auditRepo))

	authHandler := handlers.NewAuthHandler(authService)
	roleHandler := handlers.NewRoleHandler(roleService)
	groupHandler := handlers.NewGroupHandler(groupService)
	tenantHandler := handlers.NewTenantHandler(tenantService)
	userHandler := handlers.NewUserHandler(userService)

	// Public routes
	r.Group(func(r chi.Router) {
		r.Post("/login", authHandler.Login)
		// Usually, Tenant creation might be an internal admin API, but we'll expose it here for demonstration.
		r.Post("/tenants", tenantHandler.RegisterTenant)
	})

	// Protected routes (Feature D: Hydration)
	r.Group(func(r chi.Router) {
		r.Use(middleware.TenantMiddleware(tenantRepo))
		
		// Public to tenant, but requires tenant context
		r.Post("/users/register", userHandler.RegisterUser)
		
		// Fully protected routes (Require Token)
		r.Group(func(r chi.Router) {
			r.Use(middleware.RateLimitMiddleware(limiter))
			r.Use(middleware.AuthMiddleware(sessionRepo, tenantRepo))

			r.Get("/me", func(w http.ResponseWriter, r *http.Request) {
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

			// Group Management
			r.Route("/groups", func(r chi.Router) {
				r.Post("/", groupHandler.CreateGroup)
				r.Get("/", groupHandler.ListGroups)
				r.Post("/{id}/users/{userId}", groupHandler.AddUserToGroup)
				r.Post("/{id}/roles/{roleId}", groupHandler.AddRoleToGroup)
			})
		})
	})

	return r
}
