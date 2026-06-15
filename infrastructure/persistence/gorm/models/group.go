package models

import (
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/user"
	"gorm.io/gorm"
)

type GroupModel struct {
	ID       string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	TenantID string `gorm:"type:uuid;index;not null"`
	Name     string `gorm:"not null"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// Many-to-Many relationships
	Users []UserModel `gorm:"many2many:group_users;"`
	Roles []RoleModel `gorm:"many2many:group_roles;"`
}

func (GroupModel) TableName() string {
	return "groups"
}

func (m *GroupModel) ToDomain() *user.Group {
	uIDs := make([]string, len(m.Users))
	for i, u := range m.Users {
		uIDs[i] = u.ID
	}
	rIDs := make([]string, len(m.Roles))
	for i, r := range m.Roles {
		rIDs[i] = r.ID
	}

	return &user.Group{
		ID:       m.ID,
		TenantID: m.TenantID,
		Name:     m.Name,
		UserIDs:  uIDs,
		RoleIDs:  rIDs,
	}
}

func FromGroupDomain(g *user.Group) *GroupModel {
	return &GroupModel{
		ID:       g.ID,
		TenantID: g.TenantID,
		Name:     g.Name,
	}
}
