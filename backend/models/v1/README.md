# Models v1

This package contains version 1 API models for PostEaze, defining the data structures used in HTTP requests and responses. These models provide validation, type safety, and clear contracts between the frontend and backend systems.

## Architecture

The v1 models are organized by functional domain:
- **User models** (`user.go`) - Authentication, user management, and team structures
- **Token models** (`tokens.go`) - Authentication token management
- **Log models** (`log.go`) - System logging and monitoring data structures

## Key Files

- **user.go**: User authentication, profile management, and team relationship models
- **tokens.go**: Refresh token and authentication session models
- **log.go**: Application logging and log retrieval models

## User Models

### Authentication Models

#### SignupParams
Request model for user registration with validation rules:

```go
type SignupParams struct {
    Name     string   `json:"name" binding:"required,min=2"`
    Email    string   `json:"email" binding:"required,email"`
    Password string   `json:"password" binding:"required,min=8"`
    UserType UserType `json:"user_type" binding:"required,oneof=individual team"`
    TeamName string   `json:"team_name" binding:"required_if=UserType team"`
}
```

**Validation Rules:**
- Name: Required, minimum 2 characters
- Email: Required, valid email format
- Password: Required, minimum 8 characters
- UserType: Must be "individual" or "team"
- TeamName: Required only when UserType is "team"

#### LoginParams
Request model for user authentication:

```go
type LoginParams struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}
```

#### RefreshTokenParams
Request model for token refresh operations:

```go
type RefreshTokenParams struct {
    RefreshToken string `json:"refresh_token" binding:"required"`
}
```

### Core Domain Models

#### User
Primary user model for API responses:

```go
type User struct {
    ID        string    `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    Password  string    `json:"-"`                    // Hidden from JSON
    UserType  UserType  `json:"user_type"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    
    Memberships []TeamMember `json:"memberships,omitempty"`
}
```

**Key Features:**
- Password field excluded from JSON serialization (`json:"-"`)
- Optional memberships for team relationship data
- Timestamps for audit trail

#### Team
Team management model:

```go
type Team struct {
    ID        string       `json:"id"`
    Name      string       `json:"name"`
    OwnerID   string       `json:"owner_id"`
    Owner     User         `json:"owner"`
    Members   []TeamMember `json:"members"`
    CreatedAt time.Time    `json:"created_at"`
    UpdatedAt time.Time    `json:"updated_at"`
}
```

#### TeamMember
Team membership relationship model:

```go
type TeamMember struct {
    ID     string `json:"id"`
    TeamID string `json:"team_id"`
    UserID string `json:"user_id"`
    Role   Role   `json:"role"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

### Enumeration Types

#### UserType
Defines user account types:

```go
type UserType string

const (
    UserTypeIndividual UserType = "individual"
    UserTypeTeam       UserType = "team"
)
```

#### Role
Defines team member roles and permissions:

```go
type Role string

const (
    RoleAdmin   Role = "admin"    // Full team management permissions
    RoleEditor  Role = "editor"   // Content creation and editing
    RoleCreator Role = "creator"  // Content creation only
)
```

## Token Models

### RefreshToken
Database model for refresh token management:

```go
type RefreshToken struct {
    ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    UserID    uuid.UUID `gorm:"type:uuid;index"`
    User      User      `gorm:"foreignKey:UserID"`
    Token     string    `gorm:"uniqueIndex"`
    ExpiresAt time.Time
    Revoked   bool
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

**Features:**
- UUID primary key with auto-generation
- Foreign key relationship to User
- Unique token constraint
- Expiration and revocation support
- Audit timestamps

## Log Models

### LogEntry
Individual log record structure:

```go
type LogEntry struct {
    Timestamp string            `json:"timestamp"`
    Level     string            `json:"level"`
    Message   string            `json:"message"`
    LogID     string            `json:"log_id,omitempty"`
    Method    string            `json:"method,omitempty"`
    Path      string            `json:"path,omitempty"`
    Status    int               `json:"status,omitempty"`
    Duration  string            `json:"duration,omitempty"`
    IP        string            `json:"ip,omitempty"`
    UserAgent string            `json:"user_agent,omitempty"`
    Extra     map[string]string `json:"extra,omitempty"`
}
```

**Log Levels:** INFO, WARN, ERROR, DEBUG
**HTTP Context:** Method, path, status, duration, client info
**Extensibility:** Extra field for additional metadata

### LogsResponse
API response wrapper for log queries:

```go
type LogsResponse struct {
    Success bool       `json:"success"`
    Data    []LogEntry `json:"data"`
    Total   int        `json:"total"`
    Page    int        `json:"page"`
    Limit   int        `json:"limit"`
    Message string     `json:"message,omitempty"`
}
```

**Pagination Support:**
- Total record count
- Current page number
- Records per page limit
- Success/error status

## Usage Examples

### User Registration

```go
// API request
signupData := modelsv1.SignupParams{
    Name:     "John Doe",
    Email:    "john@example.com",
    Password: "securepassword123",
    UserType: modelsv1.UserTypeIndividual,
}

// Validation occurs automatically during binding
if err := c.ShouldBindJSON(&signupData); err != nil {
    // Handle validation errors
    return
}
```

### Team Creation with Owner

```go
// Create team signup
teamSignup := modelsv1.SignupParams{
    Name:     "Jane Smith",
    Email:    "jane@company.com",
    Password: "teampassword123",
    UserType: modelsv1.UserTypeTeam,
    TeamName: "Marketing Team",
}
```

### Log Query Response

```go
// API response for log retrieval
response := modelsv1.LogsResponse{
    Success: true,
    Data: []modelsv1.LogEntry{
        {
            Timestamp: "2024-01-15T10:30:00Z",
            Level:     "INFO",
            Message:   "User login successful",
            Method:    "POST",
            Path:      "/api/v1/auth/login",
            Status:    200,
            Duration:  "45ms",
        },
    },
    Total: 1,
    Page:  1,
    Limit: 50,
}
```

## Validation and Error Handling

### Binding Validation
Models use Gin's binding validation:

```go
// Automatic validation in handlers
var params modelsv1.LoginParams
if err := c.ShouldBindJSON(&params); err != nil {
    c.JSON(400, gin.H{
        "error": "Validation failed",
        "details": err.Error(),
    })
    return
}
```

### Custom Validation
For complex business rules:

```go
func (s *SignupParams) Validate() error {
    if s.UserType == UserTypeTeam && s.TeamName == "" {
        return errors.New("team name required for team accounts")
    }
    return nil
}
```

## Model Relationships

### User-Team Relationships
- **Individual Users**: Single user accounts
- **Team Owners**: Users who create and manage teams
- **Team Members**: Users belonging to teams with specific roles

### Authentication Flow
1. **Signup**: `SignupParams` → `User` creation
2. **Login**: `LoginParams` → `RefreshToken` generation
3. **Token Refresh**: `RefreshTokenParams` → New access token
4. **Logout**: Token revocation

## Related Documentation

- [Models Overview](../README.md) - General models architecture and versioning
- [Entities](../../entities/README.md) - Database entity models
- [API v1 Handlers](../../api/v1/README.md) - HTTP handlers using these models
- [Business Logic v1](../../business/v1/README.md) - Services processing these models