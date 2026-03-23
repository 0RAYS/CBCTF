package model

import (
	"CBCTF/internal/config"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
)

type FileURL string

func (a FileURL) Value() (driver.Value, error) {
	if a == "" {
		return nil, nil
	}
	return strings.TrimPrefix(string(a), config.Env.Host), nil
}

func (a *FileURL) Scan(value any) error {
	bytes, err := scanBytes(value)
	if err != nil {
		return fmt.Errorf("failed to scan FileURL: %v", value)
	}
	if len(bytes) == 0 {
		*a = ""
		return nil
	}
	if strings.HasPrefix(string(bytes), "https://") || strings.HasPrefix(string(bytes), "http://") {
		*a = FileURL(bytes)
	} else {
		*a = FileURL(config.Env.Host + string(bytes))
	}
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
	bytes, err := scanBytes(value)
	if err != nil {
		return fmt.Errorf("failed to scan StringList value")
	}
	if len(bytes) == 0 {
		*s = nil
		return nil
	}
	return json.Unmarshal(bytes, s)
}

type UintMap map[string]uint

func (u UintMap) Value() (driver.Value, error) {
	if len(u) == 0 {
		return nil, nil
	}
	return json.Marshal(u)
}

func (u *UintMap) Scan(value any) error {
	bytes, err := scanBytes(value)
	if err != nil {
		return fmt.Errorf("failed to scan UintMap value")
	}
	if len(bytes) == 0 {
		*u = nil
		return nil
	}
	return json.Unmarshal(bytes, u)
}

type StringMap map[string]string

func (s StringMap) Value() (driver.Value, error) {
	if len(s) == 0 {
		return nil, nil
	}
	return json.Marshal(s)
}

func (s *StringMap) Scan(value any) error {
	bytes, err := scanBytes(value)
	if err != nil {
		return fmt.Errorf("failed to scan StringMap value")
	}
	if len(bytes) == 0 {
		*s = nil
		return nil
	}
	return json.Unmarshal(bytes, s)
}
