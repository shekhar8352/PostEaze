# Models

Data model definitions and structures for the PostEaze backend application. This folder contains versioned data models that define the shape and validation rules for data flowing through the application.

## Architecture

The models package follows a versioned architecture pattern to support API evolution and backward compatibility. Each version contains complete model definitions for that API version, allowing for independent evolution of data structures.

## Contents

- **v1/**: Version 1 data models containing core application entities
- **README.md**: This documentation file

## Model Organization

### Versioning Strategy

Models are organized by API version to support:
- **Backward Compatibility**: Older API versions continue to work with their specific models
- **Independent Evolution**: New versions can modify data structures without breaking existing clients
- **Clear Migration Path**: Developers can see differences between versions and plan migrations

### Model Categories

The models are organized into logical categories:

1. **User Models** (`v1/user.go`): User authentication, profiles, and team management
2. **Token Models** (`v1/tokens.go`): JWT refresh tokens and authentication state
3. **Log Models** (`v1/log.go`): Application logging and audit trail structures

## Key Features

### Validation Tags
Models use Go struct tags for validation:
```go
type SignupParams struct {
    Name     string   `json:"name" binding:"required,min=2"`
    Email    string   `json:"email" binding:"required,email"`
    Password string   `json:"password" binding:"required,min=8"`
}
```

### JSON Serialization
All models include JSON tags for API serialization:
- Standard fields are serialized normally
- Sensitive fields (like passwords) use `json:"-"` to exclude from output
- Optional fields use `omitempty` to reduce payload size

### Database Integration
Models include GORM tags for database mapping:
```go
type RefreshToken struct {
    ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    UserID    uuid.UUID `gorm:"type:uuid;index"`
}
```

## Usage Patterns

### Request/Response Models
Models are used for:
- **API Request Validation**: Input parameters with validation rules
- **API Response Serialization**: Structured output with proper JSON formatting
- **Database Operations**: Entity definitions with ORM mappings

### Model Relationships
The models define clear relationships:
- Users can belong to multiple teams through TeamMember associations
- Teams have owners and members with role-based permissions
- Refresh tokens are linked to specific users for authentication

## Adding New Models

When adding new models:

1. **Choose Version**: Add to existing version or create new version folder
2. **Follow Patterns**: Use consistent struct tags and naming conventions
3. **Include Validation**: Add appropriate binding tags for input validation
4. **Document Relationships**: Clearly define foreign key relationships
5. **Update Tests**: Ensure model validation and serialization work correctly

## Related Documentation

- [Entities](../entities/README.md): Database entity implementations using these models
- [API v1](../api/v1/README.md): API endpoints that consume these models
- [Business v1](../business/v1/README.md): Business logic that operates on these models