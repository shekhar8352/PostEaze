# Application Logs

This folder contains structured application logs for the PostEaze backend service. The logging system provides comprehensive request tracking, error monitoring, and application behavior analysis through structured JSON logs with automatic daily rotation.

## Logging Strategy

PostEaze implements a structured logging approach with the following key features:
- **JSON Format**: All logs are structured as JSON for easy parsing and analysis
- **Daily Rotation**: Log files are automatically rotated daily to manage file sizes
- **Request Correlation**: Each request gets a unique log ID for tracing across multiple log entries
- **Context-Aware**: Logs include caller information (file, line, function) for debugging
- **Level-Based Filtering**: Configurable log levels (DEBUG, INFO, WARN, ERROR)

## File Organization

### Naming Convention
Log files follow a consistent naming pattern:
```
{prefix}-{date}.log
```

- **Prefix**: `app` (configurable)
- **Date**: `YYYY-MM-DD` format
- **Extension**: `.log`

### Example Files
```
logs/
├── app-2025-08-01.log    # Logs for August 1st, 2025
├── app-2025-08-02.log    # Logs for August 2nd, 2025
├── app-2025-08-04.log    # Logs for August 4th, 2025
└── README.md             # This documentation
```

### Automatic Rotation
- **Daily Rotation**: New log file created automatically at midnight
- **No Size Limits**: Files grow throughout the day without size-based rotation
- **Automatic Directory Creation**: Logs directory created if it doesn't exist
- **Graceful Handling**: Application continues logging even if file operations fail

## Log Format

### JSON Structure
Each log entry is a single-line JSON object with the following fields:

```json
{
  "timestamp": "2025-08-04T01:20:41+05:30",
  "log_id": "7b0f9e30-50ec-4c43-a7c5-a9dc260d2153",
  "level": "INFO",
  "message": "Started POST /api/v1/auth/login | IP: ::1 | User-Agent: PostmanRuntime/7.44.1",
  "file": "log_middleware.go",
  "line": 26,
  "function": "func2"
}
```

### Field Descriptions
- **timestamp**: ISO 8601 formatted timestamp with timezone
- **log_id**: UUID for correlating related log entries within a request
- **level**: Log severity (DEBUG, INFO, WARN, ERROR)
- **message**: Human-readable log message with context
- **file**: Source code filename where log was generated
- **line**: Line number in the source file
- **function**: Function name where log was generated

## Log Levels

### Available Levels
1. **DEBUG**: Detailed diagnostic information (typically disabled in production)
2. **INFO**: General application flow and important events
3. **WARN**: Warning conditions that don't prevent operation
4. **ERROR**: Error conditions that may affect functionality

### Level Configuration
- **Default Level**: INFO (DEBUG messages filtered out)
- **Runtime Changes**: Log level can be changed without restart
- **Environment-Specific**: Different levels for development vs production

### Level Usage Examples
```go
// DEBUG: Detailed diagnostic info
utils.Logger.Debug(ctx, "Processing user data: %+v", userData)

// INFO: Normal application flow
utils.Logger.Info(ctx, "User logged in successfully: %s", userEmail)

// WARN: Concerning but non-critical
utils.Logger.Warn(ctx, "Rate limit approaching for IP: %s", clientIP)

// ERROR: Error conditions
utils.Logger.Error(ctx, "Database connection failed: %v", err)
```

## Request Correlation

### Log ID System
Each HTTP request receives a unique UUID that appears in all related log entries:

```json
{"log_id": "7b0f9e30-50ec-4c43-a7c5-a9dc260d2153", "message": "Started POST /api/v1/auth/login"}
{"log_id": "7b0f9e30-50ec-4c43-a7c5-a9dc260d2153", "message": "Attempting to login user"}
{"log_id": "7b0f9e30-50ec-4c43-a7c5-a9dc260d2153", "message": "Completed POST /api/v1/auth/login"}
```

### Request Lifecycle Tracking
1. **Request Start**: Method, path, IP, User-Agent
2. **Business Logic**: Authentication attempts, database operations, errors
3. **Request Completion**: Status code, duration, final log ID

## Common Log Patterns

### HTTP Request Logging
```json
// Request start
{"level": "INFO", "message": "Started POST /api/v1/auth/login | IP: ::1 | User-Agent: PostmanRuntime/7.44.1"}

// Request completion
{"level": "INFO", "message": "Completed POST /api/v1/auth/login | Status: 401 | Duration: 269.0677ms"}
```

