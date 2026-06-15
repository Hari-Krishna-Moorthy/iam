package strategies_test

import (
	"context"

	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/application/auth/strategies"
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/role"
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/user"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"golang.org/x/crypto/bcrypt"
)

type mockUserRepo struct {
	getByUsernameFunc func(tenantID, username string) (*user.User, error)
}
func (m *mockUserRepo) GetByID(ctx context.Context, id string) (*user.User, error) { return nil, nil }
func (m *mockUserRepo) GetByEmail(ctx context.Context, tid, email string) (*user.User, error) { return nil, nil }
func (m *mockUserRepo) GetByUsername(ctx context.Context, tid, username string) (*user.User, error) { return m.getByUsernameFunc(tid, username) }
func (m *mockUserRepo) Save(ctx context.Context, u *user.User) error { return nil }

type mockRoleRepo struct {
	getByIDFunc func(id string) (*role.Role, error)
}
func (m *mockRoleRepo) GetByID(ctx context.Context, id string) (*role.Role, error) { return m.getByIDFunc(id) }
func (m *mockRoleRepo) GetByTenantID(ctx context.Context, tid string) ([]role.Role, error) { return nil, nil }
func (m *mockRoleRepo) Save(ctx context.Context, r *role.Role) error { return nil }

var _ = Describe("PasswordStrategy", func() {
	var (
		userRepo   *mockUserRepo
		roleRepo   *mockRoleRepo
		ctx        context.Context
	)

	BeforeEach(func() {
		userRepo = &mockUserRepo{}
		roleRepo = &mockRoleRepo{}
		ctx = context.WithValue(context.Background(), "tenant_id", "t1")
	})

	It("should authenticate with valid credentials", func() {
		hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		userRepo.getByUsernameFunc = func(tid, username string) (*user.User, error) {
			return &user.User{ID: "u1", TenantID: "t1", PasswordHash: string(hash), RoleID: "r1"}, nil
		}
		roleRepo.getByIDFunc = func(id string) (*role.Role, error) {
			return &role.Role{ID: "r1", Name: "admin"}, nil
		}

		s := strategies.NewPasswordStrategy(userRepo, roleRepo)
		sess, err := s.Authenticate(ctx, map[string]string{"username": "user1", "password": "password123"})

		Expect(err).NotTo(HaveOccurred())
		Expect(sess.UserID).To(Equal("u1"))
		Expect(sess.Role).To(Equal("admin"))
	})

	It("should fail with invalid password", func() {
		hash, _ := bcrypt.GenerateFromPassword([]byte("correct"), bcrypt.DefaultCost)
		userRepo.getByUsernameFunc = func(tid, username string) (*user.User, error) {
			return &user.User{PasswordHash: string(hash)}, nil
		}

		s := strategies.NewPasswordStrategy(userRepo, roleRepo)
		_, err := s.Authenticate(ctx, map[string]string{"username": "user1", "password": "wrong"})

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal("invalid credentials"))
	})
})
