# PostEaze Backend Test Execution Script (PowerShell)
# Comprehensive test runner with coverage reporting and test discovery

param(
    [string]$Type = "all",
    [string]$Suite = "",
    [switch]$Verbose,
    [switch]$Coverage,
    [string]$Database = "sqlite",
    [string]$Format = "standard",
    [int]$Threshold = 0,
    [switch]$Benchmark,
    [switch]$NoCleanup,
    [switch]$Help
)

# Colors for output
$Red = "Red"
$Green = "Green"
$Yellow = "Yellow"
$Blue = "Blue"
$Cyan = "Cyan"

# Function to print colored output
function Write-Status {
    param([string]$Message)
    Write-Host "[INFO] $Message" -ForegroundColor $Blue
}

function Write-Success {
    param([string]$Message)
    Write-Host "[SUCCESS] $Message" -ForegroundColor $Green
}

function Write-Warning {
    param([string]$Message)
    Write-Host "[WARNING] $Message" -ForegroundColor $Yellow
}

function Write-Error {
    param([string]$Message)
    Write-Host "[ERROR] $Message" -ForegroundColor $Red
}

function Write-Header {
    param([string]$Message)
    Write-Host "========================================" -ForegroundColor $Cyan
    Write-Host " $Message" -ForegroundColor $Cyan
    Write-Host "========================================" -ForegroundColor $Cyan
}

# Function to show usage
function Show-Usage {
    Write-Host "Usage: .\test.ps1 [OPTIONS]"
    Write-Host ""
    Write-Host "Test Types:"
    Write-Host "  unit                  Run unit tests only"
    Write-Host "  integration          Run integration tests only"
    Write-Host "  all                  Run all tests (default)"
    Write-Host ""
    Write-Host "Options:"
    Write-Host "  -Type TYPE            Test type (unit|integration|all)"
    Write-Host "  -Suite SUITE          Run specific test suite"
    Write-Host "  -Verbose              Enable verbose output"
    Write-Host "  -Coverage             Generate coverage report"
    Write-Host "  -Database TYPE        Database type for integration tests (sqlite|postgres)"
    Write-Host "  -Format FORMAT        Output format (standard|json|junit)"
    Write-Host "  -Threshold NUM        Coverage threshold percentage (default: 0)"
    Write-Host "  -Benchmark            Run benchmark tests"
    Write-Host "  -NoCleanup            Skip cleanup after tests"
    Write-Host "  -Help                 Show this help message"
    Write-Host ""
    Write-Host "Unit Test Suites:"
    Write-Host "  api                   API handler tests"
    Write-Host "  business              Business logic tests"
    Write-Host "  models                Model validation tests"
    Write-Host "  utils                 Utility function tests"
    Write-Host ""
    Write-Host "Integration Test Suites:"
    Write-Host "  auth                  Authentication integration tests"
    Write-Host "  log                   Log management integration tests"
    Write-Host "  database              Database integration tests"
    Write-Host "  e2e                   End-to-end integration tests"
    Write-Host ""
    Write-Host "Examples:"
    Write-Host "  .\test.ps1 -Type unit -Suite api -Verbose"
    Write-Host "  .\test.ps1 -Type integration -Coverage -Database postgres"
    Write-Host "  .\test.ps1 -Coverage -Threshold 90"
    Write-Host "  .\test.ps1 -Benchmark -Verbose"
}

# Show help if requested
if ($Help) {
    Show-Usage
    exit 0
}

# Validate test type
$ValidTypes = @("unit", "integration", "all")
if ($Type -notin $ValidTypes) {
    Write-Error "Invalid test type: $Type"
    Write-Error "Valid options: $($ValidTypes -join ', ')"
    exit 1
}

# Validate output format
$ValidFormats = @("standard", "json", "junit")
if ($Format -notin $ValidFormats) {
    Write-Error "Invalid output format: $Format"
    Write-Error "Valid options: $($ValidFormats -join ', ')"
    exit 1
}

Write-Header "PostEaze Backend Test Suite"
Write-Status "Test Type: $Type"
Write-Status "Test Suite: $(if ($Suite) { $Suite } else { 'all' })"
Write-Status "Database Type: $Database"
Write-Status "Verbose: $Verbose"
Write-Status "Coverage: $Coverage"
Write-Status "Output Format: $Format"
if ($Threshold -gt 0) {
    Write-Status "Coverage Threshold: $Threshold%"
}

# Change to backend directory
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$BackendDir = Split-Path -Parent $ScriptDir
Set-Location $BackendDir

# Check if we're in the right directory
if (-not (Test-Path "go.mod")) {
    Write-Error "go.mod not found. Please run this script from the backend directory or ensure the path is correct."
    exit 1
}

# Global variables for cleanup
$script:TempLogDir = $null
$script:CoverageDir = $null

