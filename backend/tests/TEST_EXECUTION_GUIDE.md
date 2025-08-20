# PostEaze Backend Test Execution Guide

This guide provides comprehensive information about running tests in the PostEaze backend application.

## Table of Contents

1. [Quick Start](#quick-start)
2. [Test Organization](#test-organization)
3. [Test Execution Scripts](#test-execution-scripts)
4. [Coverage Reporting](#coverage-reporting)
5. [Continuous Integration](#continuous-integration)
6. [Test Discovery](#test-discovery)
7. [Performance Testing](#performance-testing)
8. [Troubleshooting](#troubleshooting)

## Quick Start

### Run All Tests
```bash
# Using the test script (recommended)
./scripts/test.sh

# Using Make
make test

# Using Go directly
go test ./tests/...
```

### Run Specific Test Types
```bash
# Unit tests only
./scripts/test.sh --type unit
make test-unit

# Integration tests only
./scripts/test.sh --type integration
make test-integration

# With coverage
./scripts/test.sh --coverage
make test-coverage
```

### Windows Users
```powershell
# PowerShell script
.\scripts\test.ps1 -Type all

# With coverage
.\scripts\test.ps1 -Coverage

# Specific suite
.\scripts\test.ps1 -Type unit -Suite api
```

## Test Organization

### Directory Structure
```
tests/
├── config/                 # Test configuration files
│   ├── test.env            # Test environment variables
│   ├── test_config.json    # Test configuration
│   └── test_execution.json # Test execution configuration
├── integration/            # Integration tests
│   ├── auth_integration_test.go
│   ├── log_integration_test.go
│   ├── database_integration_test.go
│   └── end_to_end_integration_test.go
├── testutils/             # Test utilities and helpers
│   ├── auth.go            # Authentication test utilities
│   ├── database.go        # Database test utilities
│   ├── fixtures.go        # Test data fixtures
│   ├── http.go            # HTTP test utilities
│   ├── suite.go           # Base test suite
│   └── mocks/             # Mock implementations
└── unit/                  # Unit tests
    ├── api/               # API handler tests
    ├── business/          # Business logic tests
    ├── models/            # Model tests
    └── utils/             # Utility tests
```

### Test Types

#### Unit Tests
- **Location**: `tests/unit/`
- **Purpose**: Test individual components in isolation
- **Suites**:
  - `api`: HTTP handler tests
  - `business`: Business logic tests
  - `models`: Data model tests
  - `utils`: Utility function tests

#### Integration Tests
- **Location**: `tests/integration/`
- **Purpose**: Test component interactions and workflows
- **Suites**:
  - `auth`: Authentication workflows
  - `log`: Log management workflows
  - `database`: Database operations
  - `e2e`: End-to-end scenarios

## Test Execution Scripts

### Main Test Script (`scripts/test.sh`)

The primary test execution script with comprehensive options:

```bash
./scripts/test.sh [OPTIONS]
```

#### Options
- `--type TYPE`: Test type (unit|integration|all)
- `--suite SUITE`: Specific test suite
- `--verbose`: Enable verbose output
- `--coverage`: Generate coverage report
- `--database TYPE`: Database type (sqlite|postgres)
- `--parallel`: Run tests in parallel
- `--format FORMAT`: Output format (standard|json|junit)
- `--threshold NUM`: Coverage threshold percentage
- `--benchmark`: Run benchmark tests
- `--race`: Enable race detection
- `--short`: Run tests in short mode

#### Examples
```bash
# Run all tests with coverage
./scripts/test.sh --coverage

# Run unit tests for API handlers
./scripts/test.sh --type unit --suite api --verbose

# Run integration tests with PostgreSQL
./scripts/test.sh --type integration --database postgres

# Run tests with race detection and benchmarks
./scripts/test.sh --race --benchmark

# CI-style execution
./scripts/test.sh --coverage --race --parallel --threshold 80
```

### PowerShell Script (`scripts/test.ps1`)

Windows-compatible PowerShell version:

```powershell
.\scripts\test.ps1 [OPTIONS]
```

#### Examples
```powershell
# Run all tests with coverage
.\scripts\test.ps1 -Coverage

# Run specific suite
.\scripts\test.ps1 -Type unit -Suite api -Verbose

# Run with race detection
.\scripts\test.ps1 -Race -Benchmark
```

### Makefile Commands

Convenient Make targets for common operations:

```bash
# Basic test commands
make test                    # Run all tests
make test-unit              # Run unit tests
make test-integration       # Run integration tests
make test-coverage          # Run with coverage

# Specific suites
make test-api               # API handler tests
make test-business          # Business logic tests
make test-auth              # Auth integration tests
make test-e2e               # End-to-end tests

# Advanced options
make test-race              # With race detection
make test-bench             # Benchmark tests
make test-parallel          # Parallel execution
make test-ci                # CI-style execution

# Coverage thresholds
make test-coverage-80       # 80% threshold
make test-coverage-90       # 90% threshold

# Database-specific
make test-sqlite            # SQLite integration tests
make test-postgres          # PostgreSQL integration tests

# Output formats
make test-json              # JSON output
make test-junit             # JUnit XML output
```

## Coverage Reporting

### Coverage Generation
```bash
# Generate coverage report
./scripts/test.sh --coverage

# With specific threshold
./scripts/test.sh --coverage --threshold 90

# Coverage for specific test type
./scripts/test.sh --type unit --coverage
```

### Coverage Files Generated
- `coverage.out`: Raw coverage data
- `coverage.html`: HTML coverage report
- `coverage_func.txt`: Function-level coverage summary

### Coverage Thresholds
- **Global**: 80% (configurable)
- **Package**: 70% (configurable)
- **Function**: 60% (configurable)

### Viewing Coverage Reports
```bash
# Open HTML report in browser
open coverage.html

# View text summary
go tool cover -func=coverage.out

# Generate custom HTML report
go tool cover -html=coverage.out -o custom_coverage.html
```

## Continuous Integration

### GitHub Actions Workflow

The CI pipeline includes multiple jobs:

1. **Unit Tests**: Fast unit test execution
2. **Integration Tests (SQLite)**: Integration tests with SQLite
3. **Integration Tests (PostgreSQL)**: Integration tests with PostgreSQL
4. **Benchmark Tests**: Performance benchmarking
5. **Code Quality**: Linting and static analysis
6. **Coverage Report**: Comprehensive coverage analysis
7. **Test Matrix**: Multi-OS and Go version testing

### CI Configuration

The workflow is triggered on:
- Push to `main` or `develop` branches
- Pull requests to `main` or `develop` branches
- Changes to backend code or CI configuration

### CI Commands
```bash
# Run CI-style tests locally
make test-ci

# Unit tests for CI
make test-ci-unit

# Integration tests for CI
make test-ci-integration
```

### Artifacts Generated
- Test results (JSON/JUnit format)
- Coverage reports (HTML/XML)
- Benchmark results
- Performance profiles
- Code quality reports

## Test Discovery

### Automatic Test Discovery

The test discovery utility automatically finds and catalogs tests:

```bash
# Discover all tests
go run scripts/test-discovery.go tests/

# Output test information as JSON
go run scripts/test-discovery.go tests/ > test_catalog.json
```

### Test Discovery Output
```json
{
  "total_tests": 45,
  "total_benchmarks": 8,
  "total_suites": 4,
  "tests": [
    {
      "name": "TestAuthHandler_Signup",
      "package": "api_test",
      "file": "tests/unit/api/auth_test.go",
      "line": 25,
      "type": "test",
      "tags": ["auth", "signup"],
      "description": "Test user signup with valid data"
    }
  ],
  "suites": [
    {
      "name": "AuthIntegrationSuite",
      "package": "integration_test",
      "file": "tests/integration/auth_integration_test.go",
      "tests": [...],
      "description": "Authentication integration test suite"
    }
  ],
  "packages": ["api_test", "business_test", "models_test"]
}
```

## Performance Testing

### Benchmark Tests
```bash
# Run all benchmarks
./scripts/test.sh --benchmark

# Benchmarks with memory profiling
make perf-test

# CPU and memory profiling
make perf-profile
```

### Benchmark Output
```
BenchmarkAuthHandler_Login-8         1000    1234567 ns/op    1024 B/op    10 allocs/op
BenchmarkJWTGeneration-8             5000     234567 ns/op     512 B/op     5 allocs/op
```

### Performance Regression Testing
```bash
# Compare with baseline (in CI)
go test -bench=. -benchmem ./... > current_benchmarks.txt
# Compare with stored baseline benchmarks
```

## Troubleshooting

### Common Issues

#### 1. Database Connection Issues
```bash
# Check PostgreSQL connection
psql -h localhost -U postgres -d postease_test

# Use SQLite fallback
./scripts/test.sh --database sqlite
```

#### 2. Coverage Threshold Failures
```bash
# Check current coverage
go tool cover -func=coverage.out | tail -1

# Run with lower threshold
./scripts/test.sh --coverage --threshold 70
```

#### 3. Race Condition Detection
```bash
# Run with race detection
./scripts/test.sh --race

# Fix race conditions in code
go run -race main.go
```

#### 4. Test Timeout Issues
```bash
# Run with shorter timeout
go test -timeout 5m ./...

# Run in short mode
./scripts/test.sh --short
```

### Debug Mode

Enable verbose output for debugging:
```bash
# Verbose test output
./scripts/test.sh --verbose

# Go test verbose mode
go test -v ./tests/...

# Debug specific test
go test -v -run TestSpecificFunction ./tests/unit/api/
```

### Environment Issues

#### Missing Dependencies
```bash
# Check Go installation
go version

# Install missing tools
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
```

#### Environment Variables
```bash
# Check test environment
env | grep -E "(MODE|DATABASE|JWT)"

# Set test environment
export MODE=test
export DATABASE_DRIVER=sqlite3
export DATABASE_URL=":memory:"
```

### Performance Issues

#### Slow Tests
```bash
# Run in short mode
./scripts/test.sh --short

# Profile test execution
go test -cpuprofile=test_cpu.prof ./...

# Analyze profile
go tool pprof test_cpu.prof
```

#### Memory Issues
```bash
# Memory profiling
go test -memprofile=test_mem.prof ./...

# Check for memory leaks
go test -race -memprofile=mem.prof ./...
```

## Best Practices

### Test Execution
1. **Run tests frequently** during development
2. **Use appropriate test types** for different scenarios
3. **Monitor coverage** but don't obsess over 100%
4. **Run race detection** regularly
5. **Use parallel execution** for faster feedback

### CI/CD Integration
1. **Run tests on every commit**
2. **Use matrix testing** for multiple environments
3. **Generate artifacts** for debugging
4. **Set appropriate timeouts**
5. **Monitor performance** over time

### Debugging
1. **Use verbose output** when tests fail
2. **Run specific tests** to isolate issues
3. **Check environment variables**
4. **Verify database connections**
5. **Use profiling tools** for performance issues

## Additional Resources

- [Go Testing Documentation](https://golang.org/pkg/testing/)
- [Testify Documentation](https://github.com/stretchr/testify)
- [Go Coverage Tools](https://golang.org/doc/code.html#Testing)
- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [PostEaze Backend README](../README.md)