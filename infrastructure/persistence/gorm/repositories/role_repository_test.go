package repositories_test

import (
	"context"
	"regexp"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/role"
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/infrastructure/persistence/gorm/repositories"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var _ = Describe("RoleRepository", func() {
	var (
		db     *gorm.DB
		mock   sqlmock.Sqlmock
		repo   role.Repository
		ctx    context.Context
	)

	BeforeEach(func() {
		sqlDB, m, err := sqlmock.New()
		Expect(err).NotTo(HaveOccurred())
		mock = m

		db, err = gorm.Open(postgres.New(postgres.Config{
			Conn: sqlDB,
		}), &gorm.Config{})
		Expect(err).NotTo(HaveOccurred())

		repo = repositories.NewRoleRepository(db)
		ctx = context.Background()
	})

	Context("GetByID", func() {
		It("should return a role when found", func() {
			rid := "role-uuid"
			rows := sqlmock.NewRows([]string{"id", "tenant_id", "name", "permissions"}).
				AddRow(rid, "tenant-1", "admin", "{global:all:read}")

			mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "roles" WHERE id = $1 AND "roles"."deleted_at" IS NULL ORDER BY "roles"."id" LIMIT $2`)).
				WithArgs(rid, 1).
				WillReturnRows(rows)

			r, err := repo.GetByID(ctx, rid)
			Expect(err).NotTo(HaveOccurred())
			Expect(r.ID).To(Equal(rid))
			Expect(r.Permissions).To(HaveLen(1))
		})
	})

	Context("Delete", func() {
		It("should delete the role", func() {
			rid := "role-uuid"
			mock.ExpectBegin()
			mock.ExpectExec(regexp.QuoteMeta(`UPDATE "roles" SET "deleted_at"=$1 WHERE id = $2 AND "roles"."deleted_at" IS NULL`)).
				WithArgs(sqlmock.AnyArg(), rid).
				WillReturnResult(sqlmock.NewResult(0, 1))
			mock.ExpectCommit()

			err := repo.Delete(ctx, rid)
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
