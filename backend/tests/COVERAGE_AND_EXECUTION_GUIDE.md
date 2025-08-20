# Test Coverage and Execution Guide

This guide provides detailed information about test execution, coverage reporting, and continuous integration setup for the PostEaze backend application.

## Table of Contents

1. [Test Execution Commands](#test-execution-commands)
2. [Coverage Reporting](#coverage-reporting)
3. [Test Discovery and Organization](#test-discovery-and-organization)
4. [Continuous Integration](#continuous-integration)
5. [Performance Testing](#performance-testing)
6. [Test Environment Setup](#test-environment-setup)
7. [Troubleshooting](#troubleshooting)

## Test Execution Commands

### Basic Test Execution

#### Run All Tests
```bash
# Run all tests in the project
go test ./tests/...

# Run with verbose output
go test -v ./tests/...

# Run with race detection
go test -race ./tests/...

# Run in short mode (skips integration tests)
go test -short ./tests/...
```

#### Run Specific Test Types
```bash
# Unit tests only
go test ./tests/unit/...

# Integration tests only
go test ./tests/integration/...

# Test utilities
go test ./tests/testutils/...

# Specific package
go test ./tests/unit/api/v1/...
```

#### Run Specific Tests
```bash
# Run specific test function
go test -run TestUserHandler_CreateUser ./tests/unit/api/v1/

# Run tests matching pattern
go test -run "TestAuth.*" ./tests/unit/api/v1/

# Run specific test suite
go test -run "TestAuthHandlerTestSuite" ./tests/unit/api/v1/
```

### Advanced Test Execution

#### Parallel Execution
```bash
# Run tests in parallel (default)
go test -parallel 4 ./tests/...

# Disable parallel execution
go test -parallel 1 ./tests/...

# Run with custom parallelism
go test -parallel 8 ./tests/...
```

#### Timeout Configuration
```bash
# Set test timeout (default: 10m)
go test -timeout 30s ./tests/unit/...
go test -timeout 5m ./tests/integration/...
go test -timeout 15m ./tests/...
```

#### Test Repetition
```bash
# Run tests multiple times to catch flaky tests
go test -count 10 ./tests/unit/api/v1/

# Run until failure
go test -count 100 ./tests/... || echo "Test failed"
```

## Coverage Reporting

### Basic Coverage Commands

#### Generate Coverage Report
```bash
# Generate coverage profile
go test -coverprofile=coverage.out ./tests/...

# Generate coverage with race detection
go test -race -coverprofile=coverage.out ./tests/...

# Generate coverage for specific packages
go test -coverprofile=coverage.out ./tests/unit/...
```

#### View Coverage Reports
```bash
# View coverage summary
go tool cover -func=coverage.out

# View coverage by package
go tool cover -func=coverage.out | grep -E "(package|total)"

# Generate HTML coverage report
go tool cover -html=coverage.out -o coverage.html

# Open HTML report in browser (macOS/Linux)
open coverage.html
# Windows
start coverage.html
```

### Advanced Coverage Analysis

#### Coverage by Test Type
```bash
# Unit test coverage
go test -coverprofile=unit_coverage.out ./tests/unit/...
go tool cover -func=unit_coverage.out

# Integration test coverage
go test -coverprofile=integration_coverage.out ./tests/integration/...
go tool cover -func=integration_coverage.out

# Combined coverage
go test -coverprofile=combined_coverage.out ./tests/...
go tool cover -func=combined_coverage.out
```

#### Coverage Thresholds
```bash
# Check if coverage meets threshold (80%)
go test -cover ./tests/... | grep "coverage:" | awk -F'[:%]' '{if($2<80) exit 1}'

# More detailed threshold checking
go test -coverprofile=coverage.out ./tests/...
COVERAGE=$(go tool cover -func=coverage.out | grep "total:" | awk '{print $3}' | sed 's/%//')
if (( $(echo "$COVERAGE < 80" | bc -l) )); then
    echo "Coverage $COVERAGE% is below threshold 80%"
    exit 1
fi
```

#### Coverage Exclusions
```bash
# Exclude vendor and generated files
go test -coverprofile=coverage.out ./tests/...
grep -v -E "(vendor/|\.pb\.go|_gen\.go)" coverage.out > filtered_coverage.out
go tool cover -func=filtered_coverage.out
```

### Coverage Reporting Scripts

#### Comprehensive Coverage Script (`scripts/coverage.sh`)
```bash
#!/bin/bash

set -e

# Configuration
COVERAGE_THRESHOLD=80
OUTPUT_DIR="coverage_reports"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")

# Create output directory
mkdir -p $OUTPUT_DIR

echo "Running comprehensive coverage analysis..."

# Unit test coverage
echo "Generating unit test coverage..."
go test -coverprofile=$OUTPUT_DIR/unit_coverage_$TIMESTAMP.out ./tests/unit/...
go tool cover -func=$OUTPUT_DIR/unit_coverage_$TIMESTAMP.out > $OUTPUT_DIR/unit_coverage_$TIMESTAMP.txt
go tool cover -html=$OUTPUT_DIR/unit_coverage_$TIMESTAMP.out -o $OUTPUT_DIR/unit_coverage_$TIMESTAMP.html

# Integration test coverage
echo "Generating integration test coverage..."
go test -coverprofile=$OUTPUT_DIR/integration_coverage_$TIMESTAMP.out ./tests/integration/...
go tool cover -func=$OUTPUT_DIR/integration_coverage_$TIMESTAMP.out > $OUTPUT_DIR/integration_coverage_$TIMESTAMP.txt
go tool cover -html=$OUTPUT_DIR/integration_coverage_$TIMESTAMP.out -o $OUTPUT_DIR/integration_coverage_$TIMESTAMP.html

# Combined coverage
echo "Generating combined coverage..."
go test -coverprofile=$OUTPUT_DIR/combined_coverage_$TIMESTAMP.out ./tests/...
go tool cover -func=$OUTPUT_DIR/combined_coverage_$TIMESTAMP.out > $OUTPUT_DIR/combined_coverage_$TIMESTAMP.txt
go tool cover -html=$OUTPUT_DIR/combined_coverage_$TIMESTAMP.out -o $OUTPUT_DIR/combined_coverage_$TIMESTAMP.html

# Check coverage threshold
COVERAGE=$(go tool cover -func=$OUTPUT_DIR/combined_coverage_$TIMESTAMP.out | grep "total:" | awk '{print $3}' | sed 's/%//')
echo "Total coverage: $COVERAGE%"

if (( $(echo "$COVERAGE < $COVERAGE_THRESHOLD" | bc -l) )); then
    echo "❌ Coverage $COVERAGE% is below threshold $COVERAGE_THRESHOLD%"
    exit 1
else
    echo "✅ Coverage $COVERAGE% meets threshold $COVERAGE_THRESHOLD%"
fi

echo "Coverage reports generated in $OUTPUT_DIR/"
echo "Open $OUTPUT_DIR/combined_coverage_$TIMESTAMP.html to view detailed report"
```

## Test Discovery and Organization

### Test Discovery Commands

#### List All Tests
```bash
# List all test functions
go test -list . ./tests/...

# List tests matching pattern
go test -list "TestAuth.*" ./tests/...

# List benchmarks
go test -list "Benchmark.*" ./tests/...
```

#### Test Organization Analysis
```bash
# Count tests by package
find ./tests -name "*_test.go" -exec grep -l "func Test" {} \; | \
    xargs -I {} dirname {} | sort | uniq -c

# Count test functions
find ./tests -name "*_test.go" -exec grep -h "^func Test" {} \; | wc -l

# List test suites
find ./tests -name "*_test.go" -exec grep -h "suite.Suite" {} \; | \
    grep -o "type [A-Za-z]*TestSuite" | sort | uniq
```

### Test Categorization

#### By Test Type
```bash
# Unit tests
find ./tests/unit -name "*_test.go" | wc -l

# Integration tests  
find ./tests/integration -name "*_test.go" | wc -l

# Test utilities
find ./tests/testutils -name "*_test.go" | wc -l
```

#### By Functionality
```bash
# API tests
find ./tests -path "*/api/*" -name "*_test.go" | wc -l

# Business logic tests
find ./tests -path "*/business/*" -name "*_test.go" | wc -l

# Model tests
find ./tests -path "*/models/*" -name "*_test.go" | wc -l

# Utility tests
find ./tests -path "*/utils/*" -name "*_test.go" | wc -l
```

## Continuous Integration

### GitHub Actions Workflow

#### Basic CI Workflow (`.github/workflows/test.yml`)
```yaml
name: Test Suite

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main, develop ]

jobs:
  unit-tests:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        go-version: [1.19, 1.20, 1.21]
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v3
      with:
        go-version: ${{ matrix.go-version }}
    
    - name: Cache Go modules
      uses: actions/cache@v3
      with:
        path: ~/go/pkg/mod
        key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
        restore-keys: |
          ${{ runner.os }}-go-
    
    - name: Install dependencies
      run: go mod download
    
    - name: Run unit tests
      run: |
        go test -race -coverprofile=coverage.out ./tests/unit/...
        go tool cover -func=coverage.out
    
    - name: Upload coverage to Codecov
      uses: codecov/codecov-action@v3
      with:
        file: ./coverage.out

  integration-tests:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:13
        env:
          POSTGRES_PASSWORD: postgres
          POSTGRES_DB: postease_test
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
        ports:
          - 5432:5432
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v3
      with:
        go-version: 1.21
    
    - name: Run integration tests
      env:
        DATABASE_URL: postgres://postgres:postgres@localhost:5432/postease_test?sslmode=disable
      run: |
        go test -race ./tests/integration/...

  benchmark-tests:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v3
      with:
        go-version: 1.21
    
    - name: Run benchmarks
      run: |
        go test -bench=. -benchmem ./tests/... > benchmark_results.txt
        cat benchmark_results.txt
    
    - name: Upload benchmark results
      uses: actions/upload-artifact@v3
      with:
        name: benchmark-results
        path: benchmark_results.txt
```

### CI Test Scripts

#### Pre-commit Test Script (`scripts/pre-commit-test.sh`)
```bash
#!/bin/bash

set -e

echo "Running pre-commit tests..."

# Format check
echo "Checking code formatting..."
if [ "$(gofmt -s -l . | wc -l)" -gt 0 ]; then
    echo "Code is not formatted. Run 'gofmt -s -w .'"
    exit 1
fi

# Lint check
echo "Running linter..."
golangci-lint run ./...

# Unit tests with coverage
echo "Running unit tests..."
go test -race -coverprofile=coverage.out ./tests/unit/...

# Check coverage threshold
COVERAGE=$(go tool cover -func=coverage.out | grep "total:" | awk '{print $3}' | sed 's/%//')
if (( $(echo "$COVERAGE < 80" | bc -l) )); then
    echo "Coverage $COVERAGE% is below threshold 80%"
    exit 1
fi

# Integration tests (if not in CI)
if [ -z "$CI" ]; then
    echo "Running integration tests..."
    go test -short ./tests/integration/...
fi

echo "All pre-commit tests passed!"
```

## Performance Testing

### Benchmark Tests

#### Running Benchmarks
```bash
# Run all benchmarks
go test -bench=. ./tests/...

# Run benchmarks with memory profiling
go test -bench=. -benchmem ./tests/...

# Run specific benchmark
go test -bench=BenchmarkUserCreation ./tests/unit/business/v1/

# Run benchmarks multiple times for accuracy
go test -bench=. -count=5 ./tests/...
```

#### Benchmark Comparison
```bash
# Generate baseline benchmarks
go test -bench=. ./tests/... > baseline_benchmarks.txt

# After changes, compare performance
go test -bench=. ./tests/... > current_benchmarks.txt

# Use benchcmp to compare (install with: go install golang.org/x/tools/cmd/benchcmp@latest)
benchcmp baseline_benchmarks.txt current_benchmarks.txt
```

### Load Testing

#### Concurrent Test Execution
```bash
# Run tests with high parallelism to simulate load
go test -parallel 20 ./tests/integration/...

# Stress test with multiple iterations
for i in {1..10}; do
    echo "Iteration $i"
    go test -race ./tests/integration/...
done
```

## Test Environment Setup

### Environment Variables

#### Test Environment Configuration
```bash
# Set test environment
export MODE=test
export DATABASE_DRIVER=sqlite3
export DATABASE_URL=":memory:"
export JWT_SECRET=test-secret-key
export LOG_LEVEL=debug

# Run tests with environment
go test ./tests/...
```

#### Environment Setup Script (`scripts/setup-test-env.sh`)
```bash
#!/bin/bash

# Test environment setup
export MODE=test
export DATABASE_DRIVER=sqlite3
export DATABASE_URL=":memory:"
export JWT_SECRET=test-secret-key-for-testing-only
export JWT_EXPIRATION=1h
export REFRESH_TOKEN_EXPIRATION=24h
export LOG_LEVEL=debug
export LOG_OUTPUT=stdout

# Test database setup (for integration tests)
export TEST_DATABASE_DRIVER=postgres
export TEST_DATABASE_URL=postgres://postgres:postgres@localhost:5432/postease_test?sslmode=disable

# External service mocking
export MOCK_EXTERNAL_SERVICES=true
export EXTERNAL_API_TIMEOUT=5s

echo "Test environment configured"
echo "MODE: $MODE"
echo "DATABASE_DRIVER: $DATABASE_DRIVER"
echo "LOG_LEVEL: $LOG_LEVEL"
```

### Database Setup

#### Test Database Initialization
```bash
# Create test database (PostgreSQL)
createdb postease_test

# Run migrations
go run migrations/migrate.go -env=test -action=up

# Seed test data
go run scripts/seed-test-data.go
```

#### Docker Test Environment
```yaml
# docker-compose.test.yml
version: '3.8'
services:
  postgres-test:
    image: postgres:13
    environment:
      POSTGRES_DB: postease_test
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
    ports:
      - "5433:5432"
    volumes:
      - postgres_test_data:/var/lib/postgresql/data

volumes:
  postgres_test_data:
```

```bash
# Start test database
docker-compose -f docker-compose.test.yml up -d

# Run tests against Docker database
export TEST_DATABASE_URL=postgres://postgres:postgres@localhost:5433/postease_test?sslmode=disable
go test ./tests/integration/...

# Cleanup
docker-compose -f docker-compose.test.yml down -v
```

## Troubleshooting

### Common Issues

#### Test Failures

##### Race Conditions
```bash
# Run with race detection to identify race conditions
go test -race ./tests/...

# If race conditions are found, fix the code and re-test
go test -race -count=10 ./tests/unit/api/v1/
```

##### Flaky Tests
```bash
# Run tests multiple times to identify flaky tests
go test -count=100 ./tests/unit/api/v1/ | grep FAIL

# Run specific flaky test repeatedly
go test -run TestFlakyFunction -count=50 ./tests/unit/api/v1/
```

##### Database Connection Issues
```bash
# Check database connectivity
psql -h localhost -U postgres -d postease_test -c "SELECT 1;"

# Use SQLite for faster unit tests
export DATABASE_DRIVER=sqlite3
export DATABASE_URL=":memory:"
go test ./tests/unit/...
```

#### Coverage Issues

##### Low Coverage
```bash
# Identify uncovered code
go test -coverprofile=coverage.out ./tests/...
go tool cover -html=coverage.out -o coverage.html
# Open coverage.html to see uncovered lines

# Generate coverage by function
go tool cover -func=coverage.out | grep -E "0\.0%|[0-4][0-9]\.[0-9]%"
```

##### Coverage Calculation Problems
```bash
# Exclude test files from coverage
go test -coverprofile=coverage.out ./tests/...
grep -v "_test.go" coverage.out > filtered_coverage.out
go tool cover -func=filtered_coverage.out
```

#### Performance Issues

##### Slow Tests
```bash
# Profile test execution
go test -cpuprofile=test_cpu.prof ./tests/...
go tool pprof test_cpu.prof

# Memory profiling
go test -memprofile=test_mem.prof ./tests/...
go tool pprof test_mem.prof

# Identify slow tests
go test -v ./tests/... 2>&1 | grep -E "PASS|FAIL" | sort -k2 -nr
```

### Debug Mode

#### Verbose Test Output
```bash
# Run with maximum verbosity
go test -v -x ./tests/...

# Show test output even for passing tests
go test -v ./tests/... | grep -E "(RUN|PASS|FAIL)"

# Debug specific test
go test -v -run TestSpecificFunction ./tests/unit/api/v1/
```

#### Test Debugging
```bash
# Run single test with debugging
go test -v -run TestDebugMe ./tests/unit/api/v1/ -args -debug

# Use delve debugger
dlv test ./tests/unit/api/v1/ -- -test.run TestDebugMe
```

### Monitoring and Metrics

#### Test Execution Metrics
```bash
# Track test execution time
time go test ./tests/...

# Generate test execution report
go test -json ./tests/... > test_results.json

# Parse test results
cat test_results.json | jq '.Action, .Test, .Elapsed' | grep -E "(pass|fail)"
```

#### Coverage Trends
```bash
# Track coverage over time
echo "$(date),$(go test -cover ./tests/... | grep coverage | awk '{print $5}')" >> coverage_history.csv

# Generate coverage trend report
gnuplot -e "
set terminal png;
set output 'coverage_trend.png';
set datafile separator ',';
set xdata time;
set timefmt '%Y-%m-%d';
plot 'coverage_history.csv' using 1:2 with lines title 'Coverage %'
"
```

This comprehensive guide provides all the tools and knowledge needed to effectively execute tests, monitor coverage, and maintain a robust testing pipeline for the PostEaze backend application.