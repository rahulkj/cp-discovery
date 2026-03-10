package model

// Config represents the main configuration structure
type Config struct {
	Clusters []ClusterConfig `yaml:"clusters"`
	Output   OutputConfig    `yaml:"output"`
}

// ClusterConfig represents configuration for a single Confluent Platform cluster
type ClusterConfig struct {
	Name           string                  `yaml:"name"`
	Kafka          KafkaConfig             `yaml:"kafka"`
	SharedAuth     *SharedAuthConfig       `yaml:"shared_auth,omitempty"`
	SchemaRegistry SchemaRegistryConfig    `yaml:"schema_registry,omitempty"`
	KafkaConnect   KafkaConnectConfig      `yaml:"kafka_connect,omitempty"`
	KsqlDB         KsqlDBConfig            `yaml:"ksqldb,omitempty"`
	RestProxy      RestProxyConfig         `yaml:"rest_proxy,omitempty"`
	ControlCenter  ControlCenterConfig     `yaml:"control_center,omitempty"`
	Prometheus     PrometheusConfig        `yaml:"prometheus,omitempty"`
	Alertmanager   AlertmanagerConfig      `yaml:"alertmanager,omitempty"`
	Overrides      *ComponentOverrides     `yaml:"overrides,omitempty"`
}

// SharedAuthConfig contains authentication credentials shared across components
type SharedAuthConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// KafkaConfig represents Kafka cluster connection configuration
type KafkaConfig struct {
	BootstrapServers      string `yaml:"bootstrap_servers"`
	SecurityProtocol      string `yaml:"security_protocol,omitempty"`
	SaslMechanism         string `yaml:"sasl_mechanism,omitempty"`
	SaslUsername          string `yaml:"sasl_username,omitempty"`
	SaslPassword          string `yaml:"sasl_password,omitempty"`
	SslCaLocation         string `yaml:"ssl_ca_location,omitempty"`
	SslCertLocation       string `yaml:"ssl_cert_location,omitempty"`
	SslKeyLocation        string `yaml:"ssl_key_location,omitempty"`
	SslKeyPassword        string `yaml:"ssl_key_password,omitempty"`
	SslEndpointIdentification string `yaml:"ssl_endpoint_identification,omitempty"`
}

// SchemaRegistryConfig represents Schema Registry connection configuration
type SchemaRegistryConfig struct {
	URL               string `yaml:"url,omitempty"`
	BasicAuthUsername string `yaml:"basic_auth_username,omitempty"`
	BasicAuthPassword string `yaml:"basic_auth_password,omitempty"`
	BearerToken       string `yaml:"bearer_token,omitempty"`
	APIKey            string `yaml:"api_key,omitempty"`
	APIKeyHeader      string `yaml:"api_key_header,omitempty"`
}

// KafkaConnectConfig represents Kafka Connect connection configuration
type KafkaConnectConfig struct {
	URL               string `yaml:"url,omitempty"`
	BasicAuthUsername string `yaml:"basic_auth_username,omitempty"`
	BasicAuthPassword string `yaml:"basic_auth_password,omitempty"`
	BearerToken       string `yaml:"bearer_token,omitempty"`
	APIKey            string `yaml:"api_key,omitempty"`
	APIKeyHeader      string `yaml:"api_key_header,omitempty"`
}

// KsqlDBConfig represents ksqlDB connection configuration
type KsqlDBConfig struct {
	URL               string `yaml:"url,omitempty"`
	BasicAuthUsername string `yaml:"basic_auth_username,omitempty"`
	BasicAuthPassword string `yaml:"basic_auth_password,omitempty"`
	BearerToken       string `yaml:"bearer_token,omitempty"`
	APIKey            string `yaml:"api_key,omitempty"`
	APIKeyHeader      string `yaml:"api_key_header,omitempty"`
}

