# App Configuration

This folder contains the core Redux store configuration and typed hooks for the PostEaze React application.

## Contents

- **store.ts**: Redux store configuration using Redux Toolkit
- **hooks.ts**: Typed Redux hooks for type-safe state management

## Redux Store Configuration

The store is configured using Redux Toolkit's `configureStore` function, which provides good defaults for Redux setup including:

- Redux DevTools integration
- Immutability checks
- Serializable state checks
- Thunk middleware for async actions

### Store Structure

```typescript
export const store = configureStore({
  reducer: {
    auth: authReducer,
  },
});
```

The store currently includes:
- **auth**: Authentication state management (login, user data, tokens)

### Type Definitions

The store exports TypeScript types for type-safe Redux usage:

```typescript
export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;
```

## Typed Redux Hooks

The `hooks.ts` file provides typed versions of the standard Redux hooks to ensure type safety throughout the application:

### useAppDispatch

A typed version of `useDispatch` that knows about our store's dispatch type:

```typescript
const dispatch = useAppDispatch();
// dispatch is now typed as AppDispatch
```

### useAppSelector

A typed version of `useSelector` that provides autocomplete and type checking for the state:

```typescript
const user = useAppSelector((state) => state.auth.user);
// state is typed as RootState, providing full autocomplete
```

## Usage Patterns

### In Components

```typescript
import { useAppSelector, useAppDispatch } from '@/app/hooks';
import { loginUser } from '@/features/auth/authSlice';

function MyComponent() {
  const dispatch = useAppDispatch();
  const { user, isLoading } = useAppSelector((state) => state.auth);
  
  const handleLogin = () => {
    dispatch(loginUser({ email, password }));
  };
  
  return (
    // Component JSX
  );
}
```

### Adding New Reducers

When adding new features with Redux state:

1. Import the reducer in `store.ts`
2. Add it to the reducer configuration
3. The types will automatically update to include the new state

```typescript
import newFeatureReducer from '@/features/newFeature/newFeatureSlice';

export const store = configureStore({
  reducer: {
    auth: authReducer,
    newFeature: newFeatureReducer, // Add new reducer here
  },
});
```

## Architecture Benefits

- **Type Safety**: Full TypeScript support with autocomplete and error checking
- **Developer Experience**: Redux DevTools integration for debugging
- **Performance**: Built-in optimizations from Redux Toolkit
- **Maintainability**: Centralized state management with clear patterns

## Related Documentation

- [Features Documentation](../features/README.md) - Feature-based architecture and state slices
- [Services Documentation](../services/README.md) - API integration with Redux