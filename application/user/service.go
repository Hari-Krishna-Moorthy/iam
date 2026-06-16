package user

import (
	"context"
	"fmt"

	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/user"
	"golang.org/x/crypto/bcrypt"
)

type RegistrationRequest struct {
	TenantID string
	Username string
	Email    string
	Password string
	RoleID   string
}

type Service interface {
	RegisterUser(ctx context.Context, req RegistrationRequest) (*user.User, error)
	ListUsers(ctx context.Context, tenantID string) ([]user.User, error)
}

type userService struct {
	userRepo   user.Repository
	policyRepo user.PasswordPolicyRepository
}

func NewService(userRepo user.Repository, policyRepo user.PasswordPolicyRepository) Service {
	return &userService{
		userRepo:   userRepo,
		policyRepo: policyRepo,
	}
}

func (s *userService) RegisterUser(ctx context.Context, req RegistrationRequest) (*user.User, error) {
	// 1. Get and Validate Policy
	policy, err := s.policyRepo.GetByTenantID(ctx, req.TenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get password policy: %w", err)
	}

	if err := policy.Validate(req.Password); err != nil {
		return nil, err
	}

	// 2. Hash Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	u := &user.User{
		TenantID:     req.TenantID,
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		RoleID:       req.RoleID,
		IsActive:     true,
	}

	if err := s.userRepo.Save(ctx, u); err != nil {
		return nil, err
	}

	return u, nil
}

func (s *userService) ListUsers(ctx context.Context, tenantID string) ([]user.User, error) {
	return s.userRepo.GetByTenantID(ctx, tenantID)
}
