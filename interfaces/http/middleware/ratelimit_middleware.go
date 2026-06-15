package middleware

import (
	"net/http"

	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/ratelimit"
)

// RateLimitMiddleware enforces tenant-specific rate limits.
func RateLimitMiddleware(limiter ratelimit.Limiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tenantID, ok := r.Context().Value(TenantIDKey).(string)
			if !ok {
				// If tenant is not identified, we might want to skip or reject
				next.ServeHTTP(w, r)
				return
			}

			allowed, err := limiter.Allow(r.Context(), tenantID)
			if err != nil {
				// Handle internal error (log it)
				http.Error(w, "Internal rate limit error", http.StatusInternalServerError)
				return
			}

			if !allowed {
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