// RestProxyConfig represents REST Proxy connection configuration
type RestProxyConfig struct {
	URL               string `yaml:"url,omitempty"`
	BasicAuthUsername string `yaml:"basic_auth_username,omitempty"`
	BasicAuthPassword string `yaml:"basic_auth_password,omitempty"`
	BearerToken       string `yaml:"bearer_token,omitempty"`
	APIKey            string `yaml:"api_key,omitempty"`
	APIKeyHeader      string `yaml:"api_key_header,omitempty"`
}

// ControlCenterConfig represents Confluent Control Center connection configuration
type ControlCenterConfig struct {
	URL               string `yaml:"url,omitempty"`
	BasicAuthUsername string `yaml:"basic_auth_username,omitempty"`
	BasicAuthPassword string `yaml:"basic_auth_password,omitempty"`
	BearerToken       string `yaml:"bearer_token,omitempty"`
	APIKey            string `yaml:"api_key,omitempty"`
	APIKeyHeader      string `yaml:"api_key_header,omitempty"`
}

// PrometheusConfig represents Prometheus connection configuration
type PrometheusConfig struct {
	URL               string `yaml:"url,omitempty"`
	BasicAuthUsername string `yaml:"basic_auth_username,omitempty"`
	BasicAuthPassword string `yaml:"basic_auth_password,omitempty"`
	BearerToken       string `yaml:"bearer_token,omitempty"`
	APIKey            string `yaml:"api_key,omitempty"`
	APIKeyHeader      string `yaml:"api_key_header,omitempty"`
}

// AlertmanagerConfig represents Alertmanager connection configuration
type AlertmanagerConfig struct {
	URL               string `yaml:"url,omitempty"`
	BasicAuthUsername string `yaml:"basic_auth_username,omitempty"`
	BasicAuthPassword string `yaml:"basic_auth_password,omitempty"`
	BearerToken       string `yaml:"bearer_token,omitempty"`
	APIKey            string `yaml:"api_key,omitempty"`
	APIKeyHeader      string `yaml:"api_key_header,omitempty"`
}

// OutputConfig represents output configuration
type OutputConfig struct {
	Format   string `yaml:"format"`
	File     string `yaml:"file"`
	Detailed bool   `yaml:"detailed"`
}

// ComponentOverrides allows disabling specific component discovery
type ComponentOverrides struct {
	DisableSchemaRegistry bool `yaml:"disable_schema_registry"`
	DisableKafkaConnect   bool `yaml:"disable_kafka_connect"`
	DisableKsqlDB         bool `yaml:"disable_ksqldb"`
	DisableRestProxy      bool `yaml:"disable_rest_proxy"`
	DisableControlCenter  bool `yaml:"disable_control_center"`
	DisablePrometheus     bool `yaml:"disable_prometheus"`
	DisableAlertmanager   bool `yaml:"disable_alertmanager"`
}

// DiscoveryReport represents the complete discovery report
type DiscoveryReport struct {
	Timestamp     string          `json:"timestamp"`
	TotalClusters int             `json:"total_clusters"`
	Clusters      []ClusterReport `json:"clusters"`
}

// ClusterReport represents discovery results for a single cluster
type ClusterReport struct {
	Name                      string                  `json:"name"`
	Status                    string                  `json:"status"`
	ConfluentPlatformVersion  string                  `json:"confluent_platform_version,omitempty"`
	Errors                    []string                `json:"errors,omitempty"`
	Kafka                     KafkaReport             `json:"kafka"`
	SchemaRegistry            SchemaRegistryReport    `json:"schema_registry"`
	KafkaConnect              KafkaConnectReport      `json:"kafka_connect"`
	KsqlDB                    KsqlDBReport            `json:"ksqldb"`
	RestProxy                 RestProxyReport         `json:"rest_proxy"`
	ControlCenter             ControlCenterReport     `json:"control_center"`
	Prometheus                PrometheusReport        `json:"prometheus"`
	Alertmanager              AlertmanagerReport      `json:"alertmanager"`
}

