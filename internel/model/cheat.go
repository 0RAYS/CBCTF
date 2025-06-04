package model

const (
	Suspect string = "suspect"
	Cheater string = "cheater"
	None    string = "none"

	MagicNotMatch string = "Magic in token and request headers are not matched"
	SameFlag      string = "Flag is same"
)

type Cheat struct {
	ID        string     `gorm:"type:varchar(36);primaryKey" json:"id"`
	UserID    uint       `json:"user_id"`
	User      User       `json:"-"`
	TeamID    uint       `json:"team_id"`
	Team      Team       `json:"-"`
	ContestID uint       `json:"contest_id"`
	Contest   Contest    `json:"-"`
	Reason    string     `json:"reason"`
	Type      string     `json:"type"`
	Checked   bool       `json:"checked"`
	Cheated   bool       `json:"cheated"`
	Cheats    StringList `gorm:"type:json" json:"cheats"`
	BaseModel
}
