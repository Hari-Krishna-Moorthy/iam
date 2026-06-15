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
	Authenticate(ctx context.Context, domain string, strategy string, credentials map[string]string) (string, error)
}

type authService struct {
	tenantRepo     tenant.Repository
	sessionRepo    session.Repository
	tokenProvider  TokenProvider
	authStrategies map[string]session.AuthStrategy
}

func NewService(
	tenantRepo tenant.Repository,
	sessionRepo session.Repository,
	tokenProvider TokenProvider,
	strategies map[string]session.AuthStrategy,
) Service {
	return &authService{
		tenantRepo:     tenantRepo,
		sessionRepo:    sessionRepo,
		tokenProvider:  tokenProvider,
		authStrategies: strategies,
	}
}

func (s *authService) Authenticate(ctx context.Context, domain string, strategyName string, credentials map[string]string) (string, error) {
	// 1. Identify Tenant by Origin (Feature A)
	t, err := s.tenantRepo.GetByDomain(ctx, domain)
	if err != nil {
		return "", ErrTenantNotFound
	}

	// 2. Select Strategy (Feature B)
	strategy, ok := s.authStrategies[strategyName]
	if !ok {
		return "", ErrInvalidStrategy
	}

	// 3. Authenticate
	// Inject TenantID into context for the strategy to use
	ctx = context.WithValue(ctx, "tenant_id", t.ID)
	sess, err := s.AuthenticateWithStrategy(ctx, strategy, credentials)
	if err != nil {
		return "", err
	}

	// 4. Save Session (Feature C)
	if err := s.sessionRepo.Save(ctx, sess); err != nil {
		return "", fmt.Errorf("failed to save session: %w", err)
	}

	// 5. Generate Token (JWT)
	token, err := s.tokenProvider.GenerateToken(ctx, sess)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}

// AuthenticateWithStrategy is a helper to isolate strategy call for testing
func (s *authService) AuthenticateWithStrategy(ctx context.Context, strategy session.AuthStrategy, credentials map[string]string) (*session.Session, error) {
	return strategy.Authenticate(ctx, credentials)
}
