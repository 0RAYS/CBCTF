package model

import (
	"CBCTF/internel/config"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
)

type AvatarURL string

func (a AvatarURL) Value() (driver.Value, error) {
	if a == "" {
		return nil, nil
	}
	return strings.TrimPrefix(string(a), strings.Trim(config.Env.Backend, "/")), nil
}

func (a *AvatarURL) Scan(value any) error {
	if value == nil || value.(string) == "" {
		*a = ""
		return nil
	}
	path, ok := value.(string)
	if !ok {
		return fmt.Errorf("failed to scan AvatarURL: %v", value)
	}
	*a = AvatarURL(strings.Trim(config.Env.Backend, "/") + path)
	return nil
}

type StringList []string

func (s StringList) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *StringList) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan StringList value")
	}
	return json.Unmarshal(bytes, s)
}

type UintList []uint

func (u UintList) Value() (driver.Value, error) {
	return json.Marshal(u)
}

func (u *UintList) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan UintList value")
	}
	return json.Unmarshal(bytes, u)
}
