package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

const (
	AsyncQLogLevelSettingKey    = "asyncq.log.level"
	AsyncQConcurrencySettingKey = "asyncq.concurrency"

	GinRateLimitGlobalSettingKey           = "gin.ratelimit.global"
	GinRateLimitWhitelistSettingKey string = "gin.ratelimit.whitelist"
	GinLogWhitelistSettingKey       string = "gin.log.whitelist"
)

var DefaultSettings = []Setting{
	{Key: AsyncQLogLevelSettingKey, Value: SettingValue{V: "INFO"}},
	{Key: AsyncQConcurrencySettingKey, Value: SettingValue{V: 50}},

	{Key: GinRateLimitGlobalSettingKey, Value: SettingValue{V: 120}},
	{Key: GinRateLimitWhitelistSettingKey, Value: SettingValue{V: []string{"127.0.0.1"}}},
	{Key: GinLogWhitelistSettingKey, Value: SettingValue{V: []string{"/metrics"}}},
}

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
