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
	// LDAP Authentication
	LDAPEnabled       bool   `yaml:"ldap_enabled,omitempty"`
	LDAPServer        string `yaml:"ldap_server,omitempty"`
	LDAPUsername      string `yaml:"ldap_username,omitempty"`
	LDAPPassword      string `yaml:"ldap_password,omitempty"`
	LDAPBaseDN        string `yaml:"ldap_base_dn,omitempty"`
	// OAuth/SSO Authentication
	OAuthEnabled      bool   `yaml:"oauth_enabled,omitempty"`
	OAuthClientID     string `yaml:"oauth_client_id,omitempty"`
	OAuthClientSecret string `yaml:"oauth_client_secret,omitempty"`
	OAuthTokenURL     string `yaml:"oauth_token_url,omitempty"`
	OAuthScopes       string `yaml:"oauth_scopes,omitempty"`
}

// KafkaConnectConfig represents Kafka Connect connection configuration
type KafkaConnectConfig struct {
	URL               string `yaml:"url,omitempty"`
	BasicAuthUsername string `yaml:"basic_auth_username,omitempty"`
	BasicAuthPassword string `yaml:"basic_auth_password,omitempty"`
	BearerToken       string `yaml:"bearer_token,omitempty"`
	APIKey            string `yaml:"api_key,omitempty"`
	APIKeyHeader      string `yaml:"api_key_header,omitempty"`
	// LDAP Authentication
	LDAPEnabled       bool   `yaml:"ldap_enabled,omitempty"`
	LDAPServer        string `yaml:"ldap_server,omitempty"`
	LDAPUsername      string `yaml:"ldap_username,omitempty"`
	LDAPPassword      string `yaml:"ldap_password,omitempty"`
	LDAPBaseDN        string `yaml:"ldap_base_dn,omitempty"`
	// OAuth/SSO Authentication
	OAuthEnabled      bool   `yaml:"oauth_enabled,omitempty"`
	OAuthClientID     string `yaml:"oauth_client_id,omitempty"`
	OAuthClientSecret string `yaml:"oauth_client_secret,omitempty"`
	OAuthTokenURL     string `yaml:"oauth_token_url,omitempty"`
	OAuthScopes       string `yaml:"oauth_scopes,omitempty"`
}

// KsqlDBConfig represents ksqlDB connection configuration
type KsqlDBConfig struct {
	URL               string `yaml:"url,omitempty"`
	BasicAuthUsername string `yaml:"basic_auth_username,omitempty"`
	BasicAuthPassword string `yaml:"basic_auth_password,omitempty"`
	BearerToken       string `yaml:"bearer_token,omitempty"`
	APIKey            string `yaml:"api_key,omitempty"`
	APIKeyHeader      string `yaml:"api_key_header,omitempty"`
	// LDAP Authentication
	LDAPEnabled       bool   `yaml:"ldap_enabled,omitempty"`
	LDAPServer        string `yaml:"ldap_server,omitempty"`
	LDAPUsername      string `yaml:"ldap_username,omitempty"`
	LDAPPassword      string `yaml:"ldap_password,omitempty"`
	LDAPBaseDN        string `yaml:"ldap_base_dn,omitempty"`
	// OAuth/SSO Authentication
	OAuthEnabled      bool   `yaml:"oauth_enabled,omitempty"`
	OAuthClientID     string `yaml:"oauth_client_id,omitempty"`
	OAuthClientSecret string `yaml:"oauth_client_secret,omitempty"`
	OAuthTokenURL     string `yaml:"oauth_token_url,omitempty"`
	OAuthScopes       string `yaml:"oauth_scopes,omitempty"`
}

