# Authentication Feature

The authentication feature handles user registration, login, and authentication state management for the PostEaze application. This feature follows a modular architecture with Redux Toolkit for state management, React Router for routing, and a clean separation between API calls, business logic, and UI components.

## Architecture

The auth feature is organized using a feature-based architecture pattern:

```
auth/
├── components/          # Reusable auth-related components
│   └── signUp/         # SignUp form component
├── pages/              # Auth page components
│   └── SignUp/         # SignUp page wrapper
├── authApi.ts          # API service functions
├── authRoutes.tsx      # Route definitions
├── authSlice.ts        # Redux slice for auth state
├── thunks.ts          # Async Redux thunks
└── index.tsx          # Feature exports
```

## Key Components

### State Management (authSlice.ts)
- **Purpose**: Manages authentication state using Redux Toolkit
- **State Shape**: 
  - `loading`: Boolean indicating async operation status
  - `error`: Error message string or null
  - `user`: User object or null
- **Actions**: Handles signup flow with pending, fulfilled, and rejected states

### API Layer (authApi.ts)
- **Purpose**: Contains HTTP client functions for authentication endpoints
- **Functions**:
  - `signup(data)`: Registers a new user with name, email, and password

### Async Actions (thunks.ts)
- **Purpose**: Redux async thunks for handling side effects
- **Thunks**:
  - `signupUser`: Async thunk that calls the signup API and updates state

### Routing (authRoutes.tsx)
- **Purpose**: Defines authentication-related routes
- **Routes**:
  - `/signup`: SignUp page route

## Components Structure

### SignUp Component (`components/signUp/`)
- **Purpose**: Reusable signup form component
- **Files**:
  - `index.tsx`: Main component implementation
  - `style.module.css`: Component-specific styles
  - `index.test.tsx`: Unit tests

### SignUp Page (`pages/SignUp/`)
- **Purpose**: Page-level wrapper for signup functionality
- **Files**:
  - `index.tsx`: Page component that uses SignUp component
  - `style.module.css`: Page-specific styles

## Usage Examples

### Using the Auth Slice
```typescript
import { useAppSelector, useAppDispatch } from '@/app/hooks';
import { signupUser } from '@/features/auth';

const MyComponent = () => {
  const dispatch = useAppDispatch();
  const { loading, error, user } = useAppSelector(state => state.auth);

  const handleSignup = (userData) => {
    dispatch(signupUser(userData));
  };

  return (
    // Component JSX
  );
};
```

### Adding New Auth Routes
```typescript
// In authRoutes.tsx
const authRoutes: RouteObject[] = [
  { path: "/signup", element: <SignupPage /> },
  { path: "/login", element: <LoginPage /> }, // New route
];
```

## State Flow

1. **User Interaction**: User fills out signup form and submits
2. **Thunk Dispatch**: Component dispatches `signupUser` thunk
3. **API Call**: Thunk calls `signup` function from authApi
4. **State Update**: Slice handles pending/fulfilled/rejected states
5. **UI Update**: Components re-render based on updated state

## Dependencies

- **Redux Toolkit**: State management and async thunks
- **React Router**: Route definitions and navigation
- **Axios**: HTTP client for API calls (via services layer)
- **Mantine**: UI components (used in signup form)

## Testing

The feature includes unit tests for components:
- `components/signUp/index.test.tsx`: Tests for SignUp component

## Future Enhancements

The current implementation supports signup functionality. Future enhancements may include:
- Login functionality
- Password reset
- Email verification
- Social authentication
- JWT token management
- Protected route handling

## Related Documentation

- [Features Overview](../README.md) - Feature-based architecture explanation
- [Services](../../services/README.md) - API service layer documentation
- [App Store](../../app/README.md) - Redux store configuration
- [Routes](../../routes/README.md) - Application routing setup