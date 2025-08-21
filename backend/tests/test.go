package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// TestRunner handles test execution with simple configuration
type TestRunner struct {
	Type              string  // unit, integration, benchmark, all
	Package           string  // specific package to test
	Coverage          bool    // generate coverage
	CoverageThreshold float64 // minimum coverage threshold
	Verbose           bool    // verbose output
	ReportFormat      string  // text, html, json
	ReportFile        string  // output file for reports
	Config            *TestConfig
}

// TestConfig holds basic test configuration
type TestConfig struct {
	DatabaseURL string
	JWTSecret   string
	LogLevel    string
}

// TestResult represents the result of test execution
type TestResult struct {
	Package    string
	Tests      int
	Passed     int
	Failed     int
	Skipped    int
	Coverage   float64
	Duration   time.Duration
	Output     string
	Benchmarks []BenchmarkResult
}

// BenchmarkResult represents benchmark test results
type BenchmarkResult struct {
	Name        string
	Iterations  int
	NsPerOp     int64
	BytesPerOp  int64
	AllocsPerOp int64
}

// TestSummary aggregates all test results
type TestSummary struct {
	Results         []TestResult
	TotalTests      int
	TotalPassed     int
	TotalFailed     int
	TotalSkipped    int
	OverallCoverage float64
	TotalDuration   time.Duration
	Success         bool
	FailedTests     []FailedTest
	Timestamp       time.Time
}

// FailedTest represents a failed test with details
type FailedTest struct {
	Package     string
	TestName    string
	Error       string
	Output      string
	Suggestions []string
}

// CoverageReport represents coverage information for reporting
type CoverageReport struct {
	Summary     *TestSummary
	PackageData []PackageCoverage
	Timestamp   time.Time
}

// PackageCoverage represents coverage data for a single package
type PackageCoverage struct {
	Package    string
	Coverage   float64
	Files      []FileCoverage
	Statements int
	Covered    int
}

// FileCoverage represents coverage data for a single file
type FileCoverage struct {
	File       string
	Coverage   float64
	Statements int
	Covered    int
}

// BenchmarkComparison represents benchmark comparison data
type BenchmarkComparison struct {
	Current  BenchmarkResult
	Previous *BenchmarkResult
	Change   BenchmarkChange
}

// BenchmarkChange represents the change in benchmark performance
type BenchmarkChange struct {
	NsPerOpChange     float64 // percentage change
	BytesPerOpChange  float64
	AllocsPerOpChange float64
	Improved          bool
}

func main() {
	runner := &TestRunner{}
	var help bool
	
	// Define command-line flags
	flag.StringVar(&runner.Type, "type", "all", "Test type to run: unit, integration, benchmark, all")
	flag.StringVar(&runner.Package, "package", "", "Specific package to test (optional)")
	flag.BoolVar(&runner.Coverage, "coverage", false, "Generate coverage report")
	flag.Float64Var(&runner.CoverageThreshold, "coverage-threshold", 0.0, "Minimum coverage threshold (0.0-100.0)")
	flag.BoolVar(&runner.Verbose, "verbose", false, "Verbose output for debugging")
	flag.StringVar(&runner.ReportFormat, "report-format", "text", "Report format: text, html, json")
	flag.StringVar(&runner.ReportFile, "report-file", "", "Output file for reports (optional)")
	flag.BoolVar(&help, "help", false, "Show usage information")
	flag.Parse()
	
	if help {
		printUsage()
		return
	}

	// Validate coverage threshold
	if runner.CoverageThreshold < 0 || runner.CoverageThreshold > 100 {
		fmt.Printf("Error: coverage threshold must be between 0.0 and 100.0\n")
		os.Exit(1)
	}

	// Load configuration
	config, err := loadTestConfig()
	if err != nil {
		fmt.Printf("Error loading test configuration: %v\n", err)
		os.Exit(1)
	}
	runner.Config = config

	// Run tests and get summary
	summary, err := runner.Run()
	if err != nil {
		fmt.Printf("Test execution failed: %v\n", err)
		os.Exit(1)
	}

	// Print summary and exit with appropriate code
	runner.printSummary(summary)
	if !summary.Success {
		os.Exit(1)
	}
}

