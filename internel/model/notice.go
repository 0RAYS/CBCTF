package model

import "CBCTF/internel/i18n"

const (
	NoticeTypeNormal    = "normal"
	NoticeTypeImportant = "important"
	NoticeTypeUpdate    = "update"
)

// Notice
// BelongsTo Contest
type Notice struct {
	ContestID uint    `json:"contest_id"`
	Contest   Contest `json:"-"`
	Title     string  `json:"title"`
	Content   string  `json:"content"`
	Type      string  `gorm:"default:'normal'" json:"type"`
	Basic
}

func (n Notice) GetModelName() string {
	return "Notice"
}

func (n Notice) GetID() uint {
	return n.ID
}

func (n Notice) GetVersion() uint {
	return n.Version
}

func (n Notice) CreateErrorString() string {
	return i18n.CreateNoticeError
}

func (n Notice) DeleteErrorString() string {
	return i18n.DeleteNoticeError
}

func (n Notice) GetErrorString() string {
	return i18n.GetNoticeError
}

func (n Notice) NotFoundErrorString() string {
	return i18n.NoticeNotFound
}

func (n Notice) UpdateErrorString() string {
	return i18n.UpdateNoticeError
}

func (n Notice) GetUniqueKey() []string {
	return []string{"id"}
}
