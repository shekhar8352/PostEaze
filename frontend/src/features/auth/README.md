# Authentication feature

Login, registration, forgot password, and email verification. **Firebase** handles client identity; the Go API **`POST /api/v1/auth/authenticate`** (and refresh/logout) issues **JWTs** stored and sent via **`services/api/interceptors.ts`**.

## Key pieces

- **`AuthProvider`** — loads session on startup, exposes auth context to the tree.
- **`authRoutes.tsx`** — `/login`, `/register`, `/forgot-password`, `/email-verify` wrapped in **`PublicRoute`** where applicable.
- **`services/authService.ts`**, **`authQueries.ts`** — API calls and TanStack Query hooks.
- **`firebase.config.ts`**, **`firebaseHelper.ts`** — Firebase client setup.
- **`authSlice.ts`** — Redux state retained for flows that still use it; prefer hooks/context + Query for new work.

## Using auth in components

```typescript
import { useAuth } from "@/features/auth";

const { user, loading, logout } = useAuth();
```

For Redux selectors, use typed hooks from **`@/app/store/hooks`**.

## Related documentation

- [`../README.md`](../README.md) — features overview
- [`../../app/routes/README.md`](../../app/routes/README.md) — `ProtectedLayout` / `PublicRoute`
- [`../../services/README.md`](../../services/README.md) — `apiClient` and interceptors
