package repositories_test

import (
	"context"
	"database/sql"
	"regexp"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/tenant"
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/infrastructure/persistence/gorm/repositories"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var _ = Describe("TenantRepository", func() {
	var (
		db     *gorm.DB
		mock   sqlmock.Sqlmock
		repo   tenant.Repository
		ctx    context.Context
		sqlDB  *sql.DB
	)

	BeforeEach(func() {
		var err error
		sqlDB, mock, err = sqlmock.New()
		Expect(err).NotTo(HaveOccurred())

		db, err = gorm.Open(postgres.New(postgres.Config{
			Conn: sqlDB,
		}), &gorm.Config{})
		Expect(err).NotTo(HaveOccurred())

		repo = repositories.NewTenantRepository(db)
		ctx = context.Background()
	})

	AfterEach(func() {
		sqlDB.Close()
	})

	Context("GetByDomain", func() {
		It("should return a tenant when domain exists", func() {
			domain := "app.example.com"
			rows := sqlmock.NewRows([]string{"id", "name", "domains", "is_active"}).
				AddRow("tenant-uuid", "Example Tenant", "{app.example.com}", true)

			mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tenants" WHERE $1 = ANY(domains) AND "tenants"."deleted_at" IS NULL ORDER BY "tenants"."id" LIMIT $2`)).
				WithArgs(domain, 1).
				WillReturnRows(rows)

			res, err := repo.GetByDomain(ctx, domain)
			Expect(err).NotTo(HaveOccurred())
			Expect(res.ID).To(Equal("tenant-uuid"))
			Expect(res.Domains).To(ContainElement(domain))
		})

		It("should return error when tenant not found", func() {
			domain := "unknown.com"
			mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tenants" WHERE $1 = ANY(domains) AND "tenants"."deleted_at" IS NULL ORDER BY "tenants"."id" LIMIT $2`)).
				WithArgs(domain, 1).
				WillReturnError(gorm.ErrRecordNotFound)

			_, err := repo.GetByDomain(ctx, domain)
			Expect(err).To(Equal(gorm.ErrRecordNotFound))
		})
	})
})
