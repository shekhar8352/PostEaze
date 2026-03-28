# Features

Feature-based folders: each module owns pages, components, routes, and usually **TanStack Query** hooks/services for its API calls. Redux slices are used only where shared UI state warrants it (e.g. Instagram channel slice).

## Current features

| Feature | Role |
|---------|------|
| **auth** | Login, register, forgot password, email verification; Firebase + backend JWT; `AuthProvider`, `PublicRoute` |
| **layout** | `MainLayout` (AppShell), sidebar, header, nav constants |
| **landing** | In-app home page component (routed at `/` inside the protected shell) |
| **dashboard** | Dashboard route `/dashboard` |
| **channels** | Instagram connect/list/OAuth callback; routes under `/channels/...` |
| **analytics** | Instagram analytics UI at `/analytics` |
| **calendar** | Scheduled posts calendar at `/calendar`; uses `/api/v1/scheduled-posts` |

## Typical feature layout

Not every feature uses every file, but common pieces:

```
feature/
├── pages/
├── components/
├── services/          # API functions, often paired with *Queries.ts for TanStack Query
├── *Routes.tsx
├── index.ts / index.tsx
└── README.md          # optional
```

## Integrating a new feature

1. Add the feature folder with routes and pages.
2. Export routes from the feature barrel.
3. Append routes in `app/routes/index.tsx` inside the appropriate `ProtectedLayout` / `MainLayout` (or public) branch.
4. Add TanStack Query keys/hooks or Redux slice only if needed.

## Related documentation

- [`../app/routes/README.md`](../app/routes/README.md)
- [`../app/store/README.md`](../app/store/README.md)
- [`../services/README.md`](../services/README.md)
- [`layout/README.md`](./layout/README.md)
