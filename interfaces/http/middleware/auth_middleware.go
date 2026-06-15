package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/session"
)

// AuthMiddleware validates the session and hydrates downstream headers.
func AuthMiddleware(sessionRepo session.Repository) func(http.Handler) http.Handler {
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

			// Hydrate downstream headers (Feature D)
			r.Header.Set("X-User-ID", sess.UserID)
			r.Header.Set("X-Tenant-ID", sess.TenantID)
			r.Header.Set("X-Role", sess.Role)

			var perms []string
			for _, p := range sess.Permissions {
				perms = append(perms, p.String())
			}
			r.Header.Set("X-Permissions", strings.Join(perms, ","))

			// Inject session into context for potential internal use
			ctx := context.WithValue(r.Context(), "session", sess)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
