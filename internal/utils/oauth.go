package utils

import (
	"CBCTF/internal/log"
	"encoding/json"
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

// GetClaimValue resolves a single {key} or {nested.key} placeholder in field,
// looks up the path in resp, and returns the raw map value converted to T via
// JSON round-trip. No string formatting or template substitution is performed;
// the value is returned exactly as it appears in the map.
func GetClaimValue[T any](resp map[string]any, field string) (T, bool) {
	var zero T
	keys := GetClaimKeys(field)
	if len(keys) != 1 {
		return zero, false
	}
	ks := strings.Split(keys[0], ".")
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
			log.Logger.Warningf("Failed to marshal claim value for key %s: %s", field, err.Error())
			return zero, false
		}
		var out T
		if err = json.Unmarshal(b, &out); err != nil {
			log.Logger.Warningf("Failed to unmarshal claim value for key %s: %s", field, err.Error())
			return zero, false
		}
		return out, true
	}
	return zero, false
}
