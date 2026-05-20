package model

// Notice
// BelongsTo Contest
type Notice struct {
	ContestID uint   `gorm:"index" json:"contest_id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Type      string `gorm:"default:'normal'" json:"type"`
	BaseModel
}
