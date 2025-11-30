# Tasks Package (Asynq Integration)

This package handles asynchronous task processing and periodic (cron) jobs using [Asynq](https://github.com/hibiken/asynq), a robust Redis-based task queue library.

## 🏗 Architecture

The system is decoupled into **Producers** and **Consumers**:

1.  **Client (Producer)**:
    -   Embedded in the main API server (`backend/main.go`).
    -   Enqueues tasks to Redis.
    -   Can schedule tasks for immediate execution, future execution, or periodic execution.

2.  **Worker (Consumer)**:
    -   Runs as a standalone process (`backend/cmd/worker/main.go`).
    -   Pulls tasks from Redis and executes them.
    -   Scalable: You can run multiple worker replicas to handle higher loads.

### 📂 Directory Structure

-   `definitions.go`: **Start here.** Defines Task Type constants and Payload structs.
-   `client.go`: Wrapper for the Asynq Client. Used to enqueue tasks.
-   `server.go`: Configures the Asynq Server (Worker), including queue priorities and concurrency.
-   `handlers.go`: Contains the actual business logic (handler functions) for each task.
-   `scheduler.go`: Configures the Asynq Scheduler for periodic (cron) jobs.

### 🚦 Queues & Priorities

We utilize a weighted priority system. Workers check queues in the following order of frequency:

| Queue Name | Priority | Use Case |
| :--- | :--- | :--- |
| `fast` | **6** | Critical, time-sensitive tasks (e.g., OTPs, real-time notifications). |
| `medium` | **3** | Standard tasks (e.g., welcome emails, data processing). |
| `slow` | **1** | Low priority, background maintenance, heavy reports. |

---

## 🚀 Implementation Guide

Follow these steps to add a new background task.

### Step 1: Define the Task
Open `backend/tasks/definitions.go` and add a unique **Task Type** string and a **Payload Struct**.

```go
// definitions.go

const (
    TypeEmailDelivery = "email:deliver"
    TypeLogMessage    = "log:message"
    TypeInstagramComment = "instagram:comment"
    TypeInstagramMention = "instagram:mention"
    TypeInstagramStoryInsight = "instagram:story_insight"
)

type GenerateReportPayload struct { // [NEW]
    UserID    string
    DateRange string
}
```

### Step 2: Create the Handler
Open `backend/tasks/handlers.go` and write the function that processes the task.

```go
// handlers.go

func HandleGenerateReportTask(ctx context.Context, t *asynq.Task) error {
    var p GenerateReportPayload
    if err := json.Unmarshal(t.Payload(), &p); err != nil {
        // Return error to retry, or wrap in fmt.Errorf(..., asynq.SkipRetry) to fail permanently
        return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
    }

    log.Printf("Generating report for user: %s", p.UserID)
    
    // ... Your business logic here ...

    return nil
}
```

### Step 3: Register the Handler
Still in `backend/tasks/handlers.go`, register your new handler in the `RegisterHandlers` function.

```go
func RegisterHandlers(mux *asynq.ServeMux) {
    mux.HandleFunc(TypeEmailDelivery, HandleEmailDeliveryTask)
    mux.HandleFunc(TypeInstagramComment, HandleInstagramCommentTask)
    mux.HandleFunc(TypeInstagramMention, HandleInstagramMentionTask)
    mux.HandleFunc(TypeInstagramStoryInsight, HandleInstagramStoryInsightTask)
}
```

### Step 4: Enqueue the Task
Use the client to enqueue the task from anywhere in your application (e.g., inside an API controller).

```go
import (
    "encoding/json"
    "github.com/hibiken/asynq"
    "github.com/shekhar8352/PostEaze/tasks"
)

func TriggerReport(userID string) {
    // 1. Create Payload
    payload, _ := json.Marshal(tasks.GenerateReportPayload{
        UserID: userID, 
        DateRange: "last_week",
    })

    // 2. Create Task
    task := asynq.NewTask(tasks.TypeGenerateReport, payload)

    // 3. Enqueue
    // Simple enqueue to default 'medium' queue:
    // tasks.EnqueueTask(task, tasks.QueueMedium)

    // OR with Options (Delay, MaxRetry, etc.):
    tasks.EnqueueTask(task, tasks.QueueSlow, 
        asynq.ProcessIn(10*time.Minute), // Delay execution
        asynq.MaxRetry(3),               // Retry up to 3 times
        asynq.Timeout(20*time.Minute),   // Fail if it takes longer than 20m
    )
}
```

---

## ⏰ Cron Jobs (Periodic Tasks)

To run tasks automatically on a schedule (e.g., "every day at midnight"), use the **Scheduler**.

Open `backend/tasks/scheduler.go` and add your entry in `InitScheduler`:

```go
// scheduler.go

func InitScheduler() error {
    // ... initialization code ...

    // Example: Run 'TypeLogMessage' every minute
    if _, err := scheduler.Register("* * * * *", asynq.NewTask(tasks.TypeLogMessage, []byte("{\"Message\":\"Cron Tick\"}"))); err != nil {
        return err
    }

    // Example: Run 'TypeGenerateReport' every day at midnight
    // payload, _ := json.Marshal(tasks.GenerateReportPayload{...})
    // if _, err := scheduler.Register("@daily", asynq.NewTask(tasks.TypeGenerateReport, payload)); err != nil {
    //     return err
    // }

    return nil
}
```

---

## 🛠 Advanced Usage & Best Practices

### 1. Error Handling & Retries
-   **Automatic Retries**: If your handler returns an error, Asynq automatically retries the task with exponential backoff.
-   **Fatal Errors**: If an error is unrecoverable (e.g., invalid JSON, invalid user ID), wrap the error with `asynq.SkipRetry` to stop retrying immediately.
    ```go
    return fmt.Errorf("critical error: %w", asynq.SkipRetry)
    ```

### 2. Task Options
When enqueueing, you can pass various options:
-   `asynq.ProcessIn(d)`: Schedule task to run in the future.
-   `asynq.MaxRetry(n)`: Set custom retry limit (default is 25).
-   `asynq.Queue(q)`: Send to a specific queue (`fast`, `medium`, `slow`).
-   `asynq.TaskID(id)`: Set a custom ID. Useful for **deduplication** (if a task with this ID already exists, the new one is ignored).

### 3. Payload Size
Keep payloads **small**. Store only IDs (e.g., `UserID`, `OrderID`) in the payload and fetch the full data from the database within the handler. Do not pass large JSON objects or files in the payload.

---

## 🏃‍♂️ Running the Worker

### Local Development (Docker Compose)
The worker is configured as a service in `docker-compose.local.yml`. It starts automatically with `docker-compose up`.

```bash
# View worker logs
docker-compose -f docker-compose.local.yml logs -f worker

# Restart worker (needed if you change handler code)
docker-compose -f docker-compose.local.yml restart worker
```

### Manual Execution
If you want to run it outside Docker:

```bash
cd backend
go run cmd/worker/main.go
```
