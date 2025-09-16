# Firebase Authentication Implementation

## Overview

This document describes the Firebase authentication integration implemented for the PostEaze API. The system supports authentication through multiple platforms (Google, Facebook, Microsoft, Email) using Firebase tokens while maintaining local user data.

## Features

- **Multi-platform Authentication**: Supports Google, Facebook, Microsoft, and Email authentication
- **Firebase Token Validation**: Server-side validation using Firebase Admin SDK
- **User Management**: Automatic user creation and platform tracking
- **Backward Compatibility**: Maintains existing password-based authentication
- **Platform-specific Validation**: Email requirements based on authentication platform

## API Endpoints

### POST /api/v1/auth/authenticate

Authenticates users using Firebase tokens and creates/updates local user records.

**Request Body:**
```json
{
  "local_id": "firebase-uid-123",
  "firebase_token": "eyJhbGciOiJSUzI1NiIs...",
  "platform": "google"
}
```

**Parameters:**
- `local_id` (required): Firebase UID
- `firebase_token` (required): Firebase ID token
- `platform` (required): Authentication platform (`email`, `google`, `facebook`, `microsoft`)

**Success Response (200):**
```json
{
  "status": "success",
  "msg": "Authenticated successfully",
  "data": {
    "user": {
      "id": "firebase-uid-123",
      "name": "John Doe",
      "email": "john@example.com",
      "platforms": ["google"],
      "created_at": "2025-01-01T00:00:00Z",
      "updated_at": "2025-01-01T00:00:00Z"
    },
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
  }
}
```

**Error Responses:**
- `400 Bad Request`: Invalid request payload or missing email for non-Facebook platforms
- `401 Unauthorized`: Invalid or expired Firebase token
- `500 Internal Server Error`: Firebase SDK or database errors

### GET /api/v1/auth/authenticate

Validates existing access tokens (existing functionality).

**Headers:**
- `Authorization: Bearer <access_token>`

## User Model Changes

The User model has been updated to support Firebase authentication:

### Removed Fields
- `password`: Authentication handled by Firebase
- `user_type`: No longer needed with Firebase auth

### Modified Fields
- `email`: Now optional (Facebook may not provide email)

### Added Fields
- `platforms`: Array of authentication platforms used by the user

### Updated User Structure
```go
type User struct {
    ID        string    `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email,omitempty"`
    Platforms []string  `json:"platforms"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

## Platform-Specific Validation

### Email Requirements
- **Google**: Email required (Google always provides email)
- **Microsoft**: Email required (Microsoft always provides email)
- **Email**: Email required (email-based authentication)
- **Facebook**: Email optional (Facebook may not provide email based on user privacy settings)

### Platform Tracking
- Users can authenticate through multiple platforms
- The `platforms` array tracks all authentication methods used
- Duplicate platforms are automatically prevented

## Database Schema Changes

### Migration Script
```sql
-- Add platforms column
ALTER TABLE users ADD COLUMN platforms TEXT[] DEFAULT '{}';

-- Make email column nullable
ALTER TABLE users ALTER COLUMN email DROP NOT NULL;

-- Drop password column
ALTER TABLE users DROP COLUMN IF EXISTS password;

-- Drop user_type column
ALTER TABLE users DROP COLUMN IF EXISTS user_type;

-- Create index for efficient platform queries
CREATE INDEX IF NOT EXISTS idx_users_platforms ON users USING GIN(platforms);
```

## Firebase Configuration

### Environment Variables
```env
FIREBASE_PROJECT_ID=your-project-id
FIREBASE_PRIVATE_KEY_ID=your-private-key-id
FIREBASE_PRIVATE_KEY="-----BEGIN PRIVATE KEY-----\nYOUR_PRIVATE_KEY_HERE\n-----END PRIVATE KEY-----\n"
FIREBASE_CLIENT_EMAIL=firebase-adminsdk-xxx@your-project.iam.gserviceaccount.com
FIREBASE_CLIENT_ID=your-client-id
```

### Service Account Setup
1. Create a Firebase project in the Firebase Console
2. Generate a service account key
3. Download the JSON key file
4. Extract the required fields to environment variables

## Authentication Flow

### New User Flow
1. Client sends Firebase token and platform information
2. Server validates Firebase token using Admin SDK
3. Server extracts user information from Firebase
4. Server creates new user record with Firebase data
5. Server generates local access and refresh tokens
6. Server returns user object and tokens

### Existing User Flow
1. Client sends Firebase token and platform information
2. Server validates Firebase token using Admin SDK
3. Server finds existing user by local_id
4. Server updates user's platforms array if needed
5. Server generates new local tokens
6. Server returns updated user object and tokens

## Error Handling

### Firebase Token Validation
- Invalid tokens return 401 Unauthorized
- Expired tokens return 401 Unauthorized
- Firebase SDK errors return 500 Internal Server Error

### Email Validation
- Missing email for non-Facebook platforms returns 400 Bad Request
- Platform-specific validation enforced server-side

### Database Operations
- Transaction rollback on user creation failures
- Proper error logging for debugging

## Testing

### Unit Tests
- Firebase token validation logic
- Platform-specific email validation
- User creation and update flows
- Error handling scenarios

### Integration Tests
- Complete authentication flow
- API endpoint validation
- Error response testing
- Multi-platform authentication

### Test Files
- `tests/business/firebase_auth_test.go`: Business logic tests
- `tests/api/v1/firebase_auth_test.go`: API handler tests
- `tests/integration/firebase_auth_integration_test.go`: Integration tests

## Security Considerations

### Token Validation
- Server-side Firebase token verification using Admin SDK
- No client-side token validation
- Proper token signature and expiration checking

### Data Privacy
- Minimal user data storage
- No Firebase tokens stored locally
- Platform-specific privacy handling (Facebook email optional)

### Access Control
- Local token generation for API access
- Refresh token management
- Proper session handling

## Backward Compatibility

### Legacy Authentication
- Existing password-based signup/login maintained
- Legacy users automatically assigned "email" platform
- Gradual migration path to Firebase authentication

### API Compatibility
- Existing endpoints unchanged
- New Firebase endpoint added alongside existing auth
- No breaking changes to existing functionality

## Deployment

### Prerequisites
1. Firebase project setup
2. Service account configuration
3. Database migration execution
4. Environment variable configuration

### Deployment Steps
1. Run database migration
2. Update environment variables
3. Deploy updated API
4. Test Firebase authentication
5. Monitor error logs

### Rollback Plan
- Database migration can be reversed
- Firebase service can be disabled
- Legacy authentication remains functional

## Monitoring and Logging

### Key Metrics
- Firebase authentication success/failure rates
- Platform usage distribution
- User creation vs. existing user authentication
- Token validation performance

### Log Events
- Firebase token validation attempts
- User creation events
- Platform updates
- Authentication errors

### Error Tracking
- Firebase SDK errors
- Database operation failures
- Token validation failures
- Platform validation errors