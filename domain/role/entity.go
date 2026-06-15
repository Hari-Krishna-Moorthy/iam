package role

import (
	"context"
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/permission"
)

// Role represents a collection of permissions within a tenant.
type Role struct {
	ID          string
	TenantID    string
	Name        string
	Permissions []permission.Permission
}

// Repository defines the persistence interface for Role.
type Repository interface {
	GetByID(ctx context.Context, id string) (*Role, error)
	GetByTenantID(ctx context.Context, tenantID string) ([]Role, error)
	Save(ctx context.Context, role *Role) error
}
