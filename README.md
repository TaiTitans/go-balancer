# go-balancer

A flexible and efficient load balancer implementation in Go with multiple load balancing strategies and health checking capabilities.

## Features

- **Multiple Load Balancing Strategies**

  - Round Robin
  - Least Connections
  - Random
  - Easily extensible for custom strategies

- **Health Checking**

  - Automatic backend health monitoring
  - Configurable health check intervals
  - Automatic removal of unhealthy backends from rotation

- **Production Ready**
  - Thread-safe operations
  - Graceful shutdown support
  - Statistics and monitoring endpoint
  - Connection tracking per backend

## Installation

```bash
go get github.com/TaiTitans/go-balancer
```

## Quick Start

### 1. Start Backend Servers

First, start some backend servers to load balance across:

```bash
# Terminal 1
go run examples/backend-server/main.go -port 8081

# Terminal 2
go run examples/backend-server/main.go -port 8082

# Terminal 3
go run examples/backend-server/main.go -port 8083
```

### 2. Start the Load Balancer

```bash
go run examples/simple/main.go
```

The load balancer will start on `http://localhost:8080` and distribute requests across your backend servers.

### 3. Test the Load Balancer

```bash
# Send requests to the load balancer
curl http://localhost:8080

# Check statistics
curl http://localhost:8080/stats
```

## Usage

### Basic Example

```go
package main

import (
    "context"
    "log"
    "net/http"
    "time"

    "github.com/TaiTitans/go-balancer/balancer"
    "github.com/TaiTitans/go-balancer/strategy"
)

func main() {
    // Configure the load balancer
    config := balancer.Config{
        BackendURLs: []string{
            "http://localhost:8081",
            "http://localhost:8082",
            "http://localhost:8083",
        },
        Strategy:            strategy.NewRoundRobin(),
        HealthCheckInterval: 10 * time.Second,
        HealthCheckTimeout:  5 * time.Second,
    }

    // Create and start load balancer
    lb, err := balancer.NewLoadBalancer(config)
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()
    lb.Start(ctx)

    // Setup HTTP server
    http.Handle("/", lb)
    http.Handle("/stats", lb.HandleStats())

    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

### Using Different Strategies

```go
// Round Robin (default)
config.Strategy = strategy.NewRoundRobin()

// Least Connections
config.Strategy = strategy.NewLeastConnections()

// Random
config.Strategy = strategy.NewRandom()
```

### Custom Strategy

Implement the `Strategy` interface to create your own load balancing algorithm:

```go
type Strategy interface {
    SelectBackend(backends []*backend.Backend) *backend.Backend
    Name() string
}
```

## Project Structure

```
go-balancer/
├── backend/          # Backend server representation
├── balancer/         # Main load balancer implementation
├── strategy/         # Load balancing strategies
├── healthcheck/      # Health checking functionality
├── examples/         # Example applications
│   ├── simple/       # Simple load balancer example
│   └── backend-server/ # Example backend server
├── cmd/              # Command-line tools
├── go.mod
├── Makefile
└── README.md
```

## Building

```bash
# Build the load balancer
make build

# Build the backend server
make backend

# Run tests
make test

# Run tests with coverage
make test-coverage

# Run the load balancer
make run
```

## Testing

The project includes comprehensive tests for all components:

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Generate coverage report
make test-coverage
```

## API Endpoints

### Main Endpoint

- `GET /` - Proxy requests to backend servers

### Stats Endpoint

- `GET /stats` - View load balancer statistics including:
  - Current strategy
  - Backend health status
  - Active connections per backend
  - Response times

## Configuration Options

| Option              | Type          | Default  | Description                       |
| ------------------- | ------------- | -------- | --------------------------------- |
| BackendURLs         | []string      | required | List of backend server URLs       |
| Strategy            | Strategy      | required | Load balancing strategy to use    |
| HealthCheckInterval | time.Duration | 10s      | Interval between health checks    |
| HealthCheckTimeout  | time.Duration | 5s       | Timeout for health check requests |

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Roadmap

- [ ] Weighted load balancing
- [ ] Sticky sessions support
- [ ] Circuit breaker pattern
- [ ] Metrics export (Prometheus)
- [ ] TLS/HTTPS support
- [ ] Configuration file support
- [ ] Dynamic backend addition/removal
- [ ] Rate limiting
- [ ] Request/Response middleware support
