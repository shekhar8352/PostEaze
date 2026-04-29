# Studio feature

The **Studio** is a Kanban-style **content pipeline** for teams: customizable **phases** (columns) and **pieces** (cards) that track content from idea through publish. It integrates with the **media workspace** (linked assets on each piece) and **scheduled posts** (link or create publications from a piece).

## User-facing behaviour

| Area | Description |
|------|-------------|
| **`/studio`** | Board view: phases as columns, pieces as draggable cards (`@dnd-kit`). Open a piece in a **drawer** for detail. |
| **`/studio/settings`** | Rename studio and **piece label** (singular term shown in UI, e.g. “Piece”, “Drop”); CRUD/reorder **phases** (add, edit colours, WIP limits, delete when empty). |
| **Dashboard** | **`DashboardStudioPipeline`** summarises non-terminal phases, due-soon pieces, and links into Studio. |

### Piece drawer (tabs)

When a piece is selected, **`PieceDetailDrawer`** loads **`PieceDetail`** (`usePieceDetail`):

1. **Overview** — Edit title, description, content type, assignee, due date (`PieceOverviewTab`).
2. **Assets** — Link/unlink workspace media assets with roles (`PieceAssetsTab`).
3. **Publish** — Link scheduled posts or flows tied to publishing (`PiecePublishTab`).

Below the tabs: **comments** (`PieceCommentsPanel`) and **activity** (`PieceActivityFeed`).

## Routes

Defined in **`studioRoutes.tsx`**:

| Path | Page |
|------|------|
| `/studio` | `StudioBoard` |
| `/studio/settings` | `StudioSettings` |

Merged in **`app/routes/index.tsx`** under `ProtectedLayout` → `MainLayout`.

## API client

**`services/studioApi.ts`** wraps `/api/v1`:

- Studios: `ensureDefault`, `update`, `getBoard`
- Phases: `listPhases`, `createPhase`, `reorderPhases`, `updatePhase`, `deletePhase`
- Pieces: `createPiece`, `getPiece`, `updatePiece`, `movePiece`, `setPieceStatus`, `deletePiece`
- Links: `linkAsset`, `unlinkAsset`, `linkScheduledPost`, `unlinkScheduledPost`
- Social: `listActivities`, `listComments`, `createComment`, `deleteComment`

## TanStack Query

**`hooks/useStudioQueries.ts`** exports keyed queries/mutations (`studioKeys`), invalidation after mutations, and hooks such as `useDefaultStudio`, `useStudioBoard`, `usePhases`, `useCreatePiece`, `usePieceDetail`, etc.

## Types

**`types.ts`** mirrors **`backend/models/v1/studio.go`** contracts: `Studio`, `Phase`, `Piece`, `PieceDetail`, and request payloads (`MovePiecePayload`, asset roles, scheduled-post roles, etc.).

## Key components (non-exhaustive)

| Component | Role |
|-----------|------|
| `StudioBoard.tsx` | Main Kanban surface |
| `PhaseColumn.tsx`, `PhaseRow.tsx` | Column chrome; create piece in phase |
| `PieceCard.tsx` | Card summary in the board |
| `PieceDetailDrawer.tsx` | Tabs + comments + activity |
| `CreatePhaseModal.tsx`, `LinkToPieceModal.tsx` | Modals |

## Related documentation

- Backend pipeline and REST map: [`../../../../backend/docs/studio-pipeline.md`](../../../../backend/docs/studio-pipeline.md)
- Media workspace (assets): [`../media-workspace/README.md`](../media-workspace/README.md)
- Feature index: [`../README.md`](../README.md)
- Routes composition: [`../../app/routes/README.md`](../../app/routes/README.md)
