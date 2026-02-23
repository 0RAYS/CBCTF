package utils

import (
	"fmt"
	"reflect"
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

// GetClaimValue extracts a value from an OAuth userinfo response map.
//
// The field string may contain one or more {key} or {nested.key} placeholders.
// When T is string, each placeholder is resolved and substituted back into the
// field template, matching the original behaviour.
// When T is any other type, the field must be a single bare placeholder and the
// resolved map value is returned via a direct type assertion to T.
func GetClaimValue[T any](resp map[string]any, field string) (T, bool) {
	var zero T
	keys := GetClaimKeys(field)

	// Non-string path: single placeholder, return the raw value as T.
	var isString bool
	switch any(zero).(type) {
	case string:
		isString = true
	}

	if !isString {
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
			if i == len(ks)-1 {
				typed, ok := v.(T)
				if !ok {
					return zero, false
				}
				return typed, true
			}
			subData, ok := v.(map[string]any)
			if !ok {
				return zero, false
			}
			data = subData
		}
		return zero, false
	}

	// String path: substitute all placeholders into the template.
	result := field
	for _, key := range keys {
		ks := strings.Split(key, ".")
		data := resp
		for i, k := range ks {
			v, ok := data[k]
			if !ok {
				return zero, false
			}
			if i == len(ks)-1 {
				var s string
				if v == nil {
					s = "<nil>"
				} else if reflect.TypeOf(v).Kind() == reflect.Float64 {
					s = fmt.Sprintf("%f", v)
				} else {
					s = fmt.Sprintf("%v", v)
				}
				result = strings.ReplaceAll(result, fmt.Sprintf("{%s}", key), s)
			} else {
				subData, ok := v.(map[string]any)
				if !ok {
					return zero, false
				}
				data = subData
			}
		}
	}
	return any(result).(T), true
}
