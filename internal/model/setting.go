package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

const (
	AsyncQLogLevelSettingKey    = "asyncq.log.level"
	AsyncQConcurrencySettingKey = "asyncq.concurrency"

	GinRateLimitGlobalSettingKey    = "gin.ratelimit.global"
	GinRateLimitWhitelistSettingKey = "gin.ratelimit.whitelist"
	GinCORSSettingKey               = "gin.cors"
	GinLogWhitelistSettingKey       = "gin.log.whitelist"
)

type Setting struct {
	Key   string       `gorm:"size:50;uniqueIndex" json:"key"`
	Value SettingValue `gorm:"type:json;size:255" json:"value"`
	BaseModel
}

func (s Setting) ModelName() string {
	return "Setting"
}

func (s Setting) GetBaseModel() BaseModel {
	return s.BaseModel
}

func (s Setting) UniqueFields() []string {
	return []string{"key"}
}

func (s Setting) QueryFields() []string {
	return []string{"id", "key", "value"}
}

type SettingValue struct {
	V any
}

func (s SettingValue) Value() (driver.Value, error) {
	if s.V == nil {
		return nil, nil
	}
	return json.Marshal(s.V)
}

func (s *SettingValue) Scan(value any) error {
	if value == nil {
		s.V = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan SettingValue value")
	}
	if len(bytes) == 0 {
		s.V = nil
		return nil
	}
	var data any
	if err := json.Unmarshal(bytes, &data); err != nil {
		return err
	}
	s.V = data
	return nil
}
