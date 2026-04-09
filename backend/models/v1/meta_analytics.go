package modelsv1

// SyncMetaAnalyticsRequest optionally limits sync to specific channel IDs (must be owned by the user).
type SyncMetaAnalyticsRequest struct {
	ChannelIDs []int64 `json:"channel_ids"`
}

// SyncMetaAnalyticsResult is the outcome for one channel.
type SyncMetaAnalyticsResult struct {
	ChannelID int64  `json:"channel_id"`
	Provider  string `json:"provider"`
	Status    string `json:"status"`
	Message   string `json:"message,omitempty"`
}

// SyncMetaAnalyticsResponse aggregates per-channel sync results.
type SyncMetaAnalyticsResponse struct {
	Results []SyncMetaAnalyticsResult `json:"results"`
}
