package mediapublish

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/shekhar8352/PostEaze/entities"
	"github.com/shekhar8352/PostEaze/entities/repositories"
)

type scheduledMediaPayload struct {
	Items []struct {
		MediaAssetID *int64 `json:"media_asset_id"`
	} `json:"items"`
}

// MarkLinkedMediaAssetsPublished sets media_assets.status = published for each distinct
// media_asset_id present on scheduled post media JSON (workspace calendar flow).
func MarkLinkedMediaAssetsPublished(ctx context.Context, ownerID uuid.UUID, mediaJSON []byte) {
	var p scheduledMediaPayload
	if err := json.Unmarshal(mediaJSON, &p); err != nil {
		return
	}
	seen := make(map[int64]struct{})
	for _, it := range p.Items {
		if it.MediaAssetID == nil || *it.MediaAssetID <= 0 {
			continue
		}
		id := *it.MediaAssetID
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		_ = repositories.UpdateMediaAssetStatus(ctx, id, ownerID, string(entities.MediaAssetStatusPublished))
	}
}
