package socialcomments

import (
	"bytes"
	"encoding/json"
)

// DecodeJSONMap decodes JSON into map[string]interface{} with UseNumber so large
// Instagram IDs are kept as json.Number instead of float64 (avoids precision loss).
func DecodeJSONMap(raw []byte) (map[string]interface{}, error) {
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	var m map[string]interface{}
	if err := dec.Decode(&m); err != nil {
		return nil, err
	}
	return m, nil
}
