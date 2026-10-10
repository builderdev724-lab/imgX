#!/bin/sh
set -eu

REPO="${IMGX_REPO:-priyojeet/imgx}"
VERSION="${IMGX_VERSION:-latest}"
INSTALL_DIR="${IMGX_INSTALL_DIR:-}"

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$OS" in
  linux|darwin) ;;
  *) echo "imgx installer: unsupported OS: $OS" >&2; exit 1 ;;
esac

case "$ARCH" in
  x86_64|amd64) ARCH=amd64 ;;
  arm64|aarch64) ARCH=arm64 ;;
  riscv64) ARCH=riscv64 ;;
  *) echo "imgx installer: unsupported architecture: $ARCH" >&2; exit 1 ;;
esac

if [ "$OS" = "darwin" ] && [ "$ARCH" = "riscv64" ]; then
  echo "imgx installer: riscv64 is not supported on macOS" >&2
  exit 1
fi

if [ -z "$INSTALL_DIR" ]; then
  if [ "$(id -u)" -eq 0 ]; then
    INSTALL_DIR=/usr/local/bin
  else
    INSTALL_DIR="$HOME/.local/bin"
    mkdir -p "$INSTALL_DIR"
  fi
fi

if [ "$VERSION" = "latest" ]; then
  TAG=$(curl -fsSL -o /dev/null -w '%{url_effective}' "https://github.com/$REPO/releases/latest" | sed 's#.*/tag/##')
else
  TAG="$VERSION"
fi

VERSION_NO_V=${TAG#v}
TARBALL="imgx_${VERSION_NO_V}_${OS}_${ARCH}.tar.gz"
CHECKSUMS="imgx_${VERSION_NO_V}_checksums.txt"
BASE="https://github.com/$REPO/releases/download/${TAG}"
TMP=$(mktemp -d)
trap 'rm -rf "$TMP"' EXIT

curl -fsSL "$BASE/$TARBALL" -o "$TMP/$TARBALL"
curl -fsSL "$BASE/$CHECKSUMS" -o "$TMP/$CHECKSUMS"

EXPECTED=$(awk -v file="$TARBALL" '$2 == file { print $1 }' "$TMP/$CHECKSUMS")
if [ -z "$EXPECTED" ]; then
  echo "imgx installer: checksum entry not found for $TARBALL" >&2
  exit 1
fi

if command -v sha256sum >/dev/null 2>&1; then
  ACTUAL=$(sha256sum "$TMP/$TARBALL" | awk '{print $1}')
elif command -v shasum >/dev/null 2>&1; then
  ACTUAL=$(shasum -a 256 "$TMP/$TARBALL" | awk '{print $1}')
else
  echo "imgx installer: sha256sum or shasum is required" >&2
  exit 1
fi

if [ "$EXPECTED" != "$ACTUAL" ]; then
  echo "imgx installer: SHA-256 verification failed" >&2
  exit 1
fi

tar -xzf "$TMP/$TARBALL" -C "$TMP"
mkdir -p "$INSTALL_DIR"
if [ ! -w "$INSTALL_DIR" ] && [ "$(id -u)" -ne 0 ]; then
  echo "imgx installer: $INSTALL_DIR is not writable; rerun with sudo or set IMGX_INSTALL_DIR" >&2
  exit 1
fi

install -m 0755 "$TMP/imgx" "$INSTALL_DIR/imgx"
echo "imgx installed to $INSTALL_DIR/imgx"
"$INSTALL_DIR/imgx" --version

case ":${PATH}:" in
  *":$INSTALL_DIR:"*) ;;
  *) echo "Note: add $INSTALL_DIR to PATH if it is not already there." >&2 ;;
esac
