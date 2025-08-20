#!/bin/bash

# PostEaze Backend Test Execution Script
# Comprehensive test runner with coverage reporting and test discovery

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Default values
TEST_TYPE="all"
TEST_SUITE=""
VERBOSE=false
COVERAGE=false
CLEANUP=true
DATABASE_TYPE="sqlite"
PARALLEL=false
OUTPUT_FORMAT="standard"
COVERAGE_THRESHOLD=80
BENCHMARK=false
RACE_DETECTION=false
SHORT=false

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_header() {
    echo -e "${CYAN}========================================${NC}"
    echo -e "${CYAN} $1${NC}"
    echo -e "${CYAN}========================================${NC}"
}

# Function to show usage
show_usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Test Types:"
    echo "  unit                  Run unit tests only"
    echo "  integration          Run integration tests only"
    echo "  all                  Run all tests (default)"
    echo ""
    echo "Options:"
    echo "  -t, --type TYPE       Test type (unit|integration|all)"
    echo "  -s, --suite SUITE     Run specific test suite"
    echo "  -v, --verbose         Enable verbose output"
    echo "  -c, --coverage        Generate coverage report"
    echo "  -d, --database TYPE   Database type for integration tests (sqlite|postgres)"
    echo "  -p, --parallel        Run tests in parallel"
    echo "  -f, --format FORMAT   Output format (standard|json|junit)"
    echo "  --threshold NUM       Coverage threshold percentage (default: 80)"
    echo "  --benchmark           Run benchmark tests"
    echo "  --race                Enable race detection"
    echo "  --short               Run tests in short mode"
    echo "  --no-cleanup          Skip cleanup after tests"
    echo "  -h, --help            Show this help message"
    echo ""
    echo "Unit Test Suites:"
    echo "  api                   API handler tests"
    echo "  business              Business logic tests"
    echo "  models                Model validation tests"
    echo "  utils                 Utility function tests"
    echo ""
    echo "Integration Test Suites:"
    echo "  auth                  Authentication integration tests"
    echo "  log                   Log management integration tests"
    echo "  database              Database integration tests"
    echo "  e2e                   End-to-end integration tests"
    echo ""
    echo "Examples:"
    echo "  $0 --type unit --suite api --verbose"
    echo "  $0 --type integration --coverage --database postgres"
    echo "  $0 --coverage --threshold 90 --parallel"
    echo "  $0 --benchmark --race"
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -t|--type)
            TEST_TYPE="$2"
            shift 2
            ;;
        -s|--suite)
            TEST_SUITE="$2"
            shift 2
            ;;
        -v|--verbose)
            VERBOSE=true
            shift
            ;;
        -c|--coverage)
            COVERAGE=true
            shift
            ;;
        -d|--database)
            DATABASE_TYPE="$2"
            shift 2
            ;;
        -p|--parallel)
            PARALLEL=true
            shift
            ;;
        -f|--format)
            OUTPUT_FORMAT="$2"
            shift 2
            ;;
        --threshold)
            COVERAGE_THRESHOLD="$2"
            shift 2
            ;;
        --benchmark)
            BENCHMARK=true
            shift
            ;;
        --race)
            RACE_DETECTION=true
            shift
            ;;
        --short)
            SHORT=true
            shift
            ;;
        --no-cleanup)
            CLEANUP=false
            shift
            ;;
        -h|--help)
            show_usage
            exit 0
            ;;
        *)
            print_error "Unknown option: $1"
            show_usage
            exit 1
            ;;
    esac
done

# Validate test type
case $TEST_TYPE in
    unit|integration|all)
        ;;
    *)
        print_error "Invalid test type: $TEST_TYPE"
        print_error "Valid options: unit, integration, all"
        exit 1
        ;;
esac

# Validate output format
case $OUTPUT_FORMAT in
    standard|json|junit)
        ;;
    *)
        print_error "Invalid output format: $OUTPUT_FORMAT"
        print_error "Valid options: standard, json, junit"
        exit 1
        ;;
esac

