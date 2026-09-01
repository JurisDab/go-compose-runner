package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name    string
		yaml    string
		missing bool
		wantErr bool
		check   func(t *testing.T, cfg *Config)
	}{
		{
			name: "valid config",
			yaml: `
project: demo
services:
  - name: backend
    workDir: ./backend
    healthUrl: http://localhost:8080/health
    healthTimeoutSeconds: 10
    port: 8080
    requiredEnv: [API_KEY]
`,
			check: func(t *testing.T, cfg *Config) {
				if cfg.Project != "demo" {
					t.Errorf("Project = %q, want %q", cfg.Project, "demo")
				}
				if len(cfg.Services) != 1 {
					t.Fatalf("len(Services) = %d, want 1", len(cfg.Services))
				}
				svc := cfg.Services[0]
				if svc.Name != "backend" || svc.Port != 8080 || svc.HealthTimeoutSeconds != 10 {
					t.Errorf("unexpected service: %+v", svc)
				}
				if len(svc.RequiredEnv) != 1 || svc.RequiredEnv[0] != "API_KEY" {
					t.Errorf("RequiredEnv = %v, want [API_KEY]", svc.RequiredEnv)
				}
			},
		},
		{
			name:    "no services",
			yaml:    `project: demo`,
			wantErr: true,
		},
		{
			name:    "malformed yaml",
			yaml:    "services: [this is not valid: yaml: at all",
			wantErr: true,
		},
		{
			name:    "missing file",
			missing: true,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), ".devtool.yaml")
			if !tt.missing {
				writeFile(t, path, tt.yaml)
			}

			cfg, err := Load(path)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Load() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if tt.check != nil {
				tt.check(t, cfg)
			}
		})
	}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("writing test fixture: %v", err)
	}
}
