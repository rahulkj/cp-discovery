#!/bin/bash

set -e

VERSION=${VERSION:-"dev"}
COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS="-s -w -X main.version=${VERSION} -X main.commit=${COMMIT} -X main.date=${DATE}"

echo "Building cp-discovery binaries..."
echo "Version: ${VERSION}"
echo "Commit: ${COMMIT}"
echo "Date: ${DATE}"
echo ""

# Clean dist directory
rm -rf dist
mkdir -p dist

# Build for macOS (native builds)
echo "Building for macOS..."

echo "  - darwin/amd64"
GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build -ldflags="${LDFLAGS}" -o dist/cp-discovery_darwin_amd64/cp-discovery ./cmd/cp-discovery

echo "  - darwin/arm64"
GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 go build -ldflags="${LDFLAGS}" -o dist/cp-discovery_darwin_arm64/cp-discovery ./cmd/cp-discovery

# Build for Linux using Docker
echo ""
echo "Building for Linux (using Docker)..."

# Check if Docker is available
if ! command -v docker &> /dev/null; then
    echo "Docker not found. Skipping Linux builds."
    echo "Install Docker to build Linux binaries."
else
    # Build Linux AMD64
    echo "  - linux/amd64"
    mkdir -p dist/cp-discovery_linux_amd64
    docker run --rm \
        --platform linux/amd64 \
        -v "$PWD":/workspace \
        -w /workspace \
        -e CGO_ENABLED=1 \
        -e GOOS=linux \
        -e GOARCH=amd64 \
        golang:1.21 \
        sh -c "apt-get update -qq && apt-get install -y -qq librdkafka-dev > /dev/null 2>&1 && go build -ldflags='${LDFLAGS}' -o dist/cp-discovery_linux_amd64/cp-discovery ./cmd/cp-discovery"

    # Build Linux ARM64
    echo "  - linux/arm64"
    mkdir -p dist/cp-discovery_linux_arm64
    docker run --rm \
        --platform linux/arm64 \
        -v "$PWD":/workspace \
        -w /workspace \
        -e CGO_ENABLED=1 \
        -e GOOS=linux \
        -e GOARCH=arm64 \
        golang:1.21 \
        sh -c "apt-get update -qq && apt-get install -y -qq librdkafka-dev > /dev/null 2>&1 && go build -ldflags='${LDFLAGS}' -o dist/cp-discovery_linux_arm64/cp-discovery ./cmd/cp-discovery"
fi

# Build for Windows using Docker
echo ""
echo "Building for Windows (using Docker)..."

if command -v docker &> /dev/null; then
    # Build Windows AMD64
    echo "  - windows/amd64"
    mkdir -p dist/cp-discovery_windows_amd64
    docker run --rm \
        --platform linux/amd64 \
        -v "$PWD":/workspace \
        -w /workspace \
        -e CGO_ENABLED=1 \
        -e GOOS=windows \
        -e GOARCH=amd64 \
        -e CC=x86_64-w64-mingw32-gcc \
        -e CXX=x86_64-w64-mingw32-g++ \
        golang:1.21 \
        sh -c "apt-get update -qq && apt-get install -y -qq mingw-w64 wget tar > /dev/null 2>&1 && \
               wget -q https://github.com/confluentinc/librdkafka/archive/refs/tags/v2.3.0.tar.gz && \
               tar -xzf v2.3.0.tar.gz && \
               cd librdkafka-2.3.0 && \
               ./configure --host=x86_64-w64-mingw32 --prefix=/usr/local --enable-static --disable-shared > /dev/null 2>&1 && \
               make > /dev/null 2>&1 && make install > /dev/null 2>&1 && \
               cd /workspace && \
               go build -ldflags='${LDFLAGS}' -o dist/cp-discovery_windows_amd64/cp-discovery.exe ./cmd/cp-discovery" || echo "  Windows build failed - requires complex setup"

    # Build Windows ARM64
    echo "  - windows/arm64"
    echo "  Windows ARM64 build skipped (requires cross-compilation toolchain)"
else
    echo "  Docker not available - skipping Windows builds"
fi

echo ""
echo "Build complete! Binaries are in the dist/ directory:"
echo ""
find dist -type f -name "cp-discovery*" -exec ls -lh {} \; | awk '{print "  " $9, "(" $5 ")"}'
echo ""
