# redis-ratelimiter

A high-performance, Redis-based rate limiting library for Go applications. This library provides a simple interface for implementing rate limiting with configurable rates, time periods, and support for composite rate limits.

## Features

- 🚀 **Simple API**: Clean and intuitive interface for rate limiting
- ⚡ **High Performance**: Built on top of Redis with efficient algorithms
- 🔄 **Flexible Configuration**: Configurable rates, time periods, and contexts
- 🔀 **Composite Limiters**: Apply multiple rate limits simultaneously (e.g., per-second and per-minute)
- 🛡️ **Thread-Safe**: Safe for concurrent use in goroutines
- 🧪 **Well-Tested**: Comprehensive test suite with race condition detection
- 📦 **Production-Ready**: Used in production environments

## Installation

```bash
go get github.com/linhbkhn95/redis-ratelimiter
```

## Prerequisites

- Go 1.25.1 or higher
- Redis server (version 6.0+ recommended)

## Quick Start

### Basic Usage

```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/linhbkhn95/redis-ratelimiter"
    "github.com/redis/go-redis/v9"
)

func main() {
    // Create Redis client
    rdb := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })

    // Create a rate limiter: 10 requests per second
    limiter := ratelimiter.New(
        rdb,
        "user:123",  // Key for this rate limit
        10,          // Rate: 10 requests
        ratelimiter.Per(time.Second), // Per 1 second
    )

    // Use the limiter
    if _, err := limiter.Take(); err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    // Your code here - request is allowed
    fmt.Println("Request allowed!")
}
```

### Custom Time Period

```go
// 100 requests per minute
limiter := ratelimiter.New(
    rdb,
    "api:endpoint",
    100,
    ratelimiter.Per(time.Minute),
)

// 1000 requests per hour
limiter := ratelimiter.New(
    rdb,
    "api:endpoint",
    1000,
    ratelimiter.Per(time.Hour),
)
```

### With Context for Cancellation

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

limiter := ratelimiter.New(
    rdb,
    "user:123",
    10,
    ratelimiter.Per(time.Second),
    ratelimiter.WithContext(ctx),
)

_, err := limiter.Take()
if err != nil {
    // Handle error (could be context cancellation)
}
```

### Composite Rate Limiting

Apply multiple rate limits simultaneously. All limits must pass for the request to be allowed:

```go
// Create individual limiters
perSecLimiter := ratelimiter.New(
    rdb,
    "api:per_second",
    10,
    ratelimiter.Per(time.Second),
)

perMinLimiter := ratelimiter.New(
    rdb,
    "api:per_minute",
    100,
    ratelimiter.Per(time.Minute),
)

// Create composite limiter
composite := ratelimiter.NewComposite(perSecLimiter, perMinLimiter)

// Take() will check both limits
// Request is only allowed if both limits permit it
_, err := composite.Take()
if err != nil {
    // Handle error
}

// You can also add limiters dynamically
perHourLimiter := ratelimiter.New(
    rdb,
    "api:per_hour",
    1000,
    ratelimiter.Per(time.Hour),
)
composite.AddLimiter(perHourLimiter)
```

## API Reference

### Types

#### `Limiter` Interface

```go
type Limiter interface {
    Take() (time.Time, error)
}
```

- **Take()**: Blocks until the request is allowed under the rate limit, then returns the time when the request was allowed and any error that occurred.

#### `CompositeLimiter` Interface

```go
type CompositeLimiter interface {
    Limiter
    AddLimiter(limiter Limiter)
}
```

- **Take()**: Applies all rate limits and only passes if all limiters allow the request.
- **AddLimiter(limiter Limiter)**: Adds another limiter that must also pass.

### Functions

#### `New(rdb redis.UniversalClient, key string, rate int, opts ...Option) Limiter`

Creates a new rate limiter instance.

**Parameters:**

- `rdb`: Redis client (supports `*redis.Client`, `*redis.ClusterClient`, etc.)
- `key`: Unique key for this rate limit (e.g., user ID, API endpoint)
- `rate`: Maximum number of requests allowed in the time period
- `opts`: Optional configuration:
  - `Per(duration)`: Time period for the rate limit (default: 1 second)
  - `WithContext(ctx)`: Context for cancellation support

**Returns:** A `Limiter` instance

#### `NewComposite(limiters ...Limiter) CompositeLimiter`

Creates a composite limiter that applies multiple rate limits.

**Parameters:**

- `limiters`: One or more limiter instances to combine

**Returns:** A `CompositeLimiter` instance

### Options

#### `Per(d time.Duration) Option`

Sets the time period for the rate limit. Default is `time.Second`.

```go
ratelimiter.Per(time.Second)  // Per second
ratelimiter.Per(time.Minute)  // Per minute
ratelimiter.Per(time.Hour)    // Per hour
ratelimiter.Per(5 * time.Minute) // Custom period
```

#### `WithContext(ctx context.Context) Option`

Sets a context for cancellation support. The limiter will respect context cancellation.

```go
ratelimiter.WithContext(ctx)
```

## Examples

### HTTP Middleware

```go
func rateLimitMiddleware(limiter ratelimiter.Limiter) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        _, err := limiter.Take()
        if err != nil {
            http.Error(w, "Rate limit error", http.StatusInternalServerError)
            return
        }
        // Continue with request
    }
}

