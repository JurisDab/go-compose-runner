package healthcheck

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

const pollInterval = 500 * time.Millisecond

func Wait(ctx context.Context, url string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	client := &http.Client{Timeout: pollInterval}
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	if ok(ctx, client, url) {
		return nil
	}

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timed out waiting for %s to become healthy: %w", url, ctx.Err())
		case <-ticker.C:
			if ok(ctx, client, url) {
				return nil
			}
		}
	}
}

func ok(ctx context.Context, client *http.Client, url string) bool {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return false
	}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode >= 200 && resp.StatusCode < 300
}
