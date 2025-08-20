# Integration Tests

This directory contains comprehensive integration tests for the PostEaze backend application. These tests verify end-to-end scenarios with real database connections, file system operations, and complete API workflows.

## Test Suites

### 1. Authentication Integration Tests (`auth_integration_test.go`)
Tests complete authentication flows with real database operations:
- Complete authentication flow (signup → login → refresh → logout)
- Team user signup flow with team creation
- Invalid authentication scenarios
- Concurrent authentication requests
- Database transaction handling during authentication

### 2. Log Integration Tests (`log_integration_test.go`)
Tests log retrieval with real file system operations:
- Log retrieval by date with actual log files
- Log retrieval by log ID across multiple files
- Error handling for non-existent dates/log IDs
- Invalid date format handling
- Log retrieval across multiple dates
- Concurrent log retrieval operations
- File system operations testing

### 3. Database Integration Tests (`database_integration_test.go`)
Tests database operations with actual database connections:
- User CRUD operations with real database
- Team operations with database transactions
- Refresh token operations
- Concurrent database operations
- Database connection pooling
- Transaction isolation testing
- Database error handling
- Performance testing with large datasets

### 4. End-to-End Integration Tests (`end_to_end_integration_test.go`)
Tests complete user journeys combining all components:
- Complete user journey (signup → database verification → login → log access → logout)
- Team user complete workflow
- Error recovery scenarios
- Concurrent end-to-end operations

## Prerequisites

Before running integration tests, ensure you have:

1. **Test Database**: A PostgreSQL test database running (or SQLite for simpler setup)
2. **Test Configuration**: Proper test configuration files in `../config/`
3. **Log Directory**: Write permissions for creating temporary log files
4. **Environment Variables**: Required environment variables set

## Running the Tests

### Run All Integration Tests
```bash
cd backend
go test ./tests/integration/... -v
```

### Run Specific Test Suite
```bash
# Authentication integration tests
go test ./tests/integration/ -run TestAuthIntegrationSuite -v

# Log integration tests
go test ./tests/integration/ -run TestLogIntegrationSuite -v

# Database integration tests
go test ./tests/integration/ -run TestDatabaseIntegrationSuite -v

# End-to-end integration tests
go test ./tests/integration/ -run TestEndToEndIntegrationSuite -v
```

### Run Specific Test Case
```bash
# Run a specific test method
go test ./tests/integration/ -run TestAuthIntegrationSuite/TestCompleteAuthenticationFlow -v
```

### Run with Coverage
```bash
go test ./tests/integration/... -v -coverprofile=integration_coverage.out
go tool cover -html=integration_coverage.out -o integration_coverage.html
```

## Test Configuration

### Environment Variables
The tests require these environment variables:
```bash
# Database configuration
DATABASE_DRIVER=postgres  # or sqlite3
DATABASE_URL=postgres://user:password@localhost:5432/postease_test?sslmode=disable

# JWT secrets for testing
JWT_ACCESS_SECRET=test-access-secret-key
JWT_REFRESH_SECRET=test-refresh-secret-key

# Log directory (optional, tests create temp directories)
LOG_DIR=/path/to/test/logs

# Test mode
MODE=dev
BASE_CONFIG_PATH=./tests/config
```

### Test Database Setup
For PostgreSQL:
```sql
CREATE DATABASE postease_test;
CREATE USER test_user WITH PASSWORD 'test_password';
GRANT ALL PRIVILEGES ON DATABASE postease_test TO test_user;
```

For SQLite (simpler setup):
The tests can use SQLite in-memory databases for faster execution.

## Test Data

### Database Test Data
- Tests create and clean up their own test data
- Each test starts with a clean database state
- Fixtures are created programmatically during tests

### Log Test Data
- Tests create temporary log directories
- Sample log files are generated with realistic log entries
- Log files are cleaned up after tests complete

## Test Features

