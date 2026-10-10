# Releasing imgx

## 1. Verify the repository

```bash
git status
go test ./...
go vet ./...
gofmt -w .
git diff --check
goreleaser check
```

## 2. Commit

```bash
git add .
git commit -m 'release: prepare v1.0.0'
git push origin main
```

## 3. Tag

```bash
git tag -a v1.0.0 -m 'imgx v1.0.0'
git push origin v1.0.0
```

The GitHub Actions release workflow triggers from `vMAJOR.MINOR.PATCH` tags.

## 4. What the workflow publishes

GoReleaser creates:

- Linux amd64 / arm64 / riscv64 archives.
- macOS amd64 / arm64 archives.
- Windows amd64 / arm64 archives.
- SHA-256 checksums.
- GitHub Release notes based on commit history.

The project does not publish Windows riscv64 or macOS riscv64 artifacts; the release matrix is explicitly limited to targets tested and advertised by the project.

## 5. Installer contract

The shell installer expects release assets named like:

```text
imgx_1.0.0_linux_amd64.tar.gz
imgx_1.0.0_linux_arm64.tar.gz
imgx_1.0.0_darwin_arm64.tar.gz
imgx_1.0.0_windows_amd64.zip
imgx_1.0.0_windows_arm64.zip
```

This matches the GoReleaser archive template in `.goreleaser.yaml`.
