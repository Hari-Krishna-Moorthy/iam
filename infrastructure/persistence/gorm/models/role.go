package models

import (
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/permission"
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/role"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// RoleModel is the GORM representation of the Role entity.
type RoleModel struct {
	ID          string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	TenantID    string         `gorm:"type:uuid;index;not null"`
	Name        string         `gorm:"not null"`
	Permissions pq.StringArray `gorm:"type:text[];not null"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`

	_ struct{} `gorm:"uniqueIndex:idx_tenant_role_name"`
}

// TableName overrides the table name for RoleModel.
func (RoleModel) TableName() string {
	return "roles"
}

// ToDomain maps the GORM model to the pure domain entity.
func (m *RoleModel) ToDomain() *role.Role {
	perms := make([]permission.Permission, 0, len(m.Permissions))
	for _, pStr := range m.Permissions {
		p, _ := permission.Parse(pStr)
		perms = append(perms, p)
	}

	return &role.Role{
		ID:          m.ID,
		TenantID:    m.TenantID,
		Name:        m.Name,
		Permissions: perms,
	}
}

// FromRoleDomain creates a GORM model from the pure domain entity.
func FromRoleDomain(r *role.Role) *RoleModel {
	pStrs := make([]string, 0, len(r.Permissions))
	for _, p := range r.Permissions {
		pStrs = append(pStrs, p.String())
	}

	return &RoleModel{
		ID:          r.ID,
		TenantID:    r.TenantID,
		Name:        r.Name,
		Permissions: pq.StringArray(pStrs),
	}
}
