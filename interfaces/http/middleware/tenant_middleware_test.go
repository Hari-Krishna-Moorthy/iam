package middleware_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/tenant"
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/interfaces/http/middleware"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type mockTenantRepo struct {
	getFunc func(domain string) (*tenant.Tenant, error)
}
func (m *mockTenantRepo) GetByID(ctx context.Context, id string) (*tenant.Tenant, error) { return nil, nil }
func (m *mockTenantRepo) GetByDomain(ctx context.Context, domain string) (*tenant.Tenant, error) { return m.getFunc(domain) }
func (m *mockTenantRepo) GetAll(ctx context.Context) ([]tenant.Tenant, error) { return nil, nil }
func (m *mockTenantRepo) Save(ctx context.Context, t *tenant.Tenant) error { return nil }

var _ = Describe("TenantMiddleware", func() {
	var (
		repo   *mockTenantRepo
		next   http.Handler
		writer *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		repo = &mockTenantRepo{}
		writer = httptest.NewRecorder()
		next = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tenantID := r.Context().Value(middleware.TenantIDKey)
			w.Write([]byte(tenantID.(string)))
		})
	})

	It("should inject TenantID into context when origin is valid", func() {
		repo.getFunc = func(domain string) (*tenant.Tenant, error) {
			return &tenant.Tenant{ID: "t-123"}, nil
		}

		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Origin", "example.com")

		handler := middleware.TenantMiddleware(repo)(next)
		handler.ServeHTTP(writer, req)

		Expect(writer.Code).To(Equal(http.StatusOK))
		Expect(writer.Body.String()).To(Equal("t-123"))
	})

	It("should return 404 when tenant is not found", func() {
		repo.getFunc = func(domain string) (*tenant.Tenant, error) {
			return nil, errors.New("not found")
		}

		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Origin", "unknown.com")

		handler := middleware.TenantMiddleware(repo)(next)
		handler.ServeHTTP(writer, req)

		Expect(writer.Code).To(Equal(http.StatusNotFound))
	})
})
