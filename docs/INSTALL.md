# Installation

## Linux / macOS one-line installer

```bash
curl -fsSL https://raw.githubusercontent.com/priyojeet/imgx/main/scripts/install.sh | sh
```

The installer downloads the matching archive from GitHub Releases. It supports Linux/macOS on amd64, arm64, and Linux riscv64.

### Custom installation directory

```bash
IMGX_INSTALL_DIR="$HOME/.local/bin" \
curl -fsSL https://raw.githubusercontent.com/priyojeet/imgx/main/scripts/install.sh | sh
```

### Pin a version

```bash
IMGX_VERSION=v1.0.0 \
curl -fsSL https://raw.githubusercontent.com/priyojeet/imgx/main/scripts/install.sh | sh
```

### Windows

Run PowerShell as your normal user:

```powershell
iwr https://raw.githubusercontent.com/priyojeet/imgx/main/scripts/install.ps1 -UseBasicParsing | iex
```

Set `IMGX_INSTALL_DIR` to choose another directory.

## Manual installation

1. Open the project's GitHub Releases page.
2. Download the archive for your OS/architecture.
3. Extract `imgx` / `imgx.exe` into a directory on `PATH`.
4. Run `imgx --version`.

## Building from source

See `docs/BUILDING.md`.
