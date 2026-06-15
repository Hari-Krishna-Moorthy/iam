package models

import (
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/user"
)

type PasswordPolicyModel struct {
	TenantID        string `gorm:"primaryKey;type:uuid"`
	MinLength       int    `gorm:"not null;default:8"`
	RequireNumber   bool   `gorm:"not null;default:true"`
	RequireUppercase bool   `gorm:"not null;default:true"`
	RequireSpecial   bool   `gorm:"not null;default:true"`
}

func (PasswordPolicyModel) TableName() string {
	return "password_policies"
}

func (m *PasswordPolicyModel) ToDomain() *user.PasswordPolicy {
	return &user.PasswordPolicy{
		TenantID:        m.TenantID,
		MinLength:       m.MinLength,
		RequireNumber:   m.RequireNumber,
		RequireUppercase: m.RequireUppercase,
		RequireSpecial:   m.RequireSpecial,
	}
}

func FromPasswordPolicyDomain(p *user.PasswordPolicy) *PasswordPolicyModel {
	return &PasswordPolicyModel{
		TenantID:        p.TenantID,
		MinLength:       p.MinLength,
		RequireNumber:   p.RequireNumber,
		RequireUppercase: p.RequireUppercase,
		RequireSpecial:   p.RequireSpecial,
	}
}
