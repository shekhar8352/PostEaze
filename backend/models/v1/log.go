package modelsv1

// LogEntry represents a single log entry
type LogEntry struct {
	Timestamp string            `json:"timestamp"`
	Level     string            `json:"level"`
	Message   string            `json:"message"`
	LogID     string            `json:"log_id,omitempty"`
	Method    string            `json:"method,omitempty"`
	Path      string            `json:"path,omitempty"`
	Status    int               `json:"status,omitempty"`
	Duration  string            `json:"duration,omitempty"`
	IP        string            `json:"ip,omitempty"`
	UserAgent string            `json:"user_agent,omitempty"`
	Extra     map[string]string `json:"extra,omitempty"`
}

// LogsResponse represents the API response for logs
type LogsResponse struct {
	Success bool       `json:"success"`
	Data    []LogEntry `json:"data"`
	Total   int        `json:"total"`
	Page    int        `json:"page"`
	Limit   int        `json:"limit"`
	Message string     `json:"message,omitempty"`
}