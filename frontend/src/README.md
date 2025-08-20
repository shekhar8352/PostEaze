# Frontend Source Code

The main source code directory for the PostEaze React application, containing all TypeScript/React components, services, utilities, and application logic.

## Architecture Overview

This React application is built with modern frontend technologies and follows a feature-based architecture pattern:

- **React 19** with TypeScript for component development
- **Redux Toolkit** for state management
- **React Router v7** for client-side routing
- **Mantine UI** for component library and theming
- **Vite** for build tooling and development server
- **Vitest** for testing framework

## Application Entry Points

### main.tsx
The application entry point that sets up the React root and configures all providers:

```typescript
// Provider hierarchy (outer to inner):
// 1. React.StrictMode - Development mode checks
// 2. MantineProvider - UI theme and components
// 3. Redux Provider - Global state management
// 4. BrowserRouter - Client-side routing
// 5. App component - Main application
```

**Key Configurations:**
- **Mantine Theme**: Custom theme with Indigo primary color and Inter font family
- **Redux Store**: Connected to the global application store
- **Router**: BrowserRouter for SPA navigation
- **Styling**: Imports Mantine CSS and custom App.css

### App.tsx
A minimal wrapper component that renders the main routing component:

```typescript
const App = () => <AppRoutes />;
```

The App component delegates all routing logic to the `AppRoutes` component from the routes folder, keeping the main App component clean and focused.

## Directory Structure

### Core Directories

- **`app/`** - Redux store configuration and app-level hooks
  - `store.ts` - Redux Toolkit store setup
  - `hooks.ts` - Typed Redux hooks (useAppDispatch, useAppSelector)

- **`features/`** - Feature-based architecture with self-contained modules
  - Each feature contains its own components, Redux slices, API calls, and routes
  - Currently includes: `auth/`, `landing/`

- **`routes/`** - Centralized routing configuration
  - `index.ts` - Main routing component that combines all feature routes

- **`services/`** - External service integrations and API clients
  - `axios.ts` - HTTP client configuration for API calls

### Supporting Directories

- **`utils/`** - Shared utility functions and helpers
- **`assets/`** - Static assets like images, icons, and media files
- **`test/`** - Testing utilities and configuration
  - `setup.ts` - Test environment setup
  - `utils/` - Test helper functions and mocks

## Routing Approach

The application uses a **centralized routing system** with **feature-based route organization**:

1. **Main Router** (`routes/index.ts`): Combines routes from all features
2. **Feature Routes**: Each feature exports its own route configuration
3. **Route Composition**: Uses `useRoutes` hook for declarative routing

```typescript
// Example routing pattern:
const AppRoutes = () => {
  const routes = [...authRoutes, ...landingRoutes];
  return useRoutes(routes);
};
```

This approach allows features to be self-contained while maintaining a single source of truth for all application routes.

## State Management

The application uses **Redux Toolkit** for state management with the following patterns:

- **Feature-based slices**: Each feature manages its own state slice
- **RTK Query**: For API data fetching and caching (configured in services)
- **Typed hooks**: Custom hooks provide type safety for Redux operations

## Development Workflow

### Key Files to Know

- **`main.tsx`** - Application bootstrap and provider setup
- **`App.tsx`** - Main application component (routing wrapper)
- **`vite-env.d.ts`** - TypeScript environment declarations for Vite

### Adding New Features

1. Create feature folder in `features/`
2. Implement feature components, Redux slice, and routes
3. Export routes from feature and add to main router
4. Add any new services to `services/` if needed

### Styling Approach

- **Mantine Components**: Primary UI component library
- **CSS Modules**: For component-specific styling
- **Global Styles**: Defined in `App.css`
- **Theme Configuration**: Centralized in `main.tsx`

## Related Documentation

- [`app/`](./app/README.md) - Redux store and app configuration
- [`features/`](./features/README.md) - Feature-based architecture details
- [`routes/`](./routes/README.md) - Routing configuration and patterns
- [`services/`](./services/README.md) - API integration and HTTP client setup
- [`test/`](./test/README.md) - Testing setup and utilities