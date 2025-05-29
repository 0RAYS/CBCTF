package model

type User struct {
	Teams       []*Team      `gorm:"many2many:user_teams;" json:"-"`
	Contests    []*Contest   `gorm:"many2many:user_contests;" json:"-"`
	Submissions []Submission `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Victims     []Victim     `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Devices     []Device     `json:"-"`
	Cheats      []Cheat      `json:"-"`
	Name        string       `gorm:"type:varchar(255);uniqueIndex;not null" json:"name"`
	Password    string       `gorm:"not null" json:"-"`
	Email       string       `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Country     string       `gorm:"default:'CN'" json:"country"`
	Avatar      string       `json:"avatar"`
	Desc        string       `json:"desc"`
	Verified    bool         `gorm:"default:false" json:"verified"`
	Hidden      bool         `gorm:"default:false" json:"hidden"`
	Banned      bool         `gorm:"default:false" json:"banned"`
	Score       float64      `gorm:"default:0" json:"score"`
	Solved      int64        `gorm:"default:0" json:"solved"`
	BaseModel
}
