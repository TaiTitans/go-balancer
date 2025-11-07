package backend

import (
	"testing"
	"time"
)

func TestNewBackend(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "valid url",
			url:     "http://localhost:8080",
			wantErr: false,
		},
		{
			name:    "invalid url",
			url:     "://invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			backend, err := NewBackend(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewBackend() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && backend == nil {
				t.Error("NewBackend() returned nil backend")
			}
			if !tt.wantErr && !backend.IsAlive() {
				t.Error("NewBackend() backend should be alive by default")
			}
		})
	}
}

func TestBackend_SetAlive(t *testing.T) {
	backend, _ := NewBackend("http://localhost:8080")

	backend.SetAlive(false)
	if backend.IsAlive() {
		t.Error("SetAlive(false) did not set backend to not alive")
	}

	backend.SetAlive(true)
	if !backend.IsAlive() {
		t.Error("SetAlive(true) did not set backend to alive")
	}
}

func TestBackend_Connections(t *testing.T) {
	backend, _ := NewBackend("http://localhost:8080")

	if backend.GetConnections() != 0 {
		t.Errorf("Initial connections should be 0, got %d", backend.GetConnections())
	}

	backend.IncrementConnections()
	if backend.GetConnections() != 1 {
		t.Errorf("After increment, connections should be 1, got %d", backend.GetConnections())
	}

	backend.DecrementConnections()
	if backend.GetConnections() != 0 {
		t.Errorf("After decrement, connections should be 0, got %d", backend.GetConnections())
	}

	// Test that decrement doesn't go below 0
	backend.DecrementConnections()
	if backend.GetConnections() != 0 {
		t.Errorf("Connections should not go below 0, got %d", backend.GetConnections())
	}
}

func TestBackend_ResponseTime(t *testing.T) {
	backend, _ := NewBackend("http://localhost:8080")

	testDuration := 100 * time.Millisecond
	backend.UpdateResponseTime(testDuration)

	if backend.GetResponseTime() != testDuration {
		t.Errorf("Expected response time %v, got %v", testDuration, backend.GetResponseTime())
	}
}
