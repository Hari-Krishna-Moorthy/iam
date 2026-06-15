package models

import (
	"time"

	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/tenant"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// TenantModel is the GORM representation of the Tenant entity.
type TenantModel struct {
	ID        string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name      string         `gorm:"not null"`
	Domains   pq.StringArray `gorm:"type:text[];index:,type:gin"`
	IsActive  bool           `gorm:"default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName overrides the table name for TenantModel.
func (TenantModel) TableName() string {
	return "tenants"
}

// ToDomain maps the GORM model to the pure domain entity.
func (m *TenantModel) ToDomain() *tenant.Tenant {
	return &tenant.Tenant{
		ID:        m.ID,
		Name:      m.Name,
		Domains:   []string(m.Domains),
		IsActive:  m.IsActive,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

// FromDomain creates a GORM model from the pure domain entity.
func FromTenantDomain(t *tenant.Tenant) *TenantModel {
	return &TenantModel{
		ID:        t.ID,
		Name:      t.Name,
		Domains:   pq.StringArray(t.Domains),
		IsActive:  t.IsActive,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}
