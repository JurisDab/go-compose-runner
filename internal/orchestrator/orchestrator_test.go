package orchestrator

import (
	"context"
	"testing"
	"time"

	"github.com/JurisDab/devtool/internal/config"
)

func TestLogs_UnknownServiceName(t *testing.T) {
	cfg := &config.Config{
		Services: []config.Service{
			{Name: "backend"},
			{Name: "frontend"},
		},
	}

	_, err := Logs(context.Background(), cfg, []string{"does-not-exist"})
	if err == nil {
		t.Fatal("Logs() error = nil, want error for unknown service name")
	}
}

func TestLogs_NoFilterStartsAllServices(t *testing.T) {
	cfg := &config.Config{
		Services: []config.Service{
			{Name: "backend"},
			{Name: "frontend"},
		},
	}

	// A short-lived context ensures any docker processes this spins up get
	// killed via exec.CommandContext instead of leaking past the test.
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	lines, err := Logs(ctx, cfg, nil)
	if err != nil {
		t.Fatalf("Logs() with no name filter returned error: %v", err)
	}
	if lines == nil {
		t.Fatal("Logs() returned nil channel for empty filter")
	}

	for range lines {
		// Drain until the context timeout stops every service goroutine
		// and the channel closes.
	}
}
