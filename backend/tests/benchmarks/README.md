# Benchmark Tests

This directory contains performance benchmark tests for the PostEaze backend application. The benchmarks measure the performance of critical operations including authentication, database operations, API endpoints, and user workflows.

## Benchmark Categories

### 1. Authentication Benchmarks (`auth_benchmark_test.go`)
- JWT token generation and validation
- Password hashing and verification
- Complete authentication workflows (signup, login, refresh token)

### 2. API Endpoint Benchmarks (`api_benchmark_test.go`)
- HTTP endpoint performance
- JSON marshaling/unmarshaling
- Concurrent API request handling
- Authentication endpoint performance

### 3. Database Operation Benchmarks (`database_benchmark_test.go`)
- User CRUD operations
- Database lookups and queries
- Transaction performance
- Concurrent database operations

### 4. User Operation Benchmarks (`user_operations_benchmark_test.go`)
- Complete user workflows
- Team creation and management
- Session management
- Data validation performance

## Running Benchmarks

### Run All Benchmarks
```bash
cd backend/tests/benchmarks
go test -bench=. -benchmem
```

### Run Specific Benchmark Category
```bash
# Authentication benchmarks only
go test -bench=BenchmarkJWT -benchmem

# Database benchmarks only
go test -bench=BenchmarkDatabase -benchmem

# API benchmarks only
go test -bench=BenchmarkAPI -benchmem
```

### Run Specific Benchmark
```bash
# Run JWT generation benchmark
go test -bench=BenchmarkJWTGeneration -benchmem

# Run with specific duration
go test -bench=BenchmarkJWTGeneration -benchtime=10s -benchmem
```

### Generate Benchmark Report
```bash
# Save benchmark results to file
go test -bench=. -benchmem > benchmark_results.txt

# Compare benchmark results
go test -bench=. -benchmem -count=5 > current_benchmarks.txt
```

## Understanding Benchmark Results

Benchmark output format:
```
BenchmarkJWTGeneration-8    	  100000	     12345 ns/op	    1024 B/op	      10 allocs/op
```

- `BenchmarkJWTGeneration-8`: Benchmark name with GOMAXPROCS value
- `100000`: Number of iterations run
- `12345 ns/op`: Nanoseconds per operation
- `1024 B/op`: Bytes allocated per operation
- `10 allocs/op`: Number of allocations per operation

## Performance Targets

### Authentication Operations
- JWT Generation: < 50,000 ns/op
- JWT Validation: < 30,000 ns/op
- Password Hashing: < 100,000 ns/op (bcrypt is intentionally slow)
- Password Validation: < 100,000 ns/op

### Database Operations
- User Creation: < 10,000 ns/op
- User Lookup: < 5,000 ns/op
- Token Operations: < 8,000 ns/op

### API Endpoints
- Simple Endpoints: < 100,000 ns/op
- Authentication Endpoints: < 200,000 ns/op
- JSON Operations: < 5,000 ns/op

## Benchmark Best Practices

### Writing Benchmarks
1. Use `b.ResetTimer()` after setup
2. Use `b.StopTimer()` and `b.StartTimer()` for setup within loops
3. Clean up data between iterations for consistent results
4. Use realistic test data
5. Test both success and failure scenarios

### Running Benchmarks
1. Run benchmarks on a quiet system
2. Run multiple times and compare results
3. Use `-benchmem` to track memory allocations
4. Use `-benchtime` to run for specific duration
5. Use `-count` to run multiple iterations

### Interpreting Results
1. Focus on ns/op for performance
2. Monitor B/op and allocs/op for memory efficiency
3. Compare results over time to detect regressions
4. Consider both average and worst-case performance

## Continuous Integration

These benchmarks can be integrated into CI/CD pipelines to:
- Detect performance regressions
- Track performance improvements
- Set performance budgets
- Generate performance reports

### Example CI Integration
```bash
# Run benchmarks and fail if performance degrades
go test -bench=. -benchmem -benchtime=5s > current.txt
# Compare with baseline and fail if regression > 20%
```

## Troubleshooting

### Common Issues
1. **Inconsistent Results**: Run on a quiet system, increase benchmark time
2. **Memory Leaks**: Check for proper cleanup in benchmark code
3. **Database Errors**: Ensure test database is properly initialized
4. **Token Errors**: Check environment variables for JWT secrets

### Environment Setup
Ensure the following environment variables are set for benchmarks:
```bash
export JWT_ACCESS_SECRET="test-secret-key-for-benchmarking"
export JWT_REFRESH_SECRET="test-refresh-secret-key-for-benchmarking"
```

## Contributing

When adding new benchmarks:
1. Follow the existing naming convention
2. Include proper setup and cleanup
3. Add documentation for new benchmark categories
4. Update performance targets if needed
5. Test benchmarks locally before submitting

## Monitoring

Regular benchmark monitoring helps:
- Identify performance bottlenecks
- Track optimization efforts
- Plan capacity requirements
- Ensure consistent user experience

Run benchmarks regularly and track trends over time to maintain optimal application performance.