print_header "PostEaze Backend Test Suite"
print_status "Test Type: $TEST_TYPE"
print_status "Test Suite: ${TEST_SUITE:-all}"
print_status "Database Type: $DATABASE_TYPE"
print_status "Verbose: $VERBOSE"
print_status "Coverage: $COVERAGE"
print_status "Parallel: $PARALLEL"
print_status "Output Format: $OUTPUT_FORMAT"
print_status "Coverage Threshold: $COVERAGE_THRESHOLD%"

# Change to backend directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKEND_DIR="$(dirname "$SCRIPT_DIR")"
cd "$BACKEND_DIR"

# Check if we're in the right directory
if [[ ! -f "go.mod" ]]; then
    print_error "go.mod not found. Please run this script from the backend directory or ensure the path is correct."
    exit 1
fi

# Setup environment variables
setup_environment() {
    print_status "Setting up test environment..."
    
    export MODE="test"
    export BASE_CONFIG_PATH="./tests/config"
    export JWT_ACCESS_SECRET="test-access-secret-key-$(date +%s)"
    export JWT_REFRESH_SECRET="test-refresh-secret-key-$(date +%s)"
    
    if [[ "$DATABASE_TYPE" == "postgres" ]]; then
        # Check if PostgreSQL is available
        if command -v psql >/dev/null 2>&1; then
            export DATABASE_DRIVER="postgres"
            export DATABASE_URL="postgres://postgres:password@localhost:5432/postease_test?sslmode=disable"
            print_status "Using PostgreSQL database"
        else
            print_warning "PostgreSQL not found, falling back to SQLite"
            DATABASE_TYPE="sqlite"
        fi
    fi
    
    if [[ "$DATABASE_TYPE" == "sqlite" ]]; then
        export DATABASE_DRIVER="sqlite3"
        export DATABASE_URL=":memory:"
        print_status "Using SQLite in-memory database"
    fi
    
    # Create temporary directories
    export TEST_LOG_DIR=$(mktemp -d -t postease_test_logs_XXXXXX)
    export COVERAGE_DIR=$(mktemp -d -t postease_coverage_XXXXXX)
    print_status "Created temporary directories"
}

# Discover test packages
discover_test_packages() {
    local test_type="$1"
    local test_suite="$2"
    local packages=()
    
    case $test_type in
        unit)
            if [[ -n "$test_suite" ]]; then
                case $test_suite in
                    api)
                        packages+=("./tests/unit/api/...")
                        ;;
                    business)
                        packages+=("./tests/unit/business/...")
                        ;;
                    models)
                        packages+=("./tests/unit/models/...")
                        ;;
                    utils)
                        packages+=("./tests/unit/utils/...")
                        ;;
                    *)
                        print_error "Invalid unit test suite: $test_suite"
                        print_error "Valid options: api, business, models, utils"
                        exit 1
                        ;;
                esac
            else
                packages+=("./tests/unit/...")
            fi
            ;;
        integration)
            if [[ -n "$test_suite" ]]; then
                case $test_suite in
                    auth)
                        packages+=("./tests/integration/ -run TestAuthIntegrationSuite")
                        ;;
                    log)
                        packages+=("./tests/integration/ -run TestLogIntegrationSuite")
                        ;;
                    database)
                        packages+=("./tests/integration/ -run TestDatabaseIntegrationSuite")
                        ;;
                    e2e)
                        packages+=("./tests/integration/ -run TestEndToEndIntegrationSuite")
                        ;;
                    *)
                        print_error "Invalid integration test suite: $test_suite"
                        print_error "Valid options: auth, log, database, e2e"
                        exit 1
                        ;;
                esac
            else
                packages+=("./tests/integration/...")
            fi
            ;;
        all)
            if [[ -n "$test_suite" ]]; then
                print_warning "Test suite specified with 'all' type. Running all tests."
            fi
            packages+=("./tests/unit/..." "./tests/integration/...")
            ;;
    esac
    
    echo "${packages[@]}"
}

