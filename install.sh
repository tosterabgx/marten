#!/bin/sh
set -eu

REPO="tosterabgx/marten"
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="marten"

os="$(uname -s)"
arch="$(uname -m)"

case "$os" in
  Linux)  goos="linux" ;;
  Darwin) goos="darwin" ;;
  *)
    echo "error: unsupported OS '$os', download manually from:" >&2
    echo "  https://github.com/$REPO/releases/latest" >&2
    exit 1
    ;;
esac

case "$arch" in
  x86_64|amd64)  goarch="amd64" ;;
  arm64|aarch64) goarch="arm64" ;;
  *)
    echo "error: unsupported architecture '$arch', download manually from:" >&2
    echo "  https://github.com/$REPO/releases/latest" >&2
    exit 1
    ;;
esac

download_url="https://github.com/$REPO/releases/latest/download/marten-$goos-$goarch"

echo "Downloading marten for $goos/$goarch..."
echo "  -> $download_url"

tmpfile="$(mktemp)"
trap 'rm -f "$tmpfile"' EXIT

if ! curl --retry 3 --retry-delay 5 --max-time 20 -fsSL "$download_url" -o "$tmpfile"; then
  echo "error: download failed" >&2
  exit 1
fi

chmod +x "$tmpfile"

if [ -w "$INSTALL_DIR" ]; then
  mv "$tmpfile" "$INSTALL_DIR/$BINARY_NAME"
else
  echo "Installing to $INSTALL_DIR requires sudo:"
  sudo mv "$tmpfile" "$INSTALL_DIR/$BINARY_NAME"
fi

echo "marten installed to $INSTALL_DIR/$BINARY_NAME"
echo
echo "Try it: marten tcp 3000"