### Real Database Operations
- Tests use actual database connections
- Database transactions are tested for consistency
- Connection pooling and concurrency are verified
- Error scenarios are tested with real database constraints

### Real File System Operations
- Tests create and read actual log files
- File system permissions and access are tested
- Concurrent file access scenarios are verified
- Error handling for missing files is tested

### Complete API Workflows
- Tests exercise complete request/response cycles
- Authentication middleware integration is tested
- Error responses are verified with actual API handlers
- JSON serialization/deserialization is tested end-to-end

### Concurrency Testing
- Multiple concurrent requests are tested
- Database connection pooling under load is verified
- File system access concurrency is tested
- Race conditions are identified and verified

## Test Isolation

### Database Isolation
- Each test suite gets a fresh database connection
- Test data is cleaned up between tests
- Transactions are used to ensure data consistency
- Database state is reset before each test

### File System Isolation
- Each test suite creates its own temporary directories
- Log files are isolated per test run
- Cleanup ensures no test artifacts remain
- Original system state is restored after tests

### Environment Isolation
- Test-specific environment variables are set
- Original environment is restored after tests
- Configuration isolation prevents test interference

## Debugging Integration Tests

### Enable Verbose Logging
```bash
go test ./tests/integration/... -v -args -test.v
```

### Debug Database Operations
Set environment variable for SQL logging:
```bash
export DEBUG_SQL=true
go test ./tests/integration/ -run TestDatabaseIntegrationSuite -v
```

### Debug Log File Operations
Tests log the paths of created temporary directories and files.

### Debug API Requests/Responses
Tests log request/response details when assertions fail.

## Performance Considerations

### Test Execution Time
- Integration tests are slower than unit tests due to real I/O operations
- Database tests may take longer with actual database connections
- File system tests depend on disk I/O performance

### Resource Usage
- Tests create temporary files and directories
- Database connections are pooled but limited
- Memory usage increases with concurrent test execution

### Optimization Tips
- Use SQLite in-memory databases for faster execution
- Run tests in parallel where possible
- Clean up resources promptly to avoid resource exhaustion

## Continuous Integration

### CI/CD Pipeline Integration
```yaml
# Example GitHub Actions configuration
- name: Run Integration Tests
  run: |
    # Setup test database
    docker run -d --name postgres-test -e POSTGRES_PASSWORD=test -p 5432:5432 postgres:13
    
    # Wait for database to be ready
    sleep 10
    
    # Run integration tests
    go test ./tests/integration/... -v
  env:
    DATABASE_URL: postgres://postgres:test@localhost:5432/postgres?sslmode=disable
    JWT_ACCESS_SECRET: ci-test-access-secret
    JWT_REFRESH_SECRET: ci-test-refresh-secret
```

### Test Reporting
- Use `-json` flag for machine-readable test output
- Generate coverage reports for integration test coverage
- Archive test artifacts (logs, coverage reports) in CI

## Troubleshooting

### Common Issues

1. **Database Connection Failures**
   - Verify database is running and accessible
   - Check connection string format
   - Ensure test database exists and has proper permissions

2. **File Permission Errors**
   - Ensure write permissions in test directories
   - Check disk space availability
   - Verify temporary directory creation permissions

3. **Port Conflicts**
   - Ensure test ports are not in use
   - Use different ports for test and development environments

4. **Environment Variable Issues**
   - Verify all required environment variables are set
   - Check for typos in variable names
   - Ensure proper escaping of special characters

### Getting Help
- Check test logs for detailed error messages
- Enable verbose logging for more debugging information
- Review test setup and teardown methods for proper resource management
- Verify test dependencies and prerequisites are met

## Contributing

When adding new integration tests:

1. Follow the existing test suite patterns
2. Ensure proper setup and teardown
3. Use real database and file system operations
4. Test both success and error scenarios
5. Include concurrent operation testing where relevant
6. Document any new test requirements or setup steps
7. Ensure tests are deterministic and can run in any order