package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/audit"
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/session"
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/tenant"
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/role"
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/user"
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/job"

	applicationRole "github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/application/role"
	applicationTenant "github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/application/tenant"
	applicationUser "github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/application/user"
	interfacesHttp "github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/interfaces/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Simplified mocks for E2E routing tests

type mockTenantRepo struct {
	getFunc func(domain string) (*tenant.Tenant, error)
}
func (m *mockTenantRepo) GetByID(ctx context.Context, id string) (*tenant.Tenant, error) { return nil, nil }
func (m *mockTenantRepo) GetByDomain(ctx context.Context, domain string) (*tenant.Tenant, error) { return m.getFunc(domain) }
func (m *mockTenantRepo) GetAll(ctx context.Context) ([]tenant.Tenant, error) { return nil, nil }
func (m *mockTenantRepo) Save(ctx context.Context, t *tenant.Tenant) error { return nil }

type mockSessionRepo struct {
	getFunc func(id string) (*session.Session, error)
}
func (m *mockSessionRepo) Save(ctx context.Context, s *session.Session) error { return nil }
func (m *mockSessionRepo) GetByID(ctx context.Context, id string) (*session.Session, error) { return m.getFunc(id) }
func (m *mockSessionRepo) Delete(ctx context.Context, id string) error { return nil }
func (m *mockSessionRepo) GetByUserID(ctx context.Context, id string) ([]*session.Session, error) { return nil, nil }

type mockAuditRepo struct{}
func (m *mockAuditRepo) Save(ctx context.Context, l *audit.AuditLog) error { return nil }
func (m *mockAuditRepo) GetByTenantID(ctx context.Context, tid string) ([]*audit.AuditLog, error) { return nil, nil }

type mockLimiter struct{}
func (m *mockLimiter) Allow(ctx context.Context, tenantID string) (bool, error) { return true, nil }

type mockAuthService struct {
	authFunc func(creds map[string]string) (string, error)
}
func (m *mockAuthService) Authenticate(ctx context.Context, domain, strategy string, creds map[string]string) (string, error) {
	return m.authFunc(creds)
}

// Dummy service stubs
type dummyRoleService struct{}
func (d *dummyRoleService) CreateRole(ctx context.Context, req applicationRole.CreateRoleRequest) (*role.Role, error) { return nil, nil }
func (d *dummyRoleService) UpdateRole(ctx context.Context, req applicationRole.UpdateRoleRequest) (*role.Role, error) { return nil, nil }
func (d *dummyRoleService) DeleteRole(ctx context.Context, id string) error { return nil }
func (d *dummyRoleService) GetRole(ctx context.Context, id string) (*role.Role, error) { return nil, nil }
func (d *dummyRoleService) ListRoles(ctx context.Context, tenantID string) ([]role.Role, error) { return nil, nil }

type dummyGroupService struct{}
func (d *dummyGroupService) CreateGroup(ctx context.Context, req applicationUser.CreateGroupRequest) (*user.Group, error) { return nil, nil }
func (d *dummyGroupService) DeleteGroup(ctx context.Context, id string) error { return nil }
func (d *dummyGroupService) AddUserToGroup(ctx context.Context, gid, uid string) error { return nil }
func (d *dummyGroupService) RemoveUserFromGroup(ctx context.Context, gid, uid string) error { return nil }
func (d *dummyGroupService) AddRoleToGroup(ctx context.Context, gid, rid string) error { return nil }
func (d *dummyGroupService) RemoveRoleFromGroup(ctx context.Context, gid, rid string) error { return nil }
func (d *dummyGroupService) ListGroups(ctx context.Context, tenantID string) ([]user.Group, error) { return nil, nil }

type dummyTenantService struct{}
func (d *dummyTenantService) RegisterTenant(ctx context.Context, req applicationTenant.RegistrationRequest) (*tenant.Tenant, error) { return nil, nil }
func (d *dummyTenantService) ListTenants(ctx context.Context) ([]tenant.Tenant, error) { return nil, nil }

type dummyUserService struct{}
func (d *dummyUserService) RegisterUser(ctx context.Context, req applicationUser.RegistrationRequest) (*user.User, error) { return nil, nil }
func (d *dummyUserService) ListUsers(ctx context.Context, tenantID string) ([]user.User, error) { return nil, nil }

type dummyBulkService struct{}
func (d *dummyBulkService) SubmitBulkCreate(ctx context.Context, tenantID string, req applicationUser.BulkCreateUsersRequest) (string, error) { return "", nil }
func (d *dummyBulkService) GetJobStatus(ctx context.Context, jobID string) (*job.Job, error) { return nil, nil }

var _ = Describe("Router E2E Tests", func() {
	var (
		ts          *httptest.Server
		authSvc     *mockAuthService
		tenantRepo  *mockTenantRepo
		sessionRepo *mockSessionRepo
	)

	BeforeEach(func() {
		tenantRepo = &mockTenantRepo{}
		sessionRepo = &mockSessionRepo{}
		authSvc = &mockAuthService{}

		router := interfacesHttp.NewRouter(
			tenantRepo,
			sessionRepo,
			&mockAuditRepo{},
			&mockLimiter{},
			authSvc,
			&dummyRoleService{},
			&dummyGroupService{},
			&dummyTenantService{},
			&dummyUserService{},
			&dummyBulkService{},
		)

		ts = httptest.NewServer(router)
	})

	AfterEach(func() {
		ts.Close()
	})

	Context("Login Flow", func() {
		It("should successfully hit the login endpoint and return a token", func() {
			authSvc.authFunc = func(creds map[string]string) (string, error) {
				return "jwt.token.here", nil
			}

			reqBody, _ := json.Marshal(map[string]interface{}{
				"strategy":    "password",
				"credentials": map[string]string{"username": "admin", "password": "password"},
			})
			req, _ := http.NewRequest("POST", ts.URL+"/login", bytes.NewBuffer(reqBody))

			client := &http.Client{}
			resp, err := client.Do(req)

			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			var body map[string]string
			json.NewDecoder(resp.Body).Decode(&body)
			Expect(body["token"]).To(Equal("jwt.token.here"))
		})
	})

	Context("Protected Routes", func() {
		It("should inject hydrated headers on /me endpoint", func() {
			// 1. Mock Tenant identification by Origin
			tenantRepo.getFunc = func(domain string) (*tenant.Tenant, error) {
				return &tenant.Tenant{ID: "t-123", Domains: []string{"test.com"}}, nil
			}

			// 2. Mock Session validation
			sessionRepo.getFunc = func(id string) (*session.Session, error) {
				return &session.Session{
					ID:       "sess-1",
					UserID:   "u-456",
					TenantID: "t-123",
					Role:     "admin",
				}, nil
			}

			req, _ := http.NewRequest("GET", ts.URL+"/me", nil)
			req.Header.Set("Origin", "test.com")
			req.Header.Set("Authorization", "Bearer valid-token")

			client := &http.Client{}
			resp, err := client.Do(req)

			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			// Router /me endpoint writes "Hello, user <userID>" based on hydrated X-User-ID header
			buf := new(bytes.Buffer)
			buf.ReadFrom(resp.Body)
			Expect(buf.String()).To(Equal("Hello, user u-456"))
		})

		It("should block requests with missing authorization", func() {
			tenantRepo.getFunc = func(domain string) (*tenant.Tenant, error) {
				return &tenant.Tenant{ID: "t-123"}, nil
			}

			req, _ := http.NewRequest("GET", ts.URL+"/me", nil)
			req.Header.Set("Origin", "test.com")

			client := &http.Client{}
			resp, err := client.Do(req)

			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized))
		})
	})
})
