# Test Utils

Test utility functions that provide common testing helpers and setup for React components in the PostEaze frontend application.

## Purpose

This directory contains utility functions that simplify testing React components, particularly those that require Redux store integration and Mantine UI provider setup.

## Available Utilities

### renderWithStore

The primary test utility for rendering React components with full provider setup.

**Location**: `renderWithStore.tsx`

**Purpose**: Renders React components with Redux store and Mantine provider context, enabling comprehensive testing of connected components.

**Features**:
- Automatic Redux store configuration with auth reducer
- Mantine UI provider setup with test theme
- Customizable preloaded state for testing different scenarios
- Option to provide custom store configuration
- Returns both the configured store and render result for assertions

## Usage Examples

### Basic Component Testing

```typescript
import { renderWithStore } from '@/test/utils/renderWithStore';
import { MyComponent } from './MyComponent';

test('renders component with default store', () => {
  const { getByText } = renderWithStore(<MyComponent />);
  expect(getByText('Hello World')).toBeInTheDocument();
});
```

### Testing with Preloaded State

```typescript
import { renderWithStore } from '@/test/utils/renderWithStore';
import { LoginForm } from './LoginForm';

test('renders login form with authenticated state', () => {
  const preloadedState = {
    auth: {
      isAuthenticated: true,
      user: { id: 1, email: 'test@example.com' }
    }
  };

  const { getByText } = renderWithStore(
    <LoginForm />, 
    { preloadedState }
  );
  
  expect(getByText('Welcome back!')).toBeInTheDocument();
});
```

### Testing with Custom Store

```typescript
import { configureStore } from '@reduxjs/toolkit';
import { renderWithStore } from '@/test/utils/renderWithStore';
import authReducer from '@/features/auth/authSlice';
import { MyConnectedComponent } from './MyConnectedComponent';

test('renders with custom store configuration', () => {
  const customStore = configureStore({
    reducer: {
      auth: authReducer,
      // Add other reducers as needed
    },
    middleware: (getDefaultMiddleware) =>
      getDefaultMiddleware({
        serializableCheck: false,
      }),
  });

  const { store, getByRole } = renderWithStore(
    <MyConnectedComponent />,
    { store: customStore }
  );

  expect(getByRole('button')).toBeInTheDocument();
  // Can also make assertions about store state
  expect(store.getState().auth.isAuthenticated).toBe(false);
});
```

### Accessing Store in Tests

```typescript
import { renderWithStore } from '@/test/utils/renderWithStore';
import { UserProfile } from './UserProfile';

test('updates store state when user interacts', async () => {
  const { store, getByRole } = renderWithStore(<UserProfile />);
  
  const button = getByRole('button', { name: 'Update Profile' });
  fireEvent.click(button);
  
  // Assert store state changes
  expect(store.getState().auth.loading).toBe(true);
});
```

## Implementation Details

### Provider Setup

The `renderWithStore` utility automatically wraps components with:

1. **Redux Provider**: Provides access to the Redux store
2. **Mantine Provider**: Provides Mantine UI theme and components context

### Default Configuration

- **Store**: Configured with `authReducer` by default
- **Theme**: Uses a basic test theme optimized for testing
- **Preloaded State**: Empty object by default, can be customized

### Flexibility

The utility supports various testing scenarios:
- Components that don't use Redux (still works with default store)
- Components that need specific auth states
- Components requiring custom store configurations
- Integration tests that need to verify store state changes

## Best Practices

1. **Use for Connected Components**: Primarily use `renderWithStore` for components that connect to Redux or use Mantine components
2. **Customize State**: Provide preloaded state that matches your test scenario
3. **Store Assertions**: Use the returned store to verify state changes in integration tests
4. **Keep Tests Focused**: Only provide the minimum state needed for each test

## Dependencies

- `@reduxjs/toolkit`: Store configuration
- `react-redux`: Provider component
- `@testing-library/react`: Base render functionality
- `@mantine/core`: UI provider and theming
- `@/features/auth/authSlice`: Default auth reducer

## Related Documentation

- [Test Setup](../README.md) - Overall testing configuration
- [Features Auth](../../features/auth/README.md) - Auth slice documentation
- [App Store](../../app/README.md) - Redux store configuration