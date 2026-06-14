# Media workspace

Browse, upload, and version **media assets** (photos/videos), **import from Google Drive**, compare versions, and **publish** the current version into a scheduled post. Uses the Go API under `/api/v1/media-assets` and TanStack Query.

## Routes

| Path | Page | Description |
|------|------|-------------|
| `/workspace` | `MediaWorkspace` | Grid/list of assets; upload; **Import from Drive** |
| `/workspace/:assetId` | `MediaDetail` | Single asset: versions, compare, set current, **Drive revisions**, publish dialog |

Routes are defined in `mediaWorkspaceRoutes.tsx` and merged in `app/routes/index.tsx` inside `MainLayout`. Sidebar label **Workspace** → `/workspace` (`features/layout/constants.ts`).

## Layout

```
media-workspace/
├── pages/
│   ├── MediaWorkspace.tsx   # List + uploader + Drive import entry
│   └── MediaDetail.tsx      # Detail, timeline, publish, Drive revisions
├── components/
│   ├── AssetCard.tsx
│   ├── MediaUploader.tsx
│   ├── DriveImportModal.tsx     # Browse Drive folders/files, import as asset
│   ├── DriveRevisionsPanel.tsx  # List/import Drive revisions as new versions
│   ├── VersionTimeline.tsx
│   ├── VersionCompare.tsx
│   ├── PublishDialog.tsx
│   └── AddVersionModal.tsx
├── services/
│   └── mediaApi.ts          # Axios calls to /v1/media-assets…
├── hooks/
│   └── useMediaQueries.ts   # Query keys + useQuery / useMutation hooks
├── types.ts                 # DTOs: storage_provider, stream_url, drive_file_id, helpers
├── mediaWorkspaceRoutes.tsx
└── index.ts                 # Exports `mediaWorkspaceRoutes`
```

## Data layer

- **`services/mediaApi.ts`** — `list`, `get`, `create` (multipart), `update`, `remove`, `addVersion`, `deleteVersion`, `setCurrentVersion`, `publish`. Longer timeouts on multipart requests.
- **`hooks/useMediaQueries.ts`** — `mediaKeys` for cache; hooks such as `useMediaAssets`, `useMediaAsset`, `useCreateMediaAsset`, and mutations that invalidate `mediaKeys.all` and show Mantine notifications.
- **Drive import** — `features/integrations/services/googleDriveApi.ts` for connect/status/browse; import endpoints `POST /v1/media-assets/import/google-drive` and `POST /v1/media-assets/:id/versions/import-drive-revision`.

## Storage providers

Versions may be **blob-backed** (Vercel Blob, default) or **google_drive** (large videos stay in Drive). Drive-backed versions expose a signed **`stream_url`** for preview and for Meta/YouTube publish paths. Helpers in `types.ts`: `versionMediaUrl()`, `currentVersionForAsset()`.

## Backend contract

Paths are relative to `apiClient` base URL (typically `…/api`), so calls use `/v1/media-assets`, etc. Public streaming (no JWT): `GET /v1/media/stream/:versionId?sig=…&exp=…`. See the backend [`api/v1/README.md`](../../../../backend/api/v1/README.md) for full HTTP semantics.

## Related documentation

- [`../integrations/`](../integrations/) — Google Drive OAuth and connect UI
- [`../README.md`](../README.md) — feature index
- [`../../app/routes/README.md`](../../app/routes/README.md) — route composition
- [`../../services/README.md`](../../services/README.md) — `apiClient` and env
- [`../../../../backend/docs/google-cloud-setup.md`](../../../../backend/docs/google-cloud-setup.md) — GCP OAuth setup
