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

## Studio Demo Seeder (`seed_studio_demo`)

Populates a user's Studio board with a hand-crafted set of demo Pieces spread
across the default phases (Ideas → Scripting → Shooting → Editing → Review →
Scheduled → Published). The Studio itself is created on demand via the normal
`EnsureStudioForUser` flow, so default phases are seeded automatically on first
run.

### Usage

From the `backend` directory:

```bash
# Seed the Studio for a given user (UUID). Creates the studio + default phases
# if they do not exist, then adds ~12 demo pieces.
go run scripts/seed_studio_demo/main.go --user-id <uuid>

# Clean slate: archive all existing active pieces in the user's studio first.
go run scripts/seed_studio_demo/main.go --user-id <uuid> --reset

# Cap the number of pieces created (useful for smoke tests).
go run scripts/seed_studio_demo/main.go --user-id <uuid> --count 4
```

### Flags
- `--user-id` (required): UUID of an existing user. Their team's default Studio
  is the seed target.
- `--reset`: Archive all existing active pieces in the Studio before seeding,
  giving you a predictable starting board every run.
- `--count`: Optional cap on how many demo pieces to create.

### Configuration

Uses the same `DatabaseConfig` as the main server and worker (loaded via
`utils/configs` + `utils/env`). In dev this typically reads from your `.env`
file (`POSTGRES_URL`, `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`,
`ENCRYPTION_KEY`).

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
