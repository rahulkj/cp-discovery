# Changelog

All notable changes to cp-discovery will be documented in this file.

## [2.0.0] - 2026-03-04

### Added - Web Viewer
- **Interactive HTML Report Viewer**
  - Modern gradient UI with responsive design
  - Three tabs: Overview, Clusters, Raw JSON
  - Summary cards showing key metrics
  - Component cards for each Confluent Platform service
  - Color-coded status badges (healthy/partial/error)
  - Built-in HTTP server with automatic browser launching
  - No external dependencies - fully self-contained

- **`-view` flag**: Open report in browser after discovery
  - Automatically starts local web server
  - Opens default browser to view results
  - **Auto temp file creation**: When used without `-output`, creates temporary JSON file automatically
  - **Auto cleanup**: Temporary file is deleted when server stops (Ctrl+C)
  - Example: `./bin/cp-discovery -view`

- **`-view-file` flag**: View existing report without running discovery
  - Load any previously generated JSON report
  - Example: `./bin/cp-discovery -view-file report.json`

- **`-port` flag**: Customize web server port
  - Default: 8080
  - Example: `./bin/cp-discovery -view -port 9000`

### Added - Command-Line Arguments
- **`-output` flag**: Specify output file path from command line
  - Override config file setting
  - Support dynamic file naming
  - Example: `./bin/cp-discovery -output /tmp/report.json`

- **`-format` flag**: Override output format
  - Values: `json` or `yaml`
  - Example: `./bin/cp-discovery -format yaml`

- **`-detailed` flag**: Enable detailed discovery mode
  - Override config file setting
  - Example: `./bin/cp-discovery -detailed`

### Added - Enhanced Console Output
- **Network Throughput Display**
  - Shows Bytes In/Out per second (MB/s)
  - Shows Messages In per second
  - Displayed for Kafka and Prometheus metrics
  - Formatted for easy reading

- **Storage Details Display**
  - Shows total disk usage in GB
  - Cluster-level aggregation
  - Broker-level details (when available)
  - Formatted for capacity planning

- **Health Metrics Display**
  - Under-replicated partitions count
  - Only shown when issues detected
  - Quick health assessment

### Changed - Data Model
- **Enhanced BrokerInfo**
  - Added `DiskUsageBytes` field
  - Enables per-broker storage tracking
  - Backward compatible (optional field)

### Changed - Project Structure
- Reorganized code into standard Go layout
- Created `cmd/` for main application
- Created `internal/` for private packages
- Moved configs to `configs/` directory
- Binary now builds to `bin/`

### Added - Documentation
- **USAGE_EXAMPLES.md**: Comprehensive usage guide
- **NEW_FEATURES.md**: Detailed feature documentation
- **PROJECT_STRUCTURE.md**: Project organization guide
- **CLEANUP_SUMMARY.md**: Restructuring summary
- **CHANGELOG.md**: This file

### Fixed
- Controller count now shows total controllers (not just active)
- Proper model types for all discovery functions
- Package organization follows Go best practices

## [1.5.0] - 2026-03-04

### Added - Prometheus Integration
- Fetch cluster metrics from Prometheus
- 15 different metrics tracked
- Displayed in detailed mode
- Real-time operational visibility

### Enhanced - Component Discovery
- REST Proxy: Consumer groups, ACLs, cluster configs
- Control Center: Comprehensive component details
- Connect: Source/sink breakdown, running connectors
- Schema Registry: Mode and subject list
- ksqlDB: Cluster ID associations

## [1.0.0] - Initial Release

### Features
- Multi-cluster Kafka discovery
- Schema Registry discovery
- Kafka Connect discovery
- ksqlDB discovery
- REST Proxy discovery
- Control Center discovery
- YAML configuration support
- JSON/YAML output formats
- Detailed and minimal modes

---

## Version Notes

### Version Format
- **Major.Minor.Patch** (Semantic Versioning)
- Major: Breaking changes
- Minor: New features (backward compatible)
- Patch: Bug fixes

### Upgrade Guide

#### From 1.x to 2.0
**No breaking changes!** All existing configurations work unchanged.

**New capabilities:**
```bash
# Old way (still works)
./cp-discovery

# New way (with CLI args)
./cp-discovery -output custom.json -detailed
```

**New output sections:**
- Network throughput automatically displayed
- Storage details automatically displayed
- Health metrics automatically displayed

**Enhanced data:**
- Broker storage information
- More detailed metrics
- Better formatted output

---

## Future Enhancements

### Planned for 2.1.0
- [ ] Real-time metrics collection (JMX)
- [ ] Trend analysis over time
- [ ] Alert threshold configuration
- [ ] Export to monitoring systems
- [ ] Per-broker detailed metrics

### Planned for 3.0.0
- [ ] Web UI for reports
- [ ] Historical data storage
- [ ] Comparison between snapshots
- [ ] Automated recommendations
- [ ] Multi-datacenter topology mapping

---

## Contributing

See feature requests and report issues at:
https://github.com/rahulkj/cp-discovery/issues
