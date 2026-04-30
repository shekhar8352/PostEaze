# PostEaze

PostEaze is a social media management platform for influencers and marketing teams. It connects **Instagram** and **Facebook** (Meta), supports **scheduled posts** with a calendar-oriented API and UI, **channel analytics**, **teams**, **Firebase-backed authentication** (JWT session with the Go API), a **media workspace** (assets, versions, publish-to-schedule), and **background jobs** (Asynq on Redis)—designed to grow with additional networks over time.

## Table of contents

- [Features](#features)
- [Tech stack](#tech-stack)
- [Repository layout](#repository-layout)
- [Documentation](#documentation)
- [Prerequisites](#prerequisites)
- [Quick start](#quick-start)
  - [Full stack (Docker)](#full-stack-docker)
  - [Hybrid local dev (Postgres + Redis in Docker)](#hybrid-local-dev-postgres--redis-in-docker)
- [Environment configuration](#environment-configuration)
- [API & tooling](#api--tooling)
- [Development workflow](#development-workflow)
- [Contributing](#contributing)

## Features

- **Auth** — Firebase client sign-in; API exchanges ID tokens for JWT access/refresh.
- **Channels** — Instagram and Facebook via Meta OAuth; webhooks for Instagram.
- **Posts & scheduling** — List posts, create and manage scheduled posts; calendar views on the frontend.
- **Analytics** — Instagram channel analytics (profile, posts, audience, stories, comparisons, etc.) with access controls.
- **Media workspace** — Uploads, versioned assets, publish flow into scheduled posts.
- **Studio** — Per-team content pipeline: customizable **phases**, **pieces** (cards) with media links, scheduled-post links, comments, and activity; Kanban board UI and dashboard pipeline widget; see [`frontend/src/features/studio/README.md`](./frontend/src/features/studio/README.md) and [`backend/docs/studio-pipeline.md`](./backend/docs/studio-pipeline.md).
- **Teams** — Multi-user team model on the API.
- **Jobs** — Asynq workers for sync and related background work.

## Tech stack

| Layer | Technologies |
|--------|----------------|
| **Frontend** | React 19, TypeScript, Vite 6, Mantine 8, Redux Toolkit, TanStack Query, React Router 7, Axios, Firebase (client), Formik + Yup, Chart.js, react-big-calendar, Vitest |
| **Backend** | Go 1.23+, Gin, PostgreSQL (`lib/pq`), Redis, Asynq, Firebase Admin, JWT, Swagger (swag) |
| **Infra** | Docker & Docker Compose, Caddy (reverse proxy in full compose), optional Vercel Blob for media (`BLOB_READ_WRITE_TOKEN`) |

## Repository layout

```text
.
├── frontend/                 # React + Vite SPA
├── backend/                  # Go API, migrations, Asynq client, cmd/worker
├── init-db/                  # Postgres init scripts (full-stack compose)
├── Caddyfile                 # Reverse proxy rules (full-stack compose)
├── docker-compose.yml        # Full stack: backend, frontend, Postgres, Redis, worker, Caddy
├── docker-compose.local.yml  # Postgres + Redis + worker only (run API/UI on host)
└── README.md                 # This file
```

Sub-project docs: [`frontend/README.md`](./frontend/README.md), [`backend/README.md`](./backend/README.md).

## Documentation

| Topic | Location |
|--------|----------|
| Frontend architecture, scripts, env | [`frontend/README.md`](./frontend/README.md) |
| Backend architecture, Swagger, migrations, env | [`backend/README.md`](./backend/README.md) |
| Firebase + API auth notes | [`backend/docs/firebase-authentication.md`](./backend/docs/firebase-authentication.md) |
| Studio pipeline (pieces, phases, API) | [`backend/docs/studio-pipeline.md`](./backend/docs/studio-pipeline.md); UI: [`frontend/src/features/studio/README.md`](./frontend/src/features/studio/README.md) |

## Prerequisites

- **Node.js** 18+ (frontend)
- **Go** 1.23+ (backend; see [`backend/go.mod`](./backend/go.mod))
- **Docker** + Docker Compose (recommended for Postgres, Redis, and optional full stack)
- **PostgreSQL** and **Redis** if you run services without Docker

## Quick start

### Full stack (Docker)

Uses root **`.env`** for services (database credentials, secrets, etc.). Backend image listens on **8000** internally; Caddy routes `/api/*` to the backend per [`Caddyfile`](./Caddyfile).

```bash
git clone https://github.com/shekhar8352/PostEaze.git
cd PostEaze
# Create and fill .env (see backend README for variables)
docker compose up --build
```

Adjust published ports and hostnames as needed for your machine; see `docker-compose.yml` and `Caddyfile`.

### Hybrid local dev (Postgres + Redis in Docker)

Run databases (and optionally the worker image) in containers, then run the API and SPA on the host for fast iteration.

1. Create a **`.env`** in the repo root with at least `POSTGRES_USER`, `POSTGRES_PASSWORD`, and `POSTGRES_DB` for [`docker-compose.local.yml`](./docker-compose.local.yml).

```bash
docker compose -f docker-compose.local.yml up -d
```

2. **Backend** — From `backend/`, configure `resources/configs/dev`, `.env`/substitution vars, run migrations, then start the API and worker ([`backend/README.md`](./backend/README.md)):

```bash
cd backend
go mod download
# Set DATABASE_URL or config as documented under backend/README.md
make migrate-up   # requires golang-migrate CLI; see backend/migrations/README.md
go run .          # default HTTP port :8080 (-port to override)

# Separate terminal — Asynq worker
go run ./cmd/worker/main.go
```

3. **Frontend** — Point the SPA at your API:

```bash
cd frontend
npm install
npm run dev   # http://localhost:5173
```

Set **`VITE_API_BASE_URL`** so it matches where the Gin server exposes the API (defaults in code to `http://localhost:8080/api`). The Vite dev server uses **5173**; backend CORS should allow that origin in dev.

Ports reference:

| Compose file | Postgres (host) | Redis (host) |
|--------------|-------------------|---------------|
| `docker-compose.yml` | 5440 → 5432 | 6390 → 6379 |
| `docker-compose.local.yml` | 5432 | 6379 |

## Environment configuration

- **Root `.env`** — Used by `docker-compose.yml` / `docker-compose.local.yml` for Postgres and services that load `env_file`.
- **Backend** — Dev configs live under [`backend/resources/configs/dev/`](./backend/resources/configs/dev/); env-specific secrets are documented in [`backend/README.md`](./backend/README.md) (e.g. `ENV`, `API_HOST`, Firebase, `ENCRYPTION_KEY`, Meta/Instagram, `REDIS_ADDR`).
- **Frontend** — `VITE_API_BASE_URL` (see [`frontend/README.md`](./frontend/README.md)).

## API & tooling

- **REST base** — Routes are under `/api` (e.g. `/api/v1/...`; health checks under `/api/health`).
- **Swagger UI** — Available when the API is running, typically at `/api/swagger/index.html` (host/scheme depend on deployment; see [`backend/README.md`](./backend/README.md)).

## Development workflow

### Branching (Git Flow–style)

- **`main`** — Production-ready releases
- **`dev`** — Integration branch for ongoing work
- **`feature/*`** — Branched from `dev`
- **`hotfix/*`** — Urgent fixes from `main`

### Commits

Use **[Conventional Commits](https://www.conventionalcommits.org/)**:

- `feat:` — New features
- `fix:` — Bug fixes
- `chore:` — Tooling, deps, non-feature changes
- `refactor:` — Behavior-preserving code changes

### Pull requests

- Branch from **`dev`** unless it is a hotfix from **`main`**
- Rebase or merge often to stay current
- Run tests and linters before opening a PR
- Use clear titles and descriptions; link issues with `Closes #123` when applicable
- Prefer **Squash and merge** unless the team agrees otherwise

## Contributing

1. Open an issue or pick an existing one when possible  
2. Work on a **`feature/`** branch from **`dev`**  
3. Follow the frontend/backend READMEs for formatters, tests (`npm run lint`, `npx vitest run`, Go tests where present)  
4. Open a PR into **`dev`** and address review feedback before merge  

---

For deeper detail on any layer, start with **[`frontend/README.md`](./frontend/README.md)** and **[`backend/README.md`](./backend/README.md)**.
