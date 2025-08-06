#!/bin/bash
set -e

cd "$(dirname "$0")/.."

# Release build script for DMX Viewer (Go backend + React frontend)
# Outputs: release/ directory with all binaries and frontend assets
# IMPORTANT: Frontend must be built BEFORE backend, because Go uses go:embed to include frontend files.
# If frontend files are missing, Go build will fail.

RELEASE_DIR="release"
GO_MAIN="backend/cmd/dmx_viewer/main.go"
FRONTEND_EMBED_DIR="backend/internal/app/embed_static"

# 1. Create release directory
mkdir -p "$RELEASE_DIR"

# 2. Build frontend first (required for Go embed)
echo "Building frontend..."
cd frontend
NODE_ENV=production npm run build -- --mode production
cd ..

# 3. Check that frontend build output exists
if [ ! -f "$FRONTEND_EMBED_DIR/index.html" ]; then
  echo "ERROR: Frontend build failed or output missing: $FRONTEND_EMBED_DIR/index.html"
  exit 1
fi

# 4. Go backend cross-builds
PLATFORMS=(
  "windows amd64 dmx_viewer_windows_amd64.exe"
  "linux amd64 dmx_viewer_linux_amd64"
  "linux arm64 dmx_viewer_linux_arm64"
  "darwin amd64 dmx_viewer_mac_amd64"
  "darwin arm64 dmx_viewer_mac_arm64"
)

for entry in "${PLATFORMS[@]}"; do
  set -- $entry
  GOOS=$1
  GOARCH=$2
  OUT=$3
  echo "Building backend for $GOOS/$GOARCH -> $OUT"
  (cd backend && env GOOS=$GOOS GOARCH=$GOARCH CGO_ENABLED=0 go build -o "../$RELEASE_DIR/$OUT" "cmd/dmx_viewer/main.go")
done

# 5. Done

echo "Release build completed!"
echo "Backend binaries:"
ls -lh "$RELEASE_DIR" | grep dmx_viewer

echo "Frontend assets (embedded in Go binary):"
ls -lh "$FRONTEND_EMBED_DIR"

echo "All artifacts are in $RELEASE_DIR"
