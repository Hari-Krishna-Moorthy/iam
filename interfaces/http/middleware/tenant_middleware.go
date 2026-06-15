package middleware

import (
	"context"
	"net/http"

	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/tenant"
)

type contextKey string

const TenantIDKey contextKey = "tenant_id"

// TenantMiddleware extracts the tenant ID based on the request origin.
func TenantMiddleware(repo tenant.Repository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin == "" {
				origin = r.Host
			}

			t, err := repo.GetByDomain(r.Context(), origin)
			if err != nil {
				http.Error(w, "Tenant not found", http.StatusNotFound)
				return
			}

			ctx := context.WithValue(r.Context(), TenantIDKey, t.ID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
