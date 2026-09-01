package doctor

import (
	"net"
	"os"
	"testing"

	"github.com/JurisDab/devtool/internal/config"
)

func TestCheckWorkDir(t *testing.T) {
	dir := t.TempDir()
	file := dir + "/not-a-dir"
	if err := os.WriteFile(file, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		workDir string
		wantOK  bool
	}{
		{name: "existing directory", workDir: dir, wantOK: true},
		{name: "missing path", workDir: dir + "/nope", wantOK: false},
		{name: "path is a file, not a directory", workDir: file, wantOK: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := checkWorkDir(config.Service{Name: "svc", WorkDir: tt.workDir})
			if got.OK != tt.wantOK {
				t.Errorf("checkWorkDir(%q).OK = %v, want %v (detail: %s)", tt.workDir, got.OK, tt.wantOK, got.Detail)
			}
		})
	}
}

func TestCheckPortFree(t *testing.T) {
	t.Run("free port", func(t *testing.T) {
		// Ask the OS for a currently-unused port, then release it immediately.
		ln, err := net.Listen("tcp", "localhost:0")
		if err != nil {
			t.Fatal(err)
		}
		port := ln.Addr().(*net.TCPAddr).Port
		ln.Close()

		got := checkPortFree(config.Service{Name: "svc", Port: port})
		if !got.OK {
			t.Errorf("checkPortFree(%d).OK = false, want true (detail: %s)", port, got.Detail)
		}
	})

	t.Run("port in use", func(t *testing.T) {
		ln, err := net.Listen("tcp", "localhost:0")
		if err != nil {
			t.Fatal(err)
		}
		defer ln.Close()
		port := ln.Addr().(*net.TCPAddr).Port

		got := checkPortFree(config.Service{Name: "svc", Port: port})
		if got.OK {
			t.Errorf("checkPortFree(%d).OK = true, want false (port is held open)", port)
		}
	})
}

func TestCheckEnvVar(t *testing.T) {
	const key = "DEVTOOL_TEST_REQUIRED_ENV"

	t.Run("set", func(t *testing.T) {
		t.Setenv(key, "value")
		got := checkEnvVar(config.Service{Name: "svc"}, key)
		if !got.OK {
			t.Errorf("checkEnvVar OK = false, want true")
		}
	})

	t.Run("unset", func(t *testing.T) {
		os.Unsetenv(key)
		got := checkEnvVar(config.Service{Name: "svc"}, key)
		if got.OK {
			t.Errorf("checkEnvVar OK = true, want false")
		}
	})
}

func TestRun_ChecksEveryService(t *testing.T) {
	dir := t.TempDir()
	cfg := &config.Config{
		Project: "demo",
		Services: []config.Service{
			{Name: "a", WorkDir: dir},
			{Name: "b", WorkDir: dir, Port: 1}, // port 1 requires privileges; expected to read as unavailable or free depending on OS, just exercise the path
		},
	}

	checks := Run(cfg)

	// The 3 global docker checks, plus 2 workDir checks, plus 1 port check.
	wantLen := 3 + 2 + 1
	if len(checks) != wantLen {
		t.Fatalf("Run() returned %d checks, want %d: %+v", len(checks), wantLen, checks)
	}
}
