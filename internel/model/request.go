package model

import (
	"CBCTF/internel/i18n"
	"time"
)

// Request
// BelongsTo Device
type Request struct {
	IP        string    `json:"ip"`
	Time      time.Time `json:"time"`
	Method    string    `json:"method"`
	Path      string    `json:"path"`
	URL       string    `json:"url"`
	UserAgent string    `json:"user_agent"`
	Status    int       `json:"status"`
	Referer   string    `json:"referer"`
	Magic     string    `json:"magic"`
	BasicModel
}

func (r Request) GetModelName() string {
	return "Request"
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

func (r Request) GetForeignKeys() []string {
	return []string{"id"}
}
