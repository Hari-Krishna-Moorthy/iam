package audit

import (
	"context"
	"time"
)

type AuditLog struct {
	ID         string
	TenantID   string
	UserID     string
	Action     string // e.g., "LOGIN_SUCCESS", "LOGIN_FAILURE", "USER_REGISTERED"
	Resource   string // e.g., "auth", "user"
	Payload    string // JSON or metadata
	IPAddress  string
	UserAgent  string
	CreatedAt  time.Time
}

type Repository interface {
	Save(ctx context.Context, log *AuditLog) error
	GetByTenantID(ctx context.Context, tenantID string) ([]*AuditLog, error)
}