// RestProxyConfig represents REST Proxy connection configuration
type RestProxyConfig struct {
	URL               string `yaml:"url,omitempty"`
	BasicAuthUsername string `yaml:"basic_auth_username,omitempty"`
	BasicAuthPassword string `yaml:"basic_auth_password,omitempty"`
	BearerToken       string `yaml:"bearer_token,omitempty"`
	APIKey            string `yaml:"api_key,omitempty"`
	APIKeyHeader      string `yaml:"api_key_header,omitempty"`
	// LDAP Authentication
	LDAPEnabled       bool   `yaml:"ldap_enabled,omitempty"`
	LDAPServer        string `yaml:"ldap_server,omitempty"`
	LDAPUsername      string `yaml:"ldap_username,omitempty"`
	LDAPPassword      string `yaml:"ldap_password,omitempty"`
	LDAPBaseDN        string `yaml:"ldap_base_dn,omitempty"`
	// OAuth/SSO Authentication
	OAuthEnabled      bool   `yaml:"oauth_enabled,omitempty"`
	OAuthClientID     string `yaml:"oauth_client_id,omitempty"`
	OAuthClientSecret string `yaml:"oauth_client_secret,omitempty"`
	OAuthTokenURL     string `yaml:"oauth_token_url,omitempty"`
	OAuthScopes       string `yaml:"oauth_scopes,omitempty"`
}

// ControlCenterConfig represents Confluent Control Center connection configuration
type ControlCenterConfig struct {
	URL               string `yaml:"url,omitempty"`
	BasicAuthUsername string `yaml:"basic_auth_username,omitempty"`
	BasicAuthPassword string `yaml:"basic_auth_password,omitempty"`
	BearerToken       string `yaml:"bearer_token,omitempty"`
	APIKey            string `yaml:"api_key,omitempty"`
	APIKeyHeader      string `yaml:"api_key_header,omitempty"`
	// LDAP Authentication
	LDAPEnabled       bool   `yaml:"ldap_enabled,omitempty"`
	LDAPServer        string `yaml:"ldap_server,omitempty"`
	LDAPUsername      string `yaml:"ldap_username,omitempty"`
	LDAPPassword      string `yaml:"ldap_password,omitempty"`
	LDAPBaseDN        string `yaml:"ldap_base_dn,omitempty"`
	// OAuth/SSO Authentication
	OAuthEnabled      bool   `yaml:"oauth_enabled,omitempty"`
	OAuthClientID     string `yaml:"oauth_client_id,omitempty"`
	OAuthClientSecret string `yaml:"oauth_client_secret,omitempty"`
	OAuthTokenURL     string `yaml:"oauth_token_url,omitempty"`
	OAuthScopes       string `yaml:"oauth_scopes,omitempty"`
}

// PrometheusConfig represents Prometheus connection configuration
type PrometheusConfig struct {
	URL               string `yaml:"url,omitempty"`
	BasicAuthUsername string `yaml:"basic_auth_username,omitempty"`
	BasicAuthPassword string `yaml:"basic_auth_password,omitempty"`
	BearerToken       string `yaml:"bearer_token,omitempty"`
	APIKey            string `yaml:"api_key,omitempty"`
	APIKeyHeader      string `yaml:"api_key_header,omitempty"`
	// LDAP Authentication
	LDAPEnabled       bool   `yaml:"ldap_enabled,omitempty"`
	LDAPServer        string `yaml:"ldap_server,omitempty"`
	LDAPUsername      string `yaml:"ldap_username,omitempty"`
	LDAPPassword      string `yaml:"ldap_password,omitempty"`
	LDAPBaseDN        string `yaml:"ldap_base_dn,omitempty"`
	// OAuth/SSO Authentication
	OAuthEnabled      bool   `yaml:"oauth_enabled,omitempty"`
	OAuthClientID     string `yaml:"oauth_client_id,omitempty"`
	OAuthClientSecret string `yaml:"oauth_client_secret,omitempty"`
	OAuthTokenURL     string `yaml:"oauth_token_url,omitempty"`
	OAuthScopes       string `yaml:"oauth_scopes,omitempty"`
}

