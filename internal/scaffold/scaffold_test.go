package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestList(t *testing.T) {
	names, err := List()
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	want := []string{"go-service", "node-service"}
	if len(names) != len(want) {
		t.Fatalf("List() = %v, want %v", names, want)
	}
	for i, n := range want {
		if names[i] != n {
			t.Errorf("List()[%d] = %q, want %q", i, names[i], n)
		}
	}
}

func TestGenerate(t *testing.T) {
	tests := []struct {
		name         string
		template     string
		wantErr      bool
		wantFiles    []string
		wantContains map[string]string // file -> substring it must contain after rendering
	}{
		{
			name:     "go-service",
			template: "go-service",
			wantFiles: []string{
				"Dockerfile",
				"docker-compose.yml",
				"main.go",
				"go.mod",
			},
			wantContains: map[string]string{
				"docker-compose.yml": `"9090:9090"`,
				"main.go":            `:9090`,
				"go.mod":             "module payments",
			},
		},
		{
			name:     "node-service",
			template: "node-service",
			wantFiles: []string{
				"Dockerfile",
				"docker-compose.yml",
				"package.json",
				"index.js",
			},
			wantContains: map[string]string{
				"package.json": `"name": "payments"`,
				"index.js":     "payments listening",
			},
		},
		{
			name:     "unknown template",
			template: "does-not-exist",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dest := filepath.Join(t.TempDir(), "payments")

			err := Generate(tt.template, dest, Data{Name: "payments", Port: 9090})
			if (err != nil) != tt.wantErr {
				t.Fatalf("Generate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}

			for _, f := range tt.wantFiles {
				path := filepath.Join(dest, f)
				data, err := os.ReadFile(path)
				if err != nil {
					t.Errorf("expected file %s to exist: %v", f, err)
					continue
				}
				if want, ok := tt.wantContains[f]; ok && !strings.Contains(string(data), want) {
					t.Errorf("%s = %q, want it to contain %q", f, data, want)
				}
			}
		})
	}
}

func TestGenerate_RefusesToOverwriteExistingDir(t *testing.T) {
	dest := filepath.Join(t.TempDir(), "payments")
	if err := os.Mkdir(dest, 0o755); err != nil {
		t.Fatal(err)
	}

	err := Generate("go-service", dest, Data{Name: "payments", Port: 8080})
	if err == nil {
		t.Fatal("Generate() error = nil, want error for existing destination")
	}
}
