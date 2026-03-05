# Multi-stage build for Confluent Platform Discovery Tool

# Stage 1: Build
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache \
    gcc \
    g++ \
    make \
    musl-dev \
    librdkafka-dev \
    pkgconfig \
    git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY cmd/ ./cmd/
COPY internal/ ./internal/

# Build the binary
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o cp-discovery ./cmd/cp-discovery

# Stage 2: Runtime
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache \
    librdkafka \
    ca-certificates \
    tzdata

# Create non-root user
RUN addgroup -g 1000 discovery && \
    adduser -D -u 1000 -G discovery discovery

# Set working directory
WORKDIR /home/discovery

# Copy binary from builder
COPY --from=builder /app/cp-discovery /usr/local/bin/cp-discovery

# Copy configuration files
COPY configs/ ./configs/

# Create output directory
RUN mkdir -p /home/discovery/reports && \
    chown -R discovery:discovery /home/discovery

# Switch to non-root user
USER discovery

# Set default command
CMD ["cp-discovery"]
