#!/bin/bash

# Control Center API Test Script
# This script queries all Control Center v2 API endpoints to show discoverable data

set -e

# Configuration
C3_URL="${C3_URL:-http://localhost:9021}"
C3_USER="${C3_USER:-}"
C3_PASS="${C3_PASS:-}"

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Function to make authenticated curl request
c3_curl() {
    local endpoint="$1"
    local url="${C3_URL}${endpoint}"

    if [ -n "$C3_USER" ]; then
        curl -s -u "${C3_USER}:${C3_PASS}" "$url"
    else
        curl -s "$url"
    fi
}

# Function to print section header
print_header() {
    echo -e "\n${BLUE}========================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}========================================${NC}"
}

# Function to print subsection
print_subsection() {
    echo -e "\n${GREEN}--- $1 ---${NC}"
}

# Check if jq is installed
if ! command -v jq &> /dev/null; then
    echo -e "${RED}Error: jq is required but not installed${NC}"
    echo "Install with: brew install jq (macOS) or apt-get install jq (Linux)"
    exit 1
fi

# Test Control Center connectivity
print_header "CONTROL CENTER CONNECTIVITY TEST"
echo "Testing connection to: $C3_URL"
if c3_curl "/health" > /dev/null 2>&1; then
    echo -e "${GREEN}✓ Control Center is accessible${NC}"
else
    echo -e "${RED}✗ Cannot connect to Control Center${NC}"
    exit 1
fi

# Get version
print_subsection "Control Center Version"
c3_curl "/api/version" | jq -r '.version // "Unknown"'

# Discover Kafka Clusters
print_header "KAFKA CLUSTERS"
kafka_clusters=$(c3_curl "/2.0/clusters/kafka")
kafka_count=$(echo "$kafka_clusters" | jq 'length')
echo "Total Kafka Clusters: $kafka_count"

echo "$kafka_clusters" | jq -r '.[] | "\nCluster: \(.clusterName) (\(.clusterId))\n  Brokers: \(.brokerCount)\n  Topics: \(.topicCount)\n  Partitions: \(.partitionCount)"'

# Get detailed broker information for each cluster
if [ "$kafka_count" -gt 0 ]; then
    print_subsection "Kafka Broker Details"
    echo "$kafka_clusters" | jq -r '.[].clusterId' | while read -r cluster_id; do
        echo -e "\n${YELLOW}Cluster: $cluster_id${NC}"
        cluster_detail=$(c3_curl "/2.0/clusters/kafka/$cluster_id")
        echo "$cluster_detail" | jq -r '.brokers[]? | "  Broker \(.id): \(.host):\(.port)"' 2>/dev/null || echo "  (No broker details available)"

        # Get cluster health
        health=$(c3_curl "/2.0/clusters/kafka/$cluster_id/health" 2>/dev/null)
        if [ -n "$health" ]; then
            status=$(echo "$health" | jq -r '.status // "Unknown"')
            echo "  Health: $status"
        fi
    done
fi

# Discover Kafka Connect Clusters
print_header "KAFKA CONNECT CLUSTERS"
connect_clusters=$(c3_curl "/2.0/clusters/connect")
connect_count=$(echo "$connect_clusters" | jq 'length')
echo "Total Connect Clusters: $connect_count"

if [ "$connect_count" -gt 0 ]; then
    echo "$connect_clusters" | jq -r '.[].clusterId' | while read -r cluster_id; do
        cluster_name=$(echo "$connect_clusters" | jq -r ".[] | select(.clusterId == \"$cluster_id\") | .clusterName")
        kafka_cluster=$(echo "$connect_clusters" | jq -r ".[] | select(.clusterId == \"$cluster_id\") | .kafkaClusterId")

        echo -e "\n${YELLOW}Connect Cluster: $cluster_name${NC}"
        echo "  Cluster ID: $cluster_id"
        echo "  Kafka Cluster: $kafka_cluster"

        # Get worker nodes
        workers=$(c3_curl "/2.0/clusters/connect/$cluster_id/workers" 2>/dev/null)
        if [ -n "$workers" ]; then
            worker_count=$(echo "$workers" | jq 'length')
            echo "  Worker Nodes: $worker_count"
            echo "$workers" | jq -r '.[]? | "    - \(.workerId)"' 2>/dev/null
        fi

        # Get connectors
        connectors=$(c3_curl "/2.0/clusters/connect/$cluster_id/connectors" 2>/dev/null)
        if [ -n "$connectors" ]; then
            total_connectors=$(echo "$connectors" | jq 'length')
            source_count=$(echo "$connectors" | jq '[.[] | select(.type == "source")] | length')
            sink_count=$(echo "$connectors" | jq '[.[] | select(.type == "sink")] | length')
            running_count=$(echo "$connectors" | jq '[.[] | select(.state == "RUNNING")] | length')
            failed_count=$(echo "$connectors" | jq '[.[] | select(.state == "FAILED")] | length')

            echo "  Connectors: $total_connectors (Source: $source_count, Sink: $sink_count)"
            echo "  Status: Running: $running_count, Failed: $failed_count"

            print_subsection "Connector Details"
            echo "$connectors" | jq -r '.[] | "    \(.name) [\(.type)]: \(.state) (Tasks: \(.tasks))"'
        fi
    done
fi

# Discover Schema Registry Clusters
print_header "SCHEMA REGISTRY CLUSTERS"
sr_clusters=$(c3_curl "/2.0/clusters/schema-registry")
sr_count=$(echo "$sr_clusters" | jq 'length')
echo "Total Schema Registry Clusters: $sr_count"

if [ "$sr_count" -gt 0 ]; then
    echo "$sr_clusters" | jq -r '.[].clusterId' | while read -r cluster_id; do
        cluster_name=$(echo "$sr_clusters" | jq -r ".[] | select(.clusterId == \"$cluster_id\") | .clusterName")
        kafka_cluster=$(echo "$sr_clusters" | jq -r ".[] | select(.clusterId == \"$cluster_id\") | .kafkaClusterId")

        echo -e "\n${YELLOW}Schema Registry: $cluster_name${NC}"
        echo "  Cluster ID: $cluster_id"
        echo "  Kafka Cluster: $kafka_cluster"

        # Get Schema Registry details
        sr_detail=$(c3_curl "/2.0/clusters/schema-registry/$cluster_id" 2>/dev/null)
        if [ -n "$sr_detail" ]; then
            version=$(echo "$sr_detail" | jq -r '.version // "Unknown"')
            mode=$(echo "$sr_detail" | jq -r '.mode // "Unknown"')
            schema_count=$(echo "$sr_detail" | jq -r '.subjectCount // 0')
            node_count=$(echo "$sr_detail" | jq -r '.nodeCount // 0')

            echo "  Version: $version"
            echo "  Mode: $mode"
            echo "  Schema Count: $schema_count"
            echo "  Node Count: $node_count"

            # Try to get nodes array
            nodes=$(echo "$sr_detail" | jq -r '.nodes[]? | "    - \(.host):\(.port)"' 2>/dev/null)
            if [ -n "$nodes" ]; then
                echo "  Nodes:"
                echo "$nodes"
            fi

            # List subjects if available
            subjects=$(echo "$sr_detail" | jq -r '.subjects[]?' 2>/dev/null)
            if [ -n "$subjects" ]; then
                subject_list=$(echo "$subjects" | head -n 5)
                echo "  Subjects (first 5):"
                echo "$subject_list" | while read -r subject; do
                    echo "    - $subject"
                done
            fi
        fi

        # Try alternative nodes endpoint
        sr_nodes=$(c3_curl "/2.0/clusters/schema-registry/$cluster_id/nodes" 2>/dev/null)
        if [ -n "$sr_nodes" ]; then
            alt_node_count=$(echo "$sr_nodes" | jq 'length')
            if [ "$alt_node_count" -gt 0 ]; then
                echo "  Nodes (from /nodes endpoint): $alt_node_count"
            fi
        fi
    done
fi

# Discover ksqlDB Clusters
print_header "KSQLDB CLUSTERS"
ksql_clusters=$(c3_curl "/2.0/clusters/ksql")
ksql_count=$(echo "$ksql_clusters" | jq 'length')
echo "Total ksqlDB Clusters: $ksql_count"

if [ "$ksql_count" -gt 0 ]; then
    echo "$ksql_clusters" | jq -r '.[].clusterId' | while read -r cluster_id; do
        cluster_name=$(echo "$ksql_clusters" | jq -r ".[] | select(.clusterId == \"$cluster_id\") | .clusterName")
        kafka_cluster=$(echo "$ksql_clusters" | jq -r ".[] | select(.clusterId == \"$cluster_id\") | .kafkaClusterId")

        echo -e "\n${YELLOW}ksqlDB Cluster: $cluster_name${NC}"
        echo "  Cluster ID: $cluster_id"
        echo "  Kafka Cluster: $kafka_cluster"

        # Get ksqlDB details
        ksql_detail=$(c3_curl "/2.0/clusters/ksql/$cluster_id" 2>/dev/null)
        if [ -n "$ksql_detail" ]; then
            query_count=$(echo "$ksql_detail" | jq -r '.queryCount // 0')
            stream_count=$(echo "$ksql_detail" | jq -r '.streamCount // 0')
            table_count=$(echo "$ksql_detail" | jq -r '.tableCount // 0')
            node_count=$(echo "$ksql_detail" | jq -r '.nodeCount // 0')

            echo "  Queries: $query_count"
            echo "  Streams: $stream_count"
            echo "  Tables: $table_count"
            echo "  Node Count: $node_count"

            # Try to get servers array
            servers=$(echo "$ksql_detail" | jq -r '.servers[]? | "    - \(.host):\(.port)"' 2>/dev/null)
            if [ -n "$servers" ]; then
                echo "  Servers:"
                echo "$servers"
            fi
        fi

        # Try alternative servers endpoint
        ksql_servers=$(c3_curl "/2.0/clusters/ksql/$cluster_id/servers" 2>/dev/null)
        if [ -n "$ksql_servers" ]; then
            alt_server_count=$(echo "$ksql_servers" | jq 'length')
            if [ "$alt_server_count" -gt 0 ]; then
                echo "  Servers (from /servers endpoint): $alt_server_count"
            fi
        fi
    done
fi

# Get Consumer Group Lag
print_header "CONSUMER GROUP LAG"
lag_data=$(c3_curl "/2.0/monitoring/consumer-groups/lag" 2>/dev/null)
if [ -n "$lag_data" ]; then
    total_lag=$(echo "$lag_data" | jq -r '.totalLag // 0')
    echo "Total Consumer Lag: $total_lag"

    # Try to get per-group lag
    group_count=$(echo "$lag_data" | jq 'length' 2>/dev/null || echo "0")
    if [ "$group_count" != "null" ] && [ "$group_count" -gt 0 ]; then
        echo "Consumer Groups: $group_count"
        echo "$lag_data" | jq -r '.[] | "  \(.groupId): Lag = \(.lag)"' 2>/dev/null | head -n 10
    fi
fi

# Summary
print_header "DISCOVERY SUMMARY"
echo "Kafka Clusters: $kafka_count"
echo "Connect Clusters: $connect_count"
echo "Schema Registry Clusters: $sr_count"
echo "ksqlDB Clusters: $ksql_count"
echo ""
echo -e "${GREEN}✓ Discovery complete!${NC}"
echo ""
echo "To get JSON output for all endpoints, run:"
echo "  export C3_URL=$C3_URL"
echo "  curl -s \$C3_URL/2.0/clusters/kafka | jq ."
