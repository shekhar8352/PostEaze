# PostEaze Frontend

React SPA for PostEaze: auth (Firebase + JWT), dashboard, Instagram and **YouTube** channels, analytics, a **calendar** for scheduled posts (Instagram + YouTube), a **media workspace** (assets, versions, **Google Drive import**, publish to schedule), and **Studio** (Kanban phases and pieces with links to assets and scheduled posts). Uses a feature-based layout with a shared shell (sidebar/header).

## Architecture

### Tech stack

- **React 19** — UI
- **TypeScript** — types
- **Vite 6** — dev server and production build
- **Mantine 8** — components, dates, notifications; **@mantine/form** for some forms
- **Redux Toolkit** — global slices where needed (e.g. channels)
- **TanStack Query** — server state, caching, and feature-level queries
- **React Router 7** — routing (`useRoutes`, lazy routes)
- **Axios** — HTTP (shared `apiClient` in `src/services/api/client.ts` + interceptors)
- **Firebase** — client auth; tokens exchanged with the Go API
- **Formik + Yup** — forms where used (e.g. auth)
- **Chart.js / react-chartjs-2** — analytics charts
- **react-big-calendar** — calendar UI for scheduled posts
- **@dnd-kit** — drag-and-drop on the Studio board (sortable columns/cards)
- **Vitest** + Testing Library — tests (see `vite.config.ts` test block)

### Source layout

```
src/
├── app/                 # App shell: providers, routes, Redux store, theme, global styles
├── features/            # Feature modules (auth, layout, dashboard, channels, analytics, calendar, media-workspace, integrations, studio, landing)
├── services/            # api/client, interceptors, legacy axios re-export
├── utils/               # Shared helpers (grow as needed)
├── assets/              # Bundled static assets
├── test/                # Vitest setup and render helpers
└── main.tsx             # Entry (renders app/App.tsx)
```

## Key behaviors

- **Providers** (`app/App.tsx`): Redux → TanStack Query → Auth → Mantine → Router.
- **API base URL**: `import.meta.env.VITE_API_BASE_URL` or default `http://localhost:8080/api` (`services/api/client.ts`). Must match your Go server port and include `/api` if that is how the backend is mounted.
- **Auth**: JWT access token attached by interceptors (`services/api/interceptors.ts`); refresh flow coordinated with the API.
- **Protected UI**: `ProtectedLayout` + `MainLayout` wrap dashboard, analytics, channels (Instagram, Facebook, YouTube), calendar, media workspace (`/workspace`), Studio (`/studio`, `/studio/settings`), and the in-app home route.

### Google OAuth (Drive + YouTube)

Separate from Firebase login. Configure `VITE_GOOGLE_CLIENT_ID` and redirect URIs per [`backend/docs/google-cloud-setup.md`](../backend/docs/google-cloud-setup.md):

| Variable | Purpose |
|----------|---------|
| `VITE_GOOGLE_CLIENT_ID` | OAuth client ID (same GCP app as backend) |
| `VITE_GOOGLE_DRIVE_REDIRECT_URI` | Popup callback, default `{origin}/oauth/google/drive/callback` |
| `VITE_GOOGLE_YOUTUBE_REDIRECT_URI` | Popup callback, default `{origin}/oauth/google/youtube/callback` |

## Getting started

### Prerequisites

- Node.js 18+
- Backend running and reachable at the URL you configure (see above)

### Commands

```bash
npm install
npm run dev          # Vite — http://localhost:5173
npm run build        # typecheck + production build
npm run lint         # ESLint
npm run preview      # preview production build
```

### Tests

Vitest is configured in `vite.config.ts`. Run:

```bash
npx vitest           # watch mode
npx vitest run       # single run (e.g. CI)
```

## Documentation index

- [`src/README.md`](src/README.md) — source tree overview
- [`src/features/README.md`](src/features/README.md) — feature modules
- [`features/media-workspace/README.md`](src/features/media-workspace/README.md) — media assets, versions, Drive import, publish
- [`src/features/studio/README.md`](src/features/studio/README.md) — Studio Kanban (phases, pieces, assets, publish links)
- [`src/app/routes/README.md`](src/app/routes/README.md) — route composition
- [`../backend/docs/google-cloud-setup.md`](../backend/docs/google-cloud-setup.md) — Google Drive + YouTube OAuth
- [`src/app/store/README.md`](src/app/store/README.md) — Redux store
- [`src/services/README.md`](src/services/README.md) — HTTP client
- [`src/test/README.md`](src/test/README.md) — testing setup
