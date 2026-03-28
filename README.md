# PostEaze

PostEaze is a social media management platform built for influencers and marketing teams. It supports connected channels (including Instagram), post listing, **scheduled posts** with a calendar-oriented API, Instagram analytics, teams, and background jobs for sync—expandable toward additional networks over time.

## Table of Contents
- [Tech Stack](#tech-stack)
- [Project Structure](#project-structure)
- [Development Guidelines](#development-guidelines)
- [Branching Strategy](#branching-strategy)
- [Pull Request Protocol](#pull-request-protocol)
- [Code Review Process](#code-review-process)
- [Local Setup](#local-setup)

## Tech Stack
- **Frontend:** React + TypeScript (Vite, Mantine, Formik, Redux Toolkit)
- **Backend:** Go (Gin)
- **Database:** PostgreSQL
- **Cache:** Redis
- **Task Queue:** Asynq (Redis-based)
- **Auth:** JWT (access & refresh token flow)
- **Containerization:** Docker, Docker Compose

## Project Structure
The project is structured as follows:

```bash
.
├── frontend/          # React + Vite SPA
├── backend/           # Go API, migrations, Asynq worker
├── init-db/           # Optional Postgres init scripts (mounted in full stack compose)
├── docker-compose.yml # Full stack (backend, frontend, Postgres, Redis, worker)
├── docker-compose.local.yml # Postgres + Redis + worker only (run API/FE on host)
└── README.md
```

- Frontend: [`frontend/`](./frontend)
- Backend: [`backend/`](./backend)
- Database init: [`init-db/`](./init-db)
- Docs: This file plus READMEs under `frontend/` and `backend/`

## Development Guidelines
- Use **feature branches** for new features.
- Follow **Conventional Commits**:
  - `feat:` for new features
  - `fix:` for bug fixes
  - `chore:` for non-functional changes
  - `refactor:` for code improvement
- Run formatters/linters before pushing if configured.

## Branching Strategy
We follow the **GitFlow** strategy:
- `main`: Production-ready code
- `dev`: Ongoing development branch
- `feature/xyz`: Features branch from `dev`
- `hotfix/xyz`: Urgent fixes from `main`

## Pull Request Protocol
- Create your branch from `dev` (e.g. `feature/signup-page`)
- Rebase frequently to keep your branch up-to-date
- Ensure all tests (if any) pass before opening a PR
- Use descriptive PR titles and body
- Link related issues with `Closes #issue`
- Assign at least one reviewer

## Code Review Process
- Reviews focus on functionality, readability, and maintainability
- Suggestions must be acknowledged or resolved before merge
- Merge via **Squash and Merge** strategy (unless otherwise discussed)

## Local Setup

### 1. Prerequisites
- Node.js (>=18)
- Go (>=1.20)
- PostgreSQL (or use Docker)
- Redis (optional, recommended)
- Docker + Docker Compose (for containerized setup)

### 2. Clone the repository
```bash
git clone https://github.com/your-username/PostEaze.git
cd PostEaze
```

### 3. Run with Docker
```bash
docker-compose -f docker-compose.local.yml up --build
```

### 4. Run without Docker

From the repo root, use separate terminals for **Postgres + Redis** (e.g. via `docker-compose.local.yml`), the **API**, the **Asynq worker**, and the **frontend**.

```bash
# Terminal 1 — dependencies only (example)
docker compose -f docker-compose.local.yml up -d

# Terminal 2 — API (from backend/, after env + migrations — see backend README)
cd backend
go run .

# Terminal 3 — Asynq worker
cd backend
go run ./cmd/worker

# Terminal 4 — frontend
cd frontend
npm install
npm run dev
```

Point the SPA at your API: set `VITE_API_BASE_URL` (defaults to `http://localhost:8080/api` in code—match your `-port` / deployment). Vite dev server runs at **http://localhost:5173**; CORS in the API allows that origin.
