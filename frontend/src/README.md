# Frontend source

TypeScript/React source for the PostEaze web app.

## Stack highlights

- **React 19** + **Vite 6**
- **Mantine 8** (core, dates, notifications, form)
- **TanStack Query** for server state; **Redux Toolkit** for selected global slices
- **React Router 7** with lazy-loaded feature routes
- **Axios** via `services/api/client.ts` and `services/api/interceptors.ts`

## Entry point

- **`main.tsx`** mounts `App` from **`app/App.tsx`**.
- **`app/App.tsx`** wires providers (Redux, TanStack Query, Firebase auth context, Mantine, notifications), imports interceptors, and renders **`app/routes`** inside `BrowserRouter`.

## Directory map

| Path | Purpose |
|------|---------|
| `app/` | Providers, **`routes/`** (including `ProtectedLayout`, `MainLayout` usage), **`store/`** Redux, **`theme/`** design tokens, shell CSS |
| `features/` | Feature modules: **auth**, **layout**, **landing**, **dashboard**, **channels**, **analytics**, **calendar** (scheduled posts), **media-workspace**, **studio**, … |
| `services/` | Shared API client, interceptors; see README there |
| `utils/` | Cross-feature helpers |
| `assets/` | Images/SVG imported through Vite |
| `test/` | Vitest setup and `renderWithStore` helpers |

## Routing model

Feature route arrays are merged in **`app/routes/index.tsx`**. Authenticated app areas use **`ProtectedLayout`** then **`MainLayout`** (sidebar + header from `features/layout`).

## Related docs

- [`app/routes/README.md`](./app/routes/README.md)
- [`app/store/README.md`](./app/store/README.md)
- [`features/README.md`](./features/README.md)
- [`features/media-workspace/README.md`](./features/media-workspace/README.md)
- [`features/studio/README.md`](./features/studio/README.md)
- [`services/README.md`](./services/README.md)
- [`test/README.md`](./test/README.md)
