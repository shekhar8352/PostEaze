# Features

The features directory implements a feature-based architecture pattern for organizing React components, state management, and routing logic. Each feature represents a distinct functional area of the PostEaze application and contains all the code necessary to implement that functionality in a self-contained, modular way.

## Architecture Overview

The feature-based architecture promotes:
- **Modularity**: Each feature is self-contained with its own components, state, and routing
- **Scalability**: New features can be added without affecting existing ones
- **Maintainability**: Related code is co-located, making it easier to understand and modify
- **Team Collaboration**: Different teams can work on different features independently

## Feature Structure Pattern

Each feature follows a consistent organizational pattern:

```
feature-name/
├── components/          # Feature-specific reusable components
│   └── ComponentName/   # Individual component folders
├── pages/              # Page-level components for routing
│   └── PageName/       # Individual page folders
├── featureApi.ts       # API service functions for the feature
├── featureSlice.ts     # Redux Toolkit slice for state management
├── featureRoutes.tsx   # React Router route definitions
├── thunks.ts          # Async Redux thunks (if needed)
├── index.tsx          # Barrel export for the feature
└── README.md          # Feature-specific documentation
```

## Current Features

### Authentication (`auth/`)
Handles user registration, login, and authentication state management.

**Key Components:**
- User signup functionality with form validation
- Redux state management for auth status
- API integration for authentication endpoints
- Route definitions for auth-related pages

**State Management:**
- `authSlice.ts`: Manages user authentication state
- `thunks.ts`: Async actions for API calls
- Integration with Redux Toolkit for predictable state updates

### Landing (`landing/`)
Provides the main entry point and homepage functionality.

**Key Components:**
- Homepage rendering and navigation
- Navbar component for site navigation
- Landing page routing configuration

## Integration Patterns

### State Management Integration
Features integrate with the global Redux store through:
```typescript
// In app/store.ts
import { authSlice } from '@/features/auth';

export const store = configureStore({
  reducer: {
    auth: authSlice.reducer,
    // Other feature reducers...
  },
});
```

### Routing Integration
Feature routes are integrated into the main application routing:
```typescript
// In routes/index.tsx
import { authRoutes } from '@/features/auth';
import { landingRoutes } from '@/features/landing';

const routes = [
  ...authRoutes,
  ...landingRoutes,
  // Other feature routes...
];
```

### API Integration
Features use the shared services layer for HTTP requests:
```typescript
// In feature API files
import { apiClient } from '@/services';

export const featureApi = {
  getData: () => apiClient.get('/feature-endpoint'),
  // Other API methods...
};
```

## Development Conventions

### Component Organization
- **Components**: Reusable UI components specific to the feature
- **Pages**: Top-level components that represent full pages/routes
- **Each component/page**: Has its own folder with index.tsx, styles, and tests

### State Management
- **Redux Toolkit**: Used for complex state that needs to be shared
- **Local State**: React useState for component-specific state
- **Async Operations**: Redux Toolkit thunks for API calls

### File Naming
- **PascalCase**: For component folders and React components
- **camelCase**: For utility functions and non-component files
- **kebab-case**: For feature folder names

### Export Pattern
Each feature uses a barrel export (index.tsx) to expose:
- Route definitions
- Key components (if needed by other features)
- Types and interfaces (if shared)

## Adding New Features

To add a new feature:

1. **Create Feature Directory**: Follow the standard structure pattern
2. **Implement Core Files**:
   - `index.tsx`: Barrel exports
   - `featureRoutes.tsx`: Route definitions
   - `featureSlice.ts`: Redux state (if needed)
   - `featureApi.ts`: API functions (if needed)
3. **Add Components and Pages**: Following the folder structure
4. **Integrate with App**: Add routes and state to main app configuration
5. **Document**: Create README.md explaining the feature

## Testing Strategy

Features should include:
- **Unit Tests**: For individual components and utility functions
- **Integration Tests**: For Redux slices and thunks
- **Route Tests**: For routing configuration
- **API Tests**: For service functions

## Dependencies

Features commonly depend on:
- **React Router**: For routing and navigation
- **Redux Toolkit**: For state management
- **Mantine**: For UI components
- **Axios**: For HTTP requests (via services layer)

## Best Practices

1. **Keep Features Independent**: Minimize dependencies between features
2. **Use Shared Services**: Leverage common services for API calls and utilities
3. **Follow Naming Conventions**: Maintain consistency across features
4. **Document Thoroughly**: Each feature should have comprehensive README
5. **Test Comprehensively**: Include tests for all major functionality
6. **Optimize Imports**: Use barrel exports to simplify imports

## Related Documentation

- [App Store](../app/README.md) - Redux store configuration and setup
- [Routes](../routes/README.md) - Application routing configuration
- [Services](../services/README.md) - Shared API services and HTTP client
- [Utils](../utils/README.md) - Shared utility functions
- [Individual Feature READMEs](./auth/README.md) - Detailed feature documentation