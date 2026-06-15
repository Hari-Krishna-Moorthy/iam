package strategies

import (
	"context"
	"errors"

	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/role"
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/session"
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/user"
	"golang.org/x/crypto/bcrypt"
)

type passwordStrategy struct {
	userRepo user.Repository
	roleRepo role.Repository
}

func NewPasswordStrategy(userRepo user.Repository, roleRepo role.Repository) session.AuthStrategy {
	return &passwordStrategy{
		userRepo: userRepo,
		roleRepo: roleRepo,
	}
}

func (s *passwordStrategy) Authenticate(ctx context.Context, credentials map[string]string) (*session.Session, error) {
	tenantID, ok := ctx.Value("tenant_id").(string)
	if !ok {
		return nil, errors.New("tenant_id not found in context")
	}

	username := credentials["username"]
	password := credentials["password"]

	u, err := s.userRepo.GetByUsername(ctx, tenantID, username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	r, err := s.roleRepo.GetByID(ctx, u.RoleID)
	if err != nil {
		return nil, errors.New("user role not found")
	}

	return &session.Session{
		ID:          "sess_" + u.ID, // Simplified session ID generation for now
		UserID:      u.ID,
		TenantID:    u.TenantID,
		Role:        r.Name,
		Permissions: r.Permissions,
	}, nil
}

func (s *passwordStrategy) ValidateToken(ctx context.Context, token string) (*session.Session, error) {
	// For basic password auth, we might just rely on the session repository.
	// This could be expanded for JWT validation if this strategy handled JWTs.
	return nil, errors.New("not implemented")
}
