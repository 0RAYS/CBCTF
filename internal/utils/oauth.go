package utils

import (
	"fmt"
	"regexp"
	"strings"
)

func GetFiledKeys(field string) []string {
	re := regexp.MustCompile(`\{([^{}]*)\}`)
	matches := re.FindAllStringSubmatch(field, -1)

	var results []string
	for _, match := range matches {
		if len(match) > 1 {
			results = append(results, match[1])
		}
	}
	return results
}

func GetFiledValue(resp map[string]any, field string) (string, bool) {
	keys := GetFiledKeys(field)
	for _, key := range keys {
		ks := strings.Split(key, ".")
		data := resp
		for i, k := range ks {
			v, ok := data[k]
			if !ok {
				return "", false
			}
			if i == len(ks)-1 {
				v = fmt.Sprintf("%v", v)
				field = strings.ReplaceAll(field, fmt.Sprintf("{%s}", key), v.(string))
			} else {
				if subData, ok := v.(map[string]any); ok {
					data = subData
				} else {
					return "", false
				}
			}
		}
	}
	return field, true
}
