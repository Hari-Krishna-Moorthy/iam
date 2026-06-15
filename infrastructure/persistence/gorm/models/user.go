package models

import (
	"time"

	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/user"
	"gorm.io/gorm"
)

// UserModel is the GORM representation of the User entity.
type UserModel struct {
	ID           string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	TenantID     string `gorm:"type:uuid;index;not null"`
	Username     string `gorm:"not null"`
	Email        string `gorm:"not null"`
	PasswordHash string `gorm:"not null"`
	RoleID       string `gorm:"type:uuid;index;not null"`
	IsActive     bool   `gorm:"default:true"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`

	// Unique index on TenantID + Email and TenantID + Username
	_ struct{} `gorm:"uniqueIndex:idx_tenant_email,priority:1"`
	_ struct{} `gorm:"uniqueIndex:idx_tenant_username,priority:1"`
}

// TableName overrides the table name for UserModel.
func (UserModel) TableName() string {
	return "users"
}

// ToDomain maps the GORM model to the pure domain entity.
func (m *UserModel) ToDomain() *user.User {
	return &user.User{
		ID:           m.ID,
		TenantID:     m.TenantID,
		Username:     m.Username,
		Email:        m.Email,
		PasswordHash: m.PasswordHash,
		RoleID:       m.RoleID,
		IsActive:     m.IsActive,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

// FromUserDomain creates a GORM model from the pure domain entity.
func FromUserDomain(u *user.User) *UserModel {
	return &UserModel{
		ID:           u.ID,
		TenantID:     u.TenantID,
		Username:     u.Username,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		RoleID:       u.RoleID,
		IsActive:     u.IsActive,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}
