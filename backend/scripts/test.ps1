# PostEaze Backend Test Execution Script (PowerShell)
# Comprehensive test runner with coverage reporting and test discovery

param(
    [string]$Type = "all",
    [string]$Suite = "",
    [switch]$Verbose,
    [switch]$Coverage,
    [string]$Database = "sqlite",
    [switch]$Parallel,
    [string]$Format = "standard",
    [int]$Threshold = 80,
    [switch]$Benchmark,
    [switch]$Race,
    [switch]$Short,
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
    Write-Host "  -Parallel             Run tests in parallel"
    Write-Host "  -Format FORMAT        Output format (standard|json|junit)"
    Write-Host "  -Threshold NUM        Coverage threshold percentage (default: 80)"
    Write-Host "  -Benchmark            Run benchmark tests"
    Write-Host "  -Race                 Enable race detection"
    Write-Host "  -Short                Run tests in short mode"
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
    Write-Host "  .\test.ps1 -Coverage -Threshold 90 -Parallel"
    Write-Host "  .\test.ps1 -Benchmark -Race"
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
Write-Status "Parallel: $Parallel"
Write-Status "Output Format: $Format"
Write-Status "Coverage Threshold: $Threshold%"

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
    $env:BASE_CONFIG_PATH = "./tests/config"
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
    
    Write-Status "Created temporary directories"
}

# Discover test packages
function Get-TestPackages {
    param(
        [string]$TestType,
        [string]$TestSuite
    )
    
    $packages = @()
    
    switch ($TestType) {
        "unit" {
            if ($TestSuite) {
                switch ($TestSuite) {
                    "api" { $packages += "./tests/unit/api/..." }
                    "business" { $packages += "./tests/unit/business/..." }
                    "models" { $packages += "./tests/unit/models/..." }
                    "utils" { $packages += "./tests/unit/utils/..." }
                    default {
                        Write-Error "Invalid unit test suite: $TestSuite"
                        Write-Error "Valid options: api, business, models, utils"
                        exit 1
                    }
                }
            }
            else {
                $packages += "./tests/unit/..."
            }
        }
        "integration" {
            if ($TestSuite) {
                switch ($TestSuite) {
                    "auth" { $packages += @("./tests/integration/", "-run", "TestAuthIntegrationSuite") }
                    "log" { $packages += @("./tests/integration/", "-run", "TestLogIntegrationSuite") }
                    "database" { $packages += @("./tests/integration/", "-run", "TestDatabaseIntegrationSuite") }
                    "e2e" { $packages += @("./tests/integration/", "-run", "TestEndToEndIntegrationSuite") }
                    default {
                        Write-Error "Invalid integration test suite: $TestSuite"
                        Write-Error "Valid options: auth, log, database, e2e"
                        exit 1
                    }
                }
            }
            else {
                $packages += "./tests/integration/..."
            }
        }
        "all" {
            if ($TestSuite) {
                Write-Warning "Test suite specified with 'all' type. Running all tests."
            }
            $packages += @("./tests/unit/...", "./tests/integration/...")
        }
    }
    
    return $packages
}

# Build test command
function Build-TestCommand {
    param([array]$Packages)
    
    $cmd = @("go", "test")
    $cmd += $Packages
    
    # Add flags
    if ($Verbose) {
        $cmd += "-v"
    }
    
    if ($Coverage) {
        $coverageFile = Join-Path $script:CoverageDir.FullName "coverage.out"
        $cmd += "-coverprofile=$coverageFile"
        $cmd += "-covermode=atomic"
    }
    
    if ($Parallel) {
        $cmd += @("-parallel", "4")
    }
    
    if ($Race) {
        $cmd += "-race"
    }
    
    if ($Short) {
        $cmd += "-short"
    }
    
    if ($Benchmark) {
        $cmd += @("-bench=.", "-benchmem")
    }
    
    # Add output format
    switch ($Format) {
        "json" { $cmd += "-json" }
        "junit" { $cmd += "-v" }  # Go doesn't natively support JUnit
    }
    
    # Add timeout
    $cmd += @("-timeout", "15m")
    
    return $cmd
}

# Run tests
function Invoke-Tests {
    $packages = Get-TestPackages -TestType $Type -TestSuite $Suite
    $testCmd = Build-TestCommand -Packages $packages
    
    Write-Status "Discovered test packages: $($packages -join ' ')"
    Write-Status "Running tests..."
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
            Write-Success "All tests passed!"
            return $true
        }
        else {
            Write-Error "Some tests failed!"
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

# Generate coverage report
function New-CoverageReport {
    $coverageFile = Join-Path $script:CoverageDir.FullName "coverage.out"
    
    if ($Coverage -and (Test-Path $coverageFile)) {
        Write-Status "Generating coverage report..."
        
        # Generate HTML coverage report
        $htmlFile = Join-Path $script:CoverageDir.FullName "coverage.html"
        & go tool cover -html=$coverageFile -o $htmlFile
        
        # Generate coverage summary
        $funcFile = Join-Path $script:CoverageDir.FullName "coverage_func.txt"
        & go tool cover -func=$coverageFile | Out-File -FilePath $funcFile
        
        # Extract total coverage percentage
        $coverageOutput = & go tool cover -func=$coverageFile
        $totalLine = $coverageOutput | Select-Object -Last 1
        $totalCoverage = [regex]::Match($totalLine, '(\d+\.\d+)%').Groups[1].Value
        
        Write-Status "Coverage Summary:"
        Write-Host "----------------------------------------"
        $coverageOutput | ForEach-Object { Write-Host $_ }
        Write-Host "----------------------------------------"
        
        # Check coverage threshold
        if ([double]$totalCoverage -ge $Threshold) {
            Write-Success "Coverage ($totalCoverage%) meets threshold ($Threshold%)"
        }
        else {
            Write-Warning "Coverage ($totalCoverage%) below threshold ($Threshold%)"
            if ($env:CI -eq "true") {
                Write-Error "Coverage threshold not met in CI environment"
                return $false
            }
        }
        
        Write-Success "Coverage reports generated:"
        Write-Success "  HTML: $htmlFile"
        Write-Success "  Text: $funcFile"
        Write-Success "  Raw: $coverageFile"
        
        # Copy coverage files to standard locations for CI
        try {
            Copy-Item $coverageFile "./coverage.out" -ErrorAction SilentlyContinue
            Copy-Item $htmlFile "./coverage.html" -ErrorAction SilentlyContinue
        }
        catch {
            # Ignore copy errors
        }
        
        return $true
    }
    
    return $true
}

# Run benchmarks
function Invoke-Benchmarks {
    if ($Benchmark) {
        Write-Status "Running benchmarks..."
        
        $benchmarkFile = Join-Path $script:CoverageDir.FullName "benchmark_results.txt"
        & go test -bench=. -benchmem ./... | Out-File -FilePath $benchmarkFile
        
        Write-Status "Benchmark Results:"
        Write-Host "----------------------------------------"
        Get-Content $benchmarkFile | ForEach-Object { Write-Host $_ }
        Write-Host "----------------------------------------"
        
        Write-Success "Benchmark results saved: $benchmarkFile"
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
}

# Main execution
function Invoke-Main {
    $exitCode = 0
    
    try {
        # Setup
        Setup-Environment
        Setup-Database
        
        # Run tests
        if (-not (Invoke-Tests)) {
            $exitCode = 1
        }
        
        # Generate coverage report
        if (-not (New-CoverageReport)) {
            $exitCode = 1
        }
        
        # Run benchmarks if requested
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