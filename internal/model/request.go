package model

import (
	"database/sql"
	"time"
)

// Request
// BelongsTo Device
type Request struct {
	IP        string         `gorm:"index:idx_requests_user_ip,priority:2;index" json:"ip"`
	Time      time.Time      `gorm:"default:null" json:"time"`
	Method    string         `json:"method"`
	Path      string         `json:"path"`
	URL       string         `json:"url"`
	UserAgent string         `json:"user_agent"`
	Status    int            `json:"status"`
	Referer   string         `json:"referer"`
	Magic     string         `json:"magic"`
	UserID    sql.Null[uint] `gorm:"index:idx_requests_user_ip,priority:1;index" json:"user_id"`
	BaseModel
}
