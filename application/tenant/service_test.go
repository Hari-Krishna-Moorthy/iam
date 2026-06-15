package tenant_test

import (
	"context"
	"errors"

	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/application/tenant"
	domainTenant "github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/tenant"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type mockTenantRepo struct {
	saveFunc func(t *domainTenant.Tenant) error
}

func (m *mockTenantRepo) GetByID(ctx context.Context, id string) (*domainTenant.Tenant, error)     { return nil, nil }
func (m *mockTenantRepo) GetByDomain(ctx context.Context, domain string) (*domainTenant.Tenant, error) { return nil, nil }
func (m *mockTenantRepo) Save(ctx context.Context, t *domainTenant.Tenant) error                  { return m.saveFunc(t) }

var _ = Describe("TenantService", func() {
	var (
		service  tenant.Service
		repo     *mockTenantRepo
		ctx      context.Context
	)

	BeforeEach(func() {
		repo = &mockTenantRepo{}
		service = tenant.NewService(repo)
		ctx = context.Background()
	})

	Context("RegisterTenant", func() {
		It("should successfully register a tenant", func() {
			req := tenant.RegistrationRequest{
				Name:    "New Tenant",
				Domains: []string{"tenant.com"},
			}

			repo.saveFunc = func(t *domainTenant.Tenant) error {
				return nil
			}

			t, err := service.RegisterTenant(ctx, req)

			Expect(err).NotTo(HaveOccurred())
			Expect(t.Name).To(Equal("New Tenant"))
			Expect(t.Domains).To(ContainElement("tenant.com"))
		})

		It("should return an error if repository fails", func() {
			req := tenant.RegistrationRequest{Name: "Fail"}
			repo.saveFunc = func(t *domainTenant.Tenant) error {
				return errors.New("save failed")
			}

			_, err := service.RegisterTenant(ctx, req)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("save failed"))
		})
	})
})
