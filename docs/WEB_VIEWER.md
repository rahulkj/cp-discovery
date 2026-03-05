# Web Viewer Feature

The Confluent Discovery Tool now includes a built-in web viewer that opens discovery reports in your browser with a modern, interactive HTML interface.

## Features

- **Modern UI**: Beautiful gradient design with responsive layout
- **Tabbed Interface**: Switch between Overview, Clusters, and Raw JSON views
- **Summary Cards**: Quick metrics showing total clusters, brokers, and topics
- **Component Cards**: Detailed information for each Confluent Platform component
- **Status Badges**: Color-coded cluster health indicators (healthy/partial/error)
- **No External Dependencies**: Self-contained HTML with embedded CSS and JavaScript

## Usage

### View an Existing Report

```bash
# View a discovery report in your browser
./bin/cp-discovery -view-file test-report.json

# Use a custom port
./bin/cp-discovery -view-file test-report.json -port 8888
```

### Run Discovery and View

```bash
# Run discovery and automatically open results in browser
# If no output file is specified, a temporary file will be created automatically
./bin/cp-discovery -view

# Specify an output file to keep the report after viewing
./bin/cp-discovery -view -output my-report.json

# Combine with other flags
./bin/cp-discovery -config configs/config-production.yaml -view -port 8888
```

**Note:** When using `-view` without specifying an `-output` file, the tool automatically:
1. Creates a temporary JSON file in your system's temp directory
2. Runs the discovery and saves the report to the temp file
3. Opens the report in your browser
4. Cleans up the temp file when you stop the server (Ctrl+C)

## Command-Line Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `-view` | bool | false | Open report in web browser after discovery |
| `-port` | int | 8080 | Port for web view server |
| `-view-file` | string | "" | View existing report file (skips discovery) |

## How It Works

1. The tool starts a local HTTP server on the specified port (default: 8080)
2. Serves an HTML page at `http://localhost:8080/`
3. Provides a JSON API endpoint at `http://localhost:8080/api/report`
4. Automatically opens the page in your default browser
5. Server runs until you press Ctrl+C

## Temporary File Behavior

When you use the `-view` flag without specifying an output file:

### Automatic Temp File Creation
```bash
# This creates a temporary file automatically
./bin/cp-discovery -view

# Output:
# Using temporary report file: /tmp/cp-discovery-1234567890.json
# [discovery runs...]
# 🌐 Web Report Viewer
# Report: /tmp/cp-discovery-1234567890.json
# Server: http://localhost:8080
```

### Automatic Cleanup
- The temporary file is automatically deleted when you stop the web server (Ctrl+C)
- You'll see a cleanup message: `Cleaning up temporary file: /tmp/cp-discovery-1234567890.json`
- No manual cleanup required

### Keep the Report
If you want to keep the report file after viewing:
```bash
# Specify an output file to preserve the report
./bin/cp-discovery -view -output saved-report.json

# Now saved-report.json persists after you stop the server
```

### Temp File Location
- **macOS/Linux**: `/tmp/cp-discovery-*.json`
- **Windows**: `%TEMP%\cp-discovery-*.json`

The temporary files are created in your system's standard temporary directory and are automatically cleaned up.

## Browser Compatibility

The web viewer uses standard HTML5, CSS3, and ES6 JavaScript. It works with:
- Chrome/Edge (recommended)
- Firefox
- Safari
- Any modern browser supporting ES6

## Screenshots

### Overview Tab
- Summary cards showing key metrics
- Total clusters, healthy clusters, brokers, and topics

### Clusters Tab
- Detailed view of each cluster
- Component cards for Kafka, Schema Registry, Connect, ksqlDB, etc.
- Individual metrics and version information

### Raw JSON Tab
- Complete discovery report in formatted JSON
- Useful for debugging and data extraction

## Examples

### Basic Usage

```bash
# Run discovery on local cluster and view results
./bin/cp-discovery -config configs/example-local.yaml -view
```

### View Previous Report

```bash
# View yesterday's report
./bin/cp-discovery -view-file discovery-report-2026-03-03.json
```

### Custom Port

```bash
# Use port 9000 instead of default 8080
./bin/cp-discovery -view -port 9000
```

### Complete Workflow

```bash
# 1. Run discovery with custom output
./bin/cp-discovery -output reports/prod-$(date +%Y%m%d).json

# 2. View the generated report
./bin/cp-discovery -view-file reports/prod-$(date +%Y%m%d).json
```

## Troubleshooting

### Port Already in Use

If port 8080 is already in use:
```bash
./bin/cp-discovery -view-file report.json -port 8888
```

### Browser Doesn't Open

If the browser doesn't open automatically:
1. Check the console output for the URL
2. Manually open `http://localhost:8080` in your browser
3. The server will continue running

### File Not Found

Ensure the report file path is correct:
```bash
# Use absolute path if needed
./bin/cp-discovery -view-file /full/path/to/report.json
```

## Technical Details

### Architecture
- **Server**: Go's `net/http` package
- **Frontend**: Single-page application with vanilla JavaScript
- **Styling**: Inline CSS with gradient theme
- **Data Loading**: Fetch API for JSON retrieval

### Security Notes
- Server binds to `localhost` only (not accessible externally)
- No authentication required (local access only)
- Temporary server - stops when you terminate the process

### Platform Support
- **macOS**: Uses `open` command
- **Linux**: Uses `xdg-open` command
- **Windows**: Uses `cmd /c start` command

## Integration Examples

### CI/CD Pipeline

```bash
# Generate report in CI
./bin/cp-discovery -output reports/ci-${BUILD_ID}.json

# View locally after downloading
./bin/cp-discovery -view-file reports/ci-${BUILD_ID}.json
```

### Monitoring Script

```bash
#!/bin/bash
# Daily discovery and view
REPORT_FILE="discovery-$(date +%Y%m%d).json"
./bin/cp-discovery -output "$REPORT_FILE"
./bin/cp-discovery -view-file "$REPORT_FILE"
```

### Scheduled Reports

```bash
# Cron job to generate daily reports
0 6 * * * /path/t./cp-discovery -output /reports/daily-$(date +\%Y\%m\%d).json

# View latest report interactively
./bin/cp-discovery -view-file /reports/daily-$(date +%Y%m%d).json
```

## Future Enhancements

Potential improvements:
- Export to PDF
- Historical comparison view
- Real-time refresh
- Custom themes
- Filtering and search
- Metrics charting

## Related Documentation

- [Command-Line Arguments](../README.md#command-line-arguments-new)
- [Usage Examples](USAGE_EXAMPLES.md)
- [New Features](NEW_FEATURES.md)
