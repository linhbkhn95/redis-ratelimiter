package ratelimiter

import "time"

type Limiter interface {
	Take() (time.Time, error)
}

// CompositeLimiter applies multiple rate limits and only passes if all limits pass
type CompositeLimiter interface {
	Limiter
	// AddLimiter adds another limiter that must also pass
	AddLimiter(limiter Limiter)
}
