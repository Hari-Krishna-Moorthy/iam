package ratelimit

import (
	"context"
	"time"
)

type Algorithm string

const (
	FixedWindow   Algorithm = "FIXED_WINDOW"
	SlidingWindow Algorithm = "SLIDING_WINDOW"
	LeakyBucket   Algorithm = "LEAKY_BUCKET"
	TokenBucket   Algorithm = "TOKEN_BUCKET"
)

type Config struct {
	TenantID  string
	Algorithm Algorithm
	Limit     int           // Max requests
	Window    time.Duration // Time window (for fixed/sliding)
	Burst     int           // Max burst (for buckets)
	Rate      float64       // Fill rate (per second, for buckets)
}

type Repository interface {
	GetByTenantID(ctx context.Context, tenantID string) (*Config, error)
	Save(ctx context.Context, config *Config) error
}

type Limiter interface {
	Allow(ctx context.Context, tenantID string) (bool, error)
}
