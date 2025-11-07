# API Documentation

## Endpoints

### Main Load Balancer Endpoint

**URL:** `/`  
**Methods:** `GET`, `POST`, `PUT`, `DELETE`, etc.  
**Description:** Main endpoint that load balances requests to backend servers

**Example:**

```bash
curl http://localhost:8080/
```

**Response:**

```
Response from backend: Backend-1
Time: 2025-11-07T10:30:00Z
Path: /
Method: GET
```

---

### Statistics Endpoint

**URL:** `/stats`  
**Method:** `GET`  
**Description:** Returns current load balancer statistics

**Example:**

```bash
curl http://localhost:8080/stats
```

**Response:**

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
```

---

### Health Check Endpoint

**URL:** `/health`  
**Method:** `GET`  
**Description:** Returns health status of the load balancer

**Example:**

```bash
curl http://localhost:8080/health
```

**Response:**

```json
{
  "status": "healthy",
  "timestamp": "2025-11-07T10:30:00Z"
}
```

---

## Load Balancing Strategies

### 1. Round Robin

Distributes requests evenly across all healthy backends in a circular order.

**Usage:**

```bash
./go-balancer -strategy roundrobin
```

**Characteristics:**

- Simple and fair distribution
- No state required
- Good for backends with similar capacity

---

### 2. Least Connections

Routes requests to the backend with the fewest active connections.

**Usage:**

```bash
./go-balancer -strategy leastconnections
```

**Characteristics:**

- Better for backends with varying response times
- Prevents overloading slow backends
- Requires connection tracking

---

### 3. Random

Randomly selects a healthy backend for each request.

**Usage:**

```bash
./go-balancer -strategy random
```

**Characteristics:**

- Simple implementation
- No state required
- Good distribution over time

---

## Backend Server Endpoints

### Health Check

**URL:** `/health`  
**Method:** `GET`  
**Description:** Health check endpoint for backend servers

**Example:**

```bash
curl http://localhost:8081/health
```

**Response:**

```json
{
  "status": "healthy",
  "server": "Backend-1",
  "timestamp": "2025-11-07T10:30:00Z"
}
```

---

### Slow Response Test

**URL:** `/slow`  
**Method:** `GET`  
**Description:** Returns a response after 2 seconds delay (for testing)

**Example:**

```bash
curl http://localhost:8081/slow
```

---

### Error Response Test

**URL:** `/error`  
**Method:** `GET`  
**Description:** Returns 500 error (for testing error handling)

**Example:**

```bash
curl http://localhost:8081/error
```

---

## Configuration

### Command Line Flags

| Flag               | Type     | Default                     | Description                  |
| ------------------ | -------- | --------------------------- | ---------------------------- |
| `-port`            | int      | 8080                        | Load balancer port           |
| `-backends`        | string   | "http://localhost:8081,..." | Comma-separated backend URLs |
| `-strategy`        | string   | "roundrobin"                | Load balancing strategy      |
| `-health-interval` | duration | 10s                         | Health check interval        |
| `-health-timeout`  | duration | 5s                          | Health check timeout         |

**Example:**

```bash
./go-balancer \
  -port 9000 \
  -backends "http://backend1:8080,http://backend2:8080" \
  -strategy leastconnections \
  -health-interval 5s \
  -health-timeout 3s
```

---

## Metrics

The load balancer tracks the following metrics:

- **Total Requests:** Total number of requests processed
- **Failed Requests:** Number of requests that failed
- **Success Rate:** Percentage of successful requests
- **Active Connections:** Current number of active connections
- **Backend Status:** Health status of each backend
- **Response Time:** Average response time per backend
- **Fail Count:** Number of consecutive failures per backend

---

## Error Handling

### No Available Backends

**Status Code:** `503 Service Unavailable`  
**Response:** `Service unavailable`

**Occurs when:**

- All backends are down
- No backends configured

---

### Backend Error

**Status Code:** `502 Bad Gateway`  
**Response:** `Bad Gateway`

**Occurs when:**

- Backend connection fails
- Backend returns error
- Backend timeout

---

## Testing

### Load Testing

```bash
# Using Apache Bench
ab -n 10000 -c 100 http://localhost:8080/

# Using Hey
hey -n 10000 -c 100 http://localhost:8080/

# Using wrk
wrk -t12 -c400 -d30s http://localhost:8080/
```

### Health Check Testing

```bash
# Check load balancer health
curl http://localhost:8080/health

# Check backend health
curl http://localhost:8081/health
curl http://localhost:8082/health
curl http://localhost:8083/health
```

### Strategy Testing

```bash
# Test different strategies
for i in {1..10}; do
  curl -s http://localhost:8080/ | grep "backend"
done
```

---

## Production Deployment

### Docker

```bash
# Build image
docker build -t go-balancer:latest .

# Run container
docker run -p 8080:8080 \
  -e BACKEND_URLS="http://backend1:8080,http://backend2:8080" \
  go-balancer:latest
```

### Docker Compose

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

### Kubernetes

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: go-balancer
spec:
  replicas: 3
  selector:
    matchLabels:
      app: go-balancer
  template:
    metadata:
      labels:
        app: go-balancer
    spec:
      containers:
        - name: go-balancer
          image: go-balancer:latest
          ports:
            - containerPort: 8080
          env:
            - name: BACKEND_URLS
              value: "http://backend1:8080,http://backend2:8080"
```

---

## Monitoring

### Prometheus Metrics (Future Enhancement)

```
# HELP lb_requests_total Total number of requests
# TYPE lb_requests_total counter
lb_requests_total{strategy="roundrobin"} 15234

# HELP lb_requests_failed Total number of failed requests
# TYPE lb_requests_failed counter
lb_requests_failed{strategy="roundrobin"} 12

# HELP lb_backend_status Backend health status (1=healthy, 0=down)
# TYPE lb_backend_status gauge
lb_backend_status{backend="http://localhost:8081"} 1
```

---

## Troubleshooting

### High Latency

1. Check backend response times in `/stats`
2. Consider using `leastconnections` strategy
3. Scale up backends

### Backends Marked as Down

1. Check backend health endpoint
2. Verify network connectivity
3. Review health check timeout settings

### Uneven Distribution

1. Verify strategy configuration
2. Check backend connection counts
3. Consider using `weighted` strategy for different capacities

---

## Security Considerations

- No authentication/authorization built-in (add middleware)
- CORS enabled by default (configure as needed)
- No rate limiting per client (add middleware)
- Use HTTPS in production
- Validate backend URLs

---

## Performance Tips

1. **Use appropriate strategy:**

   - Similar backends → Round Robin
   - Different capacities → Weighted Round Robin
   - Varying response times → Least Connections

2. **Tune health checks:**

   - Longer intervals for stable backends
   - Shorter timeouts for faster failover

3. **Monitor metrics:**

   - Watch success rate
   - Track response times
   - Monitor connection counts

4. **Scale appropriately:**
   - Add backends as load increases
   - Use multiple load balancer instances
   - Consider using a service mesh for large deployments