// AlertmanagerConfig represents Alertmanager connection configuration
type AlertmanagerConfig struct {
	URL               string `yaml:"url,omitempty"`
	BasicAuthUsername string `yaml:"basic_auth_username,omitempty"`
	BasicAuthPassword string `yaml:"basic_auth_password,omitempty"`
	BearerToken       string `yaml:"bearer_token,omitempty"`
	APIKey            string `yaml:"api_key,omitempty"`
	APIKeyHeader      string `yaml:"api_key_header,omitempty"`
	// LDAP Authentication
	LDAPEnabled       bool   `yaml:"ldap_enabled,omitempty"`
	LDAPServer        string `yaml:"ldap_server,omitempty"`
	LDAPUsername      string `yaml:"ldap_username,omitempty"`
	LDAPPassword      string `yaml:"ldap_password,omitempty"`
	LDAPBaseDN        string `yaml:"ldap_base_dn,omitempty"`
	// OAuth/SSO Authentication
	OAuthEnabled      bool   `yaml:"oauth_enabled,omitempty"`
	OAuthClientID     string `yaml:"oauth_client_id,omitempty"`
	OAuthClientSecret string `yaml:"oauth_client_secret,omitempty"`
	OAuthTokenURL     string `yaml:"oauth_token_url,omitempty"`
	OAuthScopes       string `yaml:"oauth_scopes,omitempty"`
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
	Available          bool                   `json:"available"`
	BrokerCount        int                    `json:"broker_count,omitempty"`
	ControllerType     string                 `json:"controller_type,omitempty"`
	ControllerCount    int                    `json:"controller_count,omitempty"`
	ZookeeperNodes     int                    `json:"zookeeper_nodes,omitempty"`
	TopicCount         int                    `json:"topic_count,omitempty"`
	InternalTopics     int                    `json:"internal_topics,omitempty"`
	ExternalTopics     int                    `json:"external_topics,omitempty"`
	TotalPartitions    int                    `json:"total_partitions,omitempty"`
	SecurityConfig     SecurityConfig         `json:"security_config,omitempty"`
	Topics             []TopicInfo            `json:"topics,omitempty"`
	Brokers            []BrokerInfo           `json:"brokers,omitempty"`
	ClusterMetrics     ClusterMetrics         `json:"cluster_metrics,omitempty"`
	AdditionalInfo     *KafkaAdditionalInfo   `json:"kafka_endpoint_additional_info,omitempty"`
}

// KafkaAdditionalInfo contains extended Kafka cluster information
type KafkaAdditionalInfo struct {
	ConsumerGroups      []KafkaConsumerGroup    `json:"consumer_groups,omitempty"`
	TotalConsumerGroups int                     `json:"total_consumer_groups"`
	ActiveConsumerGroups int                    `json:"active_consumer_groups"`
	DetailedPartitions  []DetailedPartitionInfo `json:"detailed_partitions,omitempty"`
	BrokerConfigs       []BrokerConfigInfo      `json:"broker_configs,omitempty"`
	ClusterID           string                  `json:"cluster_id,omitempty"`
	ControllerID        int                     `json:"controller_id,omitempty"`
	ApiVersions         []ApiVersionInfo        `json:"api_versions,omitempty"`
}

// KafkaConsumerGroup represents detailed consumer group information
type KafkaConsumerGroup struct {
	GroupID          string                       `json:"group_id"`
	State            string                       `json:"state"`
	ProtocolType     string                       `json:"protocol_type,omitempty"`
	Protocol         string                       `json:"protocol,omitempty"`
	Members          []ConsumerGroupMember        `json:"members,omitempty"`
	MemberCount      int                          `json:"member_count"`
	Coordinator      int                          `json:"coordinator,omitempty"`
	Partitions       []ConsumerGroupPartition     `json:"partitions,omitempty"`
	TotalLag         int64                        `json:"total_lag,omitempty"`
}

// ConsumerGroupMember represents a consumer group member
type ConsumerGroupMember struct {
	MemberID       string   `json:"member_id"`
	ClientID       string   `json:"client_id,omitempty"`
	ClientHost     string   `json:"client_host,omitempty"`
	AssignedTopics []string `json:"assigned_topics,omitempty"`
	AssignedPartitions int  `json:"assigned_partitions"`
}

// ConsumerGroupPartition represents consumer group partition assignment
type ConsumerGroupPartition struct {
	Topic           string `json:"topic"`
	Partition       int    `json:"partition"`
	CurrentOffset   int64  `json:"current_offset"`
	LogEndOffset    int64  `json:"log_end_offset"`
	Lag             int64  `json:"lag"`
	MemberID        string `json:"member_id,omitempty"`
}

