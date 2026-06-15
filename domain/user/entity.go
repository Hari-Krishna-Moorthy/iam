package user

import (
	"context"
	"time"
)

// User represents a person within a tenant.
type User struct {
	ID           string
	TenantID     string
	Username     string
	Email        string
	PasswordHash string
	RoleID       string
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Repository defines the persistence interface for User.
type Repository interface {
	GetByID(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, tenantID, email string) (*User, error)
	GetByUsername(ctx context.Context, tenantID, username string) (*User, error)
	Save(ctx context.Context, user *User) error
}
