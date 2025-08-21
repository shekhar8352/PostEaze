# Database Initialization

The `init-db` folder serves as the database initialization directory for the PostEaze application's PostgreSQL database. This folder is mounted as a Docker volume to PostgreSQL's `/docker-entrypoint-initdb.d` directory, allowing for automatic database setup during container initialization.

## Purpose

This directory provides a standardized location for database initialization scripts that are automatically executed when the PostgreSQL container starts for the first time. It ensures consistent database setup across all environments (development, staging, production) and provides a foundation for the PostEaze application's data layer.

## Docker Integration

### Volume Mounting
The `init-db` folder is mounted in the PostgreSQL container via Docker Compose:

```yaml
postgres:
  image: postgres:15-alpine
  volumes:
    - ./init-db:/docker-entrypoint-initdb.d
```

### Execution Order
PostgreSQL automatically executes files in `/docker-entrypoint-initdb.d` in alphabetical order during the first container startup. Supported file types include:
- `.sql` files - Executed directly by PostgreSQL
- `.sh` files - Executed as shell scripts
- `.sql.gz` files - Decompressed and executed

## Current State

The `init-db` folder is currently empty, which means the database relies on:
1. **Environment Variables** - Database name, user, and password from `.env` file
2. **Application Migrations** - Schema creation handled by migration files in `backend/migrations/`
3. **Runtime Initialization** - Database schema applied when the backend application starts

## Database Configuration

### Environment Variables
Database initialization uses the following environment variables from `.env`:

```bash
POSTGRES_DB="posteaze_db"        # Database name
POSTGRES_USER="posteaze_user"    # Database user
POSTGRES_PASSWORD="posteaze_pass" # Database password
```

### Connection Details
- **Host**: `postgres` (Docker service name) or `localhost` (development)
- **Port**: `5432` (standard PostgreSQL port)
- **Database**: `posteaze_db`
- **User**: `posteaze_user`

## Schema Management Strategy

PostEaze uses a **migration-based approach** for database schema management:

### Current Approach
1. **Empty init-db**: No initialization scripts in this folder
2. **Migration Files**: Schema defined in `backend/migrations/` folder
3. **Manual Execution**: Migrations applied manually or via application startup
4. **Version Control**: All schema changes tracked in Git via migration files

### Migration Files Location
```
backend/migrations/
├── 001_create_initial_tables.up.sql    # Schema creation
├── 001_create_initial_tables.down.sql  # Schema rollback
└── README.md                           # Migration documentation
```

## Potential Use Cases

While currently empty, the `init-db` folder could be used for:

### Database Setup Scripts
```sql
-- Example: init-db/01-setup-extensions.sql
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
```

### Initial Data Population
```sql
-- Example: init-db/02-seed-data.sql
INSERT INTO users (name, email, password, user_type) VALUES
('Admin User', 'admin@posteaze.com', 'hashed_password', 'admin');
```

### Database Configuration
```sql
-- Example: init-db/03-configure-database.sql
ALTER DATABASE posteaze_db SET timezone TO 'UTC';
```

### User and Permission Setup
```sql
-- Example: init-db/04-setup-permissions.sql
GRANT ALL PRIVILEGES ON DATABASE posteaze_db TO posteaze_user;
```

## Development Workflow

### First-Time Setup
1. **Start Services**: Run `docker-compose up` to start PostgreSQL
2. **Database Creation**: PostgreSQL creates `posteaze_db` automatically
3. **Schema Application**: Apply migrations from `backend/migrations/`
4. **Application Start**: Backend connects and is ready for use

### Adding Initialization Scripts
If you need to add initialization scripts:

1. **Create SQL File**: Add `.sql` file to `init-db/` folder
2. **Use Numeric Prefix**: Name files with numeric prefixes for execution order
   - `01-extensions.sql`
   - `02-initial-data.sql`
   - `03-permissions.sql`