// DetailedPartitionInfo represents detailed partition information
type DetailedPartitionInfo struct {
	Topic           string `json:"topic"`
	Partition       int    `json:"partition"`
	Leader          int    `json:"leader"`
	Replicas        []int  `json:"replicas"`
	ISR             []int  `json:"isr"`
	OfflineReplicas []int  `json:"offline_replicas,omitempty"`
	FirstOffset     int64  `json:"first_offset"`
	LastOffset      int64  `json:"last_offset"`
	MessageCount    int64  `json:"message_count"`
}

// BrokerConfigInfo represents broker configuration
type BrokerConfigInfo struct {
	BrokerID int                       `json:"broker_id"`
	Configs  map[string]BrokerConfigEntry `json:"configs,omitempty"`
}

// BrokerConfigEntry represents a single broker configuration entry
type BrokerConfigEntry struct {
	Name      string `json:"name"`
	Value     string `json:"value"`
	Source    string `json:"source,omitempty"`
	Sensitive bool   `json:"sensitive"`
	ReadOnly  bool   `json:"read_only"`
}

// ApiVersionInfo represents API version information
type ApiVersionInfo struct {
	ApiKey     int16  `json:"api_key"`
	MinVersion int16  `json:"min_version"`
	MaxVersion int16  `json:"max_version"`
	ApiName    string `json:"api_name,omitempty"`
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
	Available       bool                        `json:"available"`
	Version         string                      `json:"version,omitempty"`
	Mode            string                      `json:"mode,omitempty"`
	TotalSchemas    int                         `json:"total_schemas,omitempty"`
	Subjects        []string                    `json:"subjects,omitempty"`
	NodeCount       int                         `json:"node_count,omitempty"`
	SchemaExporters []SchemaLinkInfo            `json:"schema_exporters"`
	ExporterCount   int                         `json:"exporter_count"`
	AdditionalInfo  *SchemaRegistryAdditionalInfo `json:"schema_registry_endpoint_additional_info,omitempty"`
}

// SchemaRegistryAdditionalInfo contains extended Schema Registry information
type SchemaRegistryAdditionalInfo struct {
	Subjects              []SubjectDetail       `json:"subjects,omitempty"`
	CompatibilityLevels   map[string]string     `json:"compatibility_levels,omitempty"`
	GlobalCompatibility   string                `json:"global_compatibility,omitempty"`
	ClusterInfo           SRClusterInfo         `json:"cluster_info,omitempty"`
	Config                map[string]string     `json:"config,omitempty"`
	Contexts              []string              `json:"contexts,omitempty"`
}

// SubjectDetail represents detailed subject information
type SubjectDetail struct {
	Subject       string          `json:"subject"`
	Versions      []int           `json:"versions"`
	LatestVersion int             `json:"latest_version"`
	LatestSchema  SchemaDetail    `json:"latest_schema,omitempty"`
	Compatibility string          `json:"compatibility,omitempty"`
}

