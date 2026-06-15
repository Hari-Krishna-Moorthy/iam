package repositories_test

import (
	"context"
	"regexp"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/user"
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/infrastructure/persistence/gorm/repositories"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var _ = Describe("UserRepository", func() {
	var (
		db     *gorm.DB
		mock   sqlmock.Sqlmock
		repo   user.Repository
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

		repo = repositories.NewUserRepository(db)
		ctx = context.Background()
	})

	Context("GetByUsername", func() {
		It("should return a user when found", func() {
			tid := "tenant-1"
			username := "testuser"
			rows := sqlmock.NewRows([]string{"id", "tenant_id", "username", "email", "password_hash"}).
				AddRow("user-uuid", tid, username, "test@test.com", "hash")

			mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE (tenant_id = $1 AND username = $2) AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $3`)).
				WithArgs(tid, username, 1).
				WillReturnRows(rows)

			u, err := repo.GetByUsername(ctx, tid, username)
			Expect(err).NotTo(HaveOccurred())
			Expect(u.Username).To(Equal(username))
			Expect(u.TenantID).To(Equal(tid))
		})
	})
})
