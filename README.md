# Go Load Balancer

[![Go Version](https://img.shields.io/badge/Go-1.25+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

A high-performance, production-ready HTTP load balancer written in Go. Features multiple load balancing strategies, health checking, metrics, and more.

## ✨ Features

- **Multiple Load Balancing Strategies**

  - Round Robin
  - Least Connections
  - Random
  - Weighted Round Robin
  - IP Hash

- **Health Checking**

  - Automatic backend health monitoring
  - Configurable check intervals
  - Automatic failure detection and recovery

- **Production Ready**

  - Reverse proxy with proper request forwarding
  - Graceful shutdown
  - Request/response logging
  - Error handling and recovery
  - CORS support
  - Rate limiting

- **Metrics & Monitoring**

  - Real-time statistics endpoint
  - Request tracking
  - Connection monitoring
  - Response time measurement
  - Success rate calculation

- **Easy Deployment**
  - Docker support
  - Docker Compose for full stack
  - Command-line configuration
  - Minimal dependencies

## 📦 Installation

```bash
# Clone the repository
git clone https://github.com/TaiTitans/go-balancer.git
cd go-balancer

# Install dependencies
go mod download

# Build
go build -o go-balancer cmd/main.go
```

## 🚀 Quick Start

### Start Backend Servers

```bash
# Terminal 1
cd examples/backend-server
go run main.go -port 8081 -name "Backend-1"

# Terminal 2
go run main.go -port 8082 -name "Backend-2"

# Terminal 3
go run main.go -port 8083 -name "Backend-3"
```

### Start Load Balancer

```bash
# Terminal 4
go run cmd/main.go \
  -port 8080 \
  -backends "http://localhost:8081,http://localhost:8082,http://localhost:8083" \
  -strategy roundrobin
```

### Test It

```bash
# Send requests
curl http://localhost:8080

# View statistics
curl http://localhost:8080/stats

# Health check
curl http://localhost:8080/health
```

## 🐳 Docker Deployment

```bash
# Build and run with Docker Compose
docker-compose up --build

# Test
curl http://localhost:8080
```

## 📖 Usage

### Command Line Options

```bash
go-balancer [options]

Options:
  -port int
        Load balancer port (default 8080)
  -backends string
        Comma-separated list of backend URLs
        (default "http://localhost:8081,http://localhost:8082,http://localhost:8083")
  -strategy string
        Load balancing strategy: roundrobin, leastconnections, random
        (default "roundrobin")
  -health-interval duration
        Health check interval (default 10s)
  -health-timeout duration
        Health check timeout (default 5s)
```

### Programmatic Usage

```go
package main

import (
    "context"
    "time"

    "github.com/TaiTitans/go-balancer/balancer"
    "github.com/TaiTitans/go-balancer/strategy"
)

func main() {
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

    lb, err := balancer.NewLoadBalancer(config)
    if err != nil {
        panic(err)
    }

    ctx := context.Background()
    lb.Start(ctx)

    // Use lb as http.Handler
    http.ListenAndServe(":8080", lb)
}
```

## 🎯 Load Balancing Strategies

### Round Robin

Distributes requests evenly across all healthy backends in a circular order.

```go
strategy := strategy.NewRoundRobin()
```

### Least Connections

Routes requests to the backend with the fewest active connections.

```go
strategy := strategy.NewLeastConnections()
```

### Random

Randomly selects a healthy backend for each request.

```go
strategy := strategy.NewRandom()
```

### Weighted Round Robin

Distributes requests based on backend weights.

```go
weights := map[*backend.Backend]int{
    backend1: 3,  // 3x more requests
    backend2: 2,
    backend3: 1,
}
strategy := strategy.NewWeightedRoundRobin(weights)
```

## 📊 Statistics Endpoint

Access `/stats` to view load balancer statistics:

```
╔════════════════════════════════════════╗
║   Load Balancer Statistics             ║
╚════════════════════════════════════════╝

Strategy:         RoundRobin
Uptime:           1h23m45s
Total Backends:   3
Alive Backends:   3
Total Requests:   15234
Failed Requests:  12
Success Rate:     99.92%
Active Connections: 5

Backend Details:
════════════════════════════════════════

[1] http://localhost:8081
    Status:       ✓ Healthy
    Connections:  2
    Response Time: 15ms
    Fail Count:   0

[2] http://localhost:8082
    Status:       ✓ Healthy
    Connections:  1
    Response Time: 12ms
    Fail Count:   0

[3] http://localhost:8083
    Status:       ✓ Healthy
    Connections:  2
    Response Time: 18ms
    Fail Count:   0
```

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package tests
go test ./backend
go test ./strategy
go test ./balancer

# Benchmark
go test -bench=. ./...
```

## 📁 Project Structure

```
go-balancer/
├── backend/          # Backend server management
├── balancer/         # Main load balancer logic
├── cmd/              # Main application entry point
├── config/           # Configuration management
├── examples/         # Example applications
│   ├── backend-server/
│   └── simple/
├── healthcheck/      # Health checking logic
├── middleware/       # HTTP middleware
├── strategy/         # Load balancing strategies
├── Dockerfile        # Docker configuration
├── docker-compose.yml
├── go.mod
├── go.sum
├── LICENSE
├── Makefile
└── README.md
```

## 🔧 Configuration

### Environment Variables

```bash
export BACKEND_URLS="http://localhost:8081,http://localhost:8082"
export LB_STRATEGY="leastconnections"
export LB_PORT="8080"
export HEALTH_CHECK_INTERVAL="10s"
```

### Configuration File (Future)

```json
{
  "server": {
    "port": 8080,
    "readTimeout": "15s",
    "writeTimeout": "15s"
  },
  "backends": [
    { "url": "http://localhost:8081", "weight": 3 },
    { "url": "http://localhost:8082", "weight": 2 }
  ],
  "healthCheck": {
    "interval": "10s",
    "timeout": "5s"
  },
  "strategy": {
    "type": "roundrobin"
  }
}
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Inspired by production load balancers like NGINX and HAProxy
- Built with Go's powerful `net/http` and `httputil` packages
- Thanks to the Go community for excellent documentation

## 📧 Contact

TaiTitans - [@TaiTitans](https://github.com/TaiTitans)

Project Link: [https://github.com/TaiTitans/go-balancer](https://github.com/TaiTitans/go-balancer)

---

## Roadmap

- [ ] Sticky sessions support
- [ ] Circuit breaker pattern
- [ ] Rate limiting

Made with ❤️ using Go
