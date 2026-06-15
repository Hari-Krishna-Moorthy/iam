package auth_test

import (
	"context"
	"errors"

	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/application/auth"
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/session"
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/tenant"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Mock objects for testing
type mockTenantRepo struct {
	getFunc func(domain string) (*tenant.Tenant, error)
}
func (m *mockTenantRepo) GetByID(ctx context.Context, id string) (*tenant.Tenant, error) { return nil, nil }
func (m *mockTenantRepo) GetByDomain(ctx context.Context, domain string) (*tenant.Tenant, error) { return m.getFunc(domain) }
func (m *mockTenantRepo) Save(ctx context.Context, t *tenant.Tenant) error { return nil }

type mockSessionRepo struct {
	saveFunc func(s *session.Session) error
}
func (m *mockSessionRepo) Save(ctx context.Context, s *session.Session) error { return m.saveFunc(s) }
func (m *mockSessionRepo) GetByID(ctx context.Context, id string) (*session.Session, error) { return nil, nil }
func (m *mockSessionRepo) Delete(ctx context.Context, id string) error { return nil }
func (m *mockSessionRepo) GetByUserID(ctx context.Context, id string) ([]*session.Session, error) { return nil, nil }

type mockAuthStrategy struct {
	authFunc func(creds map[string]string) (*session.Session, error)
}
func (m *mockAuthStrategy) Authenticate(ctx context.Context, creds map[string]string) (*session.Session, error) { return m.authFunc(creds) }
func (m *mockAuthStrategy) ValidateToken(ctx context.Context, token string) (*session.Session, error) { return nil, nil }

var _ = Describe("AuthService", func() {
	var (
		service     auth.Service
		tenantRepo  *mockTenantRepo
		sessionRepo *mockSessionRepo
		strategy    *mockAuthStrategy
		ctx         context.Context
	)

	BeforeEach(func() {
		tenantRepo = &mockTenantRepo{}
		sessionRepo = &mockSessionRepo{}
		strategy = &mockAuthStrategy{}
		ctx = context.Background()

		strategies := map[string]session.AuthStrategy{
			"password": strategy,
		}
		service = auth.NewService(tenantRepo, sessionRepo, strategies)
	})

	Context("Authenticate", func() {
		It("should succeed when everything is valid", func() {
			tenantRepo.getFunc = func(domain string) (*tenant.Tenant, error) {
				return &tenant.Tenant{ID: "t1"}, nil
			}
			strategy.authFunc = func(creds map[string]string) (*session.Session, error) {
				return &session.Session{ID: "s1", UserID: "u1"}, nil
			}
			sessionRepo.saveFunc = func(s *session.Session) error {
				return nil
			}

			sess, err := service.Authenticate(ctx, "example.com", "password", nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(sess.ID).To(Equal("s1"))
		})

		It("should fail when tenant not found", func() {
			tenantRepo.getFunc = func(domain string) (*tenant.Tenant, error) {
				return nil, errors.New("not found")
			}

			_, err := service.Authenticate(ctx, "wrong.com", "password", nil)
			Expect(err).To(Equal(auth.ErrTenantNotFound))
		})
	})
})