// SchemaDetail represents detailed schema information
type SchemaDetail struct {
	ID         int                    `json:"id"`
	Version    int                    `json:"version"`
	Schema     string                 `json:"schema"`
	SchemaType string                 `json:"schema_type"`
	References []SchemaReference      `json:"references,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

// SchemaReference represents a schema reference
type SchemaReference struct {
	Name    string `json:"name"`
	Subject string `json:"subject"`
	Version int    `json:"version"`
}

// SRClusterInfo represents Schema Registry cluster information
type SRClusterInfo struct {
	Leader   string   `json:"leader,omitempty"`
	IsLeader bool     `json:"is_leader"`
	Members  []string `json:"members,omitempty"`
}

// KafkaConnectReport represents Kafka Connect discovery results
type KafkaConnectReport struct {
	Available        bool                       `json:"available"`
	Version          string                     `json:"version,omitempty"`
	TotalConnectors  int                        `json:"total_connectors,omitempty"`
	SinkConnectors   int                        `json:"sink_connectors,omitempty"`
	SourceConnectors int                        `json:"source_connectors,omitempty"`
	Connectors       []ConnectorInfo            `json:"connectors,omitempty"`
	WorkerCount      int                        `json:"worker_count,omitempty"`
	Replicators      []ReplicatorInfo           `json:"replicators"`
	ReplicatorCount  int                        `json:"replicator_count"`
	AdditionalInfo   *KafkaConnectAdditionalInfo `json:"kafka_connect_endpoint_additional_info,omitempty"`
}

// KafkaConnectAdditionalInfo contains extended Kafka Connect information
type KafkaConnectAdditionalInfo struct {
	Connectors      []ConnectorDetailInfo `json:"connectors,omitempty"`
	ConnectorPlugins []ConnectorPlugin    `json:"connector_plugins,omitempty"`
	Workers         []WorkerInfo          `json:"workers,omitempty"`
	ClusterInfo     ConnectClusterInfo    `json:"cluster_info,omitempty"`
}

// ConnectorDetailInfo represents detailed connector information
type ConnectorDetailInfo struct {
	Name           string                 `json:"name"`
	Type           string                 `json:"type"`
	State          string                 `json:"state"`
	WorkerID       string                 `json:"worker_id,omitempty"`
	Config         map[string]interface{} `json:"config,omitempty"`
	Tasks          []ConnectorTaskDetail  `json:"tasks,omitempty"`
	TaskCount      int                    `json:"task_count"`
	Topics         []string               `json:"topics,omitempty"`
}

// ConnectorTaskDetail represents task information
type ConnectorTaskDetail struct {
	TaskID   int                    `json:"task_id"`
	State    string                 `json:"state"`
	WorkerID string                 `json:"worker_id,omitempty"`
	Config   map[string]interface{} `json:"config,omitempty"`
	Trace    string                 `json:"trace,omitempty"`
}

// ConnectorPlugin represents an available connector plugin
type ConnectorPlugin struct {
	Class   string `json:"class"`
	Type    string `json:"type"`
	Version string `json:"version"`
}

// WorkerInfo represents Connect worker information
type WorkerInfo struct {
	ID      string `json:"id"`
	Host    string `json:"host,omitempty"`
	Port    int    `json:"port,omitempty"`
	Leader  bool   `json:"leader"`
}

// ConnectClusterInfo represents Connect cluster information
type ConnectClusterInfo struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	ClusterID string `json:"cluster_id,omitempty"`
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
	Available      bool                   `json:"available"`
	Version        string                 `json:"version,omitempty"`
	Queries        int                    `json:"queries,omitempty"`
	Streams        int                    `json:"streams,omitempty"`
	Tables         int                    `json:"tables,omitempty"`
	NodeCount      int                    `json:"node_count,omitempty"`
	AdditionalInfo *KsqlDBAdditionalInfo `json:"ksqldb_endpoint_additional_info,omitempty"`
}

// KsqlDBAdditionalInfo contains extended ksqlDB information
type KsqlDBAdditionalInfo struct {
	Queries       []KsqlQueryDetail   `json:"queries,omitempty"`
	Streams       []KsqlStreamDetail  `json:"streams,omitempty"`
	Tables        []KsqlTableDetail   `json:"tables,omitempty"`
	Topics        []string            `json:"topics,omitempty"`
	ClusterStatus KsqlClusterStatus   `json:"cluster_status,omitempty"`
	ServerInfo    KsqlServerInfo      `json:"server_info,omitempty"`
	Connectors    []string            `json:"connectors,omitempty"`
}

// KsqlQueryDetail represents detailed query information
type KsqlQueryDetail struct {
	ID            string   `json:"id"`
	QueryString   string   `json:"query_string"`
	StatementText string   `json:"statement_text,omitempty"`
	Sinks         []string `json:"sinks,omitempty"`
	Sources       []string `json:"sources,omitempty"`
	State         string   `json:"state,omitempty"`
	QueryType     string   `json:"query_type,omitempty"`
}

// KsqlStreamDetail represents stream information
type KsqlStreamDetail struct {
	Name           string   `json:"name"`
	Topic          string   `json:"topic"`
	KeyFormat      string   `json:"key_format,omitempty"`
	ValueFormat    string   `json:"value_format,omitempty"`
	IsWindowed     bool     `json:"is_windowed"`
	Fields         []KsqlField `json:"fields,omitempty"`
	QueryID        string   `json:"query_id,omitempty"`
}

// KsqlTableDetail represents table information
type KsqlTableDetail struct {
	Name           string   `json:"name"`
	Topic          string   `json:"topic"`
	KeyFormat      string   `json:"key_format,omitempty"`
	ValueFormat    string   `json:"value_format,omitempty"`
	IsWindowed     bool     `json:"is_windowed"`
	Fields         []KsqlField `json:"fields,omitempty"`
	QueryID        string   `json:"query_id,omitempty"`
}

// KsqlField represents a ksqlDB field
type KsqlField struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	IsKey  bool   `json:"is_key,omitempty"`
}

// KsqlClusterStatus represents ksqlDB cluster status
type KsqlClusterStatus struct {
	Hosts   []KsqlHostInfo `json:"hosts,omitempty"`
}

// KsqlHostInfo represents ksqlDB host information
type KsqlHostInfo struct {
	HostInfo       KsqlHost `json:"host_info"`
	IsActiveHost   bool     `json:"is_active_host"`
	LastStatusUpdateMs int64 `json:"last_status_update_ms,omitempty"`
}

// KsqlHost represents ksqlDB host details
type KsqlHost struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

// KsqlServerInfo represents ksqlDB server information
type KsqlServerInfo struct {
	Version      string `json:"version"`
	KafkaClusterID string `json:"kafka_cluster_id,omitempty"`
	KsqlServiceID  string `json:"ksql_service_id,omitempty"`
	ServerStatus   string `json:"server_status,omitempty"`
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
	AdditionalInfo        *RestProxyAdditionalInfo `json:"rest_proxy_endpoint_additional_info,omitempty"`
}

// RestProxyAdditionalInfo contains extended REST Proxy information
type RestProxyAdditionalInfo struct {
	ProducerDetails   []ProducerDetail         `json:"producer_details,omitempty"`
	ConsumerInstances []ConsumerInstanceDetail `json:"consumer_instances,omitempty"`
	BrokerConfigs     []map[string]string      `json:"broker_configs,omitempty"`
	TopicConfigs      map[string]map[string]string `json:"topic_configs,omitempty"`
}

// ProducerDetail represents producer information
type ProducerDetail struct {
	InstanceID string `json:"instance_id"`
	ClientID   string `json:"client_id,omitempty"`
}

// ConsumerInstanceDetail represents consumer instance information
type ConsumerInstanceDetail struct {
	InstanceID string   `json:"instance_id"`
	ConsumerID string   `json:"consumer_id,omitempty"`
	Format     string   `json:"format,omitempty"`
	Topics     []string `json:"topics,omitempty"`
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
	Available         bool                        `json:"available"`
	Version           string                      `json:"version,omitempty"`
	URL               string                      `json:"url,omitempty"`
	NodeCount         int                         `json:"node_count,omitempty"`
	MonitoredClusters int                         `json:"monitored_clusters,omitempty"`
	Clusters          []C3ClusterInfo             `json:"clusters,omitempty"`
	ConnectClusters   []C3ConnectClusterInfo      `json:"connect_clusters,omitempty"`
	SchemaRegistries  []C3SchemaRegistryInfo      `json:"schema_registries,omitempty"`
	KsqlClusters      []C3KsqlClusterInfo         `json:"ksql_clusters,omitempty"`
	TotalConsumerLag  int64                       `json:"total_consumer_lag,omitempty"`
	AdditionalInfo    *ControlCenterAdditionalInfo `json:"control_center_endpoint_additional_info,omitempty"`
}

// ControlCenterAdditionalInfo contains extended Control Center information
type ControlCenterAdditionalInfo struct {
	Alerts             []C3Alert           `json:"alerts,omitempty"`
	MonitoringMetrics  C3MonitoringMetrics `json:"monitoring_metrics,omitempty"`
	ClusterDetails     []C3ClusterDetail   `json:"cluster_details,omitempty"`
}

// C3Alert represents a Control Center alert
type C3Alert struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	ClusterID   string `json:"cluster_id,omitempty"`
	Severity    string `json:"severity,omitempty"`
	Status      string `json:"status,omitempty"`
	Description string `json:"description,omitempty"`
}

// C3MonitoringMetrics represents monitoring metrics from C3
type C3MonitoringMetrics struct {
	TotalBrokers        int     `json:"total_brokers"`
	TotalTopics         int     `json:"total_topics"`
	TotalPartitions     int     `json:"total_partitions"`
	TotalConnectors     int     `json:"total_connectors"`
	TotalConsumerGroups int     `json:"total_consumer_groups"`
	AverageThroughput   float64 `json:"average_throughput,omitempty"`
}

// C3ClusterDetail represents detailed cluster information from C3
type C3ClusterDetail struct {
	ClusterID   string                 `json:"cluster_id"`
	ClusterName string                 `json:"cluster_name"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
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
	Available        bool                       `json:"available"`
	Version          string                     `json:"version,omitempty"`
	URL              string                     `json:"url,omitempty"`
	NodeCount        int                        `json:"node_count,omitempty"`
	TargetsUp        int                        `json:"targets_up,omitempty"`
	TargetsDown      int                        `json:"targets_down,omitempty"`
	HighAvailability bool                       `json:"high_availability,omitempty"`
	ClusterMetrics   PrometheusClusterMetrics   `json:"cluster_metrics,omitempty"`
	AdditionalInfo   *PrometheusAdditionalInfo `json:"prometheus_endpoint_additional_info,omitempty"`
}

