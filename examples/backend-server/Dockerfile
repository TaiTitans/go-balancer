# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod ./
RUN go mod download

# Copy source code
COPY . .

# Build the backend server
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o backend-server ./examples/backend-server/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates curl

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/backend-server .

# Expose port (will be overridden by docker-compose)
EXPOSE 8081

# Run the application
ENTRYPOINT ["./backend-server"]
CMD ["-port", "8081", "-name", "Backend-1"]