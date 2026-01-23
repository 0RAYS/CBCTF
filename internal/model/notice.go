package model

import "CBCTF/internal/i18n"

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
	BaseModel
}

func (n Notice) GetModelName() string {
	return "Notice"
}

func (n Notice) GetBaseModel() BaseModel {
	return n.BaseModel
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

func (n Notice) NotFoundString() string {
	return i18n.NoticeNotFound
}

func (n Notice) UpdateErrorString() string {
	return i18n.UpdateNoticeError
}

func (n Notice) GetUniqueKey() []string {
	return []string{"id"}
}

func (n Notice) GetAllowedQueryFields() []string {
	return []string{"id", "title", "content"}
}
