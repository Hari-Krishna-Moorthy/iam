package strategies

import (
	"context"
	"errors"

	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/role"
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/session"
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/user"
)

// OAuth2Provider defines the interface for external identity providers (Google, Okta, etc.)
type OAuth2Provider interface {
	// VerifyToken takes an external access_token/id_token and returns the user's email/identifier
	VerifyToken(ctx context.Context, token string) (string, error)
}

type oauth2Strategy struct {
	userRepo user.Repository
	roleRepo role.Repository
	provider OAuth2Provider
}

func NewOAuth2Strategy(userRepo user.Repository, roleRepo role.Repository, provider OAuth2Provider) session.AuthStrategy {
	return &oauth2Strategy{
		userRepo: userRepo,
		roleRepo: roleRepo,
		provider: provider,
	}
}

func (s *oauth2Strategy) Authenticate(ctx context.Context, credentials map[string]string) (*session.Session, error) {
	tenantID, ok := ctx.Value("tenant_id").(string)
	if !ok {
		return nil, errors.New("tenant_id not found in context")
	}

	token, ok := credentials["token"]
	if !ok || token == "" {
		return nil, errors.New("missing oauth2 token")
	}

	// 1. Verify token with external provider
	email, err := s.provider.VerifyToken(ctx, token)
	if err != nil {
		return nil, errors.New("invalid oauth2 token")
	}

	// 2. Find user by email within the tenant
	u, err := s.userRepo.GetByEmail(ctx, tenantID, email)
	if err != nil {
		// In a full implementation, we might auto-provision (JIT provision) the user here.
		// For now, we reject if they don't exist.
		return nil, errors.New("user not found for this tenant")
	}

	// 3. Get user's role
	r, err := s.roleRepo.GetByID(ctx, u.RoleID)
	if err != nil {
		return nil, errors.New("user role not found")
	}

	// 4. Return internal session
	return &session.Session{
		ID:          "sess_oauth2_" + u.ID, // Simplified
		UserID:      u.ID,
		TenantID:    u.TenantID,
		Role:        r.Name,
		Permissions: r.Permissions,
	}, nil
}

func (s *oauth2Strategy) ValidateToken(ctx context.Context, token string) (*session.Session, error) {
	return nil, errors.New("not implemented")
}
