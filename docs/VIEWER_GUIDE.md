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
# Web Viewer Usage Examples

This document provides practical examples of using the web viewer feature with different scenarios.

## Basic Usage

### Quick View (Temporary File)

The simplest way to view discovery results - uses a temporary file that's automatically cleaned up:

```bash
./bin/cp-discovery -view
```

**What happens:**
1. ✅ Creates temporary file: `/tmp/cp-discovery-1234567890.json`
2. ✅ Runs discovery and saves to temp file
3. ✅ Starts web server on port 8080
4. ✅ Opens browser automatically
5. ✅ Displays: `Using temporary report file: /tmp/cp-discovery-1234567890.json`
6. ✅ When you press Ctrl+C: Cleans up temp file automatically

### View and Save

Run discovery, view results, and keep the report file:

```bash
./bin/cp-discovery -view -output my-report.json
```

**What happens:**
1. ✅ Runs discovery and saves to `my-report.json`
2. ✅ Opens web viewer
3. ✅ File persists after you stop the server
4. ✅ Can re-view later with: `./bin/cp-discovery -view-file my-report.json`

### View Existing Report

View a previously generated report without running discovery:

```bash
./bin/cp-discovery -view-file my-report.json
```

**What happens:**
1. ✅ Loads existing `my-report.json`
2. ✅ Starts web viewer
3. ✅ No discovery performed
4. ✅ Original file is unchanged

## Advanced Examples

### Custom Port

Use a different port if 8080 is already in use:

```bash
# Temporary file with custom port
./bin/cp-discovery -view -port 9000

# Existing file with custom port
./bin/cp-discovery -view-file report.json -port 9000
```

### Specific Configuration

Run discovery with a specific config file and view results:

```bash
# Temporary file (auto-cleanup)
./bin/cp-discovery -config configs/config-production.yaml -view

# Save to specific location
./bin/cp-discovery -config configs/config-production.yaml -view -output prod-report.json
```

### Detailed Mode

Enable detailed discovery mode and view results:

```bash
# Temporary file
./bin/cp-discovery -view -detailed

# Save detailed report
./bin/cp-discovery -view -detailed -output detailed-report.json
```

### Custom Output Format

While the web viewer requires JSON, you can still specify the format:

```bash
# This will automatically use JSON for the temp file
./bin/cp-discovery -view -format json

# Save as YAML (won't open in web viewer automatically)
./bin/cp-discovery -output report.yaml -format yaml
```

## Workflow Examples

### Daily Health Check

Quick daily check without cluttering disk:

```bash
#!/bin/bash
# daily-check.sh - Quick health check with web viewer

echo "Running daily Confluent Platform health check..."
./bin/cp-discovery -view

# Temp file auto-cleans when you close the browser/server
```

### Weekly Report Generation

Generate and save weekly reports:

```bash
#!/bin/bash
# weekly-report.sh - Generate and view weekly reports

REPORT_FILE="reports/weekly-$(date +%Y%m%d).json"

echo "Generating weekly report: $REPORT_FILE"
./bin/cp-discovery -config configs/config-production.yaml \
  -output "$REPORT_FILE" \
  -detailed \
  -view

echo "Report saved to: $REPORT_FILE"
```

### Multi-Environment Comparison

View reports from different environments:

```bash
# Generate production report
./bin/cp-discovery -config configs/prod-config.yaml \
  -output reports/prod-$(date +%Y%m%d).json

# Generate staging report
./bin/cp-discovery -config configs/staging-config.yaml \
  -output reports/staging-$(date +%Y%m%d).json

# View production report
./bin/cp-discovery -view-file reports/prod-$(date +%Y%m%d).json -port 8080

# In another terminal, view staging report
./bin/cp-discovery -view-file reports/staging-$(date +%Y%m%d).json -port 8081
```

### Troubleshooting Session

Quick troubleshooting with temporary file:

```bash
# Run discovery and view immediately
./bin/cp-discovery -config configs/problem-cluster.yaml -view -detailed

# Review the web interface
# Take screenshots if needed
# When done, Ctrl+C to close and auto-cleanup
```

## CI/CD Integration

### Jenkins Pipeline

```groovy
pipeline {
    agent any
    stages {
        stage('Discover') {
            steps {
                sh '''
                    ./bin/cp-discovery \
                        -config configs/ci-config.yaml \
                        -output reports/build-${BUILD_NUMBER}.json \
                        -detailed
                '''
            }
        }
        stage('Archive') {
            steps {
                archiveArtifacts artifacts: 'reports/*.json'
            }
        }
    }
}

// Later, view the report locally:
// ./bin/cp-discovery -view-file reports/build-123.json
```

