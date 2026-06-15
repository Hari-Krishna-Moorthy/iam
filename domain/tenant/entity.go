package tenant

import (
	"context"
	"time"
)

// Tenant represents a single customer or organization.
type Tenant struct {
	ID        string
	Name      string
	Domains   []string // List of hostnames/origins mapped to this tenant
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Repository defines the persistence interface for Tenant.
type Repository interface {
	GetByID(ctx context.Context, id string) (*Tenant, error)
	GetByDomain(ctx context.Context, domain string) (*Tenant, error)
	Save(ctx context.Context, tenant *Tenant) error
}
