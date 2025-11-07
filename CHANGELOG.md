# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Full documentation in docs/ directory
- API documentation with examples
- Deployment guide for various platforms
- Contributing guidelines
- Example configuration files

## [1.0.0] - 2025-11-07

### Added

- Initial release of Go Load Balancer
- Multiple load balancing strategies:
  - Round Robin
  - Least Connections
  - Random
  - Weighted Round Robin (experimental)
  - IP Hash (experimental)
- Health checking system with configurable intervals
- Reverse proxy implementation with proper request forwarding
- Statistics endpoint (`/stats`) for monitoring
- Health endpoint (`/health`) for load balancer status
- Middleware support:
  - Request logging
  - Panic recovery
  - CORS handling
  - Rate limiting (basic)
- Graceful shutdown support
- Docker and Docker Compose support
- Command-line configuration
- Metrics tracking:
  - Total requests
  - Failed requests
  - Success rate
  - Active connections per backend
  - Response times
  - Failure counts
- Thread-safe operations with mutex protection
- Atomic operations for connection counting
- Context-based cancellation support

### Changed

- Improved backend health check reliability
- Enhanced error handling and recovery
- Better logging with structured output
- Optimized connection tracking

### Fixed

- Race conditions in connection counting
- Backend selection when all backends are down
- Memory leaks in proxy connections
- Graceful shutdown timing issues

### Security

- Added CORS middleware
- Implemented proper error messages without exposing internal details
- Request validation and sanitization

## [0.2.0] - 2025-11-06 (Beta)

### Added

- Least Connections strategy
- Random strategy
- Basic health checking
- Connection tracking

### Changed

- Refactored backend management
- Improved code organization

### Fixed

- Round Robin strategy edge cases
- Backend alive status updates

## [0.1.0] - 2025-11-05 (Alpha)

### Added

- Basic Round Robin load balancing
- Simple backend pool management
- HTTP reverse proxy
- Basic error handling

---

## Version History

### Versioning Scheme

We use [Semantic Versioning](https://semver.org/):

- MAJOR version for incompatible API changes
- MINOR version for new functionality in a backward compatible manner
- PATCH version for backward compatible bug fixes

### Release Schedule

- Major releases: When significant features are added
- Minor releases: Monthly or when new features are ready
- Patch releases: As needed for bug fixes

---

## Future Plans

### v1.1.0 (Planned)

- [ ] Prometheus metrics export
- [ ] JSON configuration file support
- [ ] Dynamic backend addition/removal via API
- [ ] Session persistence (sticky sessions)
- [ ] Request/response transformation
- [ ] Advanced rate limiting per client IP

### v1.2.0 (Planned)

- [ ] TLS/SSL termination
- [ ] WebSocket support
- [ ] gRPC load balancing
- [ ] Circuit breaker pattern
- [ ] Retry logic with exponential backoff
- [ ] Request caching

### v2.0.0 (Future)

- [ ] Admin UI dashboard
- [ ] Multi-cluster support
- [ ] Service mesh integration
- [ ] Advanced analytics
- [ ] Plugin system for custom strategies
- [ ] Database connection pooling

---

## How to Upgrade

### From 0.x to 1.0

1. **Configuration Changes:**

   - Update backend URLs format (no changes needed)
   - Review health check intervals (defaults changed to 10s)

2. **API Changes:**

   - `/stats` endpoint response format enhanced
   - Added new fields: `failCount`, `successRate`, `uptime`

3. **Code Changes:**

   - Import paths remain the same
   - Strategy interface unchanged
   - Backend struct has new fields (backward compatible)

4. **Migration Steps:**

```bash
# Backup current setup
cp go-balancer go-balancer.old

# Update to v1.0
go get github.com/TaiTitans/go-balancer@v1.0.0

# Rebuild
go build -o go-balancer cmd/main.go

# Test
./go-balancer --help
```

---

## Breaking Changes

### v1.0.0

- None (first stable release)

### Future Breaking Changes

We commit to maintaining backward compatibility within major versions. Breaking changes will only occur in major version bumps and will be clearly documented here.

---

## Deprecated Features

### v1.0.0

- None currently

Future deprecations will be announced here at least one minor version before removal.

---

## Contributors

Special thanks to all contributors who helped make this project possible!

### Core Team

- [@TaiTitans](https://github.com/TaiTitans) - Creator & Maintainer

### Contributors

- (Your name here! - Contributions welcome)

---

## Support

For questions, issues, or feature requests:

- GitHub Issues: https://github.com/TaiTitans/go-balancer/issues
- Discussions: https://github.com/TaiTitans/go-balancer/discussions
- Email: [Your contact email]

---

[Unreleased]: https://github.com/TaiTitans/go-balancer/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/TaiTitans/go-balancer/releases/tag/v1.0.0
[0.2.0]: https://github.com/TaiTitans/go-balancer/releases/tag/v0.2.0
[0.1.0]: https://github.com/TaiTitans/go-balancer/releases/tag/v0.1.0
