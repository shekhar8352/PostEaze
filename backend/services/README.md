# Services Package

This package provides various services used throughout the PostEaze backend application.

## Available Services

- [Redis Service](#redis-service)
- [Email Service](#email-service)

---

## Redis Service

The `RedisService` provides an interface for interacting with the Redis datastore. It supports setting keys with optional expiration, getting values, and deleting keys.

### Initialization

To use the Redis service, create a new instance using `NewRedisService()`:

```go
import "github.com/shekhar8352/PostEaze/services/redis_service"

redisService := redis_service.NewRedisService()
```

### Methods

#### `Set`
Sets a key-value pair in Redis.
- **Signature**: `Set(ctx context.Context, key string, value interface{}, expiration ...time.Duration) error`
- **Parameters**:
  - `ctx`: Context for the operation.
  - `key`: The key to set.
  - `value`: The value to store.
  - `expiration` (Optional): Duration until the key expires. **Defaults to 24 hours if omitted.**

#### `Get`
Retrieves a value by its key.
- **Signature**: `Get(ctx context.Context, key string) (string, error)`

#### `Delete`
Removes a key from Redis.
- **Signature**: `Delete(ctx context.Context, key string) error`

### Example Usage

```go
ctx := context.Background()
service := redis_service.NewRedisService()

// Set with default 24h expiration
err := service.Set(ctx, "user_session:123", "active")

// Set with custom expiration (e.g., 5 minutes)
err := service.Set(ctx, "temp_code:456", "123456", 5*time.Minute)

// Get a value
val, err := service.Get(ctx, "user_session:123")
if err != nil {
    // Handle error (e.g., key not found)
}

// Delete a key
err = service.Delete(ctx, "user_session:123")
```

---

## Email Service

The `EmailService` handles sending emails via SMTP (specifically configured for Gmail).

### Initialization

To use the Email service, create a new instance using `NewGmailEmailService()`:

```go
import "github.com/shekhar8352/PostEaze/services/email_service"

emailService := email_service.NewGmailEmailService()
```

*Note: Requires `SMTP_HOST`, `SMTP_PORT`, `SMTP_EMAIL`, and `SMTP_PASSWORD` environment variables to be set.*

### Methods

#### `SendEmail`
Sends a generic HTML email.
- **Signature**: `SendEmail(to []string, subject string, body string) error`

#### `SendNotificationEmail`
Sends a pre-formatted notification email.
- **Signature**: `SendNotificationEmail(to string, notificationContent string) error`

#### `SendTeamInviteEmail`
Sends a pre-formatted team invitation email.
- **Signature**: `SendTeamInviteEmail(to string, inviteLink string) error`

### Example Usage

```go
service := email_service.NewGmailEmailService()

// Send a generic email
err := service.SendEmail([]string{"user@example.com"}, "Welcome!", "<h1>Hello</h1>")

// Send a notification
err := service.SendNotificationEmail("user@example.com", "You have a new message.")

// Send an invite
err := service.SendTeamInviteEmail("colleague@example.com", "https://posteaze.com/join/123")
```
