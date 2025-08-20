# Services

API service layer and HTTP client configuration for the PostEaze frontend application.

## Purpose

This directory contains the centralized HTTP client configuration and API service utilities that handle communication between the React frontend and the Go backend API. The services layer provides a consistent interface for making HTTP requests across all features.

## Contents

### Key Files

- **`axios.ts`** - Main HTTP client configuration with axios instance setup, base URL configuration, and credential handling

## Architecture

The services layer follows a centralized HTTP client pattern:

- **Single Axios Instance**: One configured axios instance used throughout the application
- **Base Configuration**: Centralized API base URL and credential settings
- **Feature Integration**: Individual features import and use the configured client

### HTTP Client Configuration

The axios instance is configured with:
- **Base URL**: `http://localhost:3000/api` - Points to the Go backend API
- **Credentials**: `withCredentials: true` - Enables cookie-based authentication
- **Consistent Interface**: Provides standard HTTP methods (GET, POST, PUT, DELETE)

## Usage Examples

### Basic Import and Usage
```typescript
import axiosInstance from '@/services/axios';

// Making API requests
const response = await axiosInstance.post('/auth/signup', userData);
const data = await axiosInstance.get('/users/profile');
```

### Feature Integration Pattern
Features create their own API modules that use the configured axios instance:

```typescript
// In feature API files (e.g., features/auth/authApi.ts)
import axiosInstance from '@/services/axios';

export const signup = async (data: SignupData) => {
  const response = await axiosInstance.post('/auth/signup', data);
  return response.data;
};
```

## API Integration Patterns

### Request/Response Flow
1. **Feature Components** call API functions from their respective API modules
2. **API Modules** use the configured axios instance to make HTTP requests
3. **Axios Instance** handles base URL, credentials, and request configuration
4. **Backend API** processes requests and returns responses
5. **Response Data** is returned to components for state updates

### Error Handling
The axios instance can be extended with interceptors for:
- Request preprocessing (adding auth tokens, logging)
- Response error handling (401 redirects, error formatting)
- Loading state management

### Authentication Integration
The `withCredentials: true` configuration enables:
- Cookie-based session management
- Automatic credential inclusion in requests
- Seamless authentication with the Go backend

## Configuration

### Environment-Specific Settings
The base URL is currently hardcoded but can be made configurable:
```typescript
const baseURL = import.meta.env.VITE_API_URL || 'http://localhost:3000/api';
```

### Future Enhancements
- Request/response interceptors for error handling
- Loading state management
- Request timeout configuration
- Retry logic for failed requests

## Related Documentation

- **Features**: `../features/` - Individual feature API implementations
- **App Store**: `../app/` - Redux store configuration for API state management
- **Utils**: `../utils/` - Utility functions that may complement API services

## Dependencies

- **axios** - HTTP client library for making API requests
- **TypeScript** - Type safety for API request/response data
- **Vite** - Build tool providing environment variable support