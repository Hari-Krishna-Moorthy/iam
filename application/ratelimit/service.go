package ratelimit

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/ratelimit"
	"github.com/redis/go-redis/v9"
)

type redisLimiter struct {
	client *redis.Client
	repo   ratelimit.Repository
}

func NewRedisLimiter(client *redis.Client, repo ratelimit.Repository) ratelimit.Limiter {
	return &redisLimiter{client: client, repo: repo}
}

func (l *redisLimiter) Allow(ctx context.Context, tenantID string) (bool, error) {
	config, err := l.repo.GetByTenantID(ctx, tenantID)
	if err != nil {
		return false, err
	}

	switch config.Algorithm {
	case ratelimit.FixedWindow:
		return l.allowFixedWindow(ctx, tenantID, config)
	case ratelimit.SlidingWindow:
		return l.allowSlidingWindow(ctx, tenantID, config)
	case ratelimit.LeakyBucket:
		return l.allowLeakyBucket(ctx, tenantID, config)
	case ratelimit.TokenBucket:
		return l.allowTokenBucket(ctx, tenantID, config)
	default:
		return l.allowFixedWindow(ctx, tenantID, config)
	}
}

func (l *redisLimiter) allowFixedWindow(ctx context.Context, tenantID string, c *ratelimit.Config) (bool, error) {
	key := fmt.Sprintf("ratelimit:fixed:%s:%d", tenantID, time.Now().UnixNano()/int64(c.Window))
	
	pipe := l.client.Pipeline()
	incr := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, c.Window)
	
	_, err := pipe.Exec(ctx)
	if err != nil {
		return false, err
	}

	return incr.Val() <= int64(c.Limit), nil
}

func (l *redisLimiter) allowSlidingWindow(ctx context.Context, tenantID string, c *ratelimit.Config) (bool, error) {
	key := fmt.Sprintf("ratelimit:sliding:%s", tenantID)
	now := time.Now().UnixNano()
	windowStart := now - int64(c.Window)

	pipe := l.client.Pipeline()
	pipe.ZRemRangeByScore(ctx, key, "0", strconv.FormatInt(windowStart, 10))
	pipe.ZAdd(ctx, key, redis.Z{Score: float64(now), Member: now})
	count := pipe.ZCard(ctx, key)
	pipe.Expire(ctx, key, c.Window)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return false, err
	}

	return count.Val() <= int64(c.Limit), nil
}

func (l *redisLimiter) allowLeakyBucket(ctx context.Context, tenantID string, c *ratelimit.Config) (bool, error) {
	key := fmt.Sprintf("ratelimit:leaky:%s", tenantID)
	now := time.Now().Unix()

	// Using a simple Redis Lua script for atomic Leaky Bucket
	script := `
		local key = KEYS[1]
		local capacity = tonumber(ARGV[1])
		local rate = tonumber(ARGV[2])
		local now = tonumber(ARGV[3])

		local bucket = redis.call("HMGET", key, "last_update", "water")
		local last_update = tonumber(bucket[1]) or now
		local water = tonumber(bucket[2]) or 0

		-- Leak water
		local elapsed = now - last_update
		water = math.max(0, water - (elapsed * rate))
		
		if water < capacity then
			water = water + 1
			redis.call("HMSET", key, "last_update", now, "water", water)
			return 1
		else
			return 0
		end
	`
	res, err := l.client.Eval(ctx, script, []string{key}, c.Burst, c.Rate, now).Int()
	if err != nil {
		return false, err
	}
	return res == 1, nil
}

func (l *redisLimiter) allowTokenBucket(ctx context.Context, tenantID string, c *ratelimit.Config) (bool, error) {
	key := fmt.Sprintf("ratelimit:token:%s", tenantID)
	now := time.Now().Unix()

	script := `
		local key = KEYS[1]
		local burst = tonumber(ARGV[1])
		local rate = tonumber(ARGV[2])
		local now = tonumber(ARGV[3])

		local bucket = redis.call("HMGET", key, "last_update", "tokens")
		local last_update = tonumber(bucket[1]) or now
		local tokens = tonumber(bucket[2]) or burst

		-- Add tokens
		local elapsed = now - last_update
		tokens = math.min(burst, tokens + (elapsed * rate))

		if tokens >= 1 then
			tokens = tokens - 1
			redis.call("HMSET", key, "last_update", now, "tokens", tokens)
			return 1
		else
			return 0
		end
	`
	res, err := l.client.Eval(ctx, script, []string{key}, c.Burst, c.Rate, now).Int()
	if err != nil {
		return false, err
	}
	return res == 1, nil
}
