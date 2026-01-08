# Contributing to redis-ratelimiter

First off, thank you for considering contributing to redis-ratelimiter! It's people like you that make this project better.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [How Can I Contribute?](#how-can-i-contribute)
  - [Reporting Bugs](#reporting-bugs)
  - [Suggesting Enhancements](#suggesting-enhancements)
  - [Pull Requests](#pull-requests)
- [Development Setup](#development-setup)
- [Project Structure](#project-structure)
- [Coding Standards](#coding-standards)
- [Testing](#testing)
- [Commit Guidelines](#commit-guidelines)
- [Pull Request Process](#pull-request-process)

## Code of Conduct

This project adheres to a code of conduct that all contributors are expected to follow. Please be respectful and constructive in all interactions.

## How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check the issue list to see if the bug has already been reported. When creating a bug report, include:

- **Clear title and description**
- **Steps to reproduce** the behavior
- **Expected behavior** vs **actual behavior**
- **Go version** and **Redis version**
- **Code snippets** that reproduce the issue (if applicable)
- **Error messages** (if any)
- **Environment details** (OS, etc.)

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion, include:

- **Clear title and description**
- **Use case**: Why is this feature useful?
- **Proposed solution** (if you have one)
- **Alternatives considered** (if any)

### Pull Requests

Pull requests are welcome! Please follow these guidelines:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add or update tests
5. Ensure all tests pass and linting checks pass
6. Commit your changes (follow the commit guidelines)
7. Push to your branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

## Development Setup

### Prerequisites

- Go 1.25.1 or higher
- Redis server (for running tests)
- Make (optional, for using Makefile commands)

### Setup Steps

1. **Fork and clone the repository:**
   ```bash
   git clone https://github.com/YOUR_USERNAME/redis-ratelimiter.git
   cd redis-ratelimiter
   ```

2. **Install dependencies:**
   ```bash
   go mod download
   ```

3. **Start Redis (required for tests):**
   ```bash
   # Using Docker
   docker run -d -p 6379:6379 redis:7-alpine

   # Or install Redis locally and start the service
   ```

4. **Run tests to verify setup:**
   ```bash
   go test ./...
   ```

### Installing Development Tools

For linting, install `golangci-lint`:

```bash
# macOS
brew install golangci-lint

# Or download from https://golangci-lint.run/usage/install/
```

## Project Structure

```
redis-ratelimiter/
├── .github/
│   └── workflows/
│       └── ci.yml          # CI/CD pipeline
├── interface.go            # Core interfaces (Limiter, CompositeLimiter)
├── redis_limiter.go        # Redis-based limiter implementation
├── redis_rate_test.go      # Test suite
├── go.mod                  # Go module definition
├── go.sum                  # Go module checksums
├── Makefile                # Common tasks
├── README.md               # Project documentation
├── CONTRIBUTING.md         # This file
└── LICENSE                 # License file
```

## Coding Standards

### Go Code Style

- Follow the standard Go formatting rules (`go fmt`)
- Use `gofmt -s` for simplified code
- Follow [Effective Go](https://go.dev/doc/effective_go) guidelines
- Use meaningful variable and function names
- Add comments for exported functions and types
- Keep functions focused and small

### Code Organization

- Keep related code together
- Export only what needs to be public
- Use interfaces for abstractions
- Avoid unnecessary dependencies

### Example Format

```go
// Good
// New creates a new rate limiter with the specified configuration.
// The limiter will allow 'rate' requests per time period.
func New(rdb redis.UniversalClient, key string, rate int, opts ...Option) Limiter {
    // implementation
}
```

## Testing

### Writing Tests

- Write tests for all new features
- Write tests for bug fixes
- Aim for high code coverage (80%+)
- Test edge cases and error conditions
- Use table-driven tests when appropriate

### Running Tests

```bash
# Run all tests
make test

# Run tests with race detection
go test -race ./...

# Run specific test
go test -run TestLimiter_AllowsWithinRate

# Run with verbose output
go test -v ./...

# Run benchmarks
make benchmark
```

### Test Coverage

```bash
# Generate coverage report
make coverage

# View coverage in browser
go tool cover -html=coverage.out
```

### Test Requirements

- All tests must pass
- Tests should be deterministic (no flaky tests)
- Tests should clean up after themselves (use `t.Cleanup()`)
- Redis tests should skip if Redis is unavailable

### Example Test

```go
func TestMyNewFeature(t *testing.T) {
    rdb := newTestRedis(t)
    
    limiter := New(rdb, "test:key", 10, Per(time.Second))
    
    // Test implementation
    _, err := limiter.Take()
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    
    // Verify behavior
}
```

## Commit Guidelines

### Commit Message Format

Follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Types

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

### Examples

```
feat: add support for sliding window rate limiting

fix(composite): handle nil limiters gracefully

docs: update README with usage examples

test: add tests for context cancellation
```

### Best Practices

- Write clear, descriptive commit messages
- Keep commits focused (one logical change per commit)
- Use present tense ("add feature" not "added feature")
- Reference issue numbers when applicable: `fix #123`

## Pull Request Process

### Before Submitting

1. **Update documentation** if you're adding features or changing behavior
2. **Add tests** for new functionality
3. **Ensure all tests pass**: `make test`
4. **Run linters**: `make lint`
5. **Check formatting**: `go fmt ./...`
6. **Verify go.mod**: `go mod tidy`

### PR Checklist

- [ ] Code follows the project's style guidelines
- [ ] Self-review completed
- [ ] Comments added for complex code
- [ ] Documentation updated
- [ ] Tests added/updated
- [ ] All tests pass locally
- [ ] Linting passes
- [ ] No new warnings or errors
- [ ] Commit messages follow guidelines

### PR Description Template

```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
How was this tested?

## Checklist
- [ ] Tests pass
- [ ] Documentation updated
- [ ] Code follows style guidelines
```

### Review Process

1. Maintainers will review your PR
2. Address any feedback or requested changes
3. Once approved, your PR will be merged
4. Thank you for contributing! 🎉

## CI/CD

The project uses GitHub Actions for CI/CD. The pipeline includes:

- **Test**: Runs tests with race detection
- **Lint**: Runs golangci-lint and format checks
- **Coverage**: Uploads coverage reports

All PRs must pass CI checks before merging.

## Questions?

If you have questions about contributing, feel free to:

- Open an issue with the `question` label
- Check existing issues and discussions
- Review the code and tests for examples

Thank you for contributing to redis-ratelimiter! 🚀
