package instagram

import (
	"encoding/json"
	"fmt"
	"strings"
)

// GraphAPIError is a failed response body from graph.instagram.com with a JSON "error" object.
type GraphAPIError struct {
	Code         int
	ErrorSubcode int
	Message      string
	Type         string
	FBTraceID    string
}

func (e *GraphAPIError) Error() string {
	if e == nil {
		return ""
	}
	return fmt.Sprintf("instagram graph API: %s (code=%d, error_subcode=%d)", e.Message, e.Code, e.ErrorSubcode)
}

type graphErrorEnvelope struct {
	Error struct {
		Code         float64 `json:"code"`
		ErrorSubcode float64 `json:"error_subcode"`
		Message      string  `json:"message"`
		Type         string  `json:"type"`
		FBTraceID    string  `json:"fbtrace_id"`
	} `json:"error"`
}

// GraphAPIErrorFromBody parses a Graph API error from a non-2xx response body.
func GraphAPIErrorFromBody(body []byte) (*GraphAPIError, bool) {
	var env graphErrorEnvelope
	if err := json.Unmarshal(body, &env); err != nil {
		return nil, false
	}
	if env.Error.Message == "" && env.Error.Code == 0 {
		return nil, false
	}
	return &GraphAPIError{
		Code:         int(env.Error.Code),
		ErrorSubcode: int(env.Error.ErrorSubcode),
		Message:      env.Error.Message,
		Type:         env.Error.Type,
		FBTraceID:    env.Error.FBTraceID,
	}, true
}

// IsInsightsUnavailableForever returns true for errors where retrying or changing metrics will not help.
func (e *GraphAPIError) IsInsightsUnavailableForever() bool {
	if e == nil {
		return false
	}
	// 2108006: media posted before the account was converted to business/creator (Meta IGApiException).
	if e.ErrorSubcode == 2108006 {
		return true
	}
	msg := strings.ToLower(e.Message)
	if strings.Contains(msg, "personal") && strings.Contains(msg, "business") {
		return true
	}
	if strings.Contains(msg, "unternehmenskonto") || strings.Contains(msg, "business account") {
		if strings.Contains(msg, "converted") || strings.Contains(msg, "umgewandelt") || strings.Contains(msg, "before") {
			return true
		}
	}
	return false
}