// Run executes tests based on the configuration and returns a summary
func (tr *TestRunner) Run() (*TestSummary, error) {
	startTime := time.Now()
	summary := &TestSummary{
		Results:     make([]TestResult, 0),
		Success:     true,
		FailedTests: make([]FailedTest, 0),
		Timestamp:   startTime,
	}

	if tr.Verbose {
		fmt.Printf("=== Starting test execution ===\n")
		fmt.Printf("Test type: %s\n", tr.Type)
		if tr.Package != "" {
			fmt.Printf("Package: %s\n", tr.Package)
		}
		if tr.Coverage {
			fmt.Printf("Coverage enabled")
			if tr.CoverageThreshold > 0 {
				fmt.Printf(" (threshold: %.1f%%)", tr.CoverageThreshold)
			}
			fmt.Println()
		}
		fmt.Println()
	}
	
	var err error
	switch tr.Type {
	case "unit":
		err = tr.runUnitTests(summary)
	case "integration":
		err = tr.runIntegrationTests(summary)
	case "benchmark":
		err = tr.runBenchmarkTests(summary)
	case "all":
		err = tr.runAllTests(summary)
	default:
		return nil, fmt.Errorf("unknown test type: %s", tr.Type)
	}

	// Calculate totals
	summary.TotalDuration = time.Since(startTime)
	for _, result := range summary.Results {
		summary.TotalTests += result.Tests
		summary.TotalPassed += result.Passed
		summary.TotalFailed += result.Failed
		summary.TotalSkipped += result.Skipped
	}

	// Calculate overall coverage
	if tr.Coverage && len(summary.Results) > 0 {
		totalCoverage := 0.0
		coverageCount := 0
		for _, result := range summary.Results {
			if result.Coverage > 0 {
				totalCoverage += result.Coverage
				coverageCount++
			}
		}
		if coverageCount > 0 {
			summary.OverallCoverage = totalCoverage / float64(coverageCount)
		}
	}

	// Check if tests passed and coverage threshold is met
	if summary.TotalFailed > 0 {
		summary.Success = false
	}
	if tr.CoverageThreshold > 0 && summary.OverallCoverage < tr.CoverageThreshold {
		summary.Success = false
		if tr.Verbose {
			fmt.Printf("Coverage threshold not met: %.1f%% < %.1f%%\n", 
				summary.OverallCoverage, tr.CoverageThreshold)
		}
	}

	return summary, err
}

// runUnitTests executes unit tests for api, business, models, and utils
func (tr *TestRunner) runUnitTests(summary *TestSummary) error {
	packages := []string{"api", "business", "models", "utils"}
	
	if tr.Package != "" {
		packages = []string{tr.Package}
	}
	
	if tr.Verbose {
		fmt.Printf("Running unit tests for packages: %v\n", packages)
	}
	
	for _, pkg := range packages {
		result, err := tr.runTestsForPackage(pkg)
		if result != nil {
			summary.Results = append(summary.Results, *result)
		}
		if err != nil {
			if tr.Verbose {
				fmt.Printf("Unit tests failed for package %s: %v\n", pkg, err)
			}
			// Continue with other packages even if one fails
		}
	}
	
	return nil
}

// runIntegrationTests executes integration tests
func (tr *TestRunner) runIntegrationTests(summary *TestSummary) error {
	if tr.Verbose {
		fmt.Println("Running integration tests...")
	}
	
	result, err := tr.runTestsForPackage("integration")
	if result != nil {
		summary.Results = append(summary.Results, *result)
	}
	return err
}

// runBenchmarkTests executes benchmark tests
func (tr *TestRunner) runBenchmarkTests(summary *TestSummary) error {
	if tr.Verbose {
		fmt.Println("Running benchmark tests...")
	}
	
	result, err := tr.runBenchmarksForPackage("benchmarks")
	if result != nil {
		summary.Results = append(summary.Results, *result)
	}
	return err
}