// KafkaReport represents Kafka cluster discovery results
type KafkaReport struct {
	Available          bool            `json:"available"`
	BrokerCount        int             `json:"broker_count,omitempty"`
	ControllerType     string          `json:"controller_type,omitempty"`
	ControllerCount    int             `json:"controller_count,omitempty"`
	ZookeeperNodes     int             `json:"zookeeper_nodes,omitempty"`
	TopicCount         int             `json:"topic_count,omitempty"`
	InternalTopics     int             `json:"internal_topics,omitempty"`
	ExternalTopics     int             `json:"external_topics,omitempty"`
	TotalPartitions    int             `json:"total_partitions,omitempty"`
	SecurityConfig     SecurityConfig  `json:"security_config,omitempty"`
	Topics             []TopicInfo     `json:"topics,omitempty"`
	Brokers            []BrokerInfo    `json:"brokers,omitempty"`
	ClusterMetrics     ClusterMetrics  `json:"cluster_metrics,omitempty"`
}

// SecurityConfig represents security configuration
type SecurityConfig struct {
	SaslMechanisms       []string `json:"sasl_mechanisms,omitempty"`
	SecurityProtocols    []string `json:"security_protocols,omitempty"`
	SslEnabled           bool     `json:"ssl_enabled"`
	SaslEnabled          bool     `json:"sasl_enabled"`
	AuthenticationMethod string   `json:"authentication_method,omitempty"`
}

// TopicInfo represents information about a Kafka topic
type TopicInfo struct {
	Name                    string   `json:"name"`
	Partitions              int      `json:"partitions"`
	ReplicationFactor       int      `json:"replication_factor"`
	RetentionMs             int64    `json:"retention_ms,omitempty"`
	RetentionBytes          int64    `json:"retention_bytes,omitempty"`
	SizeBytes               int64    `json:"size_bytes,omitempty"`
	ThroughputBytesInPerSec float64  `json:"throughput_bytes_in_per_sec,omitempty"`
	ThroughputBytesOutPerSec float64 `json:"throughput_bytes_out_per_sec,omitempty"`
	AssociatedSchemas       []string `json:"associated_schemas,omitempty"`
	IsInternal              bool     `json:"is_internal"`
}

// BrokerInfo represents information about a Kafka broker
type BrokerInfo struct {
	ID              int    `json:"id"`
	Host            string `json:"host"`
	Port            int    `json:"port"`
	Rack            string `json:"rack,omitempty"`
	DiskUsageBytes  int64  `json:"disk_usage_bytes,omitempty"`
}

// ClusterMetrics represents cluster-level metrics
type ClusterMetrics struct {
	BytesInPerSec              float64 `json:"bytes_in_per_sec"`
	BytesOutPerSec             float64 `json:"bytes_out_per_sec"`
	MessagesInPerSec           float64 `json:"messages_in_per_sec"`
	TotalDiskUsageBytes        int64   `json:"total_disk_usage_bytes"`
	UnderReplicatedPartitions  int     `json:"under_replicated_partitions"`
}

// SchemaRegistryReport represents Schema Registry discovery results
type SchemaRegistryReport struct {
	Available       bool              `json:"available"`
	Version         string            `json:"version,omitempty"`
	Mode            string            `json:"mode,omitempty"`
	TotalSchemas    int               `json:"total_schemas,omitempty"`
	Subjects        []string          `json:"subjects,omitempty"`
	NodeCount       int               `json:"node_count,omitempty"`
	SchemaExporters []SchemaLinkInfo  `json:"schema_exporters"`
	ExporterCount   int               `json:"exporter_count"`
}

// KafkaConnectReport represents Kafka Connect discovery results
type KafkaConnectReport struct {
	Available        bool             `json:"available"`
	Version          string           `json:"version,omitempty"`
	TotalConnectors  int              `json:"total_connectors,omitempty"`
	SinkConnectors   int              `json:"sink_connectors,omitempty"`
	SourceConnectors int              `json:"source_connectors,omitempty"`
	Connectors       []ConnectorInfo  `json:"connectors,omitempty"`
	WorkerCount      int              `json:"worker_count,omitempty"`
	Replicators      []ReplicatorInfo `json:"replicators"`
	ReplicatorCount  int              `json:"replicator_count"`
}

