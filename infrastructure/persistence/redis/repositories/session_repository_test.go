package repositories_test

import (
	"context"
	"time"

	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/permission"
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/session"
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/infrastructure/persistence/redis/repositories"
	"github.com/alicebob/miniredis/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/redis/go-redis/v9"
)

var _ = Describe("SessionRepository", func() {
	var (
		mr     *miniredis.Miniredis
		client *redis.Client
		repo   session.Repository
		ctx    context.Context
	)

	BeforeEach(func() {
		var err error
		mr, err = miniredis.Run()
		Expect(err).NotTo(HaveOccurred())

		client = redis.NewClient(&redis.Options{
			Addr: mr.Addr(),
		})

		repo = repositories.NewSessionRepository(client)
		ctx = context.Background()
	})

	AfterEach(func() {
		mr.Close()
	})

	Context("Save and Get", func() {
		It("should save a session and retrieve it", func() {
			s := &session.Session{
				ID:       "sess-123",
				UserID:   "user-456",
				TenantID: "tenant-789",
				Role:     "admin",
				Permissions: []permission.Permission{
					{Scope: "global", ServiceName: "billing", Action: "read"},
				},
				CreatedAt: time.Now(),
				ExpiresAt: time.Now().Add(time.Hour),
			}

			err := repo.Save(ctx, s)
			Expect(err).NotTo(HaveOccurred())

			retrieved, err := repo.GetByID(ctx, "sess-123")
			Expect(err).NotTo(HaveOccurred())
			Expect(retrieved.UserID).To(Equal("user-456"))
			Expect(retrieved.Role).To(Equal("admin"))
			Expect(retrieved.Permissions).To(HaveLen(1))
			Expect(retrieved.Permissions[0].String()).To(Equal("global:billing:read"))
		})

		It("should add session ID to user's session set", func() {
			s := &session.Session{
				ID:        "sess-1",
				UserID:    "user-1",
				ExpiresAt: time.Now().Add(time.Hour),
			}
			repo.Save(ctx, s)

			sessions, err := repo.GetByUserID(ctx, "user-1")
			Expect(err).NotTo(HaveOccurred())
			Expect(sessions).To(HaveLen(1))
			Expect(sessions[0].ID).To(Equal("sess-1"))
		})
	})
})
