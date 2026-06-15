package strategies_test

import (
	"context"
	"errors"

	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/application/auth/strategies"
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/role"
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/user"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type mockOAuth2Provider struct {
	verifyFunc func(token string) (string, error)
}
func (m *mockOAuth2Provider) VerifyToken(ctx context.Context, token string) (string, error) {
	return m.verifyFunc(token)
}

var _ = Describe("OAuth2Strategy", func() {
	var (
		userRepo *mockUserRepo
		roleRepo *mockRoleRepo
		provider *mockOAuth2Provider
		ctx      context.Context
	)

	BeforeEach(func() {
		userRepo = &mockUserRepo{}
		roleRepo = &mockRoleRepo{}
		provider = &mockOAuth2Provider{}
		ctx = context.WithValue(context.Background(), "tenant_id", "t1")
	})

	It("should authenticate with a valid external token", func() {
		provider.verifyFunc = func(token string) (string, error) {
			return "user@example.com", nil
		}
		userRepo.getByEmailFunc = func(tid, email string) (*user.User, error) {
			return &user.User{ID: "u1", TenantID: "t1", Email: email, RoleID: "r1"}, nil
		}
		roleRepo.getByIDFunc = func(id string) (*role.Role, error) {
			return &role.Role{ID: "r1", Name: "user"}, nil
		}

		s := strategies.NewOAuth2Strategy(userRepo, roleRepo, provider)
		sess, err := s.Authenticate(ctx, map[string]string{"token": "valid_oauth_token"})

		Expect(err).NotTo(HaveOccurred())
		Expect(sess.UserID).To(Equal("u1"))
		Expect(sess.Role).To(Equal("user"))
	})

	It("should fail if the external token is invalid", func() {
		provider.verifyFunc = func(token string) (string, error) {
			return "", errors.New("invalid token")
		}

		s := strategies.NewOAuth2Strategy(userRepo, roleRepo, provider)
		_, err := s.Authenticate(ctx, map[string]string{"token": "bad_token"})

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal("invalid oauth2 token"))
	})

	It("should fail if user is not found in the tenant", func() {
		provider.verifyFunc = func(token string) (string, error) {
			return "user@example.com", nil
		}
		userRepo.getByEmailFunc = func(tid, email string) (*user.User, error) {
			return nil, errors.New("not found")
		}

		s := strategies.NewOAuth2Strategy(userRepo, roleRepo, provider)
		_, err := s.Authenticate(ctx, map[string]string{"token": "valid_oauth_token"})

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal("user not found for this tenant"))
	})
})
