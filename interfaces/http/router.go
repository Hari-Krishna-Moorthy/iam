package http

import (
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/application/auth"
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
	authService auth.Service,
) *chi.Mux {
	r := chi.NewRouter()

	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)

	authHandler := handlers.NewAuthHandler(authService)

	// Public routes
	r.Group(func(r chi.Router) {
		r.Post("/login", authHandler.Login)
	})

	// Protected routes (Feature D: Hydration)
	r.Group(func(r chi.Router) {
		r.Use(middleware.TenantMiddleware(tenantRepo))
		r.Use(middleware.AuthMiddleware(sessionRepo))

		r.Get("/me", func(w http.ResponseWriter, r *http.Request) {
			// Example of accessing hydrated headers
			userID := r.Header.Get("X-User-ID")
			w.Write([]byte("Hello, user " + userID))
		})
	})

	return r
}