// runAllTests executes all test types
func (tr *TestRunner) runAllTests(summary *TestSummary) error {
	if tr.Verbose {
		fmt.Println("Running all test types...")
	}
	
	// Run unit tests
	if err := tr.runUnitTests(summary); err != nil {
		if tr.Verbose {
			fmt.Printf("Unit tests completed with errors: %v\n", err)
		}
	}
	
	// Run integration tests
	if err := tr.runIntegrationTests(summary); err != nil {
		if tr.Verbose {
			fmt.Printf("Integration tests completed with errors: %v\n", err)
		}
	}
	
	// Run benchmark tests
	if err := tr.runBenchmarkTests(summary); err != nil {
		if tr.Verbose {
			fmt.Printf("Benchmark tests completed with errors: %v\n", err)
		}
	}
	
	return nil
}

// runTestsForPackage runs tests for a specific package and returns results
func (tr *TestRunner) runTestsForPackage(pkg string) (*TestResult, error) {
	startTime := time.Now()
	result := &TestResult{
		Package: pkg,
	}
	
	testDir := filepath.Join("tests", pkg)
	
	// Check if test directory exists
	if _, err := os.Stat(testDir); os.IsNotExist(err) {
		if tr.Verbose {
			fmt.Printf("No tests found for package: %s\n", pkg)
		}
		return nil, nil
	}
	
	// Check if there are any test files in the directory
	matches, err := filepath.Glob(filepath.Join(testDir, "*_test.go"))
	if err != nil {
		return nil, fmt.Errorf("error checking for test files: %v", err)
	}
	
	// If no test files found in main directory, check subdirectories (like v1/)
	if len(matches) == 0 {
		// Check for common subdirectories like v1/
		subDirs := []string{"v1"}
		for _, subDir := range subDirs {
			subDirPath := filepath.Join(testDir, subDir)
			if _, err := os.Stat(subDirPath); err == nil {
				subMatches, err := filepath.Glob(filepath.Join(subDirPath, "*_test.go"))
				if err != nil {
					return nil, fmt.Errorf("error checking for test files in %s: %v", subDirPath, err)
				}
				if len(subMatches) > 0 {
					// Update testDir to point to the subdirectory that has tests
					testDir = subDirPath
					matches = subMatches
					break
				}
			}
		}
	}
	
	if len(matches) == 0 {
		if tr.Verbose {
			fmt.Printf("No test files found in package: %s\n", pkg)
		}
		return nil, nil
	}
	
	// Build go test command
	args := []string{"test", "-json"}
	
	if tr.Verbose {
		args = append(args, "-v")
	}
	
	if tr.Coverage {
		args = append(args, "-cover")
	}
	
	// Add the test directory as a package path
	// Convert absolute testDir back to relative path for go test
	relativeTestDir := strings.TrimPrefix(testDir, "tests/")
	if relativeTestDir == testDir {
		// testDir didn't start with "tests/", so it's already relative
		args = append(args, "./"+testDir)
	} else {
		args = append(args, "./tests/"+relativeTestDir)
	}
	
	// Execute tests
	cmd := exec.Command("go", args...)
	// cmd.Dir should be the backend directory (current working directory when test.go is run)
	
	// Capture output
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	// Set test environment variables
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("TEST_DATABASE_URL=%s", tr.Config.DatabaseURL),
		fmt.Sprintf("TEST_JWT_SECRET=%s", tr.Config.JWTSecret),
		fmt.Sprintf("TEST_LOG_LEVEL=%s", tr.Config.LogLevel),
	)
	
	if tr.Verbose {
		fmt.Printf("Running tests for package: %s\n", pkg)
	}
	
	err = cmd.Run()
	result.Duration = time.Since(startTime)
	
	// Parse test output
	tr.parseTestOutput(stdout.String(), result)
	
	// If verbose, show output
	if tr.Verbose {
		if stdout.Len() > 0 {
			fmt.Printf("=== Output for %s ===\n", pkg)
			fmt.Print(stdout.String())
		}
		if stderr.Len() > 0 {
			fmt.Printf("=== Errors for %s ===\n", pkg)
			fmt.Print(stderr.String())
		}
	}
	
	// Store raw output for later analysis
	result.Output = stdout.String() + stderr.String()
	
	return result, err
}