3. **Rebuild Container**: Remove existing PostgreSQL volume and restart
   ```bash
   docker-compose down -v  # Remove volumes
   docker-compose up       # Recreate with init scripts
   ```

## Database Initialization Process

### Container Startup Sequence
1. **PostgreSQL Container Start**: `postgres:15-alpine` image initialization
2. **Environment Variables**: Database name, user, password applied
3. **Init Scripts Execution**: Files in `/docker-entrypoint-initdb.d` executed
4. **Database Ready**: PostgreSQL accepts connections
5. **Backend Connection**: Go application connects and applies migrations

### Initialization Conditions
PostgreSQL only runs initialization scripts when:
- The database data directory is empty (first startup)
- No existing PostgreSQL data volume is mounted
- The container is starting fresh (not restarting)

## Troubleshooting

### Common Issues

#### Scripts Not Executing
- **Cause**: Existing PostgreSQL data volume prevents re-initialization
- **Solution**: Remove volumes and restart containers
  ```bash
  docker-compose down -v
  docker-compose up
  ```

#### Permission Errors
- **Cause**: Incorrect file permissions on initialization scripts
- **Solution**: Ensure scripts are readable
  ```bash
  chmod 644 init-db/*.sql
  ```

#### SQL Syntax Errors
- **Cause**: Invalid SQL in initialization scripts
- **Solution**: Test scripts manually before adding to init-db
  ```bash
  psql -h localhost -U posteaze_user -d posteaze_db -f init-db/script.sql
  ```

### Debugging Initialization
To debug initialization issues:

1. **Check Container Logs**:
   ```bash
   docker-compose logs postgres
   ```

2. **Connect to Database**:
   ```bash
   docker exec -it posteaze-postgres psql -U posteaze_user -d posteaze_db
   ```

3. **Verify Database State**:
   ```sql
   \l                    -- List databases
   \dt                   -- List tables
   \du                   -- List users
   ```

## Best Practices

### Script Organization
- **Numeric Prefixes**: Use `01-`, `02-`, etc. for execution order
- **Descriptive Names**: Clear file names indicating purpose
- **Single Responsibility**: One script per logical setup task

### SQL Guidelines
- **Idempotent Scripts**: Use `IF NOT EXISTS` clauses
- **Error Handling**: Include proper error checking
- **Comments**: Document script purpose and usage

### Security Considerations
- **No Hardcoded Secrets**: Use environment variables for sensitive data
- **Minimal Permissions**: Grant only necessary database permissions
- **Secure Defaults**: Configure database with security best practices

## Alternative Approaches

### Migration-Only Strategy (Current)
- **Pros**: Version-controlled schema changes, rollback capability
- **Cons**: Manual migration execution required

### Init Scripts + Migrations
- **Pros**: Automated setup with version-controlled changes
- **Cons**: More complex setup process

### Application-Managed Setup
- **Pros**: Full control from application code
- **Cons**: Requires application-specific database management

## Related Documentation

- [Backend Migrations](../backend/migrations/README.md) - Database schema migration procedures
- [Database Utilities](../backend/utils/database/README.md) - Database connection and query utilities
- [Backend Entities](../backend/entities/README.md) - Data models and repository patterns
- [Docker Compose Configuration](../docker-compose.yml) - Container orchestration setup

## Future Enhancements

### Potential Improvements
1. **Automated Migration Integration**: Run migrations during container initialization
2. **Environment-Specific Scripts**: Different initialization for dev/staging/prod
3. **Data Seeding**: Automated test data population for development
4. **Health Checks**: Database readiness verification scripts
5. **Backup Integration**: Automated backup setup during initialization

### Migration Integration Example
```bash
#!/bin/bash
# Example: init-db/99-run-migrations.sh
echo "Running database migrations..."
for migration in /app/migrations/*.up.sql; do
    psql -U $POSTGRES_USER -d $POSTGRES_DB -f "$migration"
done
echo "Migrations completed"
```