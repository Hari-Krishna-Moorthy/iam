package middleware_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/permission"
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/session"
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

var _ = Describe("AuthMiddleware", func() {
	var (
		repo   *mockSessionRepo
		next   http.Handler
		writer *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		repo = &mockSessionRepo{}
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

		handler := middleware.AuthMiddleware(repo)(next)
		handler.ServeHTTP(writer, req)

		Expect(writer.Code).To(Equal(http.StatusOK))
		Expect(req.Header.Get("X-User-ID")).To(Equal("u-1"))
		Expect(req.Header.Get("X-Tenant-ID")).To(Equal("t-1"))
		Expect(req.Header.Get("X-Role")).To(Equal("admin"))
		Expect(req.Header.Get("X-Permissions")).To(Equal("g:s:r"))
	})

	It("should return 401 when Authorization header is missing", func() {
		req := httptest.NewRequest("GET", "/", nil)
		handler := middleware.AuthMiddleware(repo)(next)
		handler.ServeHTTP(writer, req)

		Expect(writer.Code).To(Equal(http.StatusUnauthorized))
	})

	It("should return 401 when session is invalid", func() {
		repo.getFunc = func(id string) (*session.Session, error) {
			return nil, errors.New("invalid")
		}

		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", "Bearer bad-sess")

		handler := middleware.AuthMiddleware(repo)(next)
		handler.ServeHTTP(writer, req)

		Expect(writer.Code).To(Equal(http.StatusUnauthorized))
	})
})
