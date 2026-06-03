#!/usr/bin/env bash
# Quick installer for driftmap
set -euo pipefail

REPO="DucChau/driftmap"
BINARY="driftmap"
INSTALL_DIR="${HOME}/.local/bin"

echo "Installing driftmap..."

# Ensure Go is available
if ! command -v go &>/dev/null; then
  echo "Error: Go is not installed. Visit https://go.dev/dl/" >&2
  exit 1
fi

mkdir -p "${INSTALL_DIR}"

go install "github.com/${REPO}@latest"

echo "driftmap installed to $(go env GOPATH)/bin/driftmap"
echo "Make sure \$(go env GOPATH)/bin is in your PATH."
