package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/session"
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/interfaces/http/handlers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type mockAuthService struct {
	authFunc func(domain, strategy string, creds map[string]string) (*session.Session, error)
}
func (m *mockAuthService) Authenticate(ctx context.Context, domain, strategy string, creds map[string]string) (*session.Session, error) {
	return m.authFunc(domain, strategy, creds)
}

var _ = Describe("AuthHandler", func() {
	var (
		service *mockAuthService
		handler *handlers.AuthHandler
		writer  *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		service = &mockAuthService{}
		handler = handlers.NewAuthHandler(service)
		writer = httptest.NewRecorder()
	})

	Context("Login", func() {
		It("should return token on successful login", func() {
			service.authFunc = func(domain, strategy string, creds map[string]string) (*session.Session, error) {
				return &session.Session{ID: "token-123"}, nil
			}

			body, _ := json.Marshal(map[string]interface{}{
				"strategy": "password",
				"credentials": map[string]string{"username": "u1", "password": "p1"},
			})
			req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))

			handler.Login(writer, req)

			Expect(writer.Code).To(Equal(http.StatusOK))
			var resp map[string]string
			json.Unmarshal(writer.Body.Bytes(), &resp)
			Expect(resp["token"]).To(Equal("token-123"))
		})

		It("should return 401 on failed auth", func() {
			service.authFunc = func(domain, strategy string, creds map[string]string) (*session.Session, error) {
				return nil, errors.New("unauthorized")
			}

			body, _ := json.Marshal(map[string]interface{}{
				"strategy": "password",
			})
			req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))

			handler.Login(writer, req)

			Expect(writer.Code).To(Equal(http.StatusUnauthorized))
		})
	})
})
