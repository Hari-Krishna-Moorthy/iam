package session

import (
	"context"
	"time"
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/permission"
)

// Session represents an active user session stored in Redis.
type Session struct {
	ID          string
	UserID      string
	TenantID    string
	Role        string
	Permissions []permission.Permission
	CreatedAt   time.Time
	ExpiresAt   time.Time
}

// Repository defines the persistence interface for Session (Redis).
type Repository interface {
	Save(ctx context.Context, session *Session) error
	GetByID(ctx context.Context, sessionID string) (*Session, error)
	Delete(ctx context.Context, sessionID string) error
	GetByUserID(ctx context.Context, userID string) ([]*Session, error)
}

// AuthStrategy defines the interface for different authentication methods.
type AuthStrategy interface {
	Authenticate(ctx context.Context, credentials map[string]string) (*Session, error)
	ValidateToken(ctx context.Context, token string) (*Session, error)
}
