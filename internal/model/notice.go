package model

const (
	NoticeTypeNormal    = "normal"
	NoticeTypeImportant = "important"
	NoticeTypeUpdate    = "update"
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

func (n Notice) ModelName() string {
	return "Notice"
}

func (n Notice) GetBaseModel() BaseModel {
	return n.BaseModel
}

func (n Notice) UniqueFields() []string {
	return []string{"id"}
}

func (n Notice) QueryFields() []string {
	return []string{"id", "title", "content", "type", "contest_id"}
}
