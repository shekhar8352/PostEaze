#!/bin/bash

# PostEaze Backend Test Execution Script
# Simple wrapper around the Go test runner

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
PACKAGE=""
VERBOSE=false
COVERAGE=false
COVERAGE_THRESHOLD=0

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
    echo "  benchmark            Run benchmark tests only"
    echo "  all                  Run all tests (default)"
    echo ""
    echo "Options:"
    echo "  -t, --type TYPE       Test type (unit|integration|benchmark|all)"
    echo "  -p, --package PKG     Run specific package tests"
    echo "  -v, --verbose         Enable verbose output"
    echo "  -c, --coverage        Generate coverage report"
    echo "  --threshold NUM       Coverage threshold percentage (default: 0)"
    echo "  -h, --help            Show this help message"
    echo ""
    echo "Package Options:"
    echo "  api                   API handler tests"
    echo "  business              Business logic tests"
    echo "  models                Model validation tests"
    echo "  utils                 Utility function tests"
    echo "  helpers               Test helper functions"
    echo "  integration           Integration tests"
    echo "  benchmarks            Benchmark tests"
    echo ""
    echo "Examples:"
    echo "  $0 --type unit --package api --verbose"
    echo "  $0 --type integration --coverage"
    echo "  $0 --coverage --threshold 80"
    echo "  $0 --type benchmark --verbose"
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -t|--type)
            TEST_TYPE="$2"
            shift 2
            ;;
        -p|--package)
            PACKAGE="$2"
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
        --threshold)
            COVERAGE_THRESHOLD="$2"
            shift 2
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
    unit|integration|benchmark|all)
        ;;
    *)
        print_error "Invalid test type: $TEST_TYPE"
        print_error "Valid options: unit, integration, benchmark, all"
        exit 1
        ;;
esac

print_header "PostEaze Backend Test Suite"
print_status "Test Type: $TEST_TYPE"
print_status "Package: ${PACKAGE:-all}"
print_status "Verbose: $VERBOSE"
print_status "Coverage: $COVERAGE"
if [[ "$COVERAGE_THRESHOLD" -gt 0 ]]; then
    print_status "Coverage Threshold: $COVERAGE_THRESHOLD%"
fi

# Change to backend directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKEND_DIR="$(dirname "$SCRIPT_DIR")"
cd "$BACKEND_DIR"

# Check if we're in the right directory
if [[ ! -f "go.mod" ]]; then
    print_error "go.mod not found. Please run this script from the backend directory or ensure the path is correct."
    exit 1
fi

# Check if test.go exists
if [[ ! -f "tests/test.go" ]]; then
    print_error "tests/test.go not found. Please ensure the test runner is available."
    exit 1
fi

# Build test command for the Go test runner
build_test_command() {
    local cmd=("go" "run" "tests/test.go")
    
    # Add test type
    cmd+=("-type" "$TEST_TYPE")
    
    # Add package if specified
    if [[ -n "$PACKAGE" ]]; then
        cmd+=("-package" "$PACKAGE")
    fi
    
    # Add flags
    if [[ "$VERBOSE" == true ]]; then
        cmd+=("-verbose")
    fi
    
    if [[ "$COVERAGE" == true ]]; then
        cmd+=("-coverage")
    fi
    
    if [[ "$COVERAGE_THRESHOLD" -gt 0 ]]; then
        cmd+=("-coverage-threshold" "$COVERAGE_THRESHOLD")
    fi
    
    echo "${cmd[@]}"
}

# Run tests using the Go test runner
run_tests() {
    local test_cmd=($(build_test_command))
    
    print_status "Running tests with Go test runner..."
    print_status "Command: ${test_cmd[*]}"
    
    # Run the tests
    "${test_cmd[@]}"
    local exit_code=$?
    
    if [[ $exit_code -eq 0 ]]; then
        print_success "All tests completed successfully!"
        return 0
    else
        print_error "Test execution failed!"
        return 1
    fi
}

# Main execution
main() {
    # Run tests using the Go test runner
    if run_tests; then
        print_success "Test execution completed successfully!"
        return 0
    else
        print_error "Test execution failed!"
        return 1
    fi
}

# Check for required tools
check_dependencies() {
    if ! command -v go >/dev/null 2>&1; then
        print_error "Go is required but not installed"
        exit 1
    fi
}

# Check dependencies before running
check_dependencies

# Run main function
main
exit $?