// PrometheusAdditionalInfo contains extended Prometheus information
type PrometheusAdditionalInfo struct {
	Targets        []PrometheusTarget        `json:"targets,omitempty"`
	AlertRules     []PrometheusAlertRule     `json:"alert_rules,omitempty"`
	RecordingRules []PrometheusRecordingRule `json:"recording_rules,omitempty"`
	Config         PrometheusServerConfig    `json:"config,omitempty"`
	TSDBStats      PrometheusTSDBStats       `json:"tsdb_stats,omitempty"`
	Runtimes       PrometheusRuntimeInfo     `json:"runtimes,omitempty"`
}

// PrometheusTarget represents a scrape target
type PrometheusTarget struct {
	DiscoveredLabels map[string]string `json:"discovered_labels,omitempty"`
	Labels           map[string]string `json:"labels,omitempty"`
	ScrapePool       string            `json:"scrape_pool,omitempty"`
	ScrapeURL        string            `json:"scrape_url,omitempty"`
	Health           string            `json:"health,omitempty"`
	LastError        string            `json:"last_error,omitempty"`
	LastScrape       string            `json:"last_scrape,omitempty"`
}

// PrometheusAlertRule represents an alert rule
type PrometheusAlertRule struct {
	Name        string                 `json:"name"`
	Query       string                 `json:"query"`
	Duration    float64                `json:"duration,omitempty"`
	Labels      map[string]string      `json:"labels,omitempty"`
	Annotations map[string]string      `json:"annotations,omitempty"`
	State       string                 `json:"state,omitempty"`
	Health      string                 `json:"health,omitempty"`
	Type        string                 `json:"type"`
}