// runBenchmarksForPackage runs benchmarks for a specific package and returns results
func (tr *TestRunner) runBenchmarksForPackage(pkg string) (*TestResult, error) {
	startTime := time.Now()
	result := &TestResult{
		Package:    pkg,
		Benchmarks: make([]BenchmarkResult, 0),
	}
	
	testDir := filepath.Join("tests", pkg)
	
	// Check if benchmark directory exists
	if _, err := os.Stat(testDir); os.IsNotExist(err) {
		if tr.Verbose {
			fmt.Printf("No benchmarks found for package: %s\n", pkg)
		}
		return nil, nil
	}
	
	// Check if there are any test files in the directory
	matches, err := filepath.Glob(filepath.Join(testDir, "*_test.go"))
	if err != nil {
		return nil, fmt.Errorf("error checking for benchmark files: %v", err)
	}
	
	// If no test files found in main directory, check subdirectories (like v1/)
	if len(matches) == 0 {
		// Check for common subdirectories like v1/
		subDirs := []string{"v1"}
		for _, subDir := range subDirs {
			subDirPath := filepath.Join(testDir, subDir)
			if _, err := os.Stat(subDirPath); err == nil {
				subMatches, err := filepath.Glob(filepath.Join(subDirPath, "*_test.go"))
				if err != nil {
					return nil, fmt.Errorf("error checking for benchmark files in %s: %v", subDirPath, err)
				}
				if len(subMatches) > 0 {
					// Update testDir to point to the subdirectory that has tests
					testDir = subDirPath
					matches = subMatches
					break
				}
			}
		}
	}
	
	if len(matches) == 0 {
		if tr.Verbose {
			fmt.Printf("No benchmark files found in package: %s\n", pkg)
		}
		return nil, nil
	}
	
	// Build go test command for benchmarks
	args := []string{"test", "-bench=.", "-run=^$"} // -run=^$ ensures no regular tests run
	
	if tr.Verbose {
		args = append(args, "-v")
	}
	
	// Add the test directory as a package path
	// Convert absolute testDir back to relative path for go test
	relativeTestDir := strings.TrimPrefix(testDir, "tests/")
	if relativeTestDir == testDir {
		// testDir didn't start with "tests/", so it's already relative
		args = append(args, "./"+testDir)
	} else {
		args = append(args, "./tests/"+relativeTestDir)
	}
	
	// Execute benchmarks
	cmd := exec.Command("go", args...)
	// cmd.Dir should be the backend directory (current working directory when test.go is run)
	
	// Capture output
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	// Set test environment variables
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("TEST_DATABASE_URL=%s", tr.Config.DatabaseURL),
		fmt.Sprintf("TEST_JWT_SECRET=%s", tr.Config.JWTSecret),
		fmt.Sprintf("TEST_LOG_LEVEL=%s", tr.Config.LogLevel),
	)
	
	if tr.Verbose {
		fmt.Printf("Running benchmarks for package: %s\n", pkg)
	}
	
	err = cmd.Run()
	result.Duration = time.Since(startTime)
	
	// Parse benchmark output
	tr.parseBenchmarkOutput(stdout.String(), result)
	
	// If verbose, show output
	if tr.Verbose {
		if stdout.Len() > 0 {
			fmt.Printf("=== Benchmark Output for %s ===\n", pkg)
			fmt.Print(stdout.String())
		}
		if stderr.Len() > 0 {
			fmt.Printf("=== Benchmark Errors for %s ===\n", pkg)
			fmt.Print(stderr.String())
		}
	}
	
	// Store raw output for later analysis
	result.Output = stdout.String() + stderr.String()
	
	return result, err
}

// loadTestConfig loads basic test configuration
func loadTestConfig() (*TestConfig, error) {
	config := &TestConfig{
		DatabaseURL: getEnvOrDefault("TEST_DATABASE_URL", ":memory:"),
		JWTSecret:   getEnvOrDefault("TEST_JWT_SECRET", "test-secret-key"),
		LogLevel:    getEnvOrDefault("TEST_LOG_LEVEL", "error"),
	}
	
	return config, nil
}

