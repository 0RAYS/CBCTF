package utils

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

func GetClaimKeys(field string) []string {
	re := regexp.MustCompile(`\{([^{}]*)}`)
	matches := re.FindAllStringSubmatch(field, -1)

	var results []string
	for _, match := range matches {
		if len(match) > 1 {
			results = append(results, match[1])
		}
	}
	return results
}

// GetClaimRawValue resolves a single {key} or {nested.key} placeholder in field,
// looks up the path in resp, and returns the raw map value converted to T via
// JSON round-trip. No string formatting or template substitution is performed;
// the value is returned exactly as it appears in the map.
func GetClaimRawValue[T any](resp map[string]any, key string) (T, bool) {
	var zero T
	ks := strings.Split(key, ".")
	data := resp
	for i, k := range ks {
		v, ok := data[k]
		if !ok {
			return zero, false
		}
		if i < len(ks)-1 {
			subData, ok := v.(map[string]any)
			if !ok {
				return zero, false
			}
			data = subData
			continue
		}
		// Leaf value: convert to T via JSON round-trip.
		b, err := json.Marshal(v)
		if err != nil {
			fmt.Println(err)
			return zero, false
		}
		var out T
		if err = json.Unmarshal(b, &out); err != nil {
			fmt.Println(err)
			return zero, false
		}
		return out, true
	}
	return zero, false
}

func GetClaimStringValue(resp map[string]any, field string) (string, bool) {
	keys := GetClaimKeys(field)
	for _, key := range keys {
		out, ok := GetClaimRawValue[any](resp, key)
		if !ok {
			return "", false
		}
		field = strings.Replace(field, fmt.Sprintf("{%s}", key), fmt.Sprintf("%v", out), 1)
	}
	return field, true
}
