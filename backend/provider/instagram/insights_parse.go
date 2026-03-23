package instagram

import (
	"bytes"
	"encoding/json"
	"strconv"
)

// ParseScalarInsightValue coerces Meta insight metric values to int64.
// Handles json.Number, float64, int-like types; returns false for maps/slices/unsupported values.
func ParseScalarInsightValue(v interface{}) (int64, bool) {
	if v == nil {
		return 0, false
	}
	switch x := v.(type) {
	case json.Number:
		i, err := x.Int64()
		if err != nil {
			f, err2 := x.Float64()
			if err2 != nil {
				return 0, false
			}
			return int64(f), true
		}
		return i, true
	case float64:
		return int64(x), true
	case float32:
		return int64(x), true
	case int:
		return int64(x), true
	case int64:
		return x, true
	case uint64:
		return int64(x), true
	case string:
		i, err := strconv.ParseInt(x, 10, 64)
		if err != nil {
			f, err2 := strconv.ParseFloat(x, 64)
			if err2 != nil {
				return 0, false
			}
			return int64(f), true
		}
		return i, true
	default:
		return 0, false
	}
}

// DecodeInsightsResponseJSON decodes JSON using UseNumber so numeric metrics stay precise.
func DecodeInsightsResponseJSON(data []byte) (*InsightsResponse, error) {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	var resp InsightsResponse
	if err := dec.Decode(&resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
