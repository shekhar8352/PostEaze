# Integration Test Runner Script for PostEaze Backend (PowerShell)
# This script sets up the environment and runs integration tests

param(
    [string]$Suite = "all",
    [switch]$Verbose,
    [switch]$Coverage,
    [string]$Database = "sqlite",
    [switch]$Parallel,
    [switch]$NoCleanup,
    [switch]$Help
)

# Colors for output
$Red = "Red"
$Green = "Green"
$Yellow = "Yellow"
$Blue = "Blue"

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

# Function to show usage
function Show-Usage {
    Write-Host "Usage: .\run_integration_tests.ps1 [OPTIONS]"
    Write-Host ""
    Write-Host "Options:"
    Write-Host "  -Suite SUITE          Run specific test suite (auth|log|database|e2e|all)"
    Write-Host "  -Verbose              Enable verbose output"
    Write-Host "  -Coverage             Generate coverage report"
    Write-Host "  -Database TYPE        Database type to use (sqlite|postgres)"
    Write-Host "  -Parallel             Run tests in parallel"
    Write-Host "  -NoCleanup            Skip cleanup after tests"
    Write-Host "  -Help                 Show this help message"
    Write-Host ""
    Write-Host "Examples:"
    Write-Host "  .\run_integration_tests.ps1 -Suite auth -Verbose"
    Write-Host "  .\run_integration_tests.ps1 -Coverage -Database postgres"
    Write-Host "  .\run_integration_tests.ps1 -Suite all -Parallel"
}

# Show help if requested
if ($Help) {
    Show-Usage
    exit 0
}

# Validate test suite
$ValidSuites = @("auth", "log", "database", "e2e", "all")
if ($Suite -notin $ValidSuites) {
    Write-Error "Invalid test suite: $Suite"
    Write-Error "Valid options: $($ValidSuites -join ', ')"
    exit 1
}

Write-Status "Starting PostEaze Backend Integration Tests"
Write-Status "Test Suite: $Suite"
Write-Status "Database Type: $Database"
Write-Status "Verbose: $Verbose"
Write-Status "Coverage: $Coverage"
Write-Status "Parallel: $Parallel"

# Change to backend directory
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$BackendDir = Join-Path $ScriptDir "../.."
Set-Location $BackendDir

# Check if we're in the right directory
if (-not (Test-Path "go.mod")) {
    Write-Error "go.mod not found. Please run this script from the backend directory or ensure the path is correct."
    exit 1
}

# Setup environment variables
function Setup-Environment {
    Write-Status "Setting up test environment..."
    
    $env:MODE = "dev"
    $env:BASE_CONFIG_PATH = "./tests/config"
    $env:JWT_ACCESS_SECRET = "integration-test-access-secret-key-$(Get-Date -Format 'yyyyMMddHHmmss')"
    $env:JWT_REFRESH_SECRET = "integration-test-refresh-secret-key-$(Get-Date -Format 'yyyyMMddHHmmss')"
    
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
    
    # Create temporary log directory
    $TempLogDir = New-TemporaryFile | ForEach-Object { Remove-Item $_; New-Item -ItemType Directory -Path $_.FullName }
    $env:LOG_DIR = $TempLogDir.FullName
    Write-Status "Created temporary log directory: $($env:LOG_DIR)"
    
    return $TempLogDir
}

# Setup test database (for PostgreSQL)
function Setup-Database {
    if ($Database -eq "postgres") {
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

# Build test command
function Build-TestCommand {
    $cmd = @("go", "test")
    
    # Add test path based on suite
    switch ($Suite) {
        "auth" { $cmd += @("./tests/integration/", "-run", "TestAuthIntegrationSuite") }
        "log" { $cmd += @("./tests/integration/", "-run", "TestLogIntegrationSuite") }
        "database" { $cmd += @("./tests/integration/", "-run", "TestDatabaseIntegrationSuite") }
        "e2e" { $cmd += @("./tests/integration/", "-run", "TestEndToEndIntegrationSuite") }
        "all" { $cmd += @("./tests/integration/...") }
    }
    
    # Add flags
    if ($Verbose) {
        $cmd += "-v"
    }
    
    if ($Coverage) {
        $cmd += "-coverprofile=integration_coverage.out"
    }
    
    if ($Parallel) {
        $cmd += "-parallel", "4"
    }
    
    # Add timeout
    $cmd += "-timeout", "10m"
    
    return $cmd
}

# Run tests
function Run-Tests {
    $testCmd = Build-TestCommand
    
    Write-Status "Running integration tests..."
    Write-Status "Command: $($testCmd -join ' ')"
    
    # Run the tests
    $process = Start-Process -FilePath $testCmd[0] -ArgumentList $testCmd[1..($testCmd.Length-1)] -Wait -PassThru -NoNewWindow
    
    if ($process.ExitCode -eq 0) {
        Write-Success "All integration tests passed!"
        return $true
    }
    else {
        Write-Error "Some integration tests failed!"
        return $false
    }
}

# Generate coverage report
function Generate-CoverageReport {
    if ($Coverage -and (Test-Path "integration_coverage.out")) {
        Write-Status "Generating coverage report..."
        
        # Generate HTML coverage report
        & go tool cover -html=integration_coverage.out -o integration_coverage.html
        
        # Show coverage summary
        $coverageOutput = & go tool cover -func=integration_coverage.out
        $coverageSummary = $coverageOutput | Select-Object -Last 1
        Write-Host $coverageSummary
        
        Write-Success "Coverage report generated: integration_coverage.html"
    }
}

# Cleanup function
function Cleanup {
    param([System.IO.DirectoryInfo]$TempLogDir)
    
    if (-not $NoCleanup) {
        Write-Status "Cleaning up..."
        
        # Remove temporary log directory
        if ($TempLogDir -and (Test-Path $TempLogDir.FullName)) {
            Remove-Item -Recurse -Force $TempLogDir.FullName
            Write-Status "Removed temporary log directory"
        }
        
        # Remove coverage files if not requested
        if (-not $Coverage) {
            if (Test-Path "integration_coverage.out") { Remove-Item "integration_coverage.out" }
            if (Test-Path "integration_coverage.html") { Remove-Item "integration_coverage.html" }
        }
        
        Write-Success "Cleanup completed"
    }
}

# Main execution
try {
    $exitCode = 0
    
    # Setup
    $tempLogDir = Setup-Environment
    Setup-Database
    
    # Run tests
    if (-not (Run-Tests)) {
        $exitCode = 1
    }
    
    # Generate coverage report if requested
    Generate-CoverageReport
    
    # Print summary
    if ($exitCode -eq 0) {
        Write-Success "Integration test run completed successfully!"
    }
    else {
        Write-Error "Integration test run completed with failures!"
    }
}
catch {
    Write-Error "An error occurred: $($_.Exception.Message)"
    $exitCode = 1
}
finally {
    # Cleanup
    if ($tempLogDir) {
        Cleanup -TempLogDir $tempLogDir
    }
}

exit $exitCode