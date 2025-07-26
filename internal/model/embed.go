package model

import (
	"CBCTF/internal/config"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
)

type UintList []uint

func (u UintList) Value() (driver.Value, error) {
	if len(u) == 0 {
		return nil, nil
	}
	return json.Marshal(u)
}

func (u *UintList) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan UintList value")
	}
	if len(bytes) == 0 {
		*u = nil
		return nil
	}
	return json.Unmarshal(bytes, u)
}

type Int32List []int32

func (e Int32List) Value() (driver.Value, error) {
	tmp := make([]int32, 0)
	for _, port := range e {
		if port > 0 && port < 65535 {
			tmp = append(tmp, port)
		}
	}
	return json.Marshal(e)
}

func (e *Int32List) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan Int32List value")
	}
	return json.Unmarshal(bytes, e)
}

type AvatarURL string

func (a AvatarURL) Value() (driver.Value, error) {
	if a == "" {
		return nil, nil
	}
	return strings.TrimPrefix(string(a), config.Env.Backend), nil
}

func (a *AvatarURL) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan AvatarURL: %v", value)
	}
	if len(bytes) == 0 {
		*a = ""
		return nil
	}
	*a = AvatarURL(config.Env.Backend + string(bytes))
	return nil
}

type StringList []string

func (s StringList) Value() (driver.Value, error) {
	if len(s) == 0 {
		return nil, nil
	}
	return json.Marshal(s)
}

func (s *StringList) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan StringList value")
	}
	if len(bytes) == 0 {
		*s = nil
		return nil
	}
	return json.Unmarshal(bytes, s)
}

type StringMap map[string]string

func (s StringMap) Value() (driver.Value, error) {
	if len(s) == 0 {
		return nil, nil
	}
	return json.Marshal(s)
}

func (s *StringMap) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan StringMap value")
	}
	if len(bytes) == 0 {
		*s = nil
		return nil
	}
	return json.Unmarshal(bytes, s)
}
