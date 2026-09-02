// Package orchestrator starts, stops, and streams logs for the services
// defined in a devtool config, running each service's docker compose
// commands concurrently and fanning their output into one stream.
package orchestrator

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"sync"
	"time"

	"github.com/JurisDab/go-compose-runner/internal/config"
	"github.com/JurisDab/go-compose-runner/internal/healthcheck"
)

// LogLine is one line of output from one service, tagged for display.
type LogLine struct {
	Service  string
	Text     string
	IsErr    bool
	IsStatus bool // health-check result rather than raw compose output
}

// Up starts every service concurrently and streams their combined output
// on the returned channel until the context is cancelled or all services
// exit. Services with a HealthURL are polled in the background and report
// a status line once they're actually ready to receive traffic, rather than
// just "container started." The channel is closed once every service and
// health-check goroutine has finished.
func Up(ctx context.Context, cfg *config.Config) (<-chan LogLine, error) {
	lines := make(chan LogLine)

	var wg sync.WaitGroup
	for _, svc := range cfg.Services {
		wg.Add(1)
		go func(svc config.Service) {
			defer wg.Done()
			runCompose(ctx, svc, []string{"up"}, lines)
		}(svc)

		if svc.HealthURL != "" {
			wg.Add(1)
			go func(svc config.Service) {
				defer wg.Done()
				waitHealthy(ctx, svc, lines)
			}(svc)
		}
	}

	go func() {
		wg.Wait()
		close(lines)
	}()

	return lines, nil
}

func waitHealthy(ctx context.Context, svc config.Service, out chan<- LogLine) {
	timeout := time.Duration(svc.HealthTimeoutSeconds) * time.Second
	if timeout <= 0 {
		timeout = config.DefaultHealthTimeoutSeconds * time.Second
	}

	if err := healthcheck.Wait(ctx, svc.HealthURL, timeout); err != nil {
		out <- LogLine{Service: svc.Name, Text: fmt.Sprintf("health check failed: %v", err), IsErr: true, IsStatus: true}
		return
	}
	out <- LogLine{Service: svc.Name, Text: fmt.Sprintf("ready (%s)", svc.HealthURL), IsStatus: true}
}

// Down stops every service concurrently and waits for all of them to finish.
func Down(ctx context.Context, cfg *config.Config) error {
	lines := make(chan LogLine)
	var wg sync.WaitGroup

	go func() {
		for l := range lines {
			fmt.Printf("[%s] %s\n", l.Service, l.Text)
		}
	}()

	for _, svc := range cfg.Services {
		wg.Add(1)
		go func(svc config.Service) {
			defer wg.Done()
			runCompose(ctx, svc, []string{"down"}, lines)
		}(svc)
	}

	wg.Wait()
	close(lines)
	return nil
}

// Logs tails logs for the given services (or all services in cfg if
// names is empty), fanning them into one channel until ctx is cancelled.
func Logs(ctx context.Context, cfg *config.Config, names []string) (<-chan LogLine, error) {
	services := cfg.Services
	if len(names) > 0 {
		wanted := make(map[string]bool, len(names))
		for _, n := range names {
			wanted[n] = true
		}
		services = nil
		for _, svc := range cfg.Services {
			if wanted[svc.Name] {
				services = append(services, svc)
			}
		}
		if len(services) == 0 {
			return nil, fmt.Errorf("no matching service among: %v", names)
		}
	}

	lines := make(chan LogLine)
	var wg sync.WaitGroup
	for _, svc := range services {
		wg.Add(1)
		go func(svc config.Service) {
			defer wg.Done()
			runCompose(ctx, svc, []string{"logs", "-f"}, lines)
		}(svc)
	}

	go func() {
		wg.Wait()
		close(lines)
	}()

	return lines, nil
}

func runCompose(ctx context.Context, svc config.Service, args []string, out chan<- LogLine) {
	cmdArgs := append([]string{"compose"}, args...)
	cmd := exec.CommandContext(ctx, "docker", cmdArgs...)
	cmd.Dir = svc.WorkDir

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		out <- LogLine{Service: svc.Name, Text: err.Error(), IsErr: true}
		return
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		out <- LogLine{Service: svc.Name, Text: err.Error(), IsErr: true}
		return
	}

	if err := cmd.Start(); err != nil {
		out <- LogLine{Service: svc.Name, Text: err.Error(), IsErr: true}
		return
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go pipeLines(&wg, svc.Name, stdout, out, false)
	go pipeLines(&wg, svc.Name, stderr, out, true)
	wg.Wait()

	if err := cmd.Wait(); err != nil {
		out <- LogLine{Service: svc.Name, Text: fmt.Sprintf("exited: %v", err), IsErr: true}
	}
}

func pipeLines(wg *sync.WaitGroup, service string, r io.Reader, out chan<- LogLine, isErr bool) {
	defer wg.Done()
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		out <- LogLine{Service: service, Text: scanner.Text(), IsErr: isErr}
	}
}
