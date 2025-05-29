package model

const (
	NoticeTypeNormal    = "normal"
	NoticeTypeImportant = "important"
	NoticeTypeUpdate    = "update"
)

type Notice struct {
	ContestID uint    `json:"contest_id"`
	Contest   Contest `json:"-"`
	AdminID   uint    `json:"creator_id"`
	Admin     Admin   `json:"-"`
	Title     string  `json:"title"`
	Content   string  `json:"content"`
	Type      string  `gorm:"default:'normal'" json:"type"`
	BaseModel
}
