# imgx

Fast, classical image processing from the command line.

`imgx` is a Go CLI inspired by ImageMagick's composability, but intentionally keeps its command grammar simple:

```text
imgx <verb> <input> <output> [flags]
```

No AI/ML. No cloud service. No runtime dependency. The released binaries are standalone executables.

## Features

- Convert PNG, JPEG, GIF, BMP and TIFF.
- Decode WebP input.
- Resize with nearest, bilinear, bicubic or Lanczos3 interpolation.
- Crop, rotate, flip.
- Grayscale, brightness and contrast adjustments.
- Gaussian blur and convolution-based sharpening.
- Borders and text watermarks.
- JPEG quality, PNG compression and TIFF LZW compression controls.
- stdin/stdout pipelines.
- Batch processing with bounded worker concurrency.
- Decode-size guard with `--max-pixels` to prevent accidental huge allocations.
- Cross-platform releases for Linux (amd64/arm64/riscv64), macOS (amd64/arm64), and Windows (amd64/arm64).

## Install

### Linux / macOS

After the first GitHub release:

```bash
curl -fsSL https://raw.githubusercontent.com/priyojeet/imgx/main/scripts/install.sh | sh
```

The installer accepts environment overrides:

```bash
IMGX_VERSION=v1.0.0 curl -fsSL https://raw.githubusercontent.com/priyojeet/imgx/main/scripts/install.sh | sh
IMGX_INSTALL_DIR="$HOME/.local/bin" curl -fsSL https://raw.githubusercontent.com/priyojeet/imgx/main/scripts/install.sh | sh
```

### Windows PowerShell

```powershell
iwr https://raw.githubusercontent.com/priyojeet/imgx/main/scripts/install.ps1 -UseBasicParsing | iex
```

See `docs/INSTALL.md` for package-manager and manual installation options.

## Quick start

```bash
imgx info photo.jpg
imgx convert photo.png photo.jpg --quality 92
imgx resize photo.jpg small.jpg --w 800 --h 600
imgx resize photo.jpg small.jpg --w 800 --keep-aspect --algo lanczos3
imgx crop photo.jpg crop.jpg --x 100 --y 80 --w 500 --h 400
imgx rotate photo.jpg rotated.jpg --deg 90
imgx flip photo.jpg mirrored.jpg --horizontal
imgx grayscale photo.jpg gray.jpg
imgx brightness photo.jpg brighter.jpg --value 15
imgx contrast photo.jpg punchier.jpg --value 20
imgx blur photo.jpg soft.jpg --radius 6
imgx sharpen photo.jpg sharp.jpg --amount 1.5
imgx border photo.jpg framed.jpg --width 12 --color '#101010'
imgx watermark photo.jpg marked.jpg --text 'Copyright 2026' --pos bottom-right
imgx compress photo.jpg smaller.jpg --quality 82
```

## Pipes

Output `-` means stdout. Input `-` means stdin. stdout output requires an explicit `--format` because a stream has no extension.

```bash
cat photo.png | imgx convert - out.jpg
imgx convert photo.png - --format jpeg --quality 85 > out.jpg
imgx resize - - --w 900 --keep-aspect --format png < photo.jpg > resized.png
```

Do not add shell progress text to a binary stream: success messages go to stderr and `--quiet` disables them.

## Batch

```bash
imgx batch './photos/*.jpg' \
  --op 'resize --w 1600 --h 1600 --keep-aspect --algo lanczos3' \
  --out-dir ./out

imgx batch './photos/*.png' --op 'grayscale' --out-dir ./gray --jobs 4
```

Batch operations use the same engine as the single-image commands.

## Architecture

```text
cmd/imgx/                 CLI entry point
internal/cli/              Cobra commands and flag parsing
internal/imaging/          image engine; no CLI dependencies
internal/version/          release metadata injected by ldflags
scripts/                   installers
.github/workflows/         CI + tagged releases
.goreleaser.yaml           cross-platform release configuration
docs/                      usage, build and release documentation
```

The important boundary is `internal/imaging`: image processing is implemented as ordinary Go functions over `image.Image` / `*image.RGBA`, so the engine stays testable and reusable.

## Build from source

Requirements:

- Go 1.26 or newer.
- Git.

```bash
git clone https://github.com/priyojeet/imgx.git
cd imgx
go mod download
go test ./...
go vet ./...
go build -trimpath -ldflags='-s -w' -o ./bin/imgx ./cmd/imgx
./bin/imgx --version
```

The repository targets Go 1.26+ so it remains buildable across the current supported Go toolchains. See `docs/BUILDING.md`.

## Release

Push a SemVer tag and GitHub Actions will invoke GoReleaser:

```bash
git tag -a v1.0.0 -m 'imgx v1.0.0'
git push origin v1.0.0
```

The release workflow publishes archives plus SHA-256 checksums. See `docs/RELEASING.md`.

## License

MIT. See [LICENSE](./LICENSE).
