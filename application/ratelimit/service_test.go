package ratelimit_test

import (
	"context"
	"time"

	domainRateLimit "github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/ratelimit"
	applicationRateLimit "github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/application/ratelimit"
	"github.com/alicebob/miniredis/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/redis/go-redis/v9"
)

type mockRateLimitRepo struct {
	getFunc func(tid string) (*domainRateLimit.Config, error)
}
func (m *mockRateLimitRepo) GetByTenantID(ctx context.Context, tid string) (*domainRateLimit.Config, error) {
	return m.getFunc(tid)
}
func (m *mockRateLimitRepo) Save(ctx context.Context, c *domainRateLimit.Config) error { return nil }

var _ = Describe("RedisLimiter", func() {
	var (
		mr      *miniredis.Miniredis
		client  *redis.Client
		repo    *mockRateLimitRepo
		limiter domainRateLimit.Limiter
		ctx     context.Context
	)

	BeforeEach(func() {
		var err error
		mr, err = miniredis.Run()
		Expect(err).NotTo(HaveOccurred())

		client = redis.NewClient(&redis.Options{Addr: mr.Addr()})
		repo = &mockRateLimitRepo{}
		limiter = applicationRateLimit.NewRedisLimiter(client, repo)
		ctx = context.Background()
	})

	AfterEach(func() {
		mr.Close()
	})

	Describe("Fixed Window", func() {
		It("should allow requests up to the limit", func() {
			repo.getFunc = func(tid string) (*domainRateLimit.Config, error) {
				return &domainRateLimit.Config{
					Algorithm: domainRateLimit.FixedWindow,
					Limit:     2,
					Window:    time.Minute,
				}, nil
			}

			Expect(limiter.Allow(ctx, "t1")).To(BeTrue())
			Expect(limiter.Allow(ctx, "t1")).To(BeTrue())
			Expect(limiter.Allow(ctx, "t1")).To(BeFalse())
		})
	})

	Describe("Token Bucket", func() {
		It("should allow bursts up to capacity", func() {
			repo.getFunc = func(tid string) (*domainRateLimit.Config, error) {
				return &domainRateLimit.Config{
					Algorithm: domainRateLimit.TokenBucket,
					Burst:     3,
					Rate:      1.0, // 1 token per second
				}, nil
			}

			Expect(limiter.Allow(ctx, "t1")).To(BeTrue())
			Expect(limiter.Allow(ctx, "t1")).To(BeTrue())
			Expect(limiter.Allow(ctx, "t1")).To(BeTrue())
			Expect(limiter.Allow(ctx, "t1")).To(BeFalse())
		})
	})
})
