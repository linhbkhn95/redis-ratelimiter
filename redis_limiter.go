package ratelimiter

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"
)

var (
	ErrIntervalServer = errors.New("rate limit interval server error")
)

type Option func(*config)

type config struct {
	per time.Duration
	ctx context.Context
}

func Per(d time.Duration) Option {
	return func(c *config) {
		c.per = d
	}
}

// WithContext allows cancellation (optional)
func WithContext(ctx context.Context) Option {
	return func(c *config) {
		c.ctx = ctx
	}
}

type redisLimiter struct {
	limiter *redis_rate.Limiter
	key     string
	limit   redis_rate.Limit
	ctx     context.Context
}

func New(
	rdb redis.UniversalClient,
	key string,
	rate int,
	opts ...Option,
) Limiter {
	cfg := &config{
		per: time.Second,
		ctx: context.Background(),
	}

	for _, opt := range opts {
		opt(cfg)
	}

	return &redisLimiter{
		limiter: redis_rate.NewLimiter(rdb),
		key:     key,
		limit:   redis_rate.Limit{Rate: rate, Burst: rate, Period: cfg.per},
		ctx:     cfg.ctx,
	}
}

func (l *redisLimiter) Take() (time.Time, error) {
	for {
		now := time.Now()

		res, err := l.limiter.Allow(l.ctx, l.key, l.limit)
		if err == nil && res.Allowed > 0 {
			return now, nil
		}

		if err != nil {
			return time.Now(), err
		}

		if res == nil {
			return time.Now(), ErrIntervalServer
		}

		// Rate limit exceeded, wait for RetryAfter
		if res.RetryAfter > 0 {
			select {
			case <-time.After(res.RetryAfter):
				continue
			case <-l.ctx.Done():
				return time.Now(), nil
			}
		}

		// Should not reach here, but fail open just in case
		return now, nil
	}
}

// compositeLimiter applies multiple rate limits and only passes if all limits pass
type compositeLimiter struct {
	limiters []Limiter
	mu       sync.RWMutex
}

// NewComposite creates a new composite limiter that checks multiple rate limits.
// All limiters must pass for Take() to succeed.
//
// Example - apply both per-second and per-minute limits:
//
//	perSecLimiter := New(rdb, "aggregate_per_second", 10, Per(time.Second))
//	perMinLimiter := New(rdb, "aggregate_per_per_min", 100, Per(time.Minute))
//	composite := NewComposite(perSecLimiter, perMinLimiter)
//	composite.Take() // will check both limits, only passes if both allow
func NewComposite(limiters ...Limiter) CompositeLimiter {
	return &compositeLimiter{
		limiters: limiters,
	}
}

// AddLimiter adds another limiter that must also pass
func (c *compositeLimiter) AddLimiter(limiter Limiter) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.limiters = append(c.limiters, limiter)
}

// Take applies all rate limits and only passes if all limiters allow the request.
// If any limiter blocks, Take() will block until that limiter allows, then re-check all limiters.
func (c *compositeLimiter) Take() (time.Time, error) {
	c.mu.RLock()
	limiters := make([]Limiter, len(c.limiters))
	copy(limiters, c.limiters)
	c.mu.RUnlock()

	// If no limiters, fail open
	if len(limiters) == 0 {
		return time.Now(), nil
	}

	// Apply all limiters - each Take() will block if needed
	// The last call time is returned, representing when the request was actually allowed
	var (
		lastTime time.Time
		err      error
	)
	for _, limiter := range limiters {
		lastTime, err = limiter.Take()
		if err != nil {
			// On error, fail open
			return time.Now(), err
		}
	}

	return lastTime, nil
}