// Usage
limiter := ratelimiter.New(rdb, "api:endpoint", 100, ratelimiter.Per(time.Minute))
http.HandleFunc("/api", rateLimitMiddleware(limiter))
```

### Per-User Rate Limiting

```go
func getUserLimiter(userID string) ratelimiter.Limiter {
    return ratelimiter.New(
        rdb,
        fmt.Sprintf("user:%s", userID),
        50,
        ratelimiter.Per(time.Minute),
    )
}

// In your handler
limiter := getUserLimiter("12345")
_, err := limiter.Take()
if err != nil {
    return errors.New("rate limit exceeded")
}
```

### Multiple Tier Rate Limiting

```go
// Free tier: 10 req/sec, 100 req/min
freeTierLimiter := ratelimiter.NewComposite(
    ratelimiter.New(rdb, "free:per_sec", 10, ratelimiter.Per(time.Second)),
    ratelimiter.New(rdb, "free:per_min", 100, ratelimiter.Per(time.Minute)),
)

// Premium tier: 100 req/sec, 1000 req/min
premiumTierLimiter := ratelimiter.NewComposite(
    ratelimiter.New(rdb, "premium:per_sec", 100, ratelimiter.Per(time.Second)),
    ratelimiter.New(rdb, "premium:per_min", 1000, ratelimiter.Per(time.Minute)),
)
```

## Testing

Run the test suite:

```bash
# Run all tests
go test ./...

# Run tests with race detection
go test -race ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...

# Run benchmarks
go test -bench=. -benchmem ./...
```

Or use the Makefile:

```bash
make test        # Run tests with race detection and coverage
make lint        # Run linters
make benchmark   # Run benchmarks
make coverage    # Generate coverage report
make all         # Run all checks
```

## Performance

The library is optimized for high-throughput scenarios. Benchmarks are included in the test suite. Run them with:

```bash
go test -bench=. -benchmem ./...
```

### Benchmark Results

Benchmarks were run on an Apple M4 Pro (ARM64) with Go 1.25.1:

#### Basic Limiter Benchmarks

| Benchmark                            | Operations | Time/Op         | Memory/Op  | Allocs/Op    |
| ------------------------------------ | ---------- | --------------- | ---------- | ------------ |
| `BenchmarkLimiter_Take_WithinRate`   | 7,725      | 150,124 ns/op   | 696 B/op   | 23 allocs/op |
| `BenchmarkLimiter_Take_Concurrent`   | 41,911     | 76,203 ns/op    | 2,106 B/op | 58 allocs/op |
| `BenchmarkLimiter_Take_WithBlocking` | 6,324      | 9,842,197 ns/op | 1,592 B/op | 42 allocs/op |

#### Composite Limiter Benchmarks

| Benchmark                                      | Operations | Time/Op       | Memory/Op  | Allocs/Op    |
| ---------------------------------------------- | ---------- | ------------- | ---------- | ------------ |
| `BenchmarkCompositeLimiter_Take_SingleLimiter` | 8,149      | 141,063 ns/op | 704 B/op   | 23 allocs/op |
| `BenchmarkCompositeLimiter_Take_TwoLimiters`   | 3,958      | 273,438 ns/op | 1,408 B/op | 46 allocs/op |
| `BenchmarkCompositeLimiter_Take_ThreeLimiters` | 2,754      | 416,108 ns/op | 2,160 B/op | 70 allocs/op |
| `BenchmarkCompositeLimiter_Take_Concurrent`    | 21,792     | 58,908 ns/op  | 1,483 B/op | 46 allocs/op |
| `BenchmarkCompositeLimiter_AddLimiter`         | 12,688,269 | 104.2 ns/op   | 185 B/op   | 3 allocs/op  |

**Performance Notes:**

- Concurrent operations show better throughput (76k ns/op vs 150k ns/op for sequential)
- Composite limiters scale linearly with the number of limiters
- Memory allocation is kept minimal for high-throughput scenarios
- Blocking operations (when rate limit is exceeded) show expected higher latency

## Redis Requirements

- Redis 6.0+ recommended
- Uses Redis commands: `INCR`, `EXPIRE`, and Lua scripts for atomic operations
- The rate limiter uses Redis keys with the format you provide (e.g., `"user:123"`)

## Error Handling

The `Take()` method may return errors in the following scenarios:

- Redis connection errors
- Context cancellation (when using `WithContext`)
- Internal server errors (rare)

The limiter follows a "fail-open" policy: if Redis is unavailable or an error occurs, the limiter will attempt to allow the request rather than blocking all traffic.

## Thread Safety

All limiter implementations are thread-safe and can be safely used from multiple goroutines concurrently.

## Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines on how to contribute to this project.

## License

See [LICENSE](LICENSE) file for details.

## Author

Created and maintained by [linhbkhn95](https://github.com/linhbkhn95)

## Acknowledgments

Built on top of the excellent [redis_rate](https://github.com/go-redis/redis_rate) library and [go-redis](https://github.com/redis/go-redis).
