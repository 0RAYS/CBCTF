package model

import (
	"CBCTF/internel/i18n"
	"time"
)

// Request
// BelongsTo Device
type Request struct {
	IP        string    `gorm:"size:45;not null" json:"ip"`
	Time      time.Time `gorm:"not null" json:"time"`
	Method    string    `gorm:"size:10;not null" json:"method"`
	Path      string    `gorm:"size:255;not null" json:"path"`
	URL       string    `gorm:"size:255;not null" json:"url"`
	UserAgent string    `gorm:"size:255;not null" json:"user_agent"`
	Status    int       `gorm:"not null" json:"status"`
	Referer   string    `gorm:"size:255" json:"referer"`
	Magic     string    `json:"magic"`
	Basic
}

func (r Request) GetModelName() string {
	return "Request"
}

func (r Request) GetID() uint {
	return r.ID
}

func (r Request) GetVersion() uint {
	return r.Version
}

func (r Request) CreateErrorString() string {
	return i18n.CreateRequestError
}

func (r Request) DeleteErrorString() string {
	return i18n.DeleteRequestError
}

func (r Request) GetErrorString() string {
	return i18n.GetRequestError
}

func (r Request) NotFoundErrorString() string {
	return i18n.RequestNotFound
}

func (r Request) UpdateErrorString() string {
	return i18n.UpdateRequestError
}

func (r Request) GetUniqueKey() []string {
	return []string{"id"}
}
