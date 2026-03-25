# Tasks package (Asynq)

Background work and scheduled jobs using [Asynq](https://github.com/hibiken/asynq) (Redis). The API process enqueues tasks; `cmd/worker` runs consumers.

## Architecture

- **Client** — Initialized in `main.go` (`tasks.InitClient`), used from handlers and webhooks.
- **Worker** — `go run cmd/worker/main.go` (or Docker service) runs `ServeMux` handlers.
- **Scheduler** — `scheduler.go` registers periodic jobs (e.g. Instagram profile sync every 6 hours).

## Files

| File | Role |
|------|------|
| `definitions.go` | Task type constants and JSON payloads |
| `client.go` | Enqueue helpers |
| `server.go` | Worker server configuration |
| `handlers.go` | Task handlers registered on the mux |
| `scheduler.go` | Cron-style schedules |

## Queues and priorities

Workers prioritize queues `fast`, `medium`, `slow` (see `server.go` for weights).

## Registered handlers

`RegisterHandlers` in `handlers.go` wires:

- `TypeEmailDelivery`
- `TypeLogMessage`
- `TypeInstagramComment`, `TypeInstagramMention`, `TypeInstagramStoryInsight`
- `TypeSyncInstagramProfiles`
- `TypeSyncInstagramPosts`
- `TypePeriodSnapshot`
- `TypeSyncInstagramAnalytics`

Add a new task: define type + payload in `definitions.go`, implement `Handle...` in `handlers.go`, register in `RegisterHandlers`, enqueue from code with `tasks.EnqueueTask` (or scheduler).

## Periodic tasks

Instagram profile sync is registered in `scheduler.go` (e.g. `@every 6h` for `TypeSyncInstagramProfiles`). Adjust schedules there.

## Running the worker

```bash
cd backend
go run cmd/worker/main.go
```

With Docker Compose, use the project’s compose file that includes the worker service and Redis.

## Practices

- Return `fmt.Errorf("... %w", asynq.SkipRetry)` for non-retryable failures.
- Keep payloads small (IDs only); load full rows inside the handler.

## Related documentation

- [Backend README](../README.md)
- [Webhooks](../api/webhooks/README.md) (enqueue from Instagram events)
