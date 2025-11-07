# Quick Start Guide

Get up and running with Go Load Balancer in 5 minutes!

## Prerequisites

- Go 1.21 or higher
- Git (optional, for cloning)

## Step 1: Get the Code

```bash
# Option A: Clone from GitHub
git clone https://github.com/TaiTitans/go-balancer.git
cd go-balancer

# Option B: Download and extract ZIP
# Then navigate to the directory
cd go-balancer
```

## Step 2: Build

```bash
# Build the load balancer
go build -o go-balancer cmd/main.go

# Or use the Makefile (Linux/Mac)
make build
```

## Step 3: Start Backend Servers

Open 3 terminals and run:

**Terminal 1:**

```bash
cd examples/backend-server
go run main.go -port 8081 -name "Backend-1"
```

**Terminal 2:**

```bash
cd examples/backend-server
go run main.go -port 8082 -name "Backend-2"
```

**Terminal 3:**

```bash
cd examples/backend-server
go run main.go -port 8083 -name "Backend-3"
```

You should see:

```
Backend server [Backend-1] starting on :8081
Backend server [Backend-2] starting on :8082
Backend server [Backend-3] starting on :8083
```

## Step 4: Start Load Balancer

**Terminal 4:**

```bash
# Using the binary
./go-balancer

# Or run directly
go run cmd/main.go
```

You should see:

```
╔════════════════════════════════════════╗
║   Go Load Balancer                     ║
╚════════════════════════════════════════╝
Port:          8080
Strategy:      RoundRobin
Backends:      3
Health Check:  10s

Endpoints:
  - Load Balancer: http://localhost:8080/
  - Statistics:    http://localhost:8080/stats
  - Health:        http://localhost:8080/health

Backends:
  [1] http://localhost:8081
  [2] http://localhost:8082
  [3] http://localhost:8083
════════════════════════════════════════
```

## Step 5: Test It!

### Send Requests

```bash
# Send multiple requests
for i in {1..10}; do
  curl http://localhost:8080
done
```

You'll see responses from different backends:

```
Response from backend: Backend-1
Response from backend: Backend-2
Response from backend: Backend-3
Response from backend: Backend-1
...
```

### View Statistics

```bash
curl http://localhost:8080/stats
```

Output:

```
╔════════════════════════════════════════╗
║   Load Balancer Statistics             ║
╚════════════════════════════════════════╝

Strategy:         RoundRobin
Uptime:           2m15s
Total Backends:   3
Alive Backends:   3
Total Requests:   10
Failed Requests:  0
Success Rate:     100.00%
Active Connections: 0

Backend Details:
════════════════════════════════════════

[1] http://localhost:8081
    Status:       ✓ Healthy
    Connections:  0
    Response Time: 1.2ms
    Fail Count:   0

[2] http://localhost:8082
    Status:       ✓ Healthy
    Connections:  0
    Response Time: 1.1ms
    Fail Count:   0

[3] http://localhost:8083
    Status:       ✓ Healthy
    Connections:  0
    Response Time: 1.3ms
    Fail Count:   0
```

### Check Health

```bash
curl http://localhost:8080/health
```

Output:

```json
{ "status": "healthy", "timestamp": "2025-11-07T10:30:00Z" }
```

## Step 6: Try Different Strategies

### Least Connections

```bash
./go-balancer -strategy leastconnections
```

### Random

```bash
./go-balancer -strategy random
```

### Custom Port

```bash
./go-balancer -port 9000
```

### Custom Backends

```bash
./go-balancer -backends "http://server1:8080,http://server2:8080"
```

### Custom Health Check

```bash
./go-balancer \
  -health-interval 5s \
  -health-timeout 3s
```

## Step 7: Docker Deployment (Optional)

### Using Docker Compose

```bash
# Start everything with one command
docker-compose up --build

# Test
curl http://localhost:8080
curl http://localhost:8080/stats

# Stop
docker-compose down
```

### Using Docker

```bash
# Build image
docker build -t go-balancer .

# Run container
docker run -p 8080:8080 go-balancer
```

## Testing Features

### Test Health Checking

1. Stop one backend server (Ctrl+C in terminal)
2. Watch the load balancer logs
3. Send requests - they should only go to healthy backends
4. Restart the backend
5. It should automatically rejoin the pool

### Test Load Distribution

```bash
# Send 100 requests
for i in {1..100}; do
  curl -s http://localhost:8080 | grep "backend"
done | sort | uniq -c
```

You should see roughly equal distribution:

```
  33 Response from backend: Backend-1
  34 Response from backend: Backend-2
  33 Response from backend: Backend-3
```

### Test Error Handling

```bash
# Hit the error endpoint on a backend
curl http://localhost:8081/error

# The load balancer should mark it as unhealthy
curl http://localhost:8080/stats
```

### Test Slow Responses

```bash
# Time a slow endpoint
time curl http://localhost:8081/slow

# Load balancer should track response times
curl http://localhost:8080/stats
```

## Common Commands

```bash
# Build
go build -o go-balancer cmd/main.go

# Run tests
go test ./...

# Run with coverage
go test -cover ./...

# Format code
go fmt ./...

# Clean build
make clean build

# Run example
go run examples/simple/main.go
```

## Troubleshooting

### "address already in use"

```bash
# Find process using port 8080
lsof -i :8080  # Linux/Mac
netstat -ano | findstr :8080  # Windows

# Kill the process or use a different port
./go-balancer -port 8081
```

### "no available backends"

```bash
# Make sure backend servers are running
# Check they're on the correct ports
# Verify backend URLs in the command
```

### Backends marked as down

```bash
# Check backend health endpoints
curl http://localhost:8081/health
curl http://localhost:8082/health
curl http://localhost:8083/health

# Adjust health check timeout if needed
./go-balancer -health-timeout 10s
```

## Next Steps

1. **Read the full documentation:**

   - [API Documentation](docs/API.md)
   - [Deployment Guide](docs/DEPLOYMENT.md)
   - [Project Summary](PROJECT_SUMMARY.md)

2. **Try advanced features:**

   - Configure health check intervals
   - Test different load balancing strategies
   - Monitor metrics and statistics
   - Deploy with Docker

3. **Contribute:**

   - Read [CONTRIBUTING.md](CONTRIBUTING.md)
   - Report issues on GitHub
   - Submit pull requests

4. **Deploy to production:**
   - Review [DEPLOYMENT.md](docs/DEPLOYMENT.md)
   - Set up monitoring
   - Configure TLS/SSL
   - Use a process manager

## Quick Reference

| Command                      | Description         |
| ---------------------------- | ------------------- |
| `./go-balancer`              | Start with defaults |
| `./go-balancer --help`       | Show all options    |
| `curl localhost:8080`        | Send request        |
| `curl localhost:8080/stats`  | View statistics     |
| `curl localhost:8080/health` | Check health        |
| `docker-compose up`          | Start full stack    |
| `go test ./...`              | Run tests           |
| `make build`                 | Build binary        |

## Support

- 📖 Documentation: [docs/](docs/)
- 🐛 Issues: [GitHub Issues](https://github.com/TaiTitans/go-balancer/issues)
- 💬 Discussions: [GitHub Discussions](https://github.com/TaiTitans/go-balancer/discussions)

---

**Congratulations! You're now running a production-ready load balancer! 🎉**
