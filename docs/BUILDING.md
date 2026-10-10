# Building imgx

## Toolchain

Use Go 1.26 or newer. The module uses Go 1.26 because the current `golang.org/x/image` release requires Go 1.26.

Check your toolchain:

```bash
go version
go env GOPATH GOMODCACHE GOTOOLCHAIN
```

## Local build

```bash
git clone https://github.com/priyojeet/imgx.git
cd imgx

go mod download
go test ./...
go vet ./...

gofmt -w .
go build -trimpath -ldflags='-s -w' -o ./bin/imgx ./cmd/imgx
./bin/imgx --help
```

For a smaller production binary:

```bash
go build -trimpath -ldflags='-s -w' -o ./bin/imgx ./cmd/imgx
```

## Reproducible-ish local release build

Inject release metadata manually:

```bash
go build -trimpath \
  -ldflags='-s -w \
    -X github.com/priyojeet/imgx/internal/version.Version=v1.0.0 \
    -X github.com/priyojeet/imgx/internal/version.Commit=$(git rev-parse HEAD) \
    -X github.com/priyojeet/imgx/internal/version.Date=$(date -u +%Y-%m-%dT%H:%M:%SZ)' \
  -o ./bin/imgx ./cmd/imgx
```

## Cross-compilation without GoReleaser

Because imgx has `CGO_ENABLED=0` in release builds, normal Go cross-compilation works:

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/imgx-linux-amd64 ./cmd/imgx
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o bin/imgx-linux-arm64 ./cmd/imgx
CGO_ENABLED=0 GOOS=linux GOARCH=riscv64 go build -o bin/imgx-linux-riscv64 ./cmd/imgx
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o bin/imgx-darwin-amd64 ./cmd/imgx
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o bin/imgx-darwin-arm64 ./cmd/imgx
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/imgx-windows-amd64.exe ./cmd/imgx
CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -o bin/imgx-windows-arm64.exe ./cmd/imgx
```

## GoReleaser

Install GoReleaser, then run a local snapshot:

```bash
go install github.com/goreleaser/goreleaser/v2@latest
goreleaser check
goreleaser release --snapshot --clean
```

The generated artifacts appear under `dist/`.

## Why no cgo?

The image engine avoids cgo-backed codecs. That keeps cross-compilation straightforward and means a release binary does not depend on system image libraries.
