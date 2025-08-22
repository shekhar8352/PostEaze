# Routes

Centralized routing configuration for the PostEaze React application using React Router v6.

## Purpose

This folder contains the main routing logic that combines all feature-based routes into a single routing configuration. It serves as the central hub for navigation throughout the application.

## Architecture

The routing system follows a **feature-based architecture** where each feature defines its own routes, and this folder combines them into the main application router.

### Route Organization Pattern

```
routes/
├── index.ts          # Main route configuration and combination logic
└── README.md         # This documentation
```

### Routing Flow

1. **Feature Routes**: Each feature (auth, landing, etc.) defines its own routes in `[feature]/[feature]Routes.tsx`
2. **Route Aggregation**: The main `routes/index.ts` imports and combines all feature routes
3. **Route Rendering**: The `AppRoutes` component uses React Router's `useRoutes` hook to render the combined routes
4. **App Integration**: The main `App.tsx` component renders `AppRoutes` as the primary routing component

## Key Files

### `index.ts`
The main routing configuration file that:
- Imports route definitions from all features
- Combines them into a single routes array
- Exports the `AppRoutes` component that renders the routes using `useRoutes`

## Current Route Structure

```typescript
// Current routes (expandable)
const routes = [
  ...authRoutes,      // /signup
  ...landingRoutes,   // /
  // Future: ...dashboardRoutes, ...postsRoutes, etc.
];
```

### Active Routes

| Path | Feature | Component | Description |
|------|---------|-----------|-------------|
| `/` | landing | LandingPage | Application home/landing page |
| `/signup` | auth | SignupPage | User registration page |

## Usage Patterns

### Adding New Routes

1. **Create feature routes**: Define routes in your feature's `[feature]Routes.tsx` file
2. **Export from feature**: Export routes from the feature's `index.tsx`
3. **Import in main routes**: Add the feature routes to `routes/index.ts`
4. **Combine routes**: Add to the routes array using spread operator

Example:
```typescript
// In features/dashboard/dashboardRoutes.tsx
const dashboardRoutes: RouteObject[] = [
  { path: "/dashboard", element: <DashboardPage /> },
  { path: "/dashboard/posts", element: <PostsPage /> },
];

// In routes/index.ts
import { dashboardRoutes } from "@/features/dashboard";

const AppRoutes = () => {
  const routes = [
    ...authRoutes, 
    ...landingRoutes,
    ...dashboardRoutes  // Add new feature routes
  ];
  return useRoutes(routes);
};
```

### Route Configuration

Routes use React Router v6's `RouteObject` interface:
```typescript
interface RouteObject {
  path?: string;
  element?: React.ReactNode;
  children?: RouteObject[];
  // ... other React Router properties
}
```

## Navigation Setup

The routing system is integrated with the application through:

1. **BrowserRouter**: Configured in `main.tsx` to provide routing context
2. **App Component**: Renders the `AppRoutes` component
3. **Feature Integration**: Each feature exports its routes for central aggregation

## Best Practices

### Route Organization
- **Feature-based**: Keep routes close to their related components
- **Centralized aggregation**: Combine all routes in this folder
- **Consistent naming**: Use `[feature]Routes.tsx` naming convention

### Route Definitions
- **Type safety**: Use `RouteObject[]` type for route arrays
- **Clear paths**: Use descriptive, RESTful path names
- **Lazy loading**: Consider code splitting for large features

### Future Enhancements
- **Protected routes**: Add authentication guards
- **Route groups**: Organize routes by access level or feature area
- **Dynamic routes**: Support for parameterized routes
- **Route metadata**: Add titles, breadcrumbs, or permissions

## Related Documentation

- **Features**: `../features/README.md` - Feature-based architecture overview
- **App Configuration**: `../app/README.md` - Redux store and app setup
- **Authentication**: `../features/auth/README.md` - Auth routes and components
- **Landing**: `../features/landing/README.md` - Landing page routes

## Dependencies

- **react-router-dom**: Primary routing library
- **@types/react-router-dom**: TypeScript definitions
- **Feature modules**: Individual feature route definitions