// getEnvOrDefault returns environment variable value or default if not set
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// parseTestOutput parses JSON test output and populates TestResult
func (tr *TestRunner) parseTestOutput(output string, result *TestResult) {
	scanner := bufio.NewScanner(strings.NewReader(output))
	
	for scanner.Scan() {
		line := scanner.Text()
		
		// Try to parse as JSON
		var testEvent map[string]interface{}
		if err := json.Unmarshal([]byte(line), &testEvent); err != nil {
			// If not JSON, try to parse coverage from plain text
			if strings.Contains(line, "coverage:") {
				tr.parseCoverage(line, result)
			}
			continue
		}
		
		// Parse JSON test events
		action, ok := testEvent["Action"].(string)
		if !ok {
			continue
		}
		
		switch action {
		case "pass":
			result.Passed++
			result.Tests++
		case "fail":
			result.Failed++
			result.Tests++
			// Extract failed test details
			tr.extractFailedTestDetails(testEvent, result)
		case "skip":
			result.Skipped++
			result.Tests++
		}
	}
}

// extractFailedTestDetails extracts details about failed tests for better reporting
func (tr *TestRunner) extractFailedTestDetails(testEvent map[string]interface{}, result *TestResult) {
	testName, _ := testEvent["Test"].(string)
	output, _ := testEvent["Output"].(string)
	
	if testName != "" {
		failedTest := FailedTest{
			Package:  result.Package,
			TestName: testName,
			Output:   output,
			Error:    tr.extractErrorFromOutput(output),
		}
		
		// Add suggestions based on common error patterns
		failedTest.Suggestions = tr.generateSuggestions(failedTest.Error, failedTest.Output)
		
		// Add to global failed tests list (we'll need to pass summary reference)
		// For now, store in result and aggregate later
		if result.Package != "" {
			// This will be aggregated in the main Run function
		}
	}
}

// extractErrorFromOutput extracts the main error message from test output
func (tr *TestRunner) extractErrorFromOutput(output string) string {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "Error:") || strings.Contains(line, "FAIL:") {
			return line
		}
		if strings.HasPrefix(line, "panic:") {
			return line
		}
	}
	return "Test failed without clear error message"
}

// generateSuggestions provides actionable suggestions based on error patterns
func (tr *TestRunner) generateSuggestions(error, output string) []string {
	suggestions := make([]string, 0)
	
	errorLower := strings.ToLower(error)
	outputLower := strings.ToLower(output)
	
	// Database connection issues
	if strings.Contains(errorLower, "database") || strings.Contains(errorLower, "sql") {
		suggestions = append(suggestions, "Check database connection and ensure test database is properly initialized")
		suggestions = append(suggestions, "Verify TEST_DATABASE_URL environment variable is set correctly")
	}
	
	// Authentication issues
	if strings.Contains(errorLower, "unauthorized") || strings.Contains(errorLower, "token") {
		suggestions = append(suggestions, "Check JWT token generation and validation in test setup")
		suggestions = append(suggestions, "Verify TEST_JWT_SECRET environment variable matches application secret")
	}
	
	// HTTP request issues
	if strings.Contains(errorLower, "404") || strings.Contains(errorLower, "not found") {
		suggestions = append(suggestions, "Verify API endpoint URL and routing configuration")
		suggestions = append(suggestions, "Check if the handler is properly registered in the router")
	}
	
	// Validation errors
	if strings.Contains(errorLower, "validation") || strings.Contains(errorLower, "invalid") {
		suggestions = append(suggestions, "Check input data format and validation rules")
		suggestions = append(suggestions, "Ensure test data matches expected schema")
	}
	
	// Timeout issues
	if strings.Contains(errorLower, "timeout") || strings.Contains(errorLower, "deadline") {
		suggestions = append(suggestions, "Consider increasing test timeout or optimizing slow operations")
		suggestions = append(suggestions, "Check for potential deadlocks or infinite loops")
	}
	
	// Panic recovery
	if strings.Contains(outputLower, "panic") {
		suggestions = append(suggestions, "Review the stack trace to identify the source of the panic")
		suggestions = append(suggestions, "Add proper error handling to prevent panics in production code")
	}
	
	// Generic suggestions if no specific pattern matched
	if len(suggestions) == 0 {
		suggestions = append(suggestions, "Run the test with -verbose flag for more detailed output")
		suggestions = append(suggestions, "Check test setup and teardown procedures")
		suggestions = append(suggestions, "Verify all dependencies and mocks are properly configured")
	}
	
	return suggestions
}

