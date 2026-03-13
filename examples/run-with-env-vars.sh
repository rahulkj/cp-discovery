#!/bin/bash
# Example: Running CP Discovery with environment variables
# This script demonstrates how to set environment variables and run cp-discovery

# Exit on error
set -e

# Color output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${GREEN}CP Discovery - Environment Variables Example${NC}"
echo

# Check if .env file exists
if [ ! -f .env ]; then
    echo -e "${YELLOW}No .env file found. Creating from .env.example...${NC}"
    if [ -f .env.example ]; then
        cp .env.example .env
        echo -e "${RED}IMPORTANT: Edit .env file with your actual credentials!${NC}"
        echo -e "File location: .env"
        exit 1
    else
        echo -e "${RED}ERROR: .env.example not found${NC}"
        exit 1
    fi
fi

# Load environment variables from .env file
echo -e "${GREEN}Loading environment variables from .env...${NC}"
set -a  # automatically export all variables
source .env
set +a

# Verify required variables are set
REQUIRED_VARS=(
    "KAFKA_BOOTSTRAP_SERVERS"
    "CLUSTER_NAME"
)

missing_vars=()
for var in "${REQUIRED_VARS[@]}"; do
    if [ -z "${!var}" ]; then
        missing_vars+=("$var")
    fi
done

if [ ${#missing_vars[@]} -gt 0 ]; then
    echo -e "${RED}ERROR: Required environment variables not set:${NC}"
    printf '%s\n' "${missing_vars[@]}"
    echo
    echo "Please edit .env file and set these variables"
    exit 1
fi

# Display configuration (without sensitive values)
echo -e "${GREEN}Configuration:${NC}"
echo "  Cluster Name: ${CLUSTER_NAME}"
echo "  Kafka Bootstrap: ${KAFKA_BOOTSTRAP_SERVERS}"
echo "  Output Format: ${OUTPUT_FORMAT:-json}"
echo "  Output File: ${OUTPUT_FILE:-discovery-report.json}"
echo

# Check if cp-discovery binary exists
BINARY=""
if [ -f "bin/cp-discovery" ]; then
    BINARY="bin/cp-discovery"
elif [ -f "cp-discovery" ]; then
    BINARY="./cp-discovery"
elif command -v cp-discovery &> /dev/null; then
    BINARY="cp-discovery"
else
    # Try to find in dist directory
    DIST_BINARY=$(find dist -name "cp-discovery-*" -type f 2>/dev/null | head -n 1)
    if [ -n "$DIST_BINARY" ]; then
        BINARY="$DIST_BINARY"
    else
        echo -e "${RED}ERROR: cp-discovery binary not found${NC}"
        echo "Please build the project first:"
        echo "  make build"
        echo "  OR"
        echo "  goreleaser build --clean --snapshot"
        exit 1
    fi
fi

echo -e "${GREEN}Using binary: ${BINARY}${NC}"
echo

# Run cp-discovery with environment variables
echo -e "${GREEN}Running CP Discovery...${NC}"
$BINARY -config configs/config-env-vars.yaml

echo
echo -e "${GREEN}Discovery complete!${NC}"
echo "Results saved to: ${OUTPUT_FILE:-discovery-report.json}"
