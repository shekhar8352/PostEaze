package utils

import (
	"encoding/json"
	"math"
	"sort"
	"strings"

	"github.com/shekhar8352/PostEaze/entities"
	modelsv1 "github.com/shekhar8352/PostEaze/models/v1"
)

const audienceBucketLimit = 8

var audienceMetricTitles = map[string]string{
	"audience_city":        "Top cities",
	"audience_country":     "Countries",
	"audience_gender_age":  "Gender & age",
	"audience_locale":      "Locales",
}

type audienceEnvelope struct {
	Data []audienceInsight `json:"data"`
}

type audienceInsight struct {
	Name   string           `json:"name"`
	Values []audienceValue  `json:"values"`
}

type audienceValue struct {
	Value json.RawMessage `json:"value"`
}

// AudienceDashboardFromSnapshot turns stored Meta JSON into UI-friendly buckets (latest snapshot only).
func AudienceDashboardFromSnapshot(snap *entities.InstagramAudienceSnapshot) (*modelsv1.AudienceDashboard, error) {
	if snap == nil || len(snap.Raw) == 0 {
		return nil, nil
	}

	var env audienceEnvelope
	if err := json.Unmarshal(snap.Raw, &env); err != nil {
		return nil, err
	}

	out := &modelsv1.AudienceDashboard{
		SnapshotDate: FormatAnalyticsDate(snap.SnapshotDate),
		CreatedAt:    snap.CreatedAt.UTC().Format("2006-01-02T15:04:05Z07:00"),
		Metrics:      nil,
	}

	for _, ins := range env.Data {
		if ins.Name == "" || len(ins.Values) == 0 {
			continue
		}
		buckets := rawValueToBuckets(ins.Values[0].Value, audienceBucketLimit)
		if len(buckets) == 0 {
			continue
		}
		title := audienceMetricTitles[ins.Name]
		if title == "" {
			title = strings.ReplaceAll(ins.Name, "_", " ")
		}
		out.Metrics = append(out.Metrics, modelsv1.AudienceMetricBlock{
			Key:   ins.Name,
			Title: title,
			Items: buckets,
		})
	}

	if len(out.Metrics) == 0 {
		return out, nil
	}
	return out, nil
}

func rawValueToBuckets(raw json.RawMessage, limit int) []modelsv1.AudienceBucket {
	if len(raw) == 0 || string(raw) == "null" {
		return nil
	}

	// Single number (unusual for audience_* but handle)
	var n float64
	if err := json.Unmarshal(raw, &n); err == nil {
		return []modelsv1.AudienceBucket{{Label: "Total", Value: int64(n), Pct: 100}}
	}

	var m map[string]interface{}
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil
	}

	type kv struct {
		k string
		v int64
	}
	var pairs []kv
	var total int64
	for k, v := range m {
		if c, ok := coerceAudienceCount(v); ok && c >= 0 {
			pairs = append(pairs, kv{k: k, v: c})
			total += c
		}
	}
	if len(pairs) == 0 || total == 0 {
		return nil
	}

	sort.Slice(pairs, func(i, j int) bool { return pairs[i].v > pairs[j].v })
	if len(pairs) > limit {
		pairs = pairs[:limit]
	}

	out := make([]modelsv1.AudienceBucket, 0, len(pairs))
	for _, p := range pairs {
		pct := (float64(p.v) / float64(total)) * 100
		out = append(out, modelsv1.AudienceBucket{
			Label: humanizeAudienceLabel(p.k),
			Value: p.v,
			Pct:   math.Round(pct*10) / 10,
		})
	}
	return out
}

func coerceAudienceCount(v interface{}) (int64, bool) {
	switch x := v.(type) {
	case float64:
		return int64(x), true
	case int64:
		return x, true
	case int:
		return int64(x), true
	case json.Number:
		i, err := x.Int64()
		return i, err == nil
	default:
		return 0, false
	}
}

func humanizeAudienceLabel(s string) string {
	s = strings.ReplaceAll(s, "_", " ")
	s = strings.TrimSpace(s)
	if s == "" {
		return "—"
	}
	return s
}
