package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/application/user"
	domainUser "github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/user"
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/interfaces/http/handlers"
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/interfaces/http/middleware"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type mockUserService struct {
	registerFunc func(req user.RegistrationRequest) (*domainUser.User, error)
}
func (m *mockUserService) RegisterUser(ctx context.Context, req user.RegistrationRequest) (*domainUser.User, error) {
	return m.registerFunc(req)
}

var _ = Describe("UserHandler", func() {
	var (
		service *mockUserService
		handler *handlers.UserHandler
		writer  *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		service = &mockUserService{}
		handler = handlers.NewUserHandler(service)
		writer = httptest.NewRecorder()
	})

	Context("RegisterUser", func() {
		It("should return 201 on success", func() {
			service.registerFunc = func(req user.RegistrationRequest) (*domainUser.User, error) {
				return &domainUser.User{ID: "u-123", Username: "john", PasswordHash: "secret"}, nil
			}

			body, _ := json.Marshal(user.RegistrationRequest{Username: "john", Password: "P@ssword1"})
			req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(body))
			req = req.WithContext(context.WithValue(req.Context(), middleware.TenantIDKey, "t1"))

			handler.RegisterUser(writer, req)

			Expect(writer.Code).To(Equal(http.StatusCreated))
			var resp map[string]interface{}
			json.Unmarshal(writer.Body.Bytes(), &resp)
			Expect(resp["ID"]).To(Equal("u-123"))
			Expect(resp["PasswordHash"]).To(BeEmpty()) // Ensure hash is removed
		})

		It("should return 401 if tenant context is missing", func() {
			body, _ := json.Marshal(user.RegistrationRequest{Username: "john"})
			req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(body))
			// Context has no tenant ID

			handler.RegisterUser(writer, req)

			Expect(writer.Code).To(Equal(http.StatusUnauthorized))
		})

		It("should return 400 on service error", func() {
			service.registerFunc = func(req user.RegistrationRequest) (*domainUser.User, error) {
				return nil, errors.New("policy violation")
			}

			body, _ := json.Marshal(user.RegistrationRequest{Username: "john"})
			req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(body))
			req = req.WithContext(context.WithValue(req.Context(), middleware.TenantIDKey, "t1"))

			handler.RegisterUser(writer, req)

			Expect(writer.Code).To(Equal(http.StatusBadRequest))
		})
	})
})
