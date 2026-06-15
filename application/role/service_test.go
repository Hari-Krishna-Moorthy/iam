package role_test

import (
	"context"

	applicationRole "github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/application/role"
	domainRole "github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/role"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type mockRoleRepo struct {
	getByIDFunc       func(id string) (*domainRole.Role, error)
	getByTenantIDFunc func(tid string) ([]domainRole.Role, error)
	saveFunc          func(r *domainRole.Role) error
	deleteFunc        func(id string) error
}

func (m *mockRoleRepo) GetByID(ctx context.Context, id string) (*domainRole.Role, error)     { return m.getByIDFunc(id) }
func (m *mockRoleRepo) GetByTenantID(ctx context.Context, tid string) ([]domainRole.Role, error) { return m.getByTenantIDFunc(tid) }
func (m *mockRoleRepo) Save(ctx context.Context, r *domainRole.Role) error                  { return m.saveFunc(r) }
func (m *mockRoleRepo) Delete(ctx context.Context, id string) error                         { return m.deleteFunc(id) }

var _ = Describe("RoleService", func() {
	var (
		service applicationRole.Service
		repo    *mockRoleRepo
		ctx     context.Context
	)

	BeforeEach(func() {
		repo = &mockRoleRepo{}
		service = applicationRole.NewService(repo)
		ctx = context.Background()
	})

	Context("CreateRole", func() {
		It("should successfully create a role", func() {
			req := applicationRole.CreateRoleRequest{
				TenantID:    "t1",
				Name:        "admin",
				Permissions: []string{"g:s:r"},
			}

			repo.saveFunc = func(r *domainRole.Role) error {
				return nil
			}

			r, err := service.CreateRole(ctx, req)
			Expect(err).NotTo(HaveOccurred())
			Expect(r.Name).To(Equal("admin"))
			Expect(r.Permissions).To(HaveLen(1))
			Expect(r.Permissions[0].String()).To(Equal("g:s:r"))
		})

		It("should fail if permission format is invalid", func() {
			req := applicationRole.CreateRoleRequest{
				Permissions: []string{"invalid"},
			}

			_, err := service.CreateRole(ctx, req)
			Expect(err).To(HaveOccurred())
		})
	})

	Context("UpdateRole", func() {
		It("should successfully update a role", func() {
			repo.getByIDFunc = func(id string) (*domainRole.Role, error) {
				return &domainRole.Role{ID: id, Name: "old"}, nil
			}
			repo.saveFunc = func(r *domainRole.Role) error {
				return nil
			}

			req := applicationRole.UpdateRoleRequest{
				ID:          "r1",
				Name:        "new",
				Permissions: []string{"g:s:w"},
			}

			r, err := service.UpdateRole(ctx, req)
			Expect(err).NotTo(HaveOccurred())
			Expect(r.Name).To(Equal("new"))
			Expect(r.Permissions[0].String()).To(Equal("g:s:w"))
		})
	})

	Context("DeleteRole", func() {
		It("should call delete on repository", func() {
			repo.deleteFunc = func(id string) error {
				return nil
			}

			err := service.DeleteRole(ctx, "r1")
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
