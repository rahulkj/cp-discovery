#!/bin/bash

set -e

VERSION=${VERSION:-"dev"}
COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS="-s -w -X main.version=${VERSION} -X main.commit=${COMMIT} -X main.date=${DATE}"

echo "=========================================="
echo "Building cp-discovery for all platforms"
echo "=========================================="
echo "Version: ${VERSION}"
echo "Commit: ${COMMIT}"
echo "Date: ${DATE}"
echo ""

# Clean dist directory
rm -rf dist
mkdir -p dist

#===============================
# macOS Builds (native)
#===============================
echo "Building macOS binaries (native)..."

echo "  ✓ darwin/amd64"
GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build -ldflags="${LDFLAGS}" -o dist/cp-discovery_darwin_amd64/cp-discovery ./cmd/cp-discovery 2>/dev/null || echo "    ✗ Failed"

echo "  ✓ darwin/arm64"
GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 go build -ldflags="${LDFLAGS}" -o dist/cp-discovery_darwin_arm64/cp-discovery ./cmd/cp-discovery 2>/dev/null || echo "    ✗ Failed"

#===============================
# Linux Builds (Docker)
#===============================
echo ""
echo "Building Linux binaries (Docker)..."

if ! command -v docker &> /dev/null; then
    echo "  ✗ Docker not found - skipping Linux builds"
else
    # Linux ARM64 (native to this Mac)
    echo "  ✓ linux/arm64 (building...)"
    docker build --platform linux/arm64 -f Dockerfile.build -t cp-discovery-builder-arm64 . > /dev/null 2>&1
    mkdir -p dist/cp-discovery_linux_arm64
    docker run --rm --platform linux/arm64 \
        -v "$PWD/dist:/dist" \
        -e CGO_ENABLED=1 \
        -e GOOS=linux \
        -e GOARCH=arm64 \
        cp-discovery-builder-arm64 \
        go build -ldflags="${LDFLAGS}" -o /dist/cp-discovery_linux_arm64/cp-discovery ./cmd/cp-discovery 2>/dev/null && echo "    ✓ linux/arm64 complete" || echo "    ✗ linux/arm64 failed"

    # Linux AMD64 (cross-platform - requires QEMU)
    echo "  ✓ linux/amd64 (building...)"
    docker build --platform linux/amd64 -f Dockerfile.build -t cp-discovery-builder-amd64 . > /dev/null 2>&1
    mkdir -p dist/cp-discovery_linux_amd64
    docker run --rm --platform linux/amd64 \
        -v "$PWD/dist:/dist" \
        -e CGO_ENABLED=1 \
        -e GOOS=linux \
        -e GOARCH=amd64 \
        cp-discovery-builder-amd64 \
        go build -ldflags="${LDFLAGS}" -o /dist/cp-discovery_linux_amd64/cp-discovery ./cmd/cp-discovery 2>/dev/null && echo "    ✓ linux/amd64 complete" || echo "    ✗ linux/amd64 failed"
fi

#===============================
# Windows Builds (skip - complex)
#===============================
echo ""
echo "Windows builds skipped (requires MinGW cross-compilation)"

#===============================
# Summary
#===============================
echo ""
echo "=========================================="
echo "Build Summary"
echo "=========================================="

for dir in dist/*/; do
    if [ -f "${dir}cp-discovery" ] || [ -f "${dir}cp-discovery.exe" ]; then
        binary="${dir}cp-discovery"
        [ ! -f "$binary" ] && binary="${dir}cp-discovery.exe"
        size=$(ls -lh "$binary" | awk '{print $5}')
        platform=$(basename "$dir")
        echo "✓ ${platform}: ${size}"
    fi
done

echo ""
echo "Build complete!"
