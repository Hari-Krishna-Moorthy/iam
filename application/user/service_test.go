package user_test

import (
	"context"
	"errors"

	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/application/user"
	domainUser "github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/user"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type mockUserRepo struct {
	saveFunc func(u *domainUser.User) error
}

func (m *mockUserRepo) GetByID(ctx context.Context, id string) (*domainUser.User, error)             { return nil, nil }
func (m *mockUserRepo) GetByEmail(ctx context.Context, tid, email string) (*domainUser.User, error)  { return nil, nil }
func (m *mockUserRepo) GetByUsername(ctx context.Context, tid, uname string) (*domainUser.User, error) { return nil, nil }
func (m *mockUserRepo) Save(ctx context.Context, u *domainUser.User) error                         { return m.saveFunc(u) }

var _ = Describe("UserService", func() {
	var (
		service  user.Service
		repo     *mockUserRepo
		ctx      context.Context
	)

	BeforeEach(func() {
		repo = &mockUserRepo{}
		service = user.NewService(repo)
		ctx = context.Background()
	})

	Context("RegisterUser", func() {
		It("should successfully register a user with hashed password", func() {
			req := user.RegistrationRequest{
				TenantID: "t1",
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
				RoleID:   "r1",
			}

			var savedUser *domainUser.User
			repo.saveFunc = func(u *domainUser.User) error {
				savedUser = u
				return nil
			}

			u, err := service.RegisterUser(ctx, req)

			Expect(err).NotTo(HaveOccurred())
			Expect(u.Username).To(Equal("testuser"))
			Expect(u.PasswordHash).NotTo(Equal("password123")) // Should be hashed
			Expect(savedUser).To(Equal(u))
		})

		It("should return an error if repository fails to save", func() {
			req := user.RegistrationRequest{
				TenantID: "t1",
				Username: "testuser",
				Password: "password123",
			}

			repo.saveFunc = func(u *domainUser.User) error {
				return errors.New("db error")
			}

			_, err := service.RegisterUser(ctx, req)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("db error"))
		})
	})
})
