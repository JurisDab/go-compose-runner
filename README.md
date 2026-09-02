# devtool

A CLI that replaces ad-hoc shell scripts for local multi-service development:
concurrent service startup with health checks, multiplexed log streaming
across services via goroutines/channels, project scaffolding from embedded
templates, and a single `.devtool.yaml` per project.

## Usage

```sh
go build -o devtool ./cmd/devtool

cp .devtool.example.yaml .devtool.yaml
# edit .devtool.yaml to match your project's services

./devtool doctor            # check Docker, ports, env vars before starting
./devtool up                # start every service, streamed + color-coded
./devtool logs backend      # tail one service (omit the name for all)
./devtool down              # stop everything

./devtool new go-service payments --port 9090   # scaffold a new service
./devtool new list                              # list available templates
```

## Config

`.devtool.yaml`:

```yaml
project: my-project
services:
  - name: backend
    workDir: ./backend
    healthUrl: http://localhost:8080/health
    healthTimeoutSeconds: 30 # optional, defaults to 30
    port: 8080
    requiredEnv: [API_KEY] # optional, checked by `devtool doctor`
```

## Project layout

```
cmd/devtool/             entrypoint; wires Ctrl+C to context cancellation
internal/cli/            Cobra commands: up, down, logs, doctor, new
internal/config/         .devtool.yaml parsing
internal/orchestrator/   concurrent docker compose orchestration + log fan-in
internal/healthcheck/    HTTP polling used by `up` to report "ready"
internal/doctor/         pre-flight environment checks
internal/scaffold/       embed.FS-backed project templates for `devtool new`
```

## Releasing

Tests and `go vet` run on every push/PR via GitHub Actions
([.github/workflows/ci.yml](.github/workflows/ci.yml)). Pushing a tag like
`v0.1.0` triggers
[.github/workflows/release.yml](.github/workflows/release.yml), which runs
GoReleaser ([.goreleaser.yaml](.goreleaser.yaml)) to build cross-platform
binaries and publish them as a GitHub Release.
