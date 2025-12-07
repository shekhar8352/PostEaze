# Backend Scripts

This directory contains utility scripts for testing, development, and maintenance of the PostEaze backend.

## Test Runner (`test.sh` / `test.ps1`)

A unified test runner wrapper around `go test` that provides enhanced output, coverage reporting, and test discovery.

### Usage
**Bash (Linux/macOS):**
```bash
./scripts/test.sh [OPTIONS]
```

**PowerShell (Windows):**
```powershell
.\scripts\test.ps1 [OPTIONS]
```

### Options
- `-t, --type TYPE`: Test type (`unit`, `integration`, `benchmark`, `all`). Default: `all`.
- `-p, --package PKG`: Run tests for a specific package (e.g., `api`, `business`).
- `-v, --verbose`: Enable verbose output.
- `-c, --coverage`: Generate coverage report (`coverage.out` and `coverage.html`).
- `--threshold NUM`: Fail if coverage is below this percentage.

## Decrypt Token Utility (`decrypt_token`)

A utility to decrypt channel access tokens stored in the database.

### Usage
From the `backend` directory:
```bash
go run scripts/decrypt_token/main.go
```

### Configuration
Requires the following environment variables (usually loaded from `.env`):
- `POSTGRES_USER`
- `POSTGRES_PASSWORD`
- `POSTGRES_URL`
- `POSTGRES_DB`
- `ENCRYPTION_KEY`

## Test Discovery Utility (`test_discovery`)

Analyzes the codebase to discover and list all available tests, benchmarks, and suites.

### Usage
From the `backend` directory:
```bash
go run scripts/test_discovery/main.go <directory>
# Example:
go run scripts/test_discovery/main.go ./tests
```

Outputs a JSON structure containing details about all discovered tests.
