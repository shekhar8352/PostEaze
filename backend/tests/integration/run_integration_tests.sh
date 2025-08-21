#!/bin/bash

# Integration Test Runner Script for PostEaze Backend
# This script sets up the environment and runs integration tests

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Default values
TEST_SUITE=""
VERBOSE=false
COVERAGE=false
CLEANUP=true
DATABASE_TYPE="sqlite"
PARALLEL=false

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

# Function to show usage
show_usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -s, --suite SUITE     Run specific test suite (auth|log|database|e2e|all)"
    echo "  -v, --verbose         Enable verbose output"
    echo "  -c, --coverage        Generate coverage report"
    echo "  -d, --database TYPE   Database type to use (sqlite|postgres)"
    echo "  -p, --parallel        Run tests in parallel"
    echo "  --no-cleanup          Skip cleanup after tests"
    echo "  -h, --help            Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 --suite auth --verbose"
    echo "  $0 --coverage --database postgres"
    echo "  $0 --suite all --parallel"
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
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

# Set default test suite if not specified
if [[ -z "$TEST_SUITE" ]]; then
    TEST_SUITE="all"
fi

# Validate test suite
case $TEST_SUITE in
    auth|log|database|e2e|all)
        ;;
    *)
        print_error "Invalid test suite: $TEST_SUITE"
        print_error "Valid options: auth, log, database, e2e, all"
        exit 1
        ;;
esac

print_status "Starting PostEaze Backend Integration Tests"
print_status "Test Suite: $TEST_SUITE"
print_status "Database Type: $DATABASE_TYPE"
print_status "Verbose: $VERBOSE"
print_status "Coverage: $COVERAGE"
print_status "Parallel: $PARALLEL"

# Change to backend directory
cd "$(dirname "$0")/../.."

# Check if we're in the right directory
if [[ ! -f "go.mod" ]]; then
    print_error "go.mod not found. Please run this script from the backend directory or ensure the path is correct."
    exit 1
fi

# Setup environment variables
setup_environment() {
    print_status "Setting up test environment..."
    
    export MODE="dev"
    export BASE_CONFIG_PATH="./tests/config"
    export JWT_ACCESS_SECRET="integration-test-access-secret-key-$(date +%s)"
    export JWT_REFRESH_SECRET="integration-test-refresh-secret-key-$(date +%s)"
    
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
    
    # Create temporary log directory
    export LOG_DIR=$(mktemp -d -t postease_integration_logs_XXXXXX)
    print_status "Created temporary log directory: $LOG_DIR"
}

# Setup test database (for PostgreSQL)
setup_database() {
    if [[ "$DATABASE_TYPE" == "postgres" ]]; then
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

# Build test command
build_test_command() {
    local cmd="go test"
    
    # Add test path based on suite
    case $TEST_SUITE in
        auth)
            cmd="$cmd ./tests/integration/ -run TestAuthIntegrationSuite"
            ;;
        log)
            cmd="$cmd ./tests/integration/ -run TestLogIntegrationSuite"
            ;;
        database)
            cmd="$cmd ./tests/integration/ -run TestDatabaseIntegrationSuite"
            ;;
        e2e)
            cmd="$cmd ./tests/integration/ -run TestEndToEndIntegrationSuite"
            ;;
        all)
            cmd="$cmd ./tests/integration/..."
            ;;
    esac
    
    # Add flags
    if [[ "$VERBOSE" == true ]]; then
        cmd="$cmd -v"
    fi
    
    if [[ "$COVERAGE" == true ]]; then
        cmd="$cmd -coverprofile=integration_coverage.out"
    fi
    
    if [[ "$PARALLEL" == true ]]; then
        cmd="$cmd -parallel 4"
    fi
    
    # Add timeout
    cmd="$cmd -timeout 10m"
    
    echo "$cmd"
}

# Run tests
run_tests() {
    local test_cmd=$(build_test_command)
    
    print_status "Running integration tests..."
    print_status "Command: $test_cmd"
    
    # Run the tests
    if eval "$test_cmd"; then
        print_success "All integration tests passed!"
        return 0
    else
        print_error "Some integration tests failed!"
        return 1
    fi
}

# Generate coverage report
generate_coverage_report() {
    if [[ "$COVERAGE" == true && -f "integration_coverage.out" ]]; then
        print_status "Generating coverage report..."
        
        # Generate HTML coverage report
        go tool cover -html=integration_coverage.out -o integration_coverage.html
        
        # Show coverage summary
        go tool cover -func=integration_coverage.out | tail -1
        
        print_success "Coverage report generated: integration_coverage.html"
    fi
}

# Cleanup function
cleanup() {
    if [[ "$CLEANUP" == true ]]; then
        print_status "Cleaning up..."
        
        # Remove temporary log directory
        if [[ -n "$LOG_DIR" && -d "$LOG_DIR" ]]; then
            rm -rf "$LOG_DIR"
            print_status "Removed temporary log directory"
        fi
        
        # Remove coverage files if not requested
        if [[ "$COVERAGE" != true ]]; then
            rm -f integration_coverage.out integration_coverage.html
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
    
    # Generate coverage report if requested
    generate_coverage_report
    
    # Print summary
    if [[ $exit_code -eq 0 ]]; then
        print_success "Integration test run completed successfully!"
    else
        print_error "Integration test run completed with failures!"
    fi
    
    return $exit_code
}

# Run main function
main
exit $?