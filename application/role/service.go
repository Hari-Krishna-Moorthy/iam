package role

import (
	"context"

	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/permission"
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/role"
)

type CreateRoleRequest struct {
	TenantID    string
	Name        string
	Permissions []string
}

type UpdateRoleRequest struct {
	ID          string
	Name        string
	Permissions []string
}

type Service interface {
	CreateRole(ctx context.Context, req CreateRoleRequest) (*role.Role, error)
	UpdateRole(ctx context.Context, req UpdateRoleRequest) (*role.Role, error)
	DeleteRole(ctx context.Context, id string) error
	GetRole(ctx context.Context, id string) (*role.Role, error)
	ListRoles(ctx context.Context, tenantID string) ([]role.Role, error)
}

type roleService struct {
	repo role.Repository
}

func NewService(repo role.Repository) Service {
	return &roleService{repo: repo}
}

func (s *roleService) CreateRole(ctx context.Context, req CreateRoleRequest) (*role.Role, error) {
	perms := make([]permission.Permission, 0, len(req.Permissions))
	for _, pStr := range req.Permissions {
		p, err := permission.Parse(pStr)
		if err != nil {
			return nil, err
		}
		perms = append(perms, p)
	}

	r := &role.Role{
		TenantID:    req.TenantID,
		Name:        req.Name,
		Permissions: perms,
	}

	if err := s.repo.Save(ctx, r); err != nil {
		return nil, err
	}

	return r, nil
}

func (s *roleService) UpdateRole(ctx context.Context, req UpdateRoleRequest) (*role.Role, error) {
	r, err := s.repo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	perms := make([]permission.Permission, 0, len(req.Permissions))
	for _, pStr := range req.Permissions {
		p, err := permission.Parse(pStr)
		if err != nil {
			return nil, err
		}
		perms = append(perms, p)
	}

	r.Name = req.Name
	r.Permissions = perms

	if err := s.repo.Save(ctx, r); err != nil {
		return nil, err
	}

	return r, nil
}

func (s *roleService) DeleteRole(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *roleService) GetRole(ctx context.Context, id string) (*role.Role, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *roleService) ListRoles(ctx context.Context, tenantID string) ([]role.Role, error) {
	return s.repo.GetByTenantID(ctx, tenantID)
}
