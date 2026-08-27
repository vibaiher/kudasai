#!/bin/sh
set -e

REPO="vibaiher/kudasai"
INSTALL_DIR="/usr/local/bin"

# Detect OS
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$OS" in
  linux) ;;
  darwin) ;;
  *)
    echo "Unsupported OS: $OS"
    exit 1
    ;;
esac

# Detect architecture
ARCH=$(uname -m)
case "$ARCH" in
  x86_64) ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *)
    echo "Unsupported architecture: $ARCH"
    exit 1
    ;;
esac

# Get latest release tag
LATEST=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST" ]; then
  echo "Failed to fetch latest release"
  exit 1
fi

VERSION=${LATEST#v}
TARBALL="kudasai_${OS}_${ARCH}.tar.gz"
CHECKSUMS="kudasai_${VERSION}_checksums.txt"
URL="https://github.com/${REPO}/releases/download/${LATEST}/${TARBALL}"
CHECKSUMS_URL="https://github.com/${REPO}/releases/download/${LATEST}/${CHECKSUMS}"

echo "Downloading kudasai ${LATEST} for ${OS}/${ARCH}..."

TMPDIR=$(mktemp -d)
trap 'rm -rf "$TMPDIR"' EXIT

curl -fsSL "$URL" -o "${TMPDIR}/${TARBALL}"
curl -fsSL "$CHECKSUMS_URL" -o "${TMPDIR}/${CHECKSUMS}"

if command -v sha256sum >/dev/null 2>&1; then
  ACTUAL=$(sha256sum "${TMPDIR}/${TARBALL}" | cut -d' ' -f1)
elif command -v shasum >/dev/null 2>&1; then
  ACTUAL=$(shasum -a 256 "${TMPDIR}/${TARBALL}" | cut -d' ' -f1)
else
  echo "Cannot verify the download: neither sha256sum nor shasum is available"
  exit 1
fi

EXPECTED=$(awk -v file="$TARBALL" '$2 == file { print $1 }' "${TMPDIR}/${CHECKSUMS}")

if [ -z "$EXPECTED" ]; then
  echo "No checksum published for ${TARBALL} in ${CHECKSUMS}"
  exit 1
fi

if [ "$ACTUAL" != "$EXPECTED" ]; then
  echo "Checksum mismatch for ${TARBALL}"
  echo "  expected: ${EXPECTED}"
  echo "  actual:   ${ACTUAL}"
  exit 1
fi

echo "Checksum verified"

tar -xzf "${TMPDIR}/${TARBALL}" -C "$TMPDIR"

if [ -w "$INSTALL_DIR" ]; then
  mv "${TMPDIR}/kudasai" "${INSTALL_DIR}/kudasai"
else
  echo "Installing to ${INSTALL_DIR} (requires sudo)..."
  sudo mv "${TMPDIR}/kudasai" "${INSTALL_DIR}/kudasai"
fi

echo "kudasai ${LATEST} installed to ${INSTALL_DIR}/kudasai"
