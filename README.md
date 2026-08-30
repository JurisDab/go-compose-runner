# devtool

A CLI that replaces ad-hoc shell scripts for local multi-service development:
concurrent service startup, multiplexed log streaming across services via
goroutines/channels, and a single `.devtool.yaml` per project.

## Status

Early scaffold. `up`, `down`, and `logs` wrap `docker compose` per service
defined in the config. Health-check polling, `devtool new` scaffolding, and
`devtool doctor` are not implemented yet.

## Usage

```sh
go build -o devtool ./cmd/devtool

cp .devtool.example.yaml .devtool.yaml
# edit .devtool.yaml to match your project's services

./devtool up
./devtool logs backend
./devtool down
```

## Project layout

```
cmd/devtool/        main package, entrypoint only
internal/cli/       Cobra commands (up, down, logs, ...)
internal/config/    .devtool.yaml parsing
internal/orchestrator/  concurrent docker compose orchestration + log fan-in
```
