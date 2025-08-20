# Database Migrations

This folder contains SQL migration files that define the database schema evolution for the PostEaze application. Migrations provide a version-controlled way to modify the database structure and ensure consistent schema across different environments.

## Migration Strategy

PostEaze uses a simple, file-based migration system with numbered SQL files that ensures database schema consistency across all environments (development, staging, production). This approach provides:

- **Version Control**: All schema changes are tracked in Git alongside application code
- **Reproducibility**: Identical schema can be recreated in any environment
- **Rollback Safety**: Every change can be safely reverted using down migrations
- **Team Collaboration**: Multiple developers can work on schema changes without conflicts

Each migration consists of two files:
- **Up migration** (`*.up.sql`): Applies changes to move the database forward
- **Down migration** (`*.down.sql`): Reverts changes to roll back the database

### Migration Principles
1. **Sequential Numbering**: Migrations are applied in numerical order
2. **Atomic Changes**: Each migration represents one logical schema change
3. **Backward Compatibility**: New migrations should not break existing application code
4. **Data Preservation**: Schema changes should preserve existing data when possible

## File Structure

```
migrations/
├── 001_create_initial_tables.up.sql    # Initial schema creation
├── 001_create_initial_tables.down.sql  # Initial schema rollback
└── README.md                           # This documentation
```

## Naming Convention

Migration files follow a strict naming pattern:
```
{number}_{description}.{direction}.sql
```

- **Number**: 3-digit sequential number (001, 002, 003, etc.)
- **Description**: Snake_case description of the migration purpose
- **Direction**: Either `up` (apply) or `down` (rollback)
- **Extension**: Always `.sql`

### Examples
- `001_create_initial_tables.up.sql`
- `002_add_user_preferences.up.sql`
- `003_create_posts_table.down.sql`

## Current Schema

The current migration (`001_create_initial_tables`) establishes the core database schema:

### Tables Created
- **users**: User accounts with authentication data
- **teams**: Team/organization management
- **team_members**: User-team relationships with roles
- **refresh_tokens**: JWT refresh token management

### Key Features
- **UUID Primary Keys**: Using `gen_random_uuid()` for globally unique identifiers
- **Extension Dependencies**: Requires `pgcrypto` extension for UUID generation
- **Automatic Timestamps**: All tables include `created_at` and `updated_at` fields
- **Foreign Key Constraints**: Proper relationships with cascade deletion for data integrity
- **Unique Constraints**: Prevent duplicate data (e.g., unique emails, team memberships)
- **Referential Integrity**: All foreign keys properly reference parent tables

## Migration Procedures

### Development Workflow
1. **Plan the Change**: Document what schema changes are needed and why
2. **Create Migration Files**: Generate both up and down migration files
3. **Test Locally**: Apply and rollback migrations on local development database
4. **Code Review**: Have migrations reviewed by team members
5. **Test on Staging**: Apply migrations to staging environment
6. **Deploy to Production**: Apply migrations during deployment window

### Deployment Strategy
- **Pre-deployment**: Apply schema changes that are backward compatible
- **Post-deployment**: Apply changes that require new application code
- **Rollback Plan**: Always have a tested rollback procedure ready

### Emergency Procedures
If a migration causes issues in production:
1. **Immediate Assessment**: Determine if rollback is safe
2. **Application Rollback**: Revert application to previous version if needed
3. **Schema Rollback**: Apply down migration to revert schema changes
4. **Data Recovery**: Restore from backup if data corruption occurred

## Creating New Migrations

### Step 1: Determine Migration Number
Find the highest numbered migration and increment by 1:
```bash
# If last migration is 001, next should be 002
ls backend/migrations/ | grep -E '^[0-9]+' | sort -n | tail -1
```

### Step 2: Create Migration Files
Create both up and down migration files:
```bash
# Example for adding a new posts table
touch backend/migrations/002_create_posts_table.up.sql
touch backend/migrations/002_create_posts_table.down.sql
```

### Step 3: Write Migration SQL

**Up Migration Example** (`002_create_posts_table.up.sql`):
```sql
-- Create posts table
CREATE TABLE IF NOT EXISTS posts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    title TEXT NOT NULL,
    content TEXT,
    status TEXT NOT NULL DEFAULT 'draft',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_post_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create index for better query performance
CREATE INDEX IF NOT EXISTS idx_posts_user_id ON posts(user_id);
CREATE INDEX IF NOT EXISTS idx_posts_status ON posts(status);
```