// ConnectorInfo represents information about a Kafka Connect connector
type ConnectorInfo struct {
	Name           string `json:"name"`
	Type           string `json:"type"`
	State          string `json:"state"`
	Tasks          int    `json:"tasks"`
	ConnectorClass string `json:"connector_class,omitempty"`
	Quickstart     string `json:"quickstart,omitempty"`
}

// KsqlDBReport represents ksqlDB discovery results
type KsqlDBReport struct {
	Available   bool   `json:"available"`
	Version     string `json:"version,omitempty"`
	Queries     int    `json:"queries,omitempty"`
	Streams     int    `json:"streams,omitempty"`
	Tables      int    `json:"tables,omitempty"`
	NodeCount   int    `json:"node_count,omitempty"`
}

// RestProxyReport represents REST Proxy discovery results
type RestProxyReport struct {
	Available             bool                  `json:"available"`
	Version               string                `json:"version,omitempty"`
	NodeCount             int                   `json:"node_count,omitempty"`
	BrokerCount           int                   `json:"broker_count,omitempty"`
	ControllerID          int                   `json:"controller_id,omitempty"`
	ControllerCount       int                   `json:"controller_count,omitempty"`
	ControllerMode        string                `json:"controller_mode,omitempty"`
	ClusterID             string                `json:"cluster_id,omitempty"`
	TopicCount            int                   `json:"topic_count,omitempty"`
	InternalTopics        int                   `json:"internal_topics,omitempty"`
	ExternalTopics        int                   `json:"external_topics,omitempty"`
	PartitionCount        int                   `json:"partition_count,omitempty"`
	AvgReplicationFactor  float64               `json:"avg_replication_factor,omitempty"`
	SecurityConfig        SecurityConfig        `json:"security_config,omitempty"`
	Brokers               []RestProxyBrokerInfo `json:"brokers,omitempty"`
	Topics                []RestProxyTopicInfo  `json:"topics,omitempty"`
	ConsumerGroups        []ConsumerGroupInfo   `json:"consumer_groups,omitempty"`
	ConsumerGroupCount    int                   `json:"consumer_group_count,omitempty"`
	ActiveConsumerGroups  int                   `json:"active_consumer_groups,omitempty"`
	AclCount              int                   `json:"acl_count,omitempty"`
	Acls                  []AclInfo             `json:"acls,omitempty"`
	ClusterConfig         map[string]string     `json:"cluster_config,omitempty"`
	ClusterLinks          []ClusterLinkInfo     `json:"cluster_links"`
	ClusterLinkCount      int                   `json:"cluster_link_count"`
}

// RestProxyBrokerInfo represents broker information from REST Proxy
type RestProxyBrokerInfo struct {
	BrokerID           int    `json:"broker_id"`
	Host               string `json:"host"`
	Port               int    `json:"port"`
	IsActiveController bool   `json:"is_active_controller"`
	HasControllerRole  bool   `json:"has_controller_role"`
}

// RestProxyTopicInfo represents detailed topic information from REST Proxy
type RestProxyTopicInfo struct {
	Name              string                     `json:"name"`
	IsInternal        bool                       `json:"is_internal"`
	PartitionCount    int                        `json:"partition_count"`
	ReplicationFactor int                        `json:"replication_factor"`
	Partitions        []RestProxyPartitionInfo   `json:"partitions,omitempty"`
	Configs           map[string]string          `json:"configs,omitempty"`
}

// RestProxyPartitionInfo represents partition details
type RestProxyPartitionInfo struct {
	PartitionID int   `json:"partition_id"`
	Leader      int   `json:"leader"`
	Replicas    []int `json:"replicas"`
	ISR         []int `json:"isr"`
}

// ConsumerGroupInfo represents consumer group information
type ConsumerGroupInfo struct {
	GroupID           string `json:"group_id"`
	State             string `json:"state"`
	PartitionAssignor string `json:"partition_assignor,omitempty"`
	MemberCount       int    `json:"member_count,omitempty"`
	Lag               int64  `json:"lag,omitempty"`
}