### GitHub Actions

```yaml
name: Platform Discovery

on:
  schedule:
    - cron: '0 6 * * *'  # Daily at 6 AM

jobs:
  discover:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2

      - name: Run Discovery
        run: |
          ./bin/cp-discovery \
            -config configs/ci-config.yaml \
            -output discovery-report.json \
            -detailed

      - name: Upload Report
        uses: actions/upload-artifact@v2
        with:
          name: discovery-report
          path: discovery-report.json

# Download and view locally:
# gh run download <run-id>
# ./bin/cp-discovery -view-file discovery-report/discovery-report.json
```

## Comparison: Temp vs Saved Files

### Use Temporary File When:

✅ Quick health checks
✅ Troubleshooting sessions
✅ One-time investigations
✅ Don't need to keep the report
✅ Want automatic cleanup

**Example:**
```bash
./bin/cp-discovery -view
```

### Save File When:

✅ Weekly/monthly reports
✅ Compliance documentation
✅ Historical tracking
✅ Sharing with team
✅ CI/CD pipelines

**Example:**
```bash
./bin/cp-discovery -view -output reports/$(date +%Y%m%d)-report.json
```

## Tips and Tricks

### 1. Quick Morning Check

```bash
# Add to your .bashrc or .zshrc
alias cpcheck='cd ~/cp-discovery && ./bin/cp-discovery -view'

# Then just run:
cpcheck
```

### 2. Multiple Clusters, Single View

```bash
# Generate reports for all clusters
for env in prod staging dev; do
  ./bin/cp-discovery -config configs/${env}-config.yaml \
    -output reports/${env}-report.json
done

# View each one as needed
./bin/cp-discovery -view-file reports/prod-report.json -port 8080
./bin/cp-discovery -view-file reports/staging-report.json -port 8081
```

### 3. Automated Daily Reports

```bash
# Cron job: Run daily at 8 AM, save to dated file
0 8 * * * cd /path/to/cp-discovery && ./bin/cp-discovery -output reports/daily-$(date +\%Y\%m\%d).json

# Manual review when needed:
./bin/cp-discovery -view-file reports/daily-20260304.json
```

### 4. Quick Diff Between Days

```bash
# Compare today vs yesterday using jq
./bin/cp-discovery -output today.json
diff <(jq -S . yesterday.json) <(jq -S . today.json)

# View each in browser for visual comparison
./bin/cp-discovery -view-file yesterday.json -port 8080
./bin/cp-discovery -view-file today.json -port 8081
```

## Troubleshooting

### Port Already in Use

```bash
# Check what's using port 8080
lsof -i :8080

# Use a different port
./bin/cp-discovery -view -port 8888
```

### Temp File Not Cleaned Up

If the tool crashes or is force-killed, the temp file might remain:

```bash
# Find leftover temp files
ls -lh /tmp/cp-discovery-*.json

# Clean up manually
rm /tmp/cp-discovery-*.json
```

### Browser Doesn't Open

```bash
# The server still runs, just open manually
# Look for the URL in the output:
# Server: http://localhost:8080

# Then open in your browser
```

### Want to Keep Temp File

If you started with `-view` but want to keep the report:

```bash
# Look for the temp file path in the output:
# Using temporary report file: /tmp/cp-discovery-1234567890.json

# Copy it before stopping the server:
cp /tmp/cp-discovery-1234567890.json my-saved-report.json
```

## Related Documentation

- [Web Viewer Guide](WEB_VIEWER.md) - Complete web viewer documentation
- [Usage Examples](USAGE_EXAMPLES.md) - General usage examples
- [New Features](NEW_FEATURES.md) - All v2.0.0 features
- [Configuration Reference](CONFIG_REFERENCE.md) - Configuration options

## Summary

The web viewer with automatic temporary file support makes it easy to:
- ✅ Quickly check cluster health without file clutter
- ✅ View reports in a modern, interactive interface
- ✅ Save reports when needed for compliance or history
- ✅ Integrate into workflows and automation
- ✅ Share reports with team members

**Quick Reference:**
```bash
# Temp file (auto-cleanup)
./bin/cp-discovery -view

# Save file (persistent)
./bin/cp-discovery -view -output report.json

# View existing
./bin/cp-discovery -view-file report.json
```
