package webhooks

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/shekhar8352/PostEaze/socialcomments"
	"github.com/shekhar8352/PostEaze/tasks"
	"github.com/shekhar8352/PostEaze/utils"
)

// HandleInstagramWebhookVerify godoc
// @Summary      Verify Instagram Webhook
// @Description  Handles Meta's webhook verification challenge
// @Tags         Webhooks
// @Param        hub.mode query string true "Mode"
// @Param        hub.challenge query string true "Challenge"
// @Param        hub.verify_token query string true "Verify Token"
// @Success      200 {string} string "challenge"
// @Failure      403 {object} map[string]interface{}
// @Router       /webhooks/instagram [get]
func HandleInstagramWebhookVerify(c *gin.Context) {
	mode := c.Query("hub.mode")
	token := c.Query("hub.verify_token")
	challenge := c.Query("hub.challenge")

	verifyToken := os.Getenv("INSTAGRAM_WEBHOOK_VERIFY_TOKEN")

	utils.Logger.Info(c.Request.Context(), "Received webhook verify request: ")

	if mode == "subscribe" && token == verifyToken {
		c.String(http.StatusOK, challenge)
	} else {
		c.JSON(http.StatusForbidden, gin.H{"status": "error", "msg": "Verification failed"})
	}
}

// HandleInstagramWebhookEvent godoc
// @Summary      Handle Instagram Webhook Event
// @Description  Receives and processes Instagram webhook events
// @Tags         Webhooks
// @Accept       json
// @Produce      json
// @Param        X-Hub-Signature-256 header string true "Signature"
// @Param        request body map[string]interface{} true "Webhook Payload"
// @Success      200 {string} string "EVENT_RECEIVED"
// @Failure      400 {object} map[string]interface{}
// @Failure      403 {object} map[string]interface{}
// @Router       /webhooks/instagram [post]
func HandleInstagramWebhookEvent(c *gin.Context) {
	// 1. Verify Signature
	utils.Logger.Info(c.Request.Context(), "Received webhook request: ")

	signature := c.GetHeader("X-Hub-Signature-256")
	if signature == "" {
		utils.SendError(c, http.StatusForbidden, "Missing signature")
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "Failed to read body")
		return
	}

	utils.Logger.Info(c.Request.Context(), "Received webhook body: ", body)

	appSecret := os.Getenv("INSTAGRAM_APP_SECRET")
	if !verifySignature(body, signature, appSecret) {
		utils.SendError(c, http.StatusForbidden, "Invalid signature")
		return
	}

	// 2. Parse Payload (UseNumber preserves large Instagram IDs)
	payload, err := socialcomments.DecodeJSONMap(body)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "Failed to parse JSON")
		return
	}

	utils.Logger.Info(c.Request.Context(), "Received webhook payload: ", payload)

	// 3. Process Events (Async)
	// Meta sends a list of entries
	entries, ok := payload["entry"].([]interface{})
	if !ok || len(entries) == 0 {
		utils.Logger.Warn(c.Request.Context(), "Instagram webhook: entry missing or empty (object=%v)", payload["object"])
		c.String(http.StatusOK, "EVENT_RECEIVED")
		return
	}
	for _, entry := range entries {
		entryMap, ok := entry.(map[string]interface{})
		if !ok {
			continue
		}

		// Each entry has a list of changes or messaging events
		changes, ok := entryMap["changes"].([]interface{})
		if !ok || len(changes) == 0 {
			utils.Logger.Warn(c.Request.Context(), "Instagram webhook: entry has no changes (entry_id=%v)", entryMap["id"])
			continue
		}
		for _, change := range changes {
			changeMap, ok := change.(map[string]interface{})
			if !ok {
				continue
			}

			// Enqueue task based on field
			field, _ := changeMap["field"].(string)

			// Wrap with entry id (Instagram user id) so workers can resolve channel without scanning posts.
			taskPayload := map[string]interface{}{
				"entry_id": socialcomments.StringFromAny(entryMap["id"]),
				"change":   changeMap,
			}
			changeJSON, err := json.Marshal(taskPayload)
			if err != nil {
				utils.Logger.Error(c.Request.Context(), "Failed to marshal change", err)
				continue
			}

			var taskType string
			switch field {
			case "comments":
				taskType = tasks.TypeInstagramComment
			case "mentions":
				taskType = tasks.TypeInstagramMention
			case "story_insights":
				taskType = tasks.TypeInstagramStoryInsight
			default:
				// Log unknown field
				utils.Logger.Info(c.Request.Context(), "Unknown webhook field: ", field)
				continue
			}

			task := asynq.NewTask(taskType, changeJSON)
			_, err = tasks.EnqueueTask(task, tasks.QueueMedium)
			if err != nil {
				utils.Logger.Error(c.Request.Context(), "Failed to enqueue task", err)
			}
		}
	}

	c.String(http.StatusOK, "EVENT_RECEIVED")
}

func verifySignature(payload []byte, signatureHeader string, secret string) bool {
	// Signature header format: sha256=<signature>
	parts := strings.Split(signatureHeader, "=")
	if len(parts) != 2 || parts[0] != "sha256" {
		return false
	}
	signature := parts[1]

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	expectedSignature := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}
