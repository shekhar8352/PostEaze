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

1. **User / auth** (`v1/user.go`): Firebase auth params, user DTOs, updates
2. **Token models** (`v1/tokens.go`): Refresh token payloads
3. **Log models** (`v1/log.go`): Log API shapes
4. **Team, channel, analytics** — see files under `v1/`

## Key features

### Validation tags
Example (see source for current fields):

```go
type FirebaseAuthParams struct {
    FirebaseID    string `json:"firebase_id" binding:"required"`
    FirebaseToken string `json:"firebase_token" binding:"required"`
    Platform      string `json:"platform" binding:"required,oneof=email google facebook microsoft"`
}
```

### JSON serialization
- Sensitive fields use `json:"-"` where needed
- Optional fields often use `omitempty`

### Database integration

Most persistence uses raw SQL via `entities` and `repositories`. A few structs in `models/v1` still carry `gorm` struct tags for historical or auxiliary use; the primary path is not GORM-backed in runtime code—prefer entities for DB shape.

## Usage Patterns

### Request/Response Models
Models are used for:
- **API request validation** — binding tags on request bodies
- **API response serialization** — response DTOs
- **Persistence** — usually via `entities` / SQL, not ORM

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