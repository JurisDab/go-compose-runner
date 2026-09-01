package healthcheck

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func TestWait_SucceedsImmediately(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	err := Wait(context.Background(), srv.URL, 2*time.Second)
	if err != nil {
		t.Fatalf("Wait() error = %v, want nil", err)
	}
}

func TestWait_SucceedsAfterServiceComesUp(t *testing.T) {
	var attempts atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if attempts.Add(1) < 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	err := Wait(context.Background(), srv.URL, 3*time.Second)
	if err != nil {
		t.Fatalf("Wait() error = %v, want nil", err)
	}
	if got := attempts.Load(); got < 3 {
		t.Errorf("attempts = %d, want at least 3", got)
	}
}

func TestWait_TimesOut(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer srv.Close()

	err := Wait(context.Background(), srv.URL, 1200*time.Millisecond)
	if err == nil {
		t.Fatal("Wait() error = nil, want a timeout error")
	}
}

func TestWait_RespectsCancellation(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer srv.Close()

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(200 * time.Millisecond)
		cancel()
	}()

	start := time.Now()
	err := Wait(ctx, srv.URL, 30*time.Second)
	elapsed := time.Since(start)

	if err == nil {
		t.Fatal("Wait() error = nil, want cancellation error")
	}
	if elapsed > 5*time.Second {
		t.Errorf("Wait() took %v after cancellation, want it to return promptly", elapsed)
	}
}
