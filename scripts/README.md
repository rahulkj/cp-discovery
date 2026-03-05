# Build Scripts

This directory contains build scripts for cross-platform compilation of the cp-discovery tool.

## Available Scripts

### build-all.sh

Comprehensive build script that builds binaries for all supported platforms.

**Platforms:**
- macOS (darwin/amd64, darwin/arm64) - native builds
- Linux (linux/amd64, linux/arm64) - via Docker

**Usage:**
```bash
./scripts/build-all.sh
```

**Output:**
Binaries are created in `dist/` directory with the following structure:
```
dist/
├── cp-discovery_darwin_amd64/cp-discovery
├── cp-discovery_darwin_arm64/cp-discovery
├── cp-discovery_linux_amd64/cp-discovery
└── cp-discovery_linux_arm64/cp-discovery
```

**Requirements:**
- Go 1.21+ for macOS builds
- Docker for Linux builds
- librdkafka (installed via `brew install librdkafka` on macOS)

**Environment Variables:**
- `VERSION` - Version string (default: "dev")
- Automatically captures git commit and build date

---

### build.sh

Builds binaries for macOS and Linux platforms, including experimental Windows support.

**Platforms:**
- macOS (darwin/amd64, darwin/arm64)
- Linux (linux/amd64, linux/arm64) - via Docker
- Windows (windows/amd64) - experimental, may fail

**Usage:**
```bash
./scripts/build.sh
```

**Output:**
Same structure as `build-all.sh` in the `dist/` directory.

**Requirements:**
- Go 1.21+
- Docker (for Linux and Windows builds)
- librdkafka for macOS

**Note:**
Windows builds require MinGW cross-compilation toolchain and may not succeed without proper setup.

---

### build-linux.sh

Docker-based Linux-only build script using a custom build image.

**Platforms:**
- Linux (linux/amd64, linux/arm64)

**Usage:**
```bash
./scripts/build-linux.sh
```

**Output:**
```
dist/
├── cp-discovery_linux_amd64/cp-discovery
└── cp-discovery_linux_arm64/cp-discovery
```

**Requirements:**
- Docker
- Dockerfile.build (builder image definition)

**Note:**
This script first builds a Docker image (`cp-discovery-builder`) and then uses it to compile the binaries.

---

## Build Versioning

All build scripts support version information through environment variables:

```bash
VERSION=v1.0.0 ./scripts/build-all.sh
```

This sets:
- `main.version` - from `VERSION` env var
- `main.commit` - from git commit hash
- `main.date` - build timestamp

The binaries will display this information when run with `--version` (if implemented).

## Comparison with Other Build Methods

| Method | Use Case | Platforms | Output Location |
|--------|----------|-----------|----------------|
| `make build` | Development | Current platform only | `bin/` |
| `make build-release` | Optimized single build | Current platform only | `bin/` |
| `./scripts/build-all.sh` | Multi-platform builds | macOS + Linux | `dist/` |
| `./scripts/build.sh` | Multi-platform + Windows | macOS + Linux + Windows | `dist/` |
| `./scripts/build-linux.sh` | Linux only | Linux only | `dist/` |
| `goreleaser` | Production releases | macOS + Linux | `dist/` |

## Troubleshooting

### Docker Not Found

If you see "Docker not found" errors:
```bash
# Install Docker
brew install --cask docker  # macOS
```

### librdkafka Missing (macOS)

```bash
brew install librdkafka
```

### Permission Denied

Make scripts executable:
```bash
chmod +x scripts/*.sh
```

### Cross-Platform Build Failures

Linux builds require Docker with multi-platform support. If ARM64 builds fail:
```bash
# Enable Docker multi-platform support
docker buildx create --use
docker buildx inspect --bootstrap
```
