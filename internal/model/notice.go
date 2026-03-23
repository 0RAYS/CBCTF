package model

type NoticeType string

const (
	NoticeNormalType    NoticeType = "normal"
	NoticeImportantType NoticeType = "important"
	NoticeUpdateType    NoticeType = "update"
)

// Notice
// BelongsTo Contest
type Notice struct {
	ContestID uint   `json:"contest_id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Type      string `gorm:"default:'normal'" json:"type"`
	BaseModel
}
