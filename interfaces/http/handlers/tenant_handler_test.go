package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/application/tenant"
	domainTenant "github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/tenant"
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/interfaces/http/handlers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type mockTenantService struct {
	registerFunc func(req tenant.RegistrationRequest) (*domainTenant.Tenant, error)
}
func (m *mockTenantService) RegisterTenant(ctx context.Context, req tenant.RegistrationRequest) (*domainTenant.Tenant, error) {
	return m.registerFunc(req)
}

var _ = Describe("TenantHandler", func() {
	var (
		service *mockTenantService
		handler *handlers.TenantHandler
		writer  *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		service = &mockTenantService{}
		handler = handlers.NewTenantHandler(service)
		writer = httptest.NewRecorder()
	})

	Context("RegisterTenant", func() {
		It("should return 201 on success", func() {
			service.registerFunc = func(req tenant.RegistrationRequest) (*domainTenant.Tenant, error) {
				return &domainTenant.Tenant{ID: "t-123", Name: "Acme"}, nil
			}

			body, _ := json.Marshal(tenant.RegistrationRequest{Name: "Acme"})
			req := httptest.NewRequest("POST", "/tenants", bytes.NewBuffer(body))

			handler.RegisterTenant(writer, req)

			Expect(writer.Code).To(Equal(http.StatusCreated))
			var resp map[string]interface{}
			json.Unmarshal(writer.Body.Bytes(), &resp)
			Expect(resp["ID"]).To(Equal("t-123"))
		})

		It("should return 500 on service error", func() {
			service.registerFunc = func(req tenant.RegistrationRequest) (*domainTenant.Tenant, error) {
				return nil, errors.New("internal error")
			}

			body, _ := json.Marshal(tenant.RegistrationRequest{Name: "Acme"})
			req := httptest.NewRequest("POST", "/tenants", bytes.NewBuffer(body))

			handler.RegisterTenant(writer, req)

			Expect(writer.Code).To(Equal(http.StatusInternalServerError))
		})
	})
})
