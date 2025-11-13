package constants

import "time"

const (
	RoundRobinStrategy       = "roundrobin"
	LeastConnectionsStrategy = "leastconnections"
	RandomStrategy           = "random"
)

const (
	DefaultHealthCheckInterval = 10 * time.Second
	DefaultHealthCheckTimeout  = 5 * time.Second
)
