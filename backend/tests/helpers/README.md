# Test Helpers

This directory contains simplified test helpers for the PostEaze backend testing framework.

## Database Helper (`db.go`)

The database helper provides a simple, fast way to set up test databases using in-memory SQLite.

### Features

- **In-memory SQLite**: Fast test execution without external dependencies
- **Simple fixture loading**: Easy test data management
- **Automatic cleanup**: Proper resource management
- **Transaction support**: Test isolation with transactions
- **Minimal abstraction**: Direct database access when needed

### Basic Usage

```go
func TestSomething(t *testing.T) {
    // Create test database
    db, err := NewTestDB()
    if err != nil {
        t.Fatalf("Failed to create test database: %v", err)
    }
    defer db.Cleanup() // Always cleanup

    // Load test data
    user := NewTestUser(func(u *TestUser) {
        u.Name = "Test User"
        u.Email = "test@example.com"
    })
    
    if err := db.LoadFixture(user); err != nil {
        t.Fatalf("Failed to load fixture: %v", err)
    }

    // Your test logic here...
}
```

### Loading Fixtures

#### Predefined Data
```go
// Load default users
if err := db.LoadFixture(DefaultUsers); err != nil {
    t.Fatalf("Failed to load users: %v", err)
}

// Load default teams
if err := db.LoadFixture(DefaultTeams); err != nil {
    t.Fatalf("Failed to load teams: %v", err)
}
```

#### Custom Data
```go
// Create custom user
user := NewTestUser(func(u *TestUser) {
    u.Name = "Custom Name"
    u.Email = "custom@example.com"
    u.UserType = "team"
})

// Load single fixture
if err := db.LoadFixture(user); err != nil {
    t.Fatalf("Failed to load user: %v", err)
}

// Load multiple fixtures
users := []TestUser{user1, user2, user3}
if err := db.LoadFixture(users); err != nil {
    t.Fatalf("Failed to load users: %v", err)
}
```

### Test Data Structures

#### TestUser
```go
type TestUser struct {
    ID        string
    Name      string
    Email     string
    Password  string
    UserType  string
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

#### TestTeam
```go
type TestTeam struct {
    ID        string
    Name      string
    OwnerID   string
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

#### TestToken
```go
type TestToken struct {
    ID        string
    UserID    string
    Token     string
    ExpiresAt time.Time
    Revoked   bool
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

### Helper Functions

#### Creating Test Data
```go
// Create user with defaults
user := NewTestUser()

// Create user with overrides
user := NewTestUser(func(u *TestUser) {
    u.Name = "Custom Name"
    u.Email = "custom@example.com"
})

// Create team
team := NewTestTeam("owner-id", func(t *TestTeam) {
    t.Name = "Custom Team"
})

// Create token
token := NewTestToken("user-id", func(t *TestToken) {
    t.Token = "custom_token"
    t.Revoked = true
})
```

### Database Operations

#### Direct Database Access
```go
// Query the database directly
var count int
err := db.DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
if err != nil {
    t.Fatalf("Query failed: %v", err)
}
```

#### Transactions
```go
ctx := context.Background()
tx, err := db.BeginTx(ctx)
if err != nil {
    t.Fatalf("Failed to begin transaction: %v", err)
}

// Do work in transaction...

// Rollback for test isolation
if err := tx.Rollback(); err != nil {
    t.Fatalf("Failed to rollback: %v", err)
}
```

#### Cleanup
```go
// Clean all data (keeps database open)
if err := db.CleanData(); err != nil {
    t.Fatalf("Failed to clean data: %v", err)
}

// Full cleanup (closes database)
if err := db.Cleanup(); err != nil {
    t.Fatalf("Failed to cleanup: %v", err)
}
```

### Best Practices

1. **Always use defer cleanup**:
   ```go
   db, err := NewTestDB()
   if err != nil {
       t.Fatalf("Failed to create test database: %v", err)
   }
   defer db.Cleanup()
   ```

2. **Clean data between test cases**:
   ```go
   for _, tc := range testCases {
       t.Run(tc.name, func(t *testing.T) {
           if err := db.CleanData(); err != nil {
               t.Fatalf("Failed to clean data: %v", err)
           }
           // Test logic...
       })
   }
   ```

3. **Load dependencies first**:
   ```go
   // Load user first
   user := NewTestUser()
   if err := db.LoadFixture(user); err != nil {
       t.Fatalf("Failed to load user: %v", err)
   }

   // Then load team that references the user
   team := NewTestTeam(user.ID)
   if err := db.LoadFixture(team); err != nil {
       t.Fatalf("Failed to load team: %v", err)
   }
   ```

4. **Use transactions for test isolation**:
   ```go
   tx, err := db.BeginTx(ctx)
   if err != nil {
       t.Fatalf("Failed to begin transaction: %v", err)
   }
   defer tx.Rollback() // Always rollback in tests
   ```

### Migration from Old Testing Framework

The new database helper replaces the complex `testutils/database.go` with a much simpler interface:

#### Old Way
```go
db, cleanup, err := testutils.SetupTestDB(ctx)
if err != nil {
    t.Fatalf("Failed to setup database: %v", err)
}
defer cleanup()

if err := testutils.LoadFixtures(ctx, db, fixtures...); err != nil {
    t.Fatalf("Failed to load fixtures: %v", err)
}
```

#### New Way
```go
db, err := helpers.NewTestDB()
if err != nil {
    t.Fatalf("Failed to create test database: %v", err)
}
defer db.Cleanup()

if err := db.LoadFixture(fixtures); err != nil {
    t.Fatalf("Failed to load fixtures: %v", err)
}
```

### Performance

- **In-memory SQLite**: Tests run in milliseconds
- **No external dependencies**: No need for PostgreSQL test database
- **Simple cleanup**: Fast data clearing between tests
- **Minimal overhead**: Direct database access without abstraction layers

### Examples

See `example_test.go` for complete examples of how to use the database helper in various testing scenarios.