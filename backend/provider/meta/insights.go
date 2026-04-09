package meta

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// FBInsightsResponse is the Graph API shape for /insights endpoints.
type FBInsightsResponse struct {
	Data []FBInsightMetric `json:"data"`
}

// FBInsightMetric is one metric series from insights.
type FBInsightMetric struct {
	Name   string `json:"name"`
	Period string `json:"period"`
	Values []struct {
		Value   json.RawMessage `json:"value"`
		EndTime string          `json:"end_time"`
	} `json:"values"`
}

// GraphError is a typical Graph API error payload.
type GraphError struct {
	Message   string `json:"message"`
	Type      string `json:"type"`
	Code      int    `json:"code"`
	FBTraceID string `json:"fbtrace_id"`
}

type graphErrorEnvelope struct {
	Error GraphError `json:"error"`
}

// FeedResponse is a page feed /posts listing.
type FeedResponse struct {
	Data   []FeedPost `json:"data"`
	Paging *struct {
		Cursors *struct {
			Before string `json:"before"`
			After  string `json:"after"`
		} `json:"cursors"`
		Next string `json:"next"`
	} `json:"paging"`
}

// FeedPost is a minimal post object from the feed.
type FeedPost struct {
	ID          string `json:"id"`
	Message     string `json:"message"`
	CreatedTime string `json:"created_time"`
}

// PageDailyInsightMetrics lists common day-period Page metrics (comma-separated in one request).
var PageDailyInsightMetrics = []string{
	"page_fans",
	"page_fan_adds",
	"page_impressions",
	"page_impressions_unique",
	"page_post_engagements",
	"page_engaged_users",
}

// PostInsightMetrics lists common post-level metrics.
var PostInsightMetrics = []string{
	"post_impressions",
	"post_impressions_unique",
	"post_engaged_users",
	"post_clicks",
	"post_video_views",
	"post_reactions_like_total",
}

// GetPageDailyInsights calls GET /{page-id}/insights with period=day.
func (p *MetaProviderImpl) GetPageDailyInsights(pageAccessToken, pageID string, sinceUnix, untilUnix int64, metrics []string) (*FBInsightsResponse, error) {
	if len(metrics) == 0 {
		metrics = PageDailyInsightMetrics
	}
	u, err := url.Parse(fmt.Sprintf("%s/%s/insights", GraphAPIURL, pageID))
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("access_token", pageAccessToken)
	q.Set("metric", strings.Join(metrics, ","))
	q.Set("period", "day")
	q.Set("since", strconv.FormatInt(sinceUnix, 10))
	q.Set("until", strconv.FormatInt(untilUnix, 10))
	u.RawQuery = q.Encode()

	return decodeInsightsGET(u.String())
}

// ListPageFeedPosts calls GET /{page-id}/posts with minimal fields.
func (p *MetaProviderImpl) ListPageFeedPosts(pageAccessToken, pageID string, limit int, after string) (*FeedResponse, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	u, err := url.Parse(fmt.Sprintf("%s/%s/posts", GraphAPIURL, pageID))
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("access_token", pageAccessToken)
	q.Set("fields", "id,message,created_time")
	q.Set("limit", strconv.Itoa(limit))
	if after != "" {
		q.Set("after", after)
	}
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, graphErrFromBody(body, resp.Status)
	}
	var out FeedResponse
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, fmt.Errorf("decode feed: %w", err)
	}
	return &out, nil
}

// GetPostInsights calls GET /{post-id}/insights (lifetime or default period for posts).
func (p *MetaProviderImpl) GetPostInsights(pageAccessToken, postID string, metrics []string) (*FBInsightsResponse, error) {
	if len(metrics) == 0 {
		metrics = PostInsightMetrics
	}
	u, err := url.Parse(fmt.Sprintf("%s/%s/insights", GraphAPIURL, postID))
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("access_token", pageAccessToken)
	q.Set("metric", strings.Join(metrics, ","))
	u.RawQuery = q.Encode()

	return decodeInsightsGET(u.String())
}

func decodeInsightsGET(fullURL string) (*FBInsightsResponse, error) {
	resp, err := http.Get(fullURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, graphErrFromBody(body, resp.Status)
	}
	var env graphErrorEnvelope
	if json.Unmarshal(body, &env) == nil && env.Error.Message != "" {
		return nil, fmt.Errorf("graph api: %s (code %d)", env.Error.Message, env.Error.Code)
	}
	var out FBInsightsResponse
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, fmt.Errorf("decode insights: %w", err)
	}
	return &out, nil
}

func graphErrFromBody(body []byte, status string) error {
	var env graphErrorEnvelope
	if json.Unmarshal(body, &env) == nil && env.Error.Message != "" {
		return fmt.Errorf("graph api %s: %s (code %d)", status, env.Error.Message, env.Error.Code)
	}
	return fmt.Errorf("graph api %s: %s", status, strings.TrimSpace(string(body)))
}

// ParseScalarInsightValue extracts an int from a metric value (handles JSON number or string).
func ParseScalarInsightValue(raw json.RawMessage) (int64, bool) {
	if len(raw) == 0 {
		return 0, false
	}
	var n int64
	if err := json.Unmarshal(raw, &n); err == nil {
		return n, true
	}
	var f float64
	if err := json.Unmarshal(raw, &f); err == nil {
		return int64(f), true
	}
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		if v, err := strconv.ParseInt(s, 10, 64); err == nil {
			return v, true
		}
	}
	return 0, false
}

// InsightValueByMetricEndDate merges insight series keyed by metric name into map dateKey -> metric -> value.
func InsightValueByMetricEndDate(resp *FBInsightsResponse) map[string]map[string]int64 {
	out := make(map[string]map[string]int64)
	if resp == nil {
		return out
	}
	for _, series := range resp.Data {
		name := series.Name
		for _, v := range series.Values {
			dateKey := insightEndTimeToDateKey(v.EndTime)
			if dateKey == "" {
				continue
			}
			val, ok := ParseScalarInsightValue(v.Value)
			if !ok {
				continue
			}
			if out[dateKey] == nil {
				out[dateKey] = make(map[string]int64)
			}
			out[dateKey][name] = val
		}
	}
	return out
}

func insightEndTimeToDateKey(endTime string) string {
	if endTime == "" {
		return ""
	}
	// "2024-01-15T07:00:00+0000"
	if len(endTime) >= 10 {
		return endTime[:10]
	}
	return ""
}