// PrometheusRecordingRule represents a recording rule
type PrometheusRecordingRule struct {
	Name   string            `json:"name"`
	Query  string            `json:"query"`
	Labels map[string]string `json:"labels,omitempty"`
	Health string            `json:"health,omitempty"`
}

// PrometheusServerConfig represents Prometheus server configuration
type PrometheusServerConfig struct {
	GlobalConfig    map[string]interface{} `json:"global_config,omitempty"`
	ScrapeConfigs   int                    `json:"scrape_configs_count"`
	AlertingConfig  map[string]interface{} `json:"alerting_config,omitempty"`
}

// PrometheusTSDBStats represents TSDB statistics
type PrometheusTSDBStats struct {
	NumSeries         int     `json:"num_series"`
	NumLabelPairs     int     `json:"num_label_pairs"`
	NumSamples        int64   `json:"num_samples,omitempty"`
	ChunkCount        int64   `json:"chunk_count,omitempty"`
	HeadStats         TSDBHeadStats `json:"head_stats,omitempty"`
}

// TSDBHeadStats represents head block statistics
type TSDBHeadStats struct {
	NumSeries   int   `json:"num_series"`
	ChunkCount  int64 `json:"chunk_count"`
	MinTime     int64 `json:"min_time"`
	MaxTime     int64 `json:"max_time"`
}

