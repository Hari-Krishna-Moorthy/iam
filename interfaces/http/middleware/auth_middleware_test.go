package middleware_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/permission"
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/session"
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/tenant"
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/interfaces/http/middleware"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type mockSessionRepo struct {
	getFunc func(id string) (*session.Session, error)
}
func (m *mockSessionRepo) Save(ctx context.Context, s *session.Session) error { return nil }
func (m *mockSessionRepo) GetByID(ctx context.Context, id string) (*session.Session, error) { return m.getFunc(id) }
func (m *mockSessionRepo) Delete(ctx context.Context, id string) error { return nil }
func (m *mockSessionRepo) GetByUserID(ctx context.Context, id string) ([]*session.Session, error) { return nil, nil }

type mockAuthTenantRepo struct {
	getFunc func(id string) (*tenant.Tenant, error)
}
func (m *mockAuthTenantRepo) GetByID(ctx context.Context, id string) (*tenant.Tenant, error) {
	if m.getFunc != nil {
		return m.getFunc(id)
	}
	return nil, nil
}
func (m *mockAuthTenantRepo) GetByDomain(ctx context.Context, domain string) (*tenant.Tenant, error) { return nil, nil }
func (m *mockAuthTenantRepo) GetAll(ctx context.Context) ([]tenant.Tenant, error) { return nil, nil }
func (m *mockAuthTenantRepo) Save(ctx context.Context, t *tenant.Tenant) error { return nil }

var _ = Describe("AuthMiddleware", func() {
	var (
		repo       *mockSessionRepo
		tenantRepo *mockAuthTenantRepo
		next       http.Handler
		writer     *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		repo = &mockSessionRepo{}
		tenantRepo = &mockAuthTenantRepo{}
		writer = httptest.NewRecorder()
		next = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
	})

	It("should hydrate headers when session is valid", func() {
		repo.getFunc = func(id string) (*session.Session, error) {
			return &session.Session{
				UserID:   "u-1",
				TenantID: "t-1",
				Role:     "admin",
				Permissions: []permission.Permission{
					{Scope: "g", ServiceName: "s", Action: "r"},
				},
			}, nil
		}

		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", "Bearer sess-123")

		handler := middleware.AuthMiddleware(repo, tenantRepo)(next)
		handler.ServeHTTP(writer, req)

		Expect(writer.Code).To(Equal(http.StatusOK))
		Expect(req.Header.Get("X-User-ID")).To(Equal("u-1"))
		Expect(req.Header.Get("X-Tenant-ID")).To(Equal("t-1"))
		Expect(req.Header.Get("X-Role")).To(Equal("admin"))
		Expect(req.Header.Get("X-Permissions")).To(Equal("g:s:r"))
	})

	It("should switch context for Super Admin", func() {
		repo.getFunc = func(id string) (*session.Session, error) {
			return &session.Session{
				UserID:   "u-1",
				TenantID: "system-tenant",
				Role:     "super-admin",
			}, nil
		}

		tenantRepo.getFunc = func(id string) (*tenant.Tenant, error) {
			if id == "system-tenant" {
				return &tenant.Tenant{ID: "system-tenant", IsSystem: true}, nil
			}
			if id == "target-tenant" {
				return &tenant.Tenant{ID: "target-tenant"}, nil
			}
			return nil, errors.New("not found")
		}

		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", "Bearer sess-123")
		req.Header.Set("X-Target-Tenant-ID", "target-tenant")

		handler := middleware.AuthMiddleware(repo, tenantRepo)(next)
		handler.ServeHTTP(writer, req)

		Expect(writer.Code).To(Equal(http.StatusOK))
		Expect(req.Header.Get("X-Tenant-ID")).To(Equal("target-tenant"))
	})

	It("should return 401 when Authorization header is missing", func() {
		req := httptest.NewRequest("GET", "/", nil)
		handler := middleware.AuthMiddleware(repo, tenantRepo)(next)
		handler.ServeHTTP(writer, req)

		Expect(writer.Code).To(Equal(http.StatusUnauthorized))
	})

	It("should return 401 when session is invalid", func() {
		repo.getFunc = func(id string) (*session.Session, error) {
			return nil, errors.New("invalid")
		}

		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", "Bearer bad-sess")

		handler := middleware.AuthMiddleware(repo, tenantRepo)(next)
		handler.ServeHTTP(writer, req)

		Expect(writer.Code).To(Equal(http.StatusUnauthorized))
	})
})

