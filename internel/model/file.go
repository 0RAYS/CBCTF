package model

const (
	Avatar  = "avatar"
	WriteUP = "writeup"
)

type File struct {
	ID        string `gorm:"type:varchar(36);primarykey" json:"id"`
	Filename  string `json:"filename"`
	Size      int64  `json:"size"`
	Path      string `json:"-"`
	AdminID   uint   `json:"admin_id"`
	UserID    uint   `json:"user_id"`
	TeamID    uint   `json:"team_id"`
	ContestID uint   `json:"contest_id"`
	Suffix    string `json:"suffix"`
	Hash      string `json:"hash"`
	Type      string `json:"type"`
	BaseModel
}
