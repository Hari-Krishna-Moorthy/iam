package tenant

import (
	"context"

	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/tenant"
)

type RegistrationRequest struct {
	Name    string
	Domains []string
}

type Service interface {
	RegisterTenant(ctx context.Context, req RegistrationRequest) (*tenant.Tenant, error)
}

type tenantService struct {
	repo tenant.Repository
}

func NewService(repo tenant.Repository) Service {
	return &tenantService{repo: repo}
}

func (s *tenantService) RegisterTenant(ctx context.Context, req RegistrationRequest) (*tenant.Tenant, error) {
	t := &tenant.Tenant{
		Name:     req.Name,
		Domains:  req.Domains,
		IsActive: true,
	}

	if err := s.repo.Save(ctx, t); err != nil {
		return nil, err
	}

	return t, nil
}
