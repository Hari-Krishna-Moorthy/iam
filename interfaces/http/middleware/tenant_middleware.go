package middleware

import (
	"context"
	"net/http"
	"net/url"
	"strings"

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

			// Normalize domain: remove protocol and port
			domain := normalizeDomain(origin)

			t, err := repo.GetByDomain(r.Context(), domain)
			if err != nil {
				http.Error(w, "Tenant not found for domain: "+domain, http.StatusNotFound)
				return
			}

			ctx := context.WithValue(r.Context(), TenantIDKey, t.ID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func normalizeDomain(raw string) string {
	// 1. Try parsing as URL to handle http://localhost:5173
	if strings.Contains(raw, "://") {
		u, err := url.Parse(raw)
		if err == nil {
			host := u.Hostname()
			if host != "" {
				return host
			}
		}
	}

	// 2. Handle cases like localhost:8080 or just google.com
	host := raw
	if strings.Contains(host, ":") {
		parts := strings.Split(host, ":")
		host = parts[0]
	}
	return host
}
