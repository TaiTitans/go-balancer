# Deployment Guide

## Table of Contents

- [Local Development](#local-development)
- [Docker Deployment](#docker-deployment)
- [Docker Compose](#docker-compose)
- [Production Deployment](#production-deployment)
- [Cloud Deployment](#cloud-deployment)
- [Monitoring & Logging](#monitoring--logging)

---

## Local Development

### Prerequisites

- Go 1.21 or higher
- Git

### Setup

```bash
# Clone repository
git clone https://github.com/TaiTitans/go-balancer.git
cd go-balancer

# Install dependencies
go mod download

# Build
go build -o go-balancer cmd/main.go
```

### Running Locally

```bash
# Terminal 1: Start backend servers
cd examples/backend-server
go run main.go -port 8081 -name "Backend-1" &
go run main.go -port 8082 -name "Backend-2" &
go run main.go -port 8083 -name "Backend-3" &

# Terminal 2: Start load balancer
go run cmd/main.go -port 8080

# Terminal 3: Test
curl http://localhost:8080
curl http://localhost:8080/stats
```

---

## Docker Deployment

### Build Image

```bash
docker build -t go-balancer:latest .
```

### Run Container

```bash
docker run -d \
  --name go-balancer \
  -p 8080:8080 \
  -e BACKEND_URLS="http://backend1:8080,http://backend2:8080" \
  go-balancer:latest
```

### With Custom Configuration

```bash
docker run -d \
  --name go-balancer \
  -p 8080:8080 \
  -v $(pwd)/config.json:/root/config.json \
  go-balancer:latest \
  -config /root/config.json
```

---

## Docker Compose

### Start Full Stack

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f load-balancer

# View stats
curl http://localhost:8080/stats

# Stop services
docker-compose down
```

### Custom docker-compose.yml

```yaml
version: "3.8"

services:
  load-balancer:
    build: .
    ports:
      - "8080:8080"
    environment:
      - BACKEND_URLS=http://backend1:8081,http://backend2:8082
      - LB_STRATEGY=leastconnections
    depends_on:
      - backend1
      - backend2
    networks:
      - app-network
    restart: unless-stopped

  backend1:
    build: ./examples/backend-server
    environment:
      - PORT=8081
      - SERVER_NAME=Backend-1
    networks:
      - app-network
    restart: unless-stopped

  backend2:
    build: ./examples/backend-server
    environment:
      - PORT=8082
      - SERVER_NAME=Backend-2
    networks:
      - app-network
    restart: unless-stopped

networks:
  app-network:
    driver: bridge
```

---

## Production Deployment

### Systemd Service

Create `/etc/systemd/system/go-balancer.service`:

```ini
[Unit]
Description=Go Load Balancer
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/go-balancer
ExecStart=/opt/go-balancer/go-balancer \
  -port 8080 \
  -backends "http://backend1:8080,http://backend2:8080,http://backend3:8080" \
  -strategy leastconnections \
  -health-interval 10s

Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
```

Enable and start:

```bash
sudo systemctl daemon-reload
sudo systemctl enable go-balancer
sudo systemctl start go-balancer
sudo systemctl status go-balancer
```

### Nginx Reverse Proxy

```nginx
upstream go_balancer {
    server localhost:8080;
}

server {
    listen 80;
    server_name example.com;

    location / {
        proxy_pass http://go_balancer;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### SSL/TLS with Let's Encrypt

```bash
# Install certbot
sudo apt-get install certbot python3-certbot-nginx

# Get certificate
sudo certbot --nginx -d example.com

# Auto-renewal
sudo certbot renew --dry-run
```

---

## Cloud Deployment

### AWS EC2

```bash
# Launch EC2 instance
aws ec2 run-instances \
  --image-id ami-xxxxx \
  --instance-type t2.micro \
  --key-name your-key \
  --security-groups go-balancer-sg

# SSH into instance
ssh -i your-key.pem ubuntu@ec2-xxx.compute.amazonaws.com

# Install and run
git clone https://github.com/TaiTitans/go-balancer.git
cd go-balancer
go build -o go-balancer cmd/main.go
sudo ./go-balancer -port 80
```

### AWS ECS (Fargate)

**task-definition.json:**

```json
{
  "family": "go-balancer",
  "networkMode": "awsvpc",
  "requiresCompatibilities": ["FARGATE"],
  "cpu": "256",
  "memory": "512",
  "containerDefinitions": [
    {
      "name": "go-balancer",
      "image": "your-repo/go-balancer:latest",
      "portMappings": [
        {
          "containerPort": 8080,
          "protocol": "tcp"
        }
      ],
      "environment": [
        {
          "name": "BACKEND_URLS",
          "value": "http://backend1:8080,http://backend2:8080"
        }
      ]
    }
  ]
}
```

Deploy:

```bash
aws ecs register-task-definition --cli-input-json file://task-definition.json
aws ecs create-service \
  --cluster your-cluster \
  --service-name go-balancer \
  --task-definition go-balancer \
  --desired-count 2
```

### Google Cloud Run

```bash
# Build and push image
gcloud builds submit --tag gcr.io/PROJECT_ID/go-balancer

# Deploy
gcloud run deploy go-balancer \
  --image gcr.io/PROJECT_ID/go-balancer \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated \
  --set-env-vars BACKEND_URLS="http://backend1:8080,http://backend2:8080"
```

### Kubernetes

**deployment.yaml:**

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: go-balancer
  labels:
    app: go-balancer
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
              value: "http://backend-service:8080"
          resources:
            requests:
              memory: "64Mi"
              cpu: "100m"
            limits:
              memory: "128Mi"
              cpu: "200m"
          livenessProbe:
            httpGet:
              path: /health
              port: 8080
            initialDelaySeconds: 10
            periodSeconds: 5
          readinessProbe:
            httpGet:
              path: /health
              port: 8080
            initialDelaySeconds: 5
            periodSeconds: 3
---
apiVersion: v1
kind: Service
metadata:
  name: go-balancer-service
spec:
  selector:
    app: go-balancer
  type: LoadBalancer
  ports:
    - protocol: TCP
      port: 80
      targetPort: 8080
```

Deploy:

```bash
kubectl apply -f deployment.yaml
kubectl get services
kubectl get pods
```

### Helm Chart

**values.yaml:**

```yaml
replicaCount: 3

image:
  repository: go-balancer
  tag: latest
  pullPolicy: IfNotPresent

service:
  type: LoadBalancer
  port: 80
  targetPort: 8080

config:
  backends:
    - "http://backend1:8080"
    - "http://backend2:8080"
  strategy: "roundrobin"
  healthCheckInterval: "10s"

resources:
  limits:
    cpu: 200m
    memory: 128Mi
  requests:
    cpu: 100m
    memory: 64Mi
```

Install:

```bash
helm install go-balancer ./helm-chart
helm upgrade go-balancer ./helm-chart
```

---

## Monitoring & Logging

### Prometheus Integration

Add metrics endpoint (future enhancement):

```go
import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

// In main.go
http.Handle("/metrics", promhttp.Handler())
```

**prometheus.yml:**

```yaml
scrape_configs:
  - job_name: "go-balancer"
    static_configs:
      - targets: ["localhost:8080"]
```

### Grafana Dashboard

Example metrics:

- Request rate
- Error rate
- Backend health status
- Response times
- Active connections

### Logging

Structured logging with Logrus:

```go
import log "github.com/sirupsen/logrus"

log.SetFormatter(&log.JSONFormatter{})
log.SetLevel(log.InfoLevel)
```

Send logs to centralized system:

- **ELK Stack:** Elasticsearch, Logstash, Kibana
- **Loki:** with Grafana
- **CloudWatch:** for AWS

---

## Performance Tuning

### System Limits

```bash
# Increase file descriptor limit
ulimit -n 65536

# /etc/security/limits.conf
* soft nofile 65536
* hard nofile 65536
```

### Go Runtime

```bash
# Set GOMAXPROCS
export GOMAXPROCS=4

# Garbage collection
export GOGC=100
```

### Docker Resources

```yaml
services:
  load-balancer:
    deploy:
      resources:
        limits:
          cpus: "2"
          memory: 512M
        reservations:
          cpus: "1"
          memory: 256M
```

---

## Security Best Practices

1. **Run as non-root user**
2. **Use environment variables for secrets**
3. **Enable HTTPS/TLS**
4. **Implement rate limiting**
5. **Add authentication middleware**
6. **Regular security updates**
7. **Network segmentation**
8. **Firewall configuration**

---

## Backup & Disaster Recovery

1. **Configuration backup:** Store config in version control
2. **Container images:** Tag and store in registry
3. **Monitoring:** Set up alerts for failures
4. **Automated deployments:** CI/CD pipelines
5. **Multi-region deployment:** for high availability

---

## Troubleshooting

### Container won't start

```bash
docker logs go-balancer
docker inspect go-balancer
```

### High memory usage

```bash
docker stats go-balancer
kubectl top pods
```

### Network issues

```bash
docker network inspect bridge
kubectl describe pod go-balancer-xxx
```

---

## Maintenance

### Rolling Updates

**Kubernetes:**

```bash
kubectl set image deployment/go-balancer go-balancer=go-balancer:v2
kubectl rollout status deployment/go-balancer
kubectl rollout undo deployment/go-balancer
```

**Docker Compose:**

```bash
docker-compose pull
docker-compose up -d
```

### Scaling

**Kubernetes:**

```bash
kubectl scale deployment/go-balancer --replicas=5
kubectl autoscale deployment/go-balancer --min=2 --max=10 --cpu-percent=80
```

---

## Support

For issues and questions:

- GitHub Issues: https://github.com/TaiTitans/go-balancer/issues
- Documentation: https://github.com/TaiTitans/go-balancer/docs
