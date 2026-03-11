package model

import (
	"database/sql"
	"time"
)

// Request
// BelongsTo Device
type Request struct {
	IP        string         `json:"ip"`
	Time      time.Time      `gorm:"default:null" json:"time"`
	Method    string         `json:"method"`
	Path      string         `json:"path"`
	URL       string         `json:"url"`
	UserAgent string         `json:"user_agent"`
	Status    int            `json:"status"`
	Referer   string         `json:"referer"`
	Magic     string         `json:"magic"`
	UserID    sql.Null[uint] `json:"user_id"`
	BaseModel
}

func (r Request) TableName() string {
	return "requests"
}

func (r Request) ModelName() string {
	return "Request"
}

func (r Request) GetBaseModel() BaseModel {
	return r.BaseModel
}

func (r Request) UniqueFields() []string {
	return []string{"id"}
}

func (r Request) QueryFields() []string {
	return []string{"id", "ip", "user_agent", "user_id", "method", "path", "status", "magic"}
}
