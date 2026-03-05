#!/bin/bash

set -e

VERSION=${VERSION:-"dev"}
COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS="-s -w -X main.version=${VERSION} -X main.commit=${COMMIT} -X main.date=${DATE}"

echo "Building Linux binaries using Docker..."

# Build Docker image
docker build -f scripts/Dockerfile.build -t cp-discovery-builder .

# Build for Linux AMD64
echo "Building for linux/amd64..."
mkdir -p dist/cp-discovery_linux_amd64
docker run --rm \
    -v "$PWD/dist:/dist" \
    -e CGO_ENABLED=1 \
    -e GOOS=linux \
    -e GOARCH=amd64 \
    cp-discovery-builder \
    go build -ldflags="${LDFLAGS}" -o /dist/cp-discovery_linux_amd64/cp-discovery ./cmd/cp-discovery

# Build for Linux ARM64
echo "Building for linux/arm64..."
mkdir -p dist/cp-discovery_linux_arm64
docker run --rm \
    --platform linux/arm64 \
    -v "$PWD/dist:/dist" \
    -e CGO_ENABLED=1 \
    -e GOOS=linux \
    -e GOARCH=arm64 \
    cp-discovery-builder \
    go build -ldflags="${LDFLAGS}" -o /dist/cp-discovery_linux_arm64/cp-discovery ./cmd/cp-discovery

echo "Linux builds complete!"
ls -lh dist/cp-discovery_linux_*/cp-discovery