# Setup environment variables
function Setup-Environment {
    Write-Status "Setting up test environment..."
    
    $env:MODE = "test"
    $env:JWT_ACCESS_SECRET = "test-access-secret-key-$(Get-Date -Format 'yyyyMMddHHmmss')"
    $env:JWT_REFRESH_SECRET = "test-refresh-secret-key-$(Get-Date -Format 'yyyyMMddHHmmss')"
    
    if ($Database -eq "postgres") {
        # Check if PostgreSQL is available
        try {
            $null = Get-Command psql -ErrorAction Stop
            $env:DATABASE_DRIVER = "postgres"
            $env:DATABASE_URL = "postgres://postgres:password@localhost:5432/postease_test?sslmode=disable"
            Write-Status "Using PostgreSQL database"
        }
        catch {
            Write-Warning "PostgreSQL not found, falling back to SQLite"
            $Database = "sqlite"
        }
    }
    
    if ($Database -eq "sqlite") {
        $env:DATABASE_DRIVER = "sqlite3"
        $env:DATABASE_URL = ":memory:"
        Write-Status "Using SQLite in-memory database"
    }
    
    # Create temporary directories
    $script:TempLogDir = New-TemporaryFile | ForEach-Object { 
        Remove-Item $_ 
        New-Item -ItemType Directory -Path $_.FullName 
    }
    $script:CoverageDir = New-TemporaryFile | ForEach-Object { 
        Remove-Item $_ 
        New-Item -ItemType Directory -Path $_.FullName 
    }
    
    $env:TEST_LOG_DIR = $script:TempLogDir.FullName
    $env:COVERAGE_DIR = $script:CoverageDir.FullName
    
    # Set environment variables for the Go test runner
    $env:TEST_DATABASE_URL = $env:DATABASE_URL
    $env:TEST_JWT_SECRET = $env:JWT_ACCESS_SECRET
    $env:TEST_LOG_LEVEL = "error"
    
    Write-Status "Created temporary directories and set test environment variables"
}

# Check if Go test runner exists
function Test-GoTestRunner {
    if (-not (Test-Path "tests/test.go")) {
        Write-Error "tests/test.go not found. Please ensure the Go test runner is available."
        exit 1
    }
}

# Build test command for the Go test runner
function Build-TestCommand {
    $cmd = @("go", "run", "tests/test.go")
    
    # Map test type
    $testType = switch ($Type) {
        "unit" { "unit" }
        "integration" { "integration" }
        "all" { "all" }
        default { "all" }
    }
    $cmd += @("-type", $testType)
    
    # Add package if specified (map suite to package)
    if ($Suite) {
        $packageName = switch ($Suite) {
            "api" { "api" }
            "business" { "business" }
            "models" { "models" }
            "utils" { "utils" }
            "auth" { "integration" }
            "log" { "integration" }
            "database" { "integration" }
            "e2e" { "integration" }
            default { $Suite }
        }
        $cmd += @("-package", $packageName)
    }
    
    # Add flags
    if ($Verbose) {
        $cmd += "-verbose"
    }
    
    if ($Coverage) {
        $cmd += "-coverage"
    }
    
    if ($Threshold -gt 0) {
        $cmd += @("-coverage-threshold", $Threshold.ToString())
    }
    
    # Handle benchmark mode
    if ($Benchmark) {
        $cmd[3] = "benchmark"  # Change type to benchmark
    }
    
    return $cmd
}

# Run tests using the Go test runner
function Invoke-Tests {
    $testCmd = Build-TestCommand
    
    Write-Status "Running tests with Go test runner..."
    Write-Status "Command: $($testCmd -join ' ')"
    
    # Create output file for JUnit format
    $outputFile = $null
    if ($Format -eq "junit") {
        $outputFile = Join-Path $script:CoverageDir.FullName "test_results.xml"
    }
    
    # Run the tests
    $exitCode = 0
    try {
        if ($Format -eq "junit") {
            # Capture output for JUnit conversion
            $outputPath = Join-Path $script:CoverageDir.FullName "test_output.txt"
            $process = Start-Process -FilePath $testCmd[0] -ArgumentList $testCmd[1..($testCmd.Length-1)] -Wait -PassThru -NoNewWindow -RedirectStandardOutput $outputPath -RedirectStandardError $outputPath
            $exitCode = $process.ExitCode
            
            # Convert to JUnit format (simplified)
            Convert-ToJUnit -InputFile $outputPath -OutputFile $outputFile
        }
        else {
            $process = Start-Process -FilePath $testCmd[0] -ArgumentList $testCmd[1..($testCmd.Length-1)] -Wait -PassThru -NoNewWindow
            $exitCode = $process.ExitCode
        }
        
        if ($exitCode -eq 0) {
            Write-Success "All tests completed successfully!"
            return $true
        }
        else {
            Write-Error "Test execution failed!"
            return $false
        }
    }
    catch {
        Write-Error "Error running tests: $($_.Exception.Message)"
        return $false
    }
}

