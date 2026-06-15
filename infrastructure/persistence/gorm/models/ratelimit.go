package models

import (
	"time"

	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/ratelimit"
)

type RateLimitConfigModel struct {
	TenantID  string `gorm:"primaryKey;type:uuid"`
	Algorithm string `gorm:"not null;default:'FIXED_WINDOW'"`
	Limit     int    `gorm:"not null;default:100"`
	WindowSec int64  `gorm:"not null;default:60"`
	Burst     int    `gorm:"not null;default:10"`
	Rate      float64 `gorm:"not null;default:1.0"`
}

func (RateLimitConfigModel) TableName() string {
	return "rate_limit_configs"
}

func (m *RateLimitConfigModel) ToDomain() *ratelimit.Config {
	return &ratelimit.Config{
		TenantID:  m.TenantID,
		Algorithm: ratelimit.Algorithm(m.Algorithm),
		Limit:     m.Limit,
		Window:    time.Duration(m.WindowSec) * time.Second,
		Burst:     m.Burst,
		Rate:      m.Rate,
	}
}

func FromRateLimitDomain(c *ratelimit.Config) *RateLimitConfigModel {
	return &RateLimitConfigModel{
		TenantID:  c.TenantID,
		Algorithm: string(c.Algorithm),
		Limit:     c.Limit,
		WindowSec: int64(c.Window.Seconds()),
		Burst:     c.Burst,
		Rate:      c.Rate,
	}
}
