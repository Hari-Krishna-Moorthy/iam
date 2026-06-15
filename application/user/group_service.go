package user

import (
	"context"

	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/user"
)

type CreateGroupRequest struct {
	TenantID string
	Name     string
}

type GroupService interface {
	CreateGroup(ctx context.Context, req CreateGroupRequest) (*user.Group, error)
	DeleteGroup(ctx context.Context, id string) error
	AddUserToGroup(ctx context.Context, groupID, userID string) error
	RemoveUserFromGroup(ctx context.Context, groupID, userID string) error
	AddRoleToGroup(ctx context.Context, groupID, roleID string) error
	RemoveRoleFromGroup(ctx context.Context, groupID, roleID string) error
	ListGroups(ctx context.Context, tenantID string) ([]user.Group, error)
}

type groupService struct {
	repo user.GroupRepository
}

func NewGroupService(repo user.GroupRepository) GroupService {
	return &groupService{repo: repo}
}

func (s *groupService) CreateGroup(ctx context.Context, req CreateGroupRequest) (*user.Group, error) {
	g := &user.Group{
		TenantID: req.TenantID,
		Name:     req.Name,
	}
	if err := s.repo.Save(ctx, g); err != nil {
		return nil, err
	}
	return g, nil
}

func (s *groupService) DeleteGroup(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *groupService) AddUserToGroup(ctx context.Context, groupID, userID string) error {
	return s.repo.AddUser(ctx, groupID, userID)
}

func (s *groupService) RemoveUserFromGroup(ctx context.Context, groupID, userID string) error {
	return s.repo.RemoveUser(ctx, groupID, userID)
}

func (s *groupService) AddRoleToGroup(ctx context.Context, groupID, roleID string) error {
	return s.repo.AddRole(ctx, groupID, roleID)
}

func (s *groupService) RemoveRoleFromGroup(ctx context.Context, groupID, roleID string) error {
	return s.repo.RemoveRole(ctx, groupID, roleID)
}

func (s *groupService) ListGroups(ctx context.Context, tenantID string) ([]user.Group, error) {
	return s.repo.GetByTenantID(ctx, tenantID)
}
