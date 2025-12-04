package apiv1

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/shekhar8352/PostEaze/tasks"
	"github.com/shekhar8352/PostEaze/utils"
)

// TriggerInstagramSyncHandler godoc
// @Summary      Trigger Instagram Profile Sync (Development Only)
// @Description  Manually triggers the Instagram profile sync job for immediate execution. Only available in development mode.
// @Tags         Cron Jobs
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /cron/trigger-instagram-sync [post]
func TriggerInstagramSyncHandler(c *gin.Context) {
	// Only allow in development mode
	env := os.Getenv("ENV")
	if env != "development" && env != "dev" && env != "" {
		utils.SendError(c, http.StatusForbidden, "This endpoint is only available in development mode")
		return
	}

	// Create Asynq client
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
	defer client.Close()

	// Enqueue the sync task to the slow queue (background jobs)
	task := asynq.NewTask(tasks.TypeSyncInstagramProfiles, nil)
	info, err := client.Enqueue(task, asynq.Queue(tasks.QueueSlow))
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to enqueue sync task: "+err.Error())
		return
	}

	utils.SendSuccess(c, gin.H{
		"task_id": info.ID,
		"queue":   info.Queue,
		"message": "Instagram profile sync task enqueued successfully",
	}, "Sync task triggered")
}

// TriggerInstagramPostsSyncHandler godoc
// @Summary      Trigger Instagram Posts Sync (Development Only)
// @Description  Manually triggers the Instagram posts sync job for immediate execution. Only available in development mode.
// @Tags         Cron Jobs
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /cron/trigger-instagram-posts [post]
func TriggerInstagramPostsSyncHandler(c *gin.Context) {
	// Only allow in development mode
	env := os.Getenv("ENV")
	if env != "development" && env != "dev" && env != "" {
		utils.SendError(c, http.StatusForbidden, "This endpoint is only available in development mode")
		return
	}

	// Create Asynq client
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
	defer client.Close()

	// Enqueue the posts sync task
	task := asynq.NewTask(tasks.TypeSyncInstagramPosts, nil)
	info, err := client.Enqueue(task, asynq.Queue(tasks.QueueSlow))
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to enqueue posts sync task: "+err.Error())
		return
	}

	utils.SendSuccess(c, gin.H{
		"task_id": info.ID,
		"queue":   info.Queue,
		"message": "Instagram posts sync task enqueued successfully",
	}, "Posts sync task triggered")
}

// TriggerInstagramAnalyticsSyncHandler godoc
// @Summary      Trigger Instagram Analytics Sync (Development Only)
// @Description  Manually triggers the Instagram analytics sync job for immediate execution. Only available in development mode.
// @Tags         Cron Jobs
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /cron/trigger-instagram-analytics [post]
func TriggerInstagramAnalyticsSyncHandler(c *gin.Context) {
	// Only allow in development mode
	env := os.Getenv("ENV")
	if env != "development" && env != "dev" && env != "" {
		utils.SendError(c, http.StatusForbidden, "This endpoint is only available in development mode")
		return
	}

	// Create Asynq client
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
	defer client.Close()

	// Enqueue the analytics sync task
	task := asynq.NewTask(tasks.TypeSyncInstagramAnalytics, nil)
	info, err := client.Enqueue(task, asynq.Queue(tasks.QueueSlow))
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to enqueue analytics sync task: "+err.Error())
		return
	}

	utils.SendSuccess(c, gin.H{
		"task_id": info.ID,
		"queue":   info.Queue,
		"message": "Instagram analytics sync task enqueued successfully",
	}, "Analytics sync task triggered")
}