// AclInfo represents ACL information
type AclInfo struct {
	ResourceType string `json:"resource_type"`
	ResourceName string `json:"resource_name"`
	PatternType  string `json:"pattern_type"`
	Principal    string `json:"principal"`
	Operation    string `json:"operation"`
	Permission   string `json:"permission"`
}

// ClusterLinkInfo represents cluster link information
type ClusterLinkInfo struct {
	LinkName           string              `json:"link_name"`
	LinkID             string              `json:"link_id,omitempty"`
	SourceClusterID    string              `json:"source_cluster_id,omitempty"`
	DestinationCluster string              `json:"destination_cluster,omitempty"`
	RemoteClusterID    string              `json:"remote_cluster_id,omitempty"`
	State              string              `json:"state,omitempty"`
	MirrorTopicCount   int                 `json:"mirror_topic_count,omitempty"`
	MirrorTopics       []MirrorTopicInfo   `json:"mirror_topics,omitempty"`
	Configs            map[string]string   `json:"configs,omitempty"`
}

// MirrorTopicInfo represents a mirror topic in a cluster link
type MirrorTopicInfo struct {
	MirrorTopicName  string `json:"mirror_topic_name"`
	SourceTopicName  string `json:"source_topic_name"`
	State            string `json:"state,omitempty"`
	MirrorStatus     string `json:"mirror_status,omitempty"`
	NumPartitions    int    `json:"num_partitions,omitempty"`
	MaxPerPartitionMirrorLag int64 `json:"max_per_partition_mirror_lag,omitempty"`
}

// ReplicatorInfo represents Confluent Replicator connector information
type ReplicatorInfo struct {
	Name                string `json:"name"`
	SourceCluster       string `json:"source_cluster,omitempty"`
	DestinationCluster  string `json:"destination_cluster,omitempty"`
	TopicWhitelist      string `json:"topic_whitelist,omitempty"`
	TopicBlacklist      string `json:"topic_blacklist,omitempty"`
	TopicRenameFormat   string `json:"topic_rename_format,omitempty"`
	State               string `json:"state"`
	Tasks               int    `json:"tasks"`
}

// SchemaLinkInfo represents Schema Registry linking/exporter information
type SchemaLinkInfo struct {
	ExporterName    string            `json:"exporter_name"`
	Subjects        []string          `json:"subjects,omitempty"`
	SubjectFormat   string            `json:"subject_format,omitempty"`
	ContextType     string            `json:"context_type,omitempty"`
	Context         string            `json:"context,omitempty"`
	Config          map[string]string `json:"config,omitempty"`
}

// ControlCenterReport represents Control Center discovery results
type ControlCenterReport struct {
	Available         bool                      `json:"available"`
	Version           string                    `json:"version,omitempty"`
	URL               string                    `json:"url,omitempty"`
	NodeCount         int                       `json:"node_count,omitempty"`
	MonitoredClusters int                       `json:"monitored_clusters,omitempty"`
	Clusters          []C3ClusterInfo           `json:"clusters,omitempty"`
	ConnectClusters   []C3ConnectClusterInfo    `json:"connect_clusters,omitempty"`
	SchemaRegistries  []C3SchemaRegistryInfo    `json:"schema_registries,omitempty"`
	KsqlClusters      []C3KsqlClusterInfo       `json:"ksql_clusters,omitempty"`
	TotalConsumerLag  int64                     `json:"total_consumer_lag,omitempty"`
}

// C3ClusterInfo represents a Kafka cluster monitored by Control Center
type C3ClusterInfo struct {
	ClusterID      string `json:"cluster_id"`
	ClusterName    string `json:"cluster_name"`
	BrokerCount    int    `json:"broker_count"`
	TopicCount     int    `json:"topic_count"`
	PartitionCount int    `json:"partition_count"`
	HealthStatus   string `json:"health_status"`
}

