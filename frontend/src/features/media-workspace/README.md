# Media workspace

Browse, upload, and version **media assets** (photos/videos), compare versions, and **publish** the current version into a scheduled post. Uses the Go API under `/api/v1/media-assets` and TanStack Query.

## Routes

| Path | Page | Description |
|------|------|-------------|
| `/workspace` | `MediaWorkspace` | Grid/list of assets; create upload flow |
| `/workspace/:assetId` | `MediaDetail` | Single asset: versions, compare, set current, publish dialog |

Routes are defined in `mediaWorkspaceRoutes.tsx` and merged in `app/routes/index.tsx` inside `MainLayout`. Sidebar label **Workspace** → `/workspace` (`features/layout/constants.ts`).

## Layout

```
media-workspace/
├── pages/
│   ├── MediaWorkspace.tsx   # List + uploader entry
│   └── MediaDetail.tsx      # Detail, timeline, publish
├── components/
│   ├── AssetCard.tsx
│   ├── MediaUploader.tsx
│   ├── VersionTimeline.tsx
│   ├── VersionCompare.tsx
│   ├── PublishDialog.tsx
│   └── AddVersionModal.tsx
├── services/
│   └── mediaApi.ts          # Axios calls to /v1/media-assets…
├── hooks/
│   └── useMediaQueries.ts   # Query keys + useQuery / useMutation hooks
├── types.ts                 # DTOs aligned with API envelopes
├── mediaWorkspaceRoutes.tsx
└── index.ts                 # Exports `mediaWorkspaceRoutes`
```

## Data layer

- **`services/mediaApi.ts`** — `list`, `get`, `create` (multipart), `update`, `remove`, `addVersion`, `deleteVersion`, `setCurrentVersion`, `publish`. Longer timeouts on multipart requests.
- **`hooks/useMediaQueries.ts`** — `mediaKeys` for cache; hooks such as `useMediaAssets`, `useMediaAsset`, `useCreateMediaAsset`, and mutations that invalidate `mediaKeys.all` and show Mantine notifications.

## Backend contract

Paths are relative to `apiClient` base URL (typically `…/api`), so calls use `/v1/media-assets`, etc. See the backend [`api/v1/README.md`](../../../../backend/api/v1/README.md) for full HTTP semantics.

## Related documentation

- [`../README.md`](../README.md) — feature index
- [`../../app/routes/README.md`](../../app/routes/README.md) — route composition
- [`../../services/README.md`](../../services/README.md) — `apiClient` and env
