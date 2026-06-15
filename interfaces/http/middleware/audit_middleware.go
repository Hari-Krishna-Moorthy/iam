package middleware

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"time"

	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/audit"
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/session"
)

func AuditMiddleware(repo audit.Repository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()

			// Capture request body for auditing (careful with large bodies/sensitive info)
			var bodyBytes []byte
			if r.Body != nil {
				bodyBytes, _ = io.ReadAll(r.Body)
				r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}

			// Capture Response
			rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
			next.ServeHTTP(rw, r)

			// Async audit logging
			go func() {
				tenantID, _ := r.Context().Value(TenantIDKey).(string)
				userID := "" // Extract from session if available
				if sess, ok := r.Context().Value("session").(*session.Session); ok {
					userID = sess.UserID
				}

				auditLog := &audit.AuditLog{
					TenantID:   tenantID,
					UserID:     userID,
					Action:     r.Method + " " + r.URL.Path,
					Resource:   r.URL.Path,
					Payload:    string(bodyBytes),
					IPAddress:  r.RemoteAddr,
					UserAgent:  r.UserAgent(),
					CreatedAt:  startTime,
				}
				repo.Save(context.Background(), auditLog)
			}()
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