# Build test command
build_test_command() {
    local packages=("$@")
    local cmd=("go" "test")
    
    # Add packages
    cmd+=("${packages[@]}")
    
    # Add flags
    if [[ "$VERBOSE" == true ]]; then
        cmd+=("-v")
    fi
    
    if [[ "$COVERAGE" == true ]]; then
        cmd+=("-coverprofile=$COVERAGE_DIR/coverage.out")
        cmd+=("-covermode=atomic")
    fi
    
    if [[ "$PARALLEL" == true ]]; then
        cmd+=("-parallel" "4")
    fi
    
    if [[ "$RACE_DETECTION" == true ]]; then
        cmd+=("-race")
    fi
    
    if [[ "$SHORT" == true ]]; then
        cmd+=("-short")
    fi
    
    if [[ "$BENCHMARK" == true ]]; then
        cmd+=("-bench=.")
        cmd+=("-benchmem")
    fi
    
    # Add output format
    case $OUTPUT_FORMAT in
        json)
            cmd+=("-json")
            ;;
        junit)
            # Note: Go doesn't natively support JUnit, but we can convert later
            cmd+=("-v")
            ;;
    esac
    
    # Add timeout
    cmd+=("-timeout" "15m")
    
    echo "${cmd[@]}"
}

# Run tests
run_tests() {
    local packages=($(discover_test_packages "$TEST_TYPE" "$TEST_SUITE"))
    local test_cmd=($(build_test_command "${packages[@]}"))
    
    print_status "Discovered test packages: ${packages[*]}"
    print_status "Running tests..."
    print_status "Command: ${test_cmd[*]}"
    
    # Create output file for JUnit format
    local output_file=""
    if [[ "$OUTPUT_FORMAT" == "junit" ]]; then
        output_file="$COVERAGE_DIR/test_results.xml"
    fi
    
    # Run the tests
    local exit_code=0
    if [[ "$OUTPUT_FORMAT" == "junit" ]]; then
        # Capture output for JUnit conversion
        "${test_cmd[@]}" 2>&1 | tee "$COVERAGE_DIR/test_output.txt"
        exit_code=${PIPESTATUS[0]}
        
        # Convert to JUnit format (simplified)
        convert_to_junit "$COVERAGE_DIR/test_output.txt" "$output_file"
    else
        "${test_cmd[@]}"
        exit_code=$?
    fi
    
    if [[ $exit_code -eq 0 ]]; then
        print_success "All tests passed!"
        return 0
    else
        print_error "Some tests failed!"
        return 1
    fi
}

# Convert test output to JUnit format (simplified)
convert_to_junit() {
    local input_file="$1"
    local output_file="$2"
    
    print_status "Converting test output to JUnit format..."
    
    cat > "$output_file" << 'EOF'
<?xml version="1.0" encoding="UTF-8"?>
<testsuites>
    <testsuite name="PostEaze Backend Tests" tests="0" failures="0" errors="0" time="0">
        <!-- Test results would be parsed and inserted here -->
        <!-- This is a simplified placeholder implementation -->
    </testsuite>
</testsuites>
EOF
    
    print_status "JUnit report generated: $output_file"
}

# Generate coverage report
generate_coverage_report() {
    if [[ "$COVERAGE" == true && -f "$COVERAGE_DIR/coverage.out" ]]; then
        print_status "Generating coverage report..."
        
        # Generate HTML coverage report
        go tool cover -html="$COVERAGE_DIR/coverage.out" -o "$COVERAGE_DIR/coverage.html"
        
        # Generate coverage summary
        local coverage_func="$COVERAGE_DIR/coverage_func.txt"
        go tool cover -func="$COVERAGE_DIR/coverage.out" > "$coverage_func"
        
        # Extract total coverage percentage
        local total_coverage=$(tail -1 "$coverage_func" | awk '{print $3}' | sed 's/%//')
        
        print_status "Coverage Summary:"
        echo "----------------------------------------"
        cat "$coverage_func"
        echo "----------------------------------------"
        
        # Check coverage threshold
        if (( $(echo "$total_coverage >= $COVERAGE_THRESHOLD" | bc -l) )); then
            print_success "Coverage ($total_coverage%) meets threshold ($COVERAGE_THRESHOLD%)"
        else
            print_warning "Coverage ($total_coverage%) below threshold ($COVERAGE_THRESHOLD%)"
            if [[ "$CI" == "true" ]]; then
                print_error "Coverage threshold not met in CI environment"
                return 1
            fi
        fi
        
        print_success "Coverage reports generated:"
        print_success "  HTML: $COVERAGE_DIR/coverage.html"
        print_success "  Text: $coverage_func"
        print_success "  Raw: $COVERAGE_DIR/coverage.out"
        
        # Copy coverage files to standard locations for CI
        cp "$COVERAGE_DIR/coverage.out" "./coverage.out" 2>/dev/null || true
        cp "$COVERAGE_DIR/coverage.html" "./coverage.html" 2>/dev/null || true
    fi
}

