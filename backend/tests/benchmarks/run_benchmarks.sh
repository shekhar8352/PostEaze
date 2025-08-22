#!/bin/bash

# Benchmark runner script for PostEaze backend
# Usage: ./run_benchmarks.sh [category] [options]

set -e

BENCHMARK_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$BENCHMARK_DIR"

# Default values
CATEGORY="all"
BENCHTIME="1s"
COUNT=1
OUTPUT_FILE=""
MEMORY_PROFILE=false

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print usage
usage() {
    echo "Usage: $0 [category] [options]"
    echo ""
    echo "Categories:"
    echo "  all          Run all benchmarks (default)"
    echo "  auth         Run authentication benchmarks"
    echo "  api          Run API endpoint benchmarks"
    echo "  database     Run database operation benchmarks"
    echo "  user         Run user operation benchmarks"
    echo "  jwt          Run JWT-related benchmarks"
    echo "  password     Run password hashing benchmarks"
    echo "  json         Run JSON marshaling benchmarks"
    echo ""
    echo "Options:"
    echo "  -t, --time DURATION    Benchmark time (default: 1s)"
    echo "  -c, --count COUNT      Number of benchmark runs (default: 1)"
    echo "  -o, --output FILE      Save results to file"
    echo "  -m, --memory           Include memory profiling"
    echo "  -v, --verbose          Verbose output"
    echo "  -h, --help             Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 auth -t 5s -m                    # Run auth benchmarks for 5s with memory profiling"
    echo "  $0 api -c 3 -o api_results.txt      # Run API benchmarks 3 times, save to file"
    echo "  $0 all -t 2s -v                     # Run all benchmarks for 2s with verbose output"
}

# Function to run benchmarks
run_benchmark() {
    local pattern="$1"
    local description="$2"
    
    echo -e "${BLUE}Running $description benchmarks...${NC}"
    
    local cmd="go test -bench=$pattern -benchtime=$BENCHTIME -run=^$ -count=$COUNT"
    
    if [ "$MEMORY_PROFILE" = true ]; then
        cmd="$cmd -benchmem"
    fi
    
    if [ -n "$OUTPUT_FILE" ]; then
        cmd="$cmd | tee -a $OUTPUT_FILE"
    fi
    
    echo -e "${YELLOW}Command: $cmd${NC}"
    eval $cmd
    echo ""
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        all|auth|api|database|user|jwt|password|json)
            CATEGORY="$1"
            shift
            ;;
        -t|--time)
            BENCHTIME="$2"
            shift 2
            ;;
        -c|--count)
            COUNT="$2"
            shift 2
            ;;
        -o|--output)
            OUTPUT_FILE="$2"
            shift 2
            ;;
        -m|--memory)
            MEMORY_PROFILE=true
            shift
            ;;
        -v|--verbose)
            set -x
            shift
            ;;
        -h|--help)
            usage
            exit 0
            ;;
        *)
            echo -e "${RED}Unknown option: $1${NC}"
            usage
            exit 1
            ;;
    esac
done

# Clear output file if specified
if [ -n "$OUTPUT_FILE" ]; then
    echo "# Benchmark Results - $(date)" > "$OUTPUT_FILE"
    echo "# Category: $CATEGORY, Time: $BENCHTIME, Count: $COUNT" >> "$OUTPUT_FILE"
    echo "" >> "$OUTPUT_FILE"
fi

echo -e "${GREEN}Starting benchmarks...${NC}"
echo -e "${GREEN}Category: $CATEGORY${NC}"
echo -e "${GREEN}Benchmark time: $BENCHTIME${NC}"
echo -e "${GREEN}Run count: $COUNT${NC}"
if [ "$MEMORY_PROFILE" = true ]; then
    echo -e "${GREEN}Memory profiling: enabled${NC}"
fi
if [ -n "$OUTPUT_FILE" ]; then
    echo -e "${GREEN}Output file: $OUTPUT_FILE${NC}"
fi
echo ""

# Run benchmarks based on category
case $CATEGORY in
    all)
        run_benchmark "." "all"
        ;;
    auth)
        run_benchmark "BenchmarkJWT|BenchmarkPassword|BenchmarkSignupFlow|BenchmarkLoginFlow|BenchmarkRefreshTokenFlow" "authentication"
        ;;
    api)
        run_benchmark "BenchmarkSignupEndpoint|BenchmarkLoginEndpoint|BenchmarkRefreshTokenEndpoint|BenchmarkLogoutEndpoint|BenchmarkConcurrentAPIRequests" "API endpoint"
        ;;
    database)
        run_benchmark "BenchmarkUserCreation|BenchmarkUserLookup|BenchmarkRefreshTokenInsertion|BenchmarkTeamCreation|BenchmarkTeamMemberAddition|BenchmarkConcurrentDatabaseOperations|BenchmarkDatabaseTransaction|BenchmarkBulkUserCreation" "database operation"
        ;;
    user)
        run_benchmark "BenchmarkCompleteUserSignupWorkflow|BenchmarkCompleteUserLoginWorkflow|BenchmarkUserAuthenticationFlow|BenchmarkTeamCreationWorkflow|BenchmarkUserSessionManagement|BenchmarkUserDataValidation|BenchmarkConcurrentUserOperations" "user operation"
        ;;
    jwt)
        run_benchmark "BenchmarkJWT" "JWT"
        ;;
    password)
        run_benchmark "BenchmarkPassword" "password"
        ;;
    json)
        run_benchmark "BenchmarkJSON" "JSON"
        ;;
    *)
        echo -e "${RED}Unknown category: $CATEGORY${NC}"
        usage
        exit 1
        ;;
esac

echo -e "${GREEN}Benchmarks completed!${NC}"

if [ -n "$OUTPUT_FILE" ]; then
    echo -e "${GREEN}Results saved to: $OUTPUT_FILE${NC}"
fi