// parseCoverage extracts coverage percentage from output
func (tr *TestRunner) parseCoverage(line string, result *TestResult) {
	// Look for coverage pattern like "coverage: 85.2% of statements"
	re := regexp.MustCompile(`coverage:\s+(\d+\.?\d*)%`)
	matches := re.FindStringSubmatch(line)
	if len(matches) > 1 {
		if coverage, err := strconv.ParseFloat(matches[1], 64); err == nil {
			result.Coverage = coverage
		}
	}
}

// parseBenchmarkOutput parses benchmark output and populates BenchmarkResult
func (tr *TestRunner) parseBenchmarkOutput(output string, result *TestResult) {
	scanner := bufio.NewScanner(strings.NewReader(output))
	
	for scanner.Scan() {
		line := scanner.Text()
		
		// Look for benchmark lines like "BenchmarkJWTGeneration-8   	  100000	     10234 ns/op	    1024 B/op	      10 allocs/op"
		if strings.HasPrefix(line, "Benchmark") {
			tr.parseBenchmarkLine(line, result)
		}
	}
}

// parseBenchmarkLine parses a single benchmark result line
func (tr *TestRunner) parseBenchmarkLine(line string, result *TestResult) {
	// Split by whitespace and filter empty strings
	parts := strings.Fields(line)
	if len(parts) < 3 {
		return
	}
	
	benchResult := BenchmarkResult{
		Name: parts[0],
	}
	
	// Parse iterations
	if iterations, err := strconv.Atoi(parts[1]); err == nil {
		benchResult.Iterations = iterations
	}
	
	// Parse ns/op - it might be a float with decimal places
	if len(parts) > 2 && strings.HasSuffix(parts[2], "ns/op") {
		nsStr := strings.TrimSuffix(parts[2], "ns/op")
		nsStr = strings.TrimSpace(nsStr)
		if nsFloat, err := strconv.ParseFloat(nsStr, 64); err == nil {
			benchResult.NsPerOp = int64(nsFloat)
		}
	}
	
	// Parse B/op and allocs/op if present
	for _, part := range parts {
		if strings.HasSuffix(part, "B/op") {
			bytesStr := strings.TrimSuffix(part, "B/op")
			bytesStr = strings.TrimSpace(bytesStr)
			if bytes, err := strconv.ParseInt(bytesStr, 10, 64); err == nil {
				benchResult.BytesPerOp = bytes
			}
		}
		if strings.HasSuffix(part, "allocs/op") {
			allocsStr := strings.TrimSuffix(part, "allocs/op")
			allocsStr = strings.TrimSpace(allocsStr)
			if allocs, err := strconv.ParseInt(allocsStr, 10, 64); err == nil {
				benchResult.AllocsPerOp = allocs
			}
		}
	}
	
	result.Benchmarks = append(result.Benchmarks, benchResult)
}

// printSummary prints a comprehensive test summary
func (tr *TestRunner) printSummary(summary *TestSummary) {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("TEST EXECUTION SUMMARY")
	fmt.Println(strings.Repeat("=", 60))
	
	// Overall results
	fmt.Printf("Total Duration: %v\n", summary.TotalDuration.Round(time.Millisecond))
	fmt.Printf("Total Tests: %d\n", summary.TotalTests)
	fmt.Printf("Passed: %d\n", summary.TotalPassed)
	fmt.Printf("Failed: %d\n", summary.TotalFailed)
	if summary.TotalSkipped > 0 {
		fmt.Printf("Skipped: %d\n", summary.TotalSkipped)
	}
	
	// Coverage information
	if tr.Coverage && summary.OverallCoverage > 0 {
		fmt.Printf("Overall Coverage: %.1f%%", summary.OverallCoverage)
		if tr.CoverageThreshold > 0 {
			if summary.OverallCoverage >= tr.CoverageThreshold {
				fmt.Printf(" ✓ (threshold: %.1f%%)", tr.CoverageThreshold)
			} else {
				fmt.Printf(" ✗ (threshold: %.1f%%)", tr.CoverageThreshold)
			}
		}
		fmt.Println()
	}
	
	// Package-level results
	if len(summary.Results) > 1 {
		fmt.Println("\nPer-Package Results:")
		fmt.Println(strings.Repeat("-", 60))
		for _, result := range summary.Results {
			if result.Tests > 0 || len(result.Benchmarks) > 0 {
				fmt.Printf("%-15s | Tests: %2d | Passed: %2d | Failed: %2d", 
					result.Package, result.Tests, result.Passed, result.Failed)
				if tr.Coverage && result.Coverage > 0 {
					fmt.Printf(" | Coverage: %5.1f%%", result.Coverage)
				}
				if len(result.Benchmarks) > 0 {
					fmt.Printf(" | Benchmarks: %d", len(result.Benchmarks))
				}
				fmt.Printf(" | Duration: %v\n", result.Duration.Round(time.Millisecond))
			}
		}
	}
	
	// Benchmark results summary
	totalBenchmarks := 0
	for _, result := range summary.Results {
		totalBenchmarks += len(result.Benchmarks)
	}
	
	if totalBenchmarks > 0 {
		fmt.Printf("\nBenchmark Results: %d benchmarks executed\n", totalBenchmarks)
		if tr.Verbose {
			fmt.Println(strings.Repeat("-", 60))
			for _, result := range summary.Results {
				if len(result.Benchmarks) > 0 {
					fmt.Printf("Package: %s\n", result.Package)
					for _, bench := range result.Benchmarks {
						fmt.Printf("  %-30s %8d iterations %8d ns/op", 
							bench.Name, bench.Iterations, bench.NsPerOp)
						if bench.BytesPerOp > 0 {
							fmt.Printf(" %8d B/op", bench.BytesPerOp)
						}
						if bench.AllocsPerOp > 0 {
							fmt.Printf(" %8d allocs/op", bench.AllocsPerOp)
						}
						fmt.Println()
					}
				}
			}
		}
	}
	
	// Final status
	fmt.Println(strings.Repeat("=", 60))
	if summary.Success {
		fmt.Println("✓ ALL TESTS PASSED")
	} else {
		fmt.Println("✗ TESTS FAILED")
		if summary.TotalFailed > 0 {
			fmt.Printf("  - %d test(s) failed\n", summary.TotalFailed)
		}
		if tr.CoverageThreshold > 0 && summary.OverallCoverage < tr.CoverageThreshold {
			fmt.Printf("  - Coverage below threshold: %.1f%% < %.1f%%\n", 
				summary.OverallCoverage, tr.CoverageThreshold)
		}
	}
	fmt.Println(strings.Repeat("=", 60))
}

// printUsage prints usage information
func printUsage() {
	fmt.Println("Simple Test Runner")
	fmt.Println("Usage: go run test.go [options]")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -type string")
	fmt.Println("        Test type to run: unit, integration, benchmark, all (default \"all\")")
	fmt.Println("  -package string")
	fmt.Println("        Specific package to test (optional)")
	fmt.Println("  -coverage")
	fmt.Println("        Generate coverage report")
	fmt.Println("  -coverage-threshold float")
	fmt.Println("        Minimum coverage threshold (0.0-100.0)")
	fmt.Println("  -verbose")
	fmt.Println("        Verbose output for debugging test failures")
	fmt.Println("  -help")
	fmt.Println("        Show this usage information")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  go run test.go                                    # Run all tests")
	fmt.Println("  go run test.go -type=unit                        # Run only unit tests")
	fmt.Println("  go run test.go -type=unit -package=api           # Run unit tests for API package")
	fmt.Println("  go run test.go -coverage -verbose                # Run all tests with coverage and verbose output")
	fmt.Println("  go run test.go -coverage -coverage-threshold=80  # Run tests with 80% coverage requirement")
	fmt.Println("  go run test.go -type=benchmark -verbose          # Run benchmarks with detailed output")
}