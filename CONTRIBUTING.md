# Contributing

## Development rules

- Keep CLI parsing under `internal/cli`.
- Keep image algorithms under `internal/imaging`.
- Do not import Cobra from `internal/imaging`.
- Prefer deterministic, testable pixel operations.
- Avoid cgo-backed codecs unless there is a documented portability reason.

## Before opening a pull request

```bash
gofmt -w .
go test ./...
go test -race ./...
go vet ./...
git diff --check
```

Add tests for new image operations, especially boundary conditions and dimension changes.
