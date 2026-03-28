# Redux store

Global client state for the PostEaze SPA using **Redux Toolkit**. Server state (auth session bootstrap, channels, analytics, scheduled posts) is largely loaded with **TanStack Query** inside features; this store holds slices that still use Redux (for example channel UI state).

## Files

| File | Role |
|------|------|
| `store.ts` | `configureStore` with default middleware (`serializableCheck: false` for non-serializable values where needed) |
| `rootReducer.ts` | Combines feature reducers |
| `hooks.ts` | Typed `useAppDispatch` and `useAppSelector` |

## Usage

```typescript
import { useAppDispatch, useAppSelector } from "@/app/store/hooks";
```

Wrap the tree with `StoreProvider` from `app/providers/StoreProvider.tsx` (already done in `App.tsx`).

## Related documentation

- [`../routes/README.md`](../routes/README.md) — routing and protected layouts
- [`../../features/README.md`](../../features/README.md) — feature slices and patterns
- [`../../services/README.md`](../../services/README.md) — API client and interceptors
