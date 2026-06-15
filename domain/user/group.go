package user

import (
	"context"
)

type Group struct {
	ID       string
	TenantID string
	Name     string
	RoleIDs  []string
	UserIDs  []string
}

type GroupRepository interface {
	GetByID(ctx context.Context, id string) (*Group, error)
	GetByTenantID(ctx context.Context, tenantID string) ([]Group, error)
	Save(ctx context.Context, group *Group) error
	Delete(ctx context.Context, id string) error
	AddUser(ctx context.Context, groupID, userID string) error
	RemoveUser(ctx context.Context, groupID, userID string) error
	AddRole(ctx context.Context, groupID, roleID string) error
	RemoveRole(ctx context.Context, groupID, roleID string) error
}
