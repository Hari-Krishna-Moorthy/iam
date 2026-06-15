package user_test

import (
	"context"

	applicationUser "github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/application/user"
	domainUser "github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/user"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type mockGroupRepo struct {
	saveFunc func(g *domainUser.Group) error
}

func (m *mockGroupRepo) GetByID(ctx context.Context, id string) (*domainUser.Group, error)       { return nil, nil }
func (m *mockGroupRepo) GetByTenantID(ctx context.Context, tid string) ([]domainUser.Group, error) { return nil, nil }
func (m *mockGroupRepo) Save(ctx context.Context, g *domainUser.Group) error                    { return m.saveFunc(g) }
func (m *mockGroupRepo) Delete(ctx context.Context, id string) error                           { return nil }
func (m *mockGroupRepo) AddUser(ctx context.Context, gid, uid string) error                    { return nil }
func (m *mockGroupRepo) RemoveUser(ctx context.Context, gid, uid string) error                 { return nil }
func (m *mockGroupRepo) AddRole(ctx context.Context, gid, rid string) error                    { return nil }
func (m *mockGroupRepo) RemoveRole(ctx context.Context, gid, rid string) error                 { return nil }

var _ = Describe("GroupService", func() {
	var (
		service applicationUser.GroupService
		repo    *mockGroupRepo
		ctx     context.Context
	)

	BeforeEach(func() {
		repo = &mockGroupRepo{}
		service = applicationUser.NewGroupService(repo)
		ctx = context.Background()
	})

	Context("CreateGroup", func() {
		It("should successfully create a group", func() {
			req := applicationUser.CreateGroupRequest{
				TenantID: "t1",
				Name:     "engineers",
			}
			repo.saveFunc = func(g *domainUser.Group) error {
				return nil
			}

			g, err := service.CreateGroup(ctx, req)
			Expect(err).NotTo(HaveOccurred())
			Expect(g.Name).To(Equal("engineers"))
		})
	})
})