### Authentication Events
```json
// Login attempt
{"level": "INFO", "message": "Attempting to login user with email: user@example.com"}

// Login success
{"level": "INFO", "message": "Logged in user successfully: {user_data}"}

// Login failure
{"level": "ERROR", "message": "Error validating password for user with email: user@example.com"}
```

### Error Logging
```json
// Application errors
{"level": "ERROR", "message": "Database connection failed: connection timeout"}

// Business logic errors
{"level": "ERROR", "message": "Invalid credentials provided for user: user@example.com"}
```

## Log Analysis

### Searching Logs
Use standard JSON processing tools to analyze logs:

```bash
# Find all ERROR level logs
grep '"level":"ERROR"' app-2025-08-04.log

# Extract all login attempts
grep "Attempting to login" app-2025-08-04.log

# Find logs for specific request
grep "7b0f9e30-50ec-4c43-a7c5-a9dc260d2153" app-2025-08-04.log

# Parse JSON for structured analysis
cat app-2025-08-04.log | jq '.message' | grep "login"
```

### Performance Analysis
```bash
# Extract request durations
cat app-2025-08-04.log | jq -r 'select(.message | contains("Duration")) | .message'

# Find slow requests (>1 second)
cat app-2025-08-04.log | jq 'select(.message | contains("Duration") and (. | match("Duration: ([0-9.]+)s") | .captures[0].string | tonumber) > 1)'
```

### Error Monitoring
```bash
# Count errors by type
cat app-2025-08-04.log | jq -r 'select(.level == "ERROR") | .message' | sort | uniq -c

# Find authentication failures
grep '"level":"ERROR"' app-2025-08-04.log | grep "password\|credentials"
```

## Configuration

### Logger Configuration
The logging system is configured in `backend/utils/logger.go`:

```go
type LoggerConfig struct {
    LogsDir        string    // Directory for log files (default: "logs")
    FilePrefix     string    // Log file prefix (default: "app")
    IncludeConsole bool      // Also log to console (default: true)
    MinLevel       LogLevel  // Minimum log level (default: INFO)
}
```

### Environment-Specific Settings
- **Development**: Console logging enabled, DEBUG level available
- **Production**: File-only logging, INFO level minimum
- **Testing**: Configurable for test isolation

## Docker Integration

### Volume Mounting
The logs directory is mounted as a Docker volume for persistence:

```yaml
# docker-compose.yml
volumes:
  - ./backend/logs:/app/logs
```

### Log Access
Access logs from the Docker container:

```bash
# View live logs
docker exec -it posteaze-backend tail -f /app/logs/app-$(date +%Y-%m-%d).log

# Copy logs from container
docker cp posteaze-backend:/app/logs ./local-logs
```

## Maintenance

### Log Retention
- **Manual Cleanup**: Old log files should be archived or deleted manually
- **Disk Space**: Monitor disk usage as logs accumulate daily
- **Backup Strategy**: Consider backing up logs for compliance or analysis

### Log Rotation Recommendations
```bash
# Archive logs older than 30 days
find backend/logs -name "app-*.log" -mtime +30 -exec gzip {} \;

# Delete archived logs older than 90 days
find backend/logs -name "app-*.log.gz" -mtime +90 -delete
```

### Performance Considerations
- **File I/O**: Logging is asynchronous to minimize request impact
- **JSON Parsing**: Use efficient tools (jq, grep) for log analysis
- **Storage**: Plan for approximately 10-100MB per day depending on traffic

## Troubleshooting

### Common Issues

#### Missing Log Files
- Check directory permissions (logs directory needs write access)
- Verify application has permission to create files
- Check disk space availability

#### Incomplete Logs
- Application crash may result in incomplete log entries
- Check for application restart events in logs
- Verify log file isn't being truncated by external processes

#### Performance Impact
- High log volume can impact application performance
- Consider increasing log level in production
- Monitor disk I/O if logging becomes a bottleneck

### Debug Commands
```bash
# Check log file permissions
ls -la backend/logs/

# Monitor real-time logging
tail -f backend/logs/app-$(date +%Y-%m-%d).log

# Validate JSON format
cat backend/logs/app-2025-08-04.log | jq empty
```

## Related Documentation

- [Logger Implementation](../utils/logger.go) - Core logging utility
- [Log Middleware](../middleware/README.md) - HTTP request logging middleware
- [Application Configuration](../utils/configs/README.md) - Configuration management
- [Docker Setup](../../docker-compose.yml) - Container log volume configuration