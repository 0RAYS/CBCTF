package model

import (
	"database/sql"
	"time"
)

// Request
// BelongsTo Device
type Request struct {
	IP        string         `json:"ip"`
	Time      time.Time      `json:"time"`
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

func (r Request) GetModelName() string {
	return "Request"
}

func (r Request) GetBaseModel() BaseModel {
	return r.BaseModel
}

func (r Request) GetUniqueKey() []string {
	return []string{"id"}
}

func (r Request) GetAllowedQueryFields() []string {
	return []string{"id", "ip", "user_agent", "user_id"}
}
