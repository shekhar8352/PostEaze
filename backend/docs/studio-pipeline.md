# Studio pipeline (backend)

**Studio** is a per-team Kanban-style pipeline: a **studio** has ordered **phases** (columns) and **pieces** (work items). Pieces can link to **media assets** and **scheduled posts**, support **comments** and an **activity** log, and move between phases with fractional-index **positions**. Creating a scheduled post can optionally reference a `piece_id` so publish lifecycle updates the linked piece (e.g. move toward “published”) via `scheduledpublish` hooks.

## Domain model (summary)

| Entity | Role |
|--------|------|
| **Studio** | One per team (`team_id`); `name`, `piece_label` (UI label for “Piece” / custom). |
| **Phase** | Ordered column within a studio: `name`, `slug`, `order_index`, **`kind`** (`idea`, `script`, `shoot`, `edit`, `review`, `scheduled`, `published`, `custom`), optional WIP limit, color. |
| **Piece** | Lives in one phase: `title`, `description`, `content_type`, `assignee`, `due_at`, **`position`** (fractional index string), `status` (`active` \| `archived`). |
| **Piece asset link** | Join to `media_assets` with a **role** (`script`, `raw`, `edit`, `thumbnail`, `final`, `attachment`). |
| **Piece scheduled-post link** | Join to scheduled posts with **role** (`primary`, `cross_post`, `repost`). |

Schema: migration **`009_studio_pipeline`**. Entities: **`entities/studio.go`**, **`entities/phase.go`**, **`entities/piece.go`**. Repositories: **`entities/repositories/`**.

### Utilities

**[`services/studio/`](../services/studio/README.md)** — fractional indexing for ordering pieces without rewriting whole lists.

### Integration with scheduling

- **`models/v1/scheduled_post.go`** — scheduled posts may set **`piece_id`** to associate publish with a Studio piece.
- **`scheduledpublish/`** — after Instagram publish attempts, lifecycle can move linked pieces (see `hooks.go`, `instagram.go`).

## HTTP API (`/api/v1`, JWT unless noted)

All Studio routes use **`AuthMiddleware`**.

### Studios (`/studios`)

| Method | Path | Purpose |
|--------|------|---------|
| GET | `/studios/default` | Ensure the current user’s team has a studio; create default phases if new. |
| PUT | `/studios/:id` | Update studio `name`, `piece_label`. |
| GET | `/studios/:id/board` | **Full board**: studio + phases + active pieces (for Kanban UI). |
| GET | `/studios/:id/phases` | List phases. |
| POST | `/studios/:id/phases` | Create phase. |
| POST | `/studios/:id/phases/reorder` | Reorder phases (`phase_ids`). |
| POST | `/studios/:id/pieces` | Create piece in the studio (body includes `phase_id`). |

### Phases (`/phases`)

| Method | Path | Purpose |
|--------|------|---------|
| PUT | `/phases/:phaseId` | Update phase metadata. |
| DELETE | `/phases/:phaseId` | Delete phase (typically blocked if active pieces remain). |

### Pieces (`/pieces`)

| Method | Path | Purpose |
|--------|------|---------|
| GET | `/pieces/:pieceId` | **Piece detail**: piece + assets + scheduled_post links + activities + comments. |
| PUT | `/pieces/:pieceId` | Update fields. |
| DELETE | `/pieces/:pieceId` | Delete piece. |
| POST | `/pieces/:pieceId/move` | Move to another phase with optional **`before_piece_id`** / **`after_piece_id`** (fractional positioning). |
| PUT | `/pieces/:pieceId/status` | `active` or `archived`. |
| POST | `/pieces/:pieceId/assets` | Link media asset. |
| DELETE | `/pieces/:pieceId/assets/:assetId` | Unlink asset (optional `?role=`). |
| POST | `/pieces/:pieceId/scheduled-posts` | Link scheduled post. |
| DELETE | `/pieces/:pieceId/scheduled-posts/:scheduledPostId` | Unlink. |
| GET | `/pieces/:pieceId/activities` | Activity feed (`?limit=`). |
| GET | `/pieces/:pieceId/comments` | List comments. |
| POST | `/pieces/:pieceId/comments` | Add comment (optional `parent_id`). |
| DELETE | `/pieces/:pieceId/comments/:commentId` | Delete comment. |

Handlers: **`api/v1/studio.go`**. Router registration: **`api/router.go`** (`addV1StudioRoutes`, phases group, pieces group). Swagger tag: **Studio**.

## Business logic

**`business/v1/studio.go`** — orchestration for studios, phases, pieces, ACLs (team membership), moves, comments, asset/scheduled-post links.

## Demo data

**`scripts/seed_studio_demo/main.go`** — seeds a user’s default studio with sample pieces across phases (see script header for flags).

## Related docs

- [`../README.md`](../README.md) — backend overview  
- [`../business/v1/README.md`](../business/v1/README.md) — business layer  
- [`../../frontend/src/features/studio/README.md`](../../frontend/src/features/studio/README.md) — React UI  
