package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/session"
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/tenant"
)

// AuthMiddleware validates the session and hydrates downstream headers.
func AuthMiddleware(sessionRepo session.Repository, tenantRepo tenant.Repository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			sessionID := strings.TrimPrefix(authHeader, "Bearer ")
			sess, err := sessionRepo.GetByID(r.Context(), sessionID)
			if err != nil {
				http.Error(w, "Invalid session", http.StatusUnauthorized)
				return
			}

			effectiveTenantID := sess.TenantID

			// Handle Super Admin Cross-Tenant Switching
			targetTenantID := r.Header.Get("X-Target-Tenant-ID")
			if targetTenantID != "" && targetTenantID != effectiveTenantID {
				// Check if the user's home tenant is the System Tenant
				homeTenant, err := tenantRepo.GetByID(r.Context(), sess.TenantID)
				if err == nil && homeTenant != nil && homeTenant.IsSystem {
					// Verify target tenant actually exists before allowing switch
					_, err := tenantRepo.GetByID(r.Context(), targetTenantID)
					if err == nil {
						effectiveTenantID = targetTenantID
					} else {
						http.Error(w, "Target tenant not found", http.StatusNotFound)
						return
					}
				} else {
					http.Error(w, "Forbidden: Only System Admins can switch tenants", http.StatusForbidden)
					return
				}
			}

			// Hydrate downstream headers (Feature D)
			r.Header.Set("X-User-ID", sess.UserID)
			r.Header.Set("X-Tenant-ID", effectiveTenantID)
			r.Header.Set("X-Role", sess.Role)

			var perms []string
			for _, p := range sess.Permissions {
				perms = append(perms, p.String())
			}
			r.Header.Set("X-Permissions", strings.Join(perms, ","))

			// Inject session into context for potential internal use
			ctx := context.WithValue(r.Context(), "session", sess)
			// Override the TenantID in the context with the effective one
			ctx = context.WithValue(ctx, TenantIDKey, effectiveTenantID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
