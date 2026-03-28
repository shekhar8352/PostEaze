# Services

HTTP layer for the SPA: one **Axios** instance for the Go API, plus optional legacy entry points.

## Files

| File | Role |
|------|------|
| `api/client.ts` | **`apiClient`** — `baseURL` from `import.meta.env.VITE_API_BASE_URL` or **`http://localhost:8080/api`**, JSON headers, `withCredentials: true`, 10s timeout |
| `api/interceptors.ts` | Side effects: attach access token, handle 401 / refresh; **imported from `app/App.tsx`** so they run at startup |
| `api/types.ts` | Shared API-related types |
| `axios.ts` | Legacy standalone instance with a fixed base URL; prefer `apiClient` for new code |
| `base/BaseService.ts` | Base class pattern for services (where used) |

## Usage

Prefer the shared client and relative paths under the v1 API (the base URL should already include `/api`):

```typescript
import { apiClient } from "@/services/api/client";

const { data } = await apiClient.get("/v1/channels");
```

If your backend is mounted at `/api/v1`, ensure `VITE_API_BASE_URL` is e.g. `http://localhost:8080/api` so requests hit `/api/v1/...`.

## Configuration

Set in `.env` (Vite):

```bash
VITE_API_BASE_URL=http://localhost:8080/api
```

Redevelopment CORS on the API allows `http://localhost:5173` by default.

## Related documentation

- [`../features/README.md`](../features/README.md) — feature APIs and TanStack Query usage
- [`../app/store/README.md`](../app/store/README.md) — Redux (orthogonal to HTTP)