# Run benchmarks
run_benchmarks() {
    if [[ "$BENCHMARK" == true ]]; then
        print_status "Running benchmarks..."
        
        local benchmark_output="$COVERAGE_DIR/benchmark_results.txt"
        go test -bench=. -benchmem ./... > "$benchmark_output" 2>&1
        
        print_status "Benchmark Results:"
        echo "----------------------------------------"
        cat "$benchmark_output"
        echo "----------------------------------------"
        
        print_success "Benchmark results saved: $benchmark_output"
    fi
}

# Setup test database (for PostgreSQL)
setup_database() {
    if [[ "$DATABASE_TYPE" == "postgres" && ("$TEST_TYPE" == "integration" || "$TEST_TYPE" == "all") ]]; then
        print_status "Setting up PostgreSQL test database..."
        
        # Check if database exists and create if needed
        if ! psql -lqt | cut -d \| -f 1 | grep -qw postease_test; then
            print_status "Creating test database..."
            createdb postease_test || {
                print_warning "Failed to create database, it may already exist"
            }
        fi
        
        print_success "PostgreSQL test database ready"
    fi
}

# Cleanup function
cleanup() {
    if [[ "$CLEANUP" == true ]]; then
        print_status "Cleaning up..."
        
        # Remove temporary directories
        if [[ -n "$TEST_LOG_DIR" && -d "$TEST_LOG_DIR" ]]; then
            rm -rf "$TEST_LOG_DIR"
        fi
        
        if [[ -n "$COVERAGE_DIR" && -d "$COVERAGE_DIR" ]]; then
            # Keep coverage files if coverage was requested
            if [[ "$COVERAGE" != true ]]; then
                rm -rf "$COVERAGE_DIR"
            else
                print_status "Coverage files preserved in: $COVERAGE_DIR"
            fi
        fi
        
        print_success "Cleanup completed"
    fi
}

# Trap to ensure cleanup runs on exit
trap cleanup EXIT

# Main execution
main() {
    local exit_code=0
    
    # Setup
    setup_environment
    setup_database
    
    # Run tests
    if ! run_tests; then
        exit_code=1
    fi
    
    # Generate coverage report
    if ! generate_coverage_report; then
        exit_code=1
    fi
    
    # Run benchmarks if requested
    run_benchmarks
    
    # Print summary
    print_header "Test Execution Summary"
    if [[ $exit_code -eq 0 ]]; then
        print_success "All test operations completed successfully!"
    else
        print_error "Test execution completed with failures!"
    fi
    
    return $exit_code
}

# Check for required tools
check_dependencies() {
    local missing_deps=()
    
    if ! command -v go >/dev/null 2>&1; then
        missing_deps+=("go")
    fi
    
    if [[ "$COVERAGE" == true ]] && ! command -v bc >/dev/null 2>&1; then
        missing_deps+=("bc")
    fi
    
    if [[ "$DATABASE_TYPE" == "postgres" ]] && ! command -v psql >/dev/null 2>&1; then
        print_warning "PostgreSQL tools not found, will fallback to SQLite"
    fi
    
    if [[ ${#missing_deps[@]} -gt 0 ]]; then
        print_error "Missing required dependencies: ${missing_deps[*]}"
        exit 1
    fi
}

# Check dependencies before running
check_dependencies

# Run main function
main
exit $?