# Confluent Platform API Reference

Complete reference for all API endpoints used by cp-discovery for component discovery.

## Table of Contents

1. [Control Center API](#control-center-api)
2. [REST Proxy API](#rest-proxy-api)
3. [Component API Summary](#component-api-summary)
4. [Discovery Capabilities Matrix](#discovery-capabilities-matrix)

---

## Control Center API

Control Center provides a centralized monitoring API (v2.0) that can discover multiple Confluent Platform components.

### Base URL
```
http://control-center-host:9021/2.0
```

### Authentication
All endpoints support Basic Authentication:
```bash
curl -u username:password http://control-center:9021/2.0/clusters/kafka
```

---

### Kafka Clusters

**List all Kafka clusters:**
```http
GET /2.0/clusters/kafka
```

**Response:**
```json
[
  {
    "clusterId": "abc123",
    "clusterName": "production",
    "brokerCount": 3,
    "topicCount": 150,
    "partitionCount": 450
  }
]
```

**Get cluster health:**
```http
GET /2.0/clusters/kafka/{clusterId}/health
```

**Discovered Information:**
- ✅ Broker count
- ✅ Topic count
- ✅ Partition count
- ✅ Health status
- ✅ Cluster ID and name

---

### Kafka Connect

**List all Connect clusters:**
```http
GET /2.0/clusters/connect
```

**Get connector details:**
```http
GET /2.0/clusters/connect/{clusterId}/connectors
```

**Get worker nodes:**
```http
GET /2.0/clusters/connect/{clusterId}/workers
```

**Response:**
```json
[
  {
    "workerId": "connect-worker-1:8083"
  },
  {
    "workerId": "connect-worker-2:8083"
  }
]
```

**Discovered Information:**
- ✅ **Worker node count** - Number of Connect workers
- ✅ Connector count (total, source, sink)
- ✅ Connector details (name, type, state, tasks)
- ✅ Running/failed connector counts
- ✅ Kafka cluster association

---

### Schema Registry

**List all Schema Registry clusters:**
```http
GET /2.0/clusters/schema-registry
```

**Get cluster details:**
```http
GET /2.0/clusters/schema-registry/{clusterId}
```

**Response:**
```json
{
  "version": "7.6.0",
  "subjectCount": 75,
  "mode": "READWRITE",
  "nodeCount": 3,
  "nodes": [
    {"id": "sr-1", "host": "sr-1.example.com", "port": 8081},
    {"id": "sr-2", "host": "sr-2.example.com", "port": 8081},
    {"id": "sr-3", "host": "sr-3.example.com", "port": 8081}
  ]
}
```

**Get nodes (alternative):**
```http
GET /2.0/clusters/schema-registry/{clusterId}/nodes
```

**Discovered Information:**
- ✅ **Node count** - Number of Schema Registry nodes
- ✅ Node IDs, hosts, and ports
- ✅ Schema/subject count
- ✅ Version and mode
- ✅ Subject list
- ✅ Kafka cluster association

---

### ksqlDB

**List all ksqlDB clusters:**
```http
GET /2.0/clusters/ksql
```

**Get cluster details:**
```http
GET /2.0/clusters/ksql/{clusterId}
```

**Response:**
```json
{
  "queryCount": 8,
  "streamCount": 15,
  "tableCount": 10,
  "nodeCount": 3,
  "servers": [
    {"id": "ksql-1", "host": "ksql-1.example.com", "port": 8088},
    {"id": "ksql-2", "host": "ksql-2.example.com", "port": 8088},
    {"id": "ksql-3", "host": "ksql-3.example.com", "port": 8088}
  ]
}
```

**Get servers (alternative):**
```http
GET /2.0/clusters/ksql/{clusterId}/servers
```

**Discovered Information:**
- ✅ **Node/server count** - Number of ksqlDB servers
- ✅ Server IDs, hosts, and ports
- ✅ Query count
- ✅ Stream count
- ✅ Table count
- ✅ Kafka cluster association

---

### Consumer Groups

**Get consumer lag:**
```http
GET /2.0/monitoring/consumer-groups/lag
```

**Response:**
```json
{
  "totalLag": 1250
}
```

Or per-group:
```json
[
  {"groupId": "order-processor", "lag": 500},
  {"groupId": "analytics", "lag": 750}
]
```

**Discovered Information:**
- ✅ Total consumer lag
- ✅ Per-group lag
- ✅ Consumer group IDs

---

### Control Center Limitations

**Cannot discover via Control Center:**
- ❌ REST Proxy (no API endpoints)
- ❌ Prometheus (external monitoring)
- ❌ Alertmanager (external monitoring)
- ❌ Control Center itself (limited self-discovery)

---

## REST Proxy API

REST Proxy provides HTTP access to Kafka Admin API. **Only discovers Kafka cluster components.**

### Base URL
```
http://rest-proxy-host:8082/v3
```

### Authentication
Supports Basic Authentication, Bearer Token, and API Key:
```bash
curl -u username:password http://rest-proxy:8082/v3/clusters
```

---

### Kafka Cluster Info

**List clusters:**
```http
GET /v3/clusters
```

**Response:**
```json
{
  "data": [
    {
      "cluster_id": "abc123",
      "controller": {
        "related": "/v3/clusters/abc123/brokers/1"
      }
    }
  ]
}
```

**Get cluster details:**
```http
GET /v3/clusters/{cluster_id}
```

**Discovered Information:**
- ✅ Cluster ID
- ✅ Controller information
- ✅ Cluster metadata

---

### Brokers

**List brokers:**
```http
GET /v3/clusters/{cluster_id}/brokers
```

**Response:**
```json
{
  "data": [
    {
      "broker_id": 1,
      "host": "broker-1.example.com",
      "port": 9092,
      "rack": "rack-1"
    },
    {
      "broker_id": 2,
      "host": "broker-2.example.com",
      "port": 9092,
      "rack": "rack-2"
    }
  ]
}
```

**Get broker configs:**
```http
GET /v3/clusters/{cluster_id}/brokers/{broker_id}/configs
```

**Discovered Information:**
- ✅ **Broker count and IDs**
- ✅ **Broker hosts and ports**
- ✅ **Controller ID and mode** (KRaft/ZooKeeper)
- ✅ **Controller count** (KRaft only)
- ✅ Rack assignment
- ✅ Process roles (KRaft)
- ✅ Broker configurations

---

### Topics and Partitions

**List topics:**
```http
GET /v3/clusters/{cluster_id}/topics
```

**Get topic partitions:**
```http
GET /v3/clusters/{cluster_id}/topics/{topic}/partitions
```

**Response:**
```json
{
  "data": [
    {
      "partition_id": 0,
      "leader": {"broker_id": 1},
      "replicas": [
        {"broker_id": 1},
        {"broker_id": 2},
        {"broker_id": 3}
      ],
      "isr": [
        {"broker_id": 1},
        {"broker_id": 2},
        {"broker_id": 3}
      ]
    }
  ]
}
```

**Get topic configs:**
```http
GET /v3/clusters/{cluster_id}/topics/{topic}/configs
```

**Discovered Information:**
- ✅ Topic count (total, internal, external)
- ✅ Partition count (total)
- ✅ **Partition leader broker** (per partition)
- ✅ **Replica broker IDs** (per partition)
- ✅ **ISR (In-Sync Replicas)** (per partition)
- ✅ Replication factor
- ✅ Topic configurations (retention, compression, etc.)

---

### Consumer Groups

**List consumer groups:**
```http
GET /v3/clusters/{cluster_id}/consumer-groups
```

**Response:**
```json
{
  "data": [
    {
      "consumer_group_id": "order-processor",
      "state": "Stable",
      "partition_assignor": "range"
    }
  ]
}
```

**Get group members:**
```http
GET /v3/clusters/{cluster_id}/consumer-groups/{group_id}/consumers
```

**Response:**
```json
{
  "data": [
    {
      "consumer_id": "consumer-1",
      "instance_id": "instance-1",
      "client_id": "client-1"
    }
  ]
}
```

**Discovered Information:**
- ✅ Consumer group count (total, active)
- ✅ **Consumer member count** (per group)
- ✅ Consumer group state
- ✅ Partition assignor strategy
- ✅ Consumer IDs

---

### ACLs (Access Control Lists)

**List ACLs:**
```http
GET /v3/clusters/{cluster_id}/acls
```

**Response:**
```json
{
  "data": [
    {
      "resource_type": "TOPIC",
      "resource_name": "orders",
      "pattern_type": "LITERAL",
      "principal": "User:alice",
      "operation": "READ",
      "permission": "ALLOW"
    }
  ]
}
```

**Discovered Information:**
- ✅ ACL count
- ✅ Resource types and names
- ✅ Principals and operations
- ✅ Permissions (ALLOW/DENY)

---

### Security Configuration

**Get broker configs:**
```http
GET /v3/clusters/{cluster_id}/brokers/{broker_id}/configs
```

Extracts security-related configs:
- `sasl.enabled.mechanisms`
- `listener.security.protocol.map`
- `security.inter.broker.protocol`

**Discovered Information:**
- ✅ SASL mechanisms (PLAIN, SCRAM-SHA-256, etc.)
- ✅ Security protocols (PLAINTEXT, SSL, SASL_SSL, etc.)
- ✅ SSL/TLS enabled status
- ✅ SASL enabled status

---

### Cluster Configuration

**Get cluster configs:**
```http
GET /v3/clusters/{cluster_id}/broker-configs
```

**Discovered Information:**
- ✅ Important cluster settings
- ✅ Retention policies
- ✅ Replication defaults
- ✅ Topic auto-creation settings

---

### REST Proxy Limitations

**Cannot discover via REST Proxy:**
- ❌ Schema Registry nodes
- ❌ Kafka Connect workers
- ❌ ksqlDB servers
- ❌ Control Center instances
- ❌ Prometheus/Alertmanager
- ❌ REST Proxy instances (no self-discovery)

**REST Proxy ONLY discovers Kafka cluster components.**

---

## Component API Summary

### Direct Component APIs

These components have their own APIs and must be queried directly:

#### Schema Registry (Port 8081)
```http
GET /
GET /subjects
GET /mode
GET /v1/metadata/id
```

#### Kafka Connect (Port 8083)
```http
GET /
GET /connectors
GET /connectors?expand=status
GET /connectors?expand=info
```

#### ksqlDB (Port 8088)
```http
GET /info
POST /ksql (SHOW QUERIES, SHOW STREAMS, SHOW TABLES)
```

#### Prometheus (Port 9090)
```http
GET /api/v1/status/buildinfo
GET /api/v1/targets
GET /api/v1/query
```

#### Alertmanager (Port 9093)
```http
GET /api/v2/status
GET /api/v2/alerts
```

---

## Discovery Capabilities Matrix

### Complete Node Count Discovery

| Component | Control Center API | REST Proxy API | Direct API |
|-----------|-------------------|----------------|------------|
| **Kafka Brokers** | ✅ Count only | ✅ **Full details** | N/A |
| **Kafka Controllers** | ✅ Count only | ✅ **ID, mode, count** | N/A |
| **Kafka Partitions** | ✅ Count only | ✅ **Leaders, replicas, ISR** | N/A |
| **Kafka Connect Workers** | ✅ **Full details** | ❌ | ✅ Limited |
| **Schema Registry Nodes** | ✅ **Full details** | ❌ | ✅ Limited |
| **ksqlDB Servers** | ✅ **Full details** | ❌ | ✅ Limited |
| **Consumer Groups** | ✅ Lag only | ✅ **Members, state** | N/A |
| **REST Proxy Instances** | ❌ | ❌ | ⚠️ No cluster API |
| **Control Center Nodes** | ⚠️ Limited | ❌ | ⚠️ Health only |
| **Prometheus Nodes** | ❌ | ❌ | ✅ Direct |
| **Alertmanager Nodes** | ❌ | ❌ | ✅ Direct |

### Detailed Component Information

| Information Type | Control Center | REST Proxy | Best Source |
|-----------------|----------------|------------|-------------|
| Kafka broker topology | Basic | ✅ **Detailed** | **REST Proxy** |
| Partition assignments | Basic | ✅ **Full** | **REST Proxy** |
| Connect worker count | ✅ **Yes** | ❌ | **Control Center** |
| Connect connectors | ✅ **Yes** | ❌ | **Control Center** or Direct |
| Schema Registry nodes | ✅ **Yes** | ❌ | **Control Center** |
| Schema count | ✅ **Yes** | ❌ | **Control Center** or Direct |
| ksqlDB servers | ✅ **Yes** | ❌ | **Control Center** |
| ksqlDB queries | ✅ **Yes** | ❌ | **Control Center** or Direct |
| Consumer groups | Lag only | ✅ **Members** | **REST Proxy** |
| ACLs | ❌ | ✅ **Full** | **REST Proxy** |
| Security config | ❌ | ✅ **Full** | **REST Proxy** |

---

## Recommended Discovery Strategy

### For Complete Platform Topology:

1. **Use Control Center API** (if available):
   - Kafka cluster counts
   - Connect worker counts
   - Schema Registry node counts
   - ksqlDB server counts
   - Connector and schema inventory

2. **Use REST Proxy API** (for Kafka details):
   - Complete partition topology
   - Leader and replica assignments
   - Consumer group members
   - ACLs and security config

3. **Use Direct APIs** (for components not in C3):
   - Prometheus/Alertmanager
   - Additional component-specific details

### For Kafka-Only Discovery:

Use **REST Proxy API** exclusively - it provides complete Kafka cluster topology including brokers, partitions, replicas, consumers, and security configuration.

---

## Query Examples

### Get Complete Kafka Topology via REST Proxy

```bash
# 1. Get cluster ID
CLUSTER_ID=$(curl -s http://rest-proxy:8082/v3/clusters | jq -r '.data[0].cluster_id')

# 2. Get all brokers
curl -s http://rest-proxy:8082/v3/clusters/$CLUSTER_ID/brokers | jq '.data'

# 3. Get partition assignments for a topic
curl -s http://rest-proxy:8082/v3/clusters/$CLUSTER_ID/topics/orders/partitions | jq '.data'

# 4. Get consumer groups
curl -s http://rest-proxy:8082/v3/clusters/$CLUSTER_ID/consumer-groups | jq '.data'
```

### Get All Component Nodes via Control Center

```bash
# 1. Get Kafka clusters
curl -u admin:secret http://c3:9021/2.0/clusters/kafka | jq '.'

# 2. Get Connect clusters with workers
CONNECT_ID=$(curl -s -u admin:secret http://c3:9021/2.0/clusters/connect | jq -r '.[0].clusterId')
curl -u admin:secret http://c3:9021/2.0/clusters/connect/$CONNECT_ID/workers | jq 'length'

# 3. Get Schema Registry with nodes
SR_ID=$(curl -s -u admin:secret http://c3:9021/2.0/clusters/schema-registry | jq -r '.[0].clusterId')
curl -u admin:secret http://c3:9021/2.0/clusters/schema-registry/$SR_ID | jq '.nodeCount'

# 4. Get ksqlDB servers
KSQL_ID=$(curl -s -u admin:secret http://c3:9021/2.0/clusters/ksql | jq -r '.[0].clusterId')
curl -u admin:secret http://c3:9021/2.0/clusters/ksql/$KSQL_ID | jq '.nodeCount'
```

---

## API Version Compatibility

| Component | API Version | Minimum Version | Notes |
|-----------|-------------|-----------------|-------|
| Control Center | v2.0 | 5.0.0+ | Used for multi-component discovery |
| REST Proxy | v3 | 5.0.0+ | Used for Kafka cluster details |
| Schema Registry | v1 | 3.0.0+ | Direct API |
| Kafka Connect | v1 | 0.10.0+ | Direct API |
| ksqlDB | v1 | 0.6.0+ | Direct API |
| Prometheus | v1 | 2.0.0+ | Direct API |
| Alertmanager | v2 | 0.15.0+ | Direct API |

---

## Summary

- **Control Center API**: Best for discovering multiple component node counts (Connect, SR, ksqlDB)
- **REST Proxy API**: Best for complete Kafka cluster topology (brokers, partitions, consumers)
- **Direct APIs**: Required for Prometheus, Alertmanager, and detailed component information
- **Combined Approach**: Use all three for complete Confluent Platform discovery

The cp-discovery tool leverages all available APIs to provide comprehensive platform topology and node count information.
