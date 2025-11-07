# Contributing to Go Load Balancer

First off, thank you for considering contributing to Go Load Balancer! It's people like you that make this project better.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [How Can I Contribute?](#how-can-i-contribute)
- [Development Setup](#development-setup)
- [Coding Guidelines](#coding-guidelines)
- [Commit Messages](#commit-messages)
- [Pull Request Process](#pull-request-process)

## Code of Conduct

This project and everyone participating in it is governed by a Code of Conduct. By participating, you are expected to uphold this code.

- Be respectful and inclusive
- Welcome newcomers
- Focus on what is best for the community
- Show empathy towards other community members

## How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check the existing issues to avoid duplicates. When creating a bug report, include as many details as possible:

**Bug Report Template:**

```markdown
**Describe the bug**
A clear description of what the bug is.

**To Reproduce**
Steps to reproduce the behavior:

1. Start load balancer with '...'
2. Send request to '...'
3. See error

**Expected behavior**
What you expected to happen.

**Actual behavior**
What actually happened.

**Environment:**

- OS: [e.g., Ubuntu 20.04]
- Go Version: [e.g., 1.21]
- Version: [e.g., v1.0.0]

**Additional context**
Logs, screenshots, or any other context.
```

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion, include:

- **Use a clear and descriptive title**
- **Provide a detailed description** of the suggested enhancement
- **Explain why this enhancement would be useful**
- **List some examples** of how it would be used

### Your First Code Contribution

Unsure where to begin? You can start by looking through `beginner` and `help-wanted` issues:

- `beginner` - issues that should only require a few lines of code
- `help-wanted` - issues that may be more involved

### Pull Requests

1. Fork the repo and create your branch from `main`
2. If you've added code that should be tested, add tests
3. If you've changed APIs, update the documentation
4. Ensure the test suite passes
5. Make sure your code follows the style guidelines
6. Issue that pull request!

## Development Setup

### Prerequisites

- Go 1.21+
- Git
- Docker (optional, for testing)

### Setup Steps

```bash
# Fork and clone your fork
git clone https://github.com/YOUR_USERNAME/go-balancer.git
cd go-balancer

# Add upstream remote
git remote add upstream https://github.com/TaiTitans/go-balancer.git

# Install dependencies
go mod download

# Run tests
go test ./...

# Build
go build -o go-balancer cmd/main.go
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...

# Run benchmarks
go test -bench=. ./...
```

### Running Locally

```bash
# Terminal 1: Start backend servers
cd examples/backend-server
go run main.go -port 8081 &
go run main.go -port 8082 &
go run main.go -port 8083 &

# Terminal 2: Start load balancer
go run cmd/main.go
```

## Coding Guidelines

### Go Code Style

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `go fmt` to format your code
- Use `golint` and `go vet` to catch common issues
- Keep functions small and focused
- Write clear, descriptive variable names

### Code Organization

```
package/
├── package.go        # Main package code
├── package_test.go   # Tests
└── doc.go           # Package documentation (if needed)
```

### Naming Conventions

- **Packages:** lowercase, single word
- **Files:** lowercase with underscores
- **Functions:** MixedCaps or mixedCaps
- **Constants:** MixedCaps
- **Variables:** mixedCaps

### Comments

- Package-level comments describing the package
- Function comments describing what the function does
- Complex code sections should have explanatory comments
- Use complete sentences

```go
// Package strategy implements various load balancing strategies.
package strategy

// SelectBackend selects a backend server from the pool
// based on the round-robin algorithm.
func (rr *RoundRobin) SelectBackend(backends []*backend.Backend) *backend.Backend {
    // Implementation
}
```

### Error Handling

```go
// Good
if err != nil {
    return fmt.Errorf("failed to create backend: %w", err)
}

// Bad
if err != nil {
    return err
}
```

### Testing

- Write unit tests for all new functionality
- Aim for >80% code coverage
- Use table-driven tests where appropriate
- Test edge cases and error conditions

```go
func TestSelectBackend(t *testing.T) {
    tests := []struct {
        name     string
        backends []*backend.Backend
        want     *backend.Backend
        wantErr  bool
    }{
        {
            name:     "empty backends",
            backends: []*backend.Backend{},
            want:     nil,
            wantErr:  false,
        },
        // More test cases...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := SelectBackend(tt.backends)
            if got != tt.want {
                t.Errorf("got %v, want %v", got, tt.want)
            }
        })
    }
}
```

## Commit Messages

### Format

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Types

- **feat:** A new feature
- **fix:** A bug fix
- **docs:** Documentation only changes
- **style:** Changes that don't affect code meaning
- **refactor:** Code change that neither fixes a bug nor adds a feature
- **perf:** Performance improvement
- **test:** Adding or correcting tests
- **chore:** Changes to build process or auxiliary tools

### Examples

```
feat(strategy): add weighted round robin strategy

Implement weighted round robin load balancing strategy
that distributes requests based on backend weights.

Closes #123
```

```
fix(backend): prevent race condition in connection tracking

Use atomic operations for connection counter to prevent
race conditions when multiple goroutines access the same backend.

Fixes #456
```

## Pull Request Process

### Before Submitting

- [ ] Run tests: `go test ./...`
- [ ] Run linter: `golangci-lint run`
- [ ] Format code: `go fmt ./...`
- [ ] Update documentation
- [ ] Add/update tests
- [ ] Update CHANGELOG.md

### PR Title Format

Follow the same format as commit messages:

```
feat(strategy): add IP hash strategy
fix(balancer): handle nil backend gracefully
docs(readme): update installation instructions
```

### PR Description Template

```markdown
## Description

Brief description of changes

## Motivation and Context

Why is this change needed? What problem does it solve?

## Type of Change

- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## How Has This Been Tested?

Describe the tests you ran

## Checklist

- [ ] My code follows the code style of this project
- [ ] I have performed a self-review of my own code
- [ ] I have commented my code, particularly in hard-to-understand areas
- [ ] I have made corresponding changes to the documentation
- [ ] My changes generate no new warnings
- [ ] I have added tests that prove my fix is effective or that my feature works
- [ ] New and existing unit tests pass locally with my changes
```

### Review Process

1. Automated checks must pass (tests, linting)
2. At least one maintainer approval required
3. Address all review comments
4. Keep the PR updated with main branch
5. Squash commits if requested

### After Your PR is Merged

- Delete your branch
- Update your fork's main branch
- Close related issues

## Additional Notes

### Branch Naming

- `feature/` - New features
- `fix/` - Bug fixes
- `docs/` - Documentation
- `refactor/` - Code refactoring
- `test/` - Test improvements

Example: `feature/add-ip-hash-strategy`

### Issue Labels

- `bug` - Something isn't working
- `enhancement` - New feature or request
- `documentation` - Improvements or additions to documentation
- `good first issue` - Good for newcomers
- `help wanted` - Extra attention is needed
- `question` - Further information is requested

## Getting Help

- Read the [documentation](docs/)
- Check [existing issues](https://github.com/TaiTitans/go-balancer/issues)
- Ask in discussions
- Contact maintainers

## Recognition

Contributors will be recognized in:

- README.md Contributors section
- Release notes
- Project documentation

Thank you for contributing! 🎉
