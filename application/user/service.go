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
}

type userService struct {
	userRepo user.Repository
}

func NewService(userRepo user.Repository) Service {
	return &userService{userRepo: userRepo}
}

func (s *userService) RegisterUser(ctx context.Context, req RegistrationRequest) (*user.User, error) {
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
