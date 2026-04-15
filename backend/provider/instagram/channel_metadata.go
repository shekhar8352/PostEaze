package instagram

import (
	"encoding/json"
	"fmt"
)

// UserIDFromChannelMetadata returns the Instagram user (IG) id stored on a channel row.
func UserIDFromChannelMetadata(metadata []byte) (string, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(metadata, &m); err != nil {
		return "", fmt.Errorf("invalid channel metadata")
	}
	id, ok := m["id"].(string)
	if !ok || id == "" {
		if v, ok := m["id"].(float64); ok {
			return fmt.Sprintf("%.0f", v), nil
		}
		return "", fmt.Errorf("instagram user id missing in channel metadata")
	}
	return id, nil
}