# Convert test output to JUnit format (simplified)
function Convert-ToJUnit {
    param(
        [string]$InputFile,
        [string]$OutputFile
    )
    
    Write-Status "Converting test output to JUnit format..."
    
    $junitContent = @'
<?xml version="1.0" encoding="UTF-8"?>
<testsuites>
    <testsuite name="PostEaze Backend Tests" tests="0" failures="0" errors="0" time="0">
        <!-- Test results would be parsed and inserted here -->
        <!-- This is a simplified placeholder implementation -->
    </testsuite>
</testsuites>
'@
    
    Set-Content -Path $OutputFile -Value $junitContent
    Write-Status "JUnit report generated: $OutputFile"
}

# Coverage reporting is now handled by the Go test runner
function New-CoverageReport {
    # The Go test runner handles coverage internally
    # Coverage reports are generated automatically when -coverage flag is used
    if ($Coverage) {
        Write-Status "Coverage reporting handled by Go test runner"
        
        # Check for generated coverage files in standard locations
        if (Test-Path "./coverage.out") {
            Write-Success "Coverage report available: ./coverage.out"
        }
        if (Test-Path "./coverage.html") {
            Write-Success "HTML coverage report available: ./coverage.html"
        }
    }
    
    return $true
}

# Benchmarks are now handled by the Go test runner
function Invoke-Benchmarks {
    # Benchmarks are handled by the Go test runner when -type=benchmark is used
    # No separate benchmark execution needed
    if ($Benchmark) {
        Write-Status "Benchmarks handled by Go test runner"
    }
}

# Setup test database (for PostgreSQL)
function Setup-Database {
    if ($Database -eq "postgres" -and ($Type -eq "integration" -or $Type -eq "all")) {
        Write-Status "Setting up PostgreSQL test database..."
        
        try {
            # Check if database exists and create if needed
            $dbExists = & psql -lqt | Select-String "postease_test"
            if (-not $dbExists) {
                Write-Status "Creating test database..."
                & createdb postease_test
            }
            Write-Success "PostgreSQL test database ready"
        }
        catch {
            Write-Warning "Failed to setup PostgreSQL database: $($_.Exception.Message)"
        }
    }
}

# Cleanup function
function Invoke-Cleanup {
    if (-not $NoCleanup) {
        Write-Status "Cleaning up..."
        
        # Remove temporary directories
        if ($script:TempLogDir -and (Test-Path $script:TempLogDir.FullName)) {
            Remove-Item -Recurse -Force $script:TempLogDir.FullName -ErrorAction SilentlyContinue
        }
        
        if ($script:CoverageDir -and (Test-Path $script:CoverageDir.FullName)) {
            # Keep coverage files if coverage was requested
            if (-not $Coverage) {
                Remove-Item -Recurse -Force $script:CoverageDir.FullName -ErrorAction SilentlyContinue
            }
            else {
                Write-Status "Coverage files preserved in: $($script:CoverageDir.FullName)"
            }
        }
        
        Write-Success "Cleanup completed"
    }
}

# Check for required tools
function Test-Dependencies {
    $missingDeps = @()
    
    try {
        $null = Get-Command go -ErrorAction Stop
    }
    catch {
        $missingDeps += "go"
    }
    
    if ($Database -eq "postgres") {
        try {
            $null = Get-Command psql -ErrorAction Stop
        }
        catch {
            Write-Warning "PostgreSQL tools not found, will fallback to SQLite"
        }
    }
    
    if ($missingDeps.Count -gt 0) {
        Write-Error "Missing required dependencies: $($missingDeps -join ', ')"
        exit 1
    }
    
    # Check if Go test runner exists
    Test-GoTestRunner
}

# Main execution
function Invoke-Main {
    $exitCode = 0
    
    try {
        # Setup
        Setup-Environment
        Setup-Database
        
        # Run tests using the Go test runner
        if (-not (Invoke-Tests)) {
            $exitCode = 1
        }
        
        # Check coverage reports (handled by Go test runner)
        if (-not (New-CoverageReport)) {
            $exitCode = 1
        }
        
        # Benchmarks are handled by the Go test runner
        Invoke-Benchmarks
        
        # Print summary
        Write-Header "Test Execution Summary"
        if ($exitCode -eq 0) {
            Write-Success "All test operations completed successfully!"
        }
        else {
            Write-Error "Test execution completed with failures!"
        }
        
        return $exitCode
    }
    catch {
        Write-Error "An error occurred: $($_.Exception.Message)"
        return 1
    }
    finally {
        # Cleanup
        Invoke-Cleanup
    }
}

# Check dependencies before running
Test-Dependencies

# Run main function
$result = Invoke-Main
exit $result