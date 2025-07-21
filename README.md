# PostEaze

PostEaze is a social media management platform built for influencers and marketing teams. It helps users schedule posts across platforms like Instagram, Facebook, YouTube, WhatsApp Channels, LinkedIn, and Twitter (X), with advanced analytics and team collaboration support.

## Table of Contents
- [Tech Stack](#tech-stack)
- [Project Structure](#project-structure)
- [Development Guidelines](#development-guidelines)
- [Branching Strategy](#branching-strategy)
- [Pull Request Protocol](#pull-request-protocol)
- [Code Review Process](#code-review-process)

## Tech Stack
- **Frontend:** React + TypeScript (Vite, Mantine, Formik, Redux Toolkit)
- **Backend:** Go (Gin)
- **Database:** PostgreSQL
- **Cache:** Redis
- **Auth:** JWT (access & refresh token flow)
- **Containerization:** Docker, Docker Compose

## Project Structure
The project is structured as follows:

```bash
.
├── frontend/ # React frontend
├── backend/ # Go backend
├── docker-compose.yml # Multi-container setup
├── README.md # Main documentation
```

- Frontend: [`frontend/`](./frontend)
- Backend: [`backend/`](./backend)
- Docs: This file and separate READMEs per app

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
docker-compose up --build
```

### 4. Run without Docker
```bash
cd backend
go run main.go

cd frontend
npm install
npm run dev
```
