package model

import (
	"encoding/json"
	"fmt"
)

func scanBytes(value any) ([]byte, error) {
	switch v := value.(type) {
	case nil:
		return nil, nil
	case []byte:
		return v, nil
	case string:
		return []byte(v), nil
	default:
		return nil, fmt.Errorf("failed to scan value of type %T", value)
	}
}

func scanJSON[T any](value any, dest *T) error {
	bytes, err := scanBytes(value)
	if err != nil {
		return err
	}
	if len(bytes) == 0 {
		var zero T
		*dest = zero
		return nil
	}
	return json.Unmarshal(bytes, dest)
}
