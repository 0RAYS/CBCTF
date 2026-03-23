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
