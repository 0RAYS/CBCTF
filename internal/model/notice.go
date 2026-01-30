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

func (n Notice) GetModelName() string {
	return "Notice"
}

func (n Notice) GetBaseModel() BaseModel {
	return n.BaseModel
}

func (n Notice) GetUniqueKey() []string {
	return []string{"id"}
}

func (n Notice) GetAllowedQueryFields() []string {
	return []string{"id", "title", "content"}
}