**Down Migration Example** (`002_create_posts_table.down.sql`):
```sql
-- Drop indexes first
DROP INDEX IF EXISTS idx_posts_status;
DROP INDEX IF EXISTS idx_posts_user_id;

-- Drop the posts table
DROP TABLE IF EXISTS posts;
```

## Migration Execution

### Prerequisites
Before running migrations, ensure:
- PostgreSQL server is running and accessible
- Database exists (created during initial setup)
- User has appropriate permissions (CREATE, DROP, ALTER)
- Required extensions can be installed (`pgcrypto`)

### Manual Execution
Currently, migrations are executed manually using PostgreSQL client tools:

```bash
# Apply up migration
psql -h localhost -U your_user -d posteaze_db -f backend/migrations/001_create_initial_tables.up.sql

# Rollback with down migration
psql -h localhost -U your_user -d posteaze_db -f backend/migrations/001_create_initial_tables.down.sql
```

### Docker Environment
For containerized deployments, migrations can be run against the PostgreSQL container:

```bash
# Execute migration in Docker container
docker exec -i posteaze-postgres psql -U $POSTGRES_USER -d $POSTGRES_DB < backend/migrations/001_create_initial_tables.up.sql

# Or connect interactively and run migration
docker exec -it posteaze-postgres psql -U $POSTGRES_USER -d $POSTGRES_DB
\i /path/to/migration/file.sql
```

### Environment-Specific Execution
Different environments may require different connection parameters:

```bash
# Development
psql -h localhost -p 5432 -U dev_user -d posteaze_dev -f migration.up.sql

# Staging
psql -h staging-db.example.com -p 5432 -U staging_user -d posteaze_staging -f migration.up.sql

# Production (with SSL)
psql "postgresql://prod_user:password@prod-db.example.com:5432/posteaze_prod?sslmode=require" -f migration.up.sql
```

## Best Practices

### Migration Design
- **Incremental Changes**: Each migration should represent a single, logical change
- **Reversible**: Always provide corresponding down migrations
- **Idempotent**: Use `IF NOT EXISTS` and `IF EXISTS` clauses to prevent errors on re-runs
- **Data Safety**: Consider data preservation when modifying existing tables

### SQL Guidelines
- Use `CREATE TABLE IF NOT EXISTS` for table creation
- Use `DROP TABLE IF EXISTS` for table removal
- Include proper foreign key constraints
- Add indexes for frequently queried columns
- Use consistent naming conventions (snake_case)

### Testing Migrations
Before applying migrations to production, thoroughly test them:

1. **Test Up Migration**: Verify the migration applies cleanly
   ```bash
   # Test on a copy of production data
   psql -d test_db -f new_migration.up.sql
   ```

2. **Test Down Migration**: Ensure rollback works correctly
   ```bash
   # Apply then immediately rollback
   psql -d test_db -f new_migration.up.sql
   psql -d test_db -f new_migration.down.sql
   ```

3. **Test Data Integrity**: Verify existing data remains intact
   - Check row counts before and after migration
   - Verify foreign key relationships are maintained
   - Ensure no data corruption occurs

4. **Test Performance**: Check that new indexes improve query performance
   ```sql
   -- Before migration
   EXPLAIN ANALYZE SELECT * FROM users WHERE email = 'test@example.com';
   
   -- After migration (if adding index)
   EXPLAIN ANALYZE SELECT * FROM users WHERE email = 'test@example.com';
   ```

5. **Test Application Compatibility**: Ensure the application works with the new schema
   - Run application tests against migrated database
   - Verify all queries still function correctly
   - Check that new features work as expected

## Database Connection

The application connects to the database using configuration from:
- **Config Files**: `backend/resources/configs/database.json`
- **Environment Variables**: Database URL and credentials
- **Connection Pool**: Managed by `backend/utils/database/database.go`

## Future Enhancements

### Automated Migration System
Consider implementing an automated migration runner that:
- Tracks applied migrations in a `schema_migrations` table
- Automatically applies pending migrations on application startup
- Provides CLI commands for migration management
- Supports migration rollback with safety checks

### Migration Tools
Potential tools to integrate:
- **golang-migrate**: Popular Go migration library
- **Goose**: Database migration tool for Go
- **Custom Solution**: Application-specific migration runner

## Related Documentation

- [Database Utilities](../utils/database/README.md) - Database connection and query utilities
- [Entities](../entities/README.md) - Data models and entity definitions
- [Models](../models/README.md) - Application data models
- [Database Initialization](../../init-db/README.md) - Docker database setup