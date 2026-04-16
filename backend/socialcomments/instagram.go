package socialcomments

import (
	"encoding/json"
	"fmt"
)

// InstagramCommentNormalized is a parsed Instagram "comments" webhook value.
type InstagramCommentNormalized struct {
	MediaID         string
	CommentID       string
	Text            string
	ParentCommentID *string
	AuthorUserID    *string
	AuthorUsername  *string
}

// ParseInstagramCommentValue extracts fields from Meta's comments field value object.
func ParseInstagramCommentValue(value map[string]interface{}) (*InstagramCommentNormalized, error) {
	if value == nil {
		return nil, fmt.Errorf("empty value")
	}

	commentID := firstNonEmpty(
		StringFromAny(value["id"]),
		StringFromAny(value["comment_id"]),
	)
	if commentID == "" {
		return nil, fmt.Errorf("missing comment id")
	}

	mediaID := firstNonEmpty(StringFromAny(value["media_id"]), StringFromAny(value["ig_media_id"]))
	if mediaID == "" {
		mediaID = StringFromAny(value["object_id"])
	}
	if mediaID == "" {
		if m, ok := value["media"].(map[string]interface{}); ok {
			mediaID = StringFromAny(m["id"])
		}
	}
	if mediaID == "" {
		return nil, fmt.Errorf("missing media id")
	}

	text := StringFromAny(value["text"])

	var parent *string
	if p := firstNonEmpty(StringFromAny(value["parent_id"]), StringFromAny(value["parent_comment_id"])); p != "" {
		parent = &p
	}

	n := &InstagramCommentNormalized{
		MediaID:         mediaID,
		CommentID:       commentID,
		Text:            text,
		ParentCommentID: parent,
	}

	if from, ok := value["from"].(map[string]interface{}); ok {
		if u := StringFromAny(from["id"]); u != "" {
			n.AuthorUserID = &u
		}
		if un := StringFromAny(from["username"]); un != "" {
			n.AuthorUsername = &un
		}
	}

	return n, nil
}

// MarshalRawPayload returns JSON bytes for storage, or nil if empty.
func MarshalRawPayload(v interface{}) []byte {
	if v == nil {
		return nil
	}
	b, err := json.Marshal(v)
	if err != nil || len(b) == 0 {
		return nil
	}
	return b
}

// StringFromAny converts JSON-decoded values (string, number, etc.) to string.
func StringFromAny(v interface{}) string {
	return stringFromAny(v)
}

func stringFromAny(v interface{}) string {
	if v == nil {
		return ""
	}
	switch t := v.(type) {
	case string:
		return t
	case float64:
		// JSON numbers
		return fmt.Sprintf("%.0f", t)
	case json.Number:
		return t.String()
	case int:
		return fmt.Sprintf("%d", t)
	case int64:
		return fmt.Sprintf("%d", t)
	default:
		return fmt.Sprintf("%v", t)
	}
}

func firstNonEmpty(a, b string) string {
	if a != "" {
		return a
	}
	return b
}
