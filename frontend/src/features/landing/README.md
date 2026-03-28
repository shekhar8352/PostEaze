# Landing feature

Home view for **authenticated** users at **`/`**, rendered inside **`MainLayout`** (see `app/routes/index.tsx`). Unauthenticated visitors are directed to auth routes by `ProtectedLayout`, not this page.

## Architecture

The landing feature follows the standard feature-based architecture pattern used throughout the application, with clear separation between components, pages, and routing configuration.

## Contents

### Core Files
- **`index.tsx`** - Feature barrel export that exposes the landing routes
- **`landingRoutes.tsx`** - React Router configuration for landing-related routes

### Components
- **`components/Navbar/`** - Navigation bar component used across landing pages
  - Simple navbar component providing basic navigation structure

### Pages
- **`pages/Landing Page/`** - Main landing page implementation
  - Root landing page component that serves as the application homepage
  - Integrates the Navbar component for consistent navigation

## Routing

```typescript
const landingRoutes: RouteObject[] = [
  { path: "/", element: <LandingPage /> },
];
```

This **`/`** route is registered **inside** the protected + `MainLayout` branch, so it is the post-login home, not a marketing page.

## Component Hierarchy

```
LandingPage
└── Navbar
```

The page composes the shared **Navbar** and any home content for signed-in users.

## Usage

The landing feature is automatically integrated into the application's routing system through the main routes configuration. The feature exports its routes via the barrel export pattern:

```typescript
// From index.tsx
export { default as landingRoutes } from "./landingRoutes";
```

## Related documentation

- [Features overview](../README.md)
- [Routes](../../app/routes/README.md)
- [Layout / shell](../layout/README.md)