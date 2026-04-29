# Routes

Central routing for the SPA: **`index.tsx`** composes feature `RouteObject[]` arrays with **`useRoutes`** (React Router 7).

## Files

| File | Role |
|------|------|
| `index.tsx` | Merges auth + protected tree (dashboard, analytics, landing home, channels, calendar, media workspace, studio) + catch-all 404 |
| `ProtectedRoute.tsx` | `ProtectedLayout` — requires authenticated user |
| `PublicRoute.tsx` | Redirects authenticated users away from public auth pages |
| `NotFound.tsx` | 404 |
| `LoadingFallback.tsx` | Suspense fallback for lazy routes |

## Route tree (summary)

**Public (no app shell)**

| Path | Feature |
|------|---------|
| `/login` | auth |
| `/register` | auth |
| `/forgot-password` | auth |
| `/email-verify` | auth |

**Protected** — wrapped in `ProtectedLayout` → `MainLayout` (sidebar/header)

| Path | Feature |
|------|---------|
| `/` | landing (in-app home) |
| `/dashboard` | dashboard |
| `/analytics` | analytics |
| `/calendar` | calendar (scheduled posts) |
| `/workspace` | media-workspace (asset grid / upload) |
| `/workspace/:assetId` | media-workspace (detail, versions, publish) |
| `/studio` | studio (Kanban board) |
| `/studio/settings` | studio (phases + labels) |
| `/channels/instagram` | channels |

**Special**

| Path | Notes |
|------|--------|
| `/auth/instagram/callback` | OAuth redirect (still under protected branch in `index.tsx`) |

**Catch-all**

| Path | Element |
|------|---------|
| `*` | `NotFound` |

## Adding routes

1. Define `RouteObject[]` in the feature’s `*Routes.tsx` with `lazy` + `Suspense` as needed.
2. Export from the feature barrel.
3. Import in `index.tsx` and add to the correct parent (`ProtectedLayout` children vs public `authRoutes`).

## Related documentation

- [`../../features/README.md`](../../features/README.md)
- [`../../features/media-workspace/README.md`](../../features/media-workspace/README.md)
- [`../../features/studio/README.md`](../../features/studio/README.md)
- [`../store/README.md`](../store/README.md)
