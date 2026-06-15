package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/session"
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/tenant"
)

var (
	ErrTenantNotFound = errors.New("tenant not found")
	ErrInvalidStrategy = errors.New("invalid auth strategy")
)

// Service coordinates authentication use cases.
type Service interface {
	Authenticate(ctx context.Context, domain string, strategy string, credentials map[string]string) (*session.Session, error)
}

type authService struct {
	tenantRepo     tenant.Repository
	sessionRepo    session.Repository
	authStrategies map[string]session.AuthStrategy
}

func NewService(
	tenantRepo tenant.Repository,
	sessionRepo session.Repository,
	strategies map[string]session.AuthStrategy,
) Service {
	return &authService{
		tenantRepo:     tenantRepo,
		sessionRepo:    sessionRepo,
		authStrategies: strategies,
	}
}

func (s *authService) Authenticate(ctx context.Context, domain string, strategyName string, credentials map[string]string) (*session.Session, error) {
	// 1. Identify Tenant by Origin (Feature A)
	t, err := s.tenantRepo.GetByDomain(ctx, domain)
	if err != nil {
		return nil, ErrTenantNotFound
	}

	// 2. Select Strategy (Feature B)
	strategy, ok := s.authStrategies[strategyName]
	if !ok {
		return nil, ErrInvalidStrategy
	}

	// 3. Authenticate
	// Inject TenantID into context for the strategy to use
	ctx = context.WithValue(ctx, "tenant_id", t.ID)
	sess, err := s.AuthenticateWithStrategy(ctx, strategy, credentials)
	if err != nil {
		return nil, err
	}

	// 4. Save Session (Feature C)
	if err := s.sessionRepo.Save(ctx, sess); err != nil {
		return nil, fmt.Errorf("failed to save session: %w", err)
	}

	return sess, nil
}

// AuthenticateWithStrategy is a helper to isolate strategy call for testing
func (s *authService) AuthenticateWithStrategy(ctx context.Context, strategy session.AuthStrategy, credentials map[string]string) (*session.Session, error) {
	return strategy.Authenticate(ctx, credentials)
}