// PrometheusRuntimeInfo represents Prometheus runtime information
type PrometheusRuntimeInfo struct {
	StartTime           string `json:"start_time,omitempty"`
	CWD                 string `json:"cwd,omitempty"`
	ReloadConfigSuccess bool   `json:"reload_config_success"`
	LastConfigTime      string `json:"last_config_time,omitempty"`
	StorageRetention    string `json:"storage_retention,omitempty"`
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
	Available      bool                         `json:"available"`
	Version        string                       `json:"version,omitempty"`
	URL            string                       `json:"url,omitempty"`
	ClusterSize    int                          `json:"cluster_size,omitempty"`
	ClusterPeers   []string                     `json:"cluster_peers,omitempty"`
	ActiveAlerts   int                          `json:"active_alerts,omitempty"`
	AdditionalInfo *AlertmanagerAdditionalInfo `json:"alertmanager_endpoint_additional_info,omitempty"`
}

// AlertmanagerAdditionalInfo contains extended Alertmanager information
type AlertmanagerAdditionalInfo struct {
	Alerts      []AlertmanagerAlert   `json:"alerts,omitempty"`
	Silences    []AlertmanagerSilence `json:"silences,omitempty"`
	Receivers   []string              `json:"receivers,omitempty"`
	Status      AlertmanagerStatus    `json:"status,omitempty"`
	Config      map[string]interface{} `json:"config,omitempty"`
}

// AlertmanagerAlert represents an alert
type AlertmanagerAlert struct {
	Labels       map[string]string `json:"labels,omitempty"`
	Annotations  map[string]string `json:"annotations,omitempty"`
	StartsAt     string            `json:"starts_at,omitempty"`
	EndsAt       string            `json:"ends_at,omitempty"`
	GeneratorURL string            `json:"generator_url,omitempty"`
	Status       AlertStatus       `json:"status,omitempty"`
}

// AlertStatus represents alert status
type AlertStatus struct {
	State       string   `json:"state,omitempty"`
	SilencedBy  []string `json:"silenced_by,omitempty"`
	InhibitedBy []string `json:"inhibited_by,omitempty"`
}

// AlertmanagerSilence represents a silence
type AlertmanagerSilence struct {
	ID        string            `json:"id"`
	Status    SilenceStatus     `json:"status,omitempty"`
	Matchers  []SilenceMatcher  `json:"matchers,omitempty"`
	StartsAt  string            `json:"starts_at,omitempty"`
	EndsAt    string            `json:"ends_at,omitempty"`
	CreatedBy string            `json:"created_by,omitempty"`
	Comment   string            `json:"comment,omitempty"`
}

// SilenceStatus represents silence status
type SilenceStatus struct {
	State string `json:"state,omitempty"`
}

// SilenceMatcher represents a silence matcher
type SilenceMatcher struct {
	Name    string `json:"name"`
	Value   string `json:"value"`
	IsRegex bool   `json:"is_regex"`
}

// AlertmanagerStatus represents Alertmanager status
type AlertmanagerStatus struct {
	Uptime      string                 `json:"uptime,omitempty"`
	VersionInfo map[string]interface{} `json:"version_info,omitempty"`
	Config      AlertmanagerConfigStatus `json:"config,omitempty"`
}

// AlertmanagerConfigStatus represents config status
type AlertmanagerConfigStatus struct {
	Original string `json:"original,omitempty"`
}