// C3ConnectClusterInfo represents a Connect cluster monitored by Control Center
type C3ConnectClusterInfo struct {
	ClusterName       string             `json:"cluster_name"`
	ClusterID         string             `json:"cluster_id"`
	KafkaClusterID    string             `json:"kafka_cluster_id"`
	ConnectorCount    int                `json:"connector_count"`
	WorkerCount       int                `json:"worker_count"`
	FailedConnectors  int                `json:"failed_connectors"`
	SourceConnectors  int                `json:"source_connectors"`
	SinkConnectors    int                `json:"sink_connectors"`
	RunningConnectors int                `json:"running_connectors"`
	Connectors        []C3ConnectorInfo  `json:"connectors,omitempty"`
}

// C3ConnectorInfo represents connector details from Control Center
type C3ConnectorInfo struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	State string `json:"state"`
	Tasks int    `json:"tasks"`
}

// C3SchemaRegistryInfo represents a Schema Registry monitored by Control Center
type C3SchemaRegistryInfo struct {
	ClusterName    string   `json:"cluster_name"`
	ClusterID      string   `json:"cluster_id"`
	KafkaClusterID string   `json:"kafka_cluster_id"`
	Version        string   `json:"version"`
	SchemaCount    int      `json:"schema_count"`
	Mode           string   `json:"mode"`
	NodeCount      int      `json:"node_count,omitempty"`
	Subjects       []string `json:"subjects,omitempty"`
}

// C3KsqlClusterInfo represents a ksqlDB cluster monitored by Control Center
type C3KsqlClusterInfo struct {
	ClusterName    string `json:"cluster_name"`
	ClusterID      string `json:"cluster_id"`
	KafkaClusterID string `json:"kafka_cluster_id"`
	QueryCount     int    `json:"query_count"`
	StreamCount    int    `json:"stream_count"`
	TableCount     int    `json:"table_count"`
	NodeCount      int    `json:"node_count,omitempty"`
}

// PrometheusReport represents Prometheus discovery results
type PrometheusReport struct {
	Available        bool                     `json:"available"`
	Version          string                   `json:"version,omitempty"`
	URL              string                   `json:"url,omitempty"`
	NodeCount        int                      `json:"node_count,omitempty"`
	TargetsUp        int                      `json:"targets_up,omitempty"`
	TargetsDown      int                      `json:"targets_down,omitempty"`
	HighAvailability bool                     `json:"high_availability,omitempty"`
	ClusterMetrics   PrometheusClusterMetrics `json:"cluster_metrics,omitempty"`
}

// PrometheusClusterMetrics represents cluster metrics from Prometheus
type PrometheusClusterMetrics struct {
	// Throughput Metrics
	BytesInPerSec    float64 `json:"bytes_in_per_sec"`
	BytesOutPerSec   float64 `json:"bytes_out_per_sec"`
	MessagesInPerSec float64 `json:"messages_in_per_sec"`

	// Controller Metrics
	ActiveControllerCount int `json:"active_controller_count"`

	// Partition Metrics
	UnderReplicatedPartitions int `json:"under_replicated_partitions"`
	OfflinePartitions         int `json:"offline_partitions"`
	TotalPartitions           int `json:"total_partitions"`

	// Broker Metrics
	TotalBrokers  int `json:"total_brokers"`
	OnlineBrokers int `json:"online_brokers"`
	LeaderCount   int `json:"leader_count"`

	// Consumer Metrics
	TotalConsumerLag int64 `json:"total_consumer_lag"`
	ConsumerGroups   int   `json:"consumer_groups"`

	// JVM Metrics
	AvgHeapUsedPercent float64 `json:"avg_heap_used_percent,omitempty"`
	AvgCPUUsedPercent  float64 `json:"avg_cpu_used_percent,omitempty"`
}

// AlertmanagerReport represents Alertmanager discovery results
type AlertmanagerReport struct {
	Available    bool     `json:"available"`
	Version      string   `json:"version,omitempty"`
	URL          string   `json:"url,omitempty"`
	ClusterSize  int      `json:"cluster_size,omitempty"`
	ClusterPeers []string `json:"cluster_peers,omitempty"`
	ActiveAlerts int      `json:"active_alerts,omitempty"